package scheduler

import (
	"career-log-be/models/job_satisfaction"
	"career-log-be/models/job_satisfaction/enums"
	"career-log-be/models/note/chat"
	"career-log-be/utils/chatgpt"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gofiber/fiber/v2"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

type ChatAnalyzeScheduler struct {
	scheduler *gocron.Scheduler
	db        *gorm.DB
	chatGPT   *chatgpt.Service
}

// ChatGPTAnalysisResponse ChatGPT 응답을 파싱하기 위한 구조체
type ChatGPTAnalysisResponse struct {
	Workload          float64 `json:"workload"`
	Compensation      float64 `json:"compensation"`
	Growth            float64 `json:"growth"`
	WorkEnvironment   float64 `json:"workEnvironment"`
	WorkRelationships float64 `json:"workRelationships"`
	WorkValues        float64 `json:"workValues"`
}

// getAnalysisPrompt 분석을 위한 프롬프트를 반환합니다
func getAnalysisPrompt() string {
	return `당신은 현재 상담자와 내담자의 대화를 분석하여 내담자의 직무 만족도를 평가하는 역할을 합니다.
내담자의 대화 내용을 바탕으로 다음 6가지 항목에 대한 만족도를 평가해주세요.

평가 항목:
1. workload: 업무량과 업무에서의 성취감
2. compensation: 회사에서의 금전적인 보상
3. growth: 회사에서의 커리어나 내면적 성장
4. workEnvironment: 회사의 워라벨
5. workRelationships: 회사 내 동료들과의 관계
6. workValues: 회사에서의 업무의 가치와 개인의 삶의 방향성

평가 방법:
- 각 항목에 대해 -10에서 +10 사이의 점수를 매깁니다
- 0: 언급되지 않았거나 중립적
- 양수: 긍정적인 경험 (최대 +10)
- 음수: 부정적인 경험 (최소 -10)

응답 형식:
다음과 같은 JSON 형식으로 응답해주세요:
{
    "workload": 0,
    "compensation": 0,
    "growth": 0,
    "workEnvironment": 0,
    "workRelationships": 0,
    "workValues": 0
}

주의사항:
- 상담자의 답변은 평가에 반영하지 않습니다
- 명확한 언급이 없는 항목은 0점으로 처리합니다
- 감정의 강도에 따라 적절한 점수를 배분합니다

다음은 분석할 대화 내용입니다:`
}

// analyzeChat 채팅 내용을 분석하여 JobSatisfactionUpdateEvent를 생성합니다
func (cs *ChatAnalyzeScheduler) analyzeChat(ctx context.Context, chatSet *chat.ChatSet) (*job_satisfaction.JobSatisfactionUpdateEvent, error) {
	// 대화 내용 구성
	var conversation string
	for _, msg := range chatSet.ChatData.Messages {
		conversation += fmt.Sprintf("%s: %s\n", msg.Role, msg.Content)
	}

	// ChatGPT 요청
	messages := []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: getAnalysisPrompt(),
		},
		{
			Role:    "user",
			Content: conversation,
		},
	}

	response, err := cs.chatGPT.CompleteChatRequest(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to get ChatGPT response: %v", err)
	}

	// 응답 파싱
	var analysis ChatGPTAnalysisResponse
	if err := json.Unmarshal([]byte(response), &analysis); err != nil {
		return nil, fmt.Errorf("failed to parse ChatGPT response: %v", err)
	}

	// JobSatisfactionUpdateEvent 생성
	kst, _ := time.LoadLocation("Asia/Seoul")
	event := &job_satisfaction.JobSatisfactionUpdateEvent{
		UserID:            chatSet.UserID,
		EventType:         enums.ChatAnalysisEvent,
		Workload:          int(analysis.Workload),
		Compensation:      int(analysis.Compensation),
		Growth:            int(analysis.Growth),
		WorkEnvironment:   int(analysis.WorkEnvironment),
		WorkRelationships: int(analysis.WorkRelationships),
		WorkValues:        int(analysis.WorkValues),
		SourceId:          &chatSet.ID,
		CreatedAt:         time.Now().In(kst),
	}

	return event, nil
}

// NewChatAnalyzeScheduler 새로운 ChatAnalyzeScheduler 인스턴스를 생성합니다
func NewChatAnalyzeScheduler(db *gorm.DB) (*ChatAnalyzeScheduler, error) {
	chatGPTService, err := chatgpt.NewChatGPTBuilder().Build()

	if err != nil {
		return nil, fmt.Errorf("failed to create ChatGPT service: %v", err)
	}

	kst, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		return nil, fmt.Errorf("failed to load KST timezone: %v", err)
	}

	return &ChatAnalyzeScheduler{
		scheduler: gocron.NewScheduler(kst),
		db:        db,
		chatGPT:   chatGPTService,
	}, nil
}

// Start 스케줄러를 시작합니다
func (cs *ChatAnalyzeScheduler) Start() {
	// 매일 오전 7시(UTC 기준)에 실행
	_, err := cs.scheduler.Every(1).Day().At("07:00").Do(cs.AnalyzeDailyChat)
	if err != nil {
		log.Printf("Failed to schedule daily chat count reset: %v", err)
	}

	cs.scheduler.StartAsync()
}

// Stop 스케줄러를 중지합니다
func (cs *ChatAnalyzeScheduler) Stop() {
	cs.scheduler.Stop()
}

func (cs *ChatAnalyzeScheduler) AnalyzeDailyChat() {
	kst, _ := time.LoadLocation("Asia/Seoul")
	now := time.Now().In(kst)
	// 전날 오전 6시부터
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 6, 0, 0, 0, kst).AddDate(0, 0, -1)
	// 당일 오전 6시까지
	endOfDay := startOfDay.Add(24 * time.Hour)

	// 해당 기간의 ChatSet 조회
	var chatSets []chat.ChatSet
	if err := cs.db.Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay).Find(&chatSets).Error; err != nil {
		log.Printf("Failed to retrieve chat sets: %v", err)
		return
	}

	for _, chatSet := range chatSets {
		event, err := cs.analyzeChat(context.Background(), &chatSet)
		if err != nil {
			log.Printf("Failed to analyze chat %s: %v", chatSet.ID, err)
			continue
		}

		// DB에 저장
		if err := cs.db.Create(event).Error; err != nil {
			log.Printf("Failed to save analysis result for chat %s: %v", chatSet.ID, err)
			continue
		}

		log.Printf("Successfully analyzed and saved result for chat %s", chatSet.ID)
	}

	log.Println("Daily chat analysis has been completed")
}

// InitChatAnalyzeScheduler Fiber 앱에 스케줄러를 초기화하고 등록하는 함수
func InitChatAnalyzeScheduler(app *fiber.App, db *gorm.DB) error {
	chatScheduler, err := NewChatAnalyzeScheduler(db)
	if err != nil {
		return err
	}

	chatScheduler.Start()

	// Fiber 앱이 종료될 때 스케줄러도 함께 종료
	app.Hooks().OnShutdown(func() error {
		chatScheduler.Stop()
		return nil
	})

	return nil
}
