package utils

import (
	"math"

	job_satisfaction "career-log-be/models/job_satisfaction"
)

// 로지스틱 함수 상수
const (
	// L: 로지스틱 함수의 최대값
	L = 100.0
	// k: 로지스틱 함수의 기울기 (클수록 가파름)
	k = 0.05
	// x0: 로지스틱 함수의 중간점
	x0 = 50.0
)

// CalculateWeightedScore는 만족도와 중요도를 가중하여 점수를 계산합니다.
// sum{ (중요도/100)*(만족도/100) } / sum(중요도/100) * 100
func CalculateWeightedScore(satisfaction *job_satisfaction.UserJobSatisfaction) float64 {
	// 각 항목별 가중 점수 계산
	weightedSum := 0.0
	importanceSum := 0.0

	// 워크로드
	workloadImportance := float64(satisfaction.WorkloadImportance) / 100.0
	weightedSum += workloadImportance * (float64(satisfaction.Workload) / 100.0)
	importanceSum += workloadImportance

	// 보상
	compensationImportance := float64(satisfaction.CompensationImportance) / 100.0
	weightedSum += compensationImportance * (float64(satisfaction.Compensation) / 100.0)
	importanceSum += compensationImportance

	// 성장
	growthImportance := float64(satisfaction.GrowthImportance) / 100.0
	weightedSum += growthImportance * (float64(satisfaction.Growth) / 100.0)
	importanceSum += growthImportance

	// 근무 환경
	workEnvironmentImportance := float64(satisfaction.WorkEnvironmentImportance) / 100.0
	weightedSum += workEnvironmentImportance * (float64(satisfaction.WorkEnvironment) / 100.0)
	importanceSum += workEnvironmentImportance

	// 업무 관계
	workRelationshipsImportance := float64(satisfaction.WorkRelationshipsImportance) / 100.0
	weightedSum += workRelationshipsImportance * (float64(satisfaction.WorkRelationships) / 100.0)
	importanceSum += workRelationshipsImportance

	// 업무 가치
	workValuesImportance := float64(satisfaction.WorkValuesImportance) / 100.0
	weightedSum += workValuesImportance * (float64(satisfaction.WorkValues) / 100.0)
	importanceSum += workValuesImportance

	// 중요도 총합이 0이면 0 반환
	if importanceSum == 0 {
		return 0
	}

	// 가중 평균 계산
	weightedAverage := (weightedSum / importanceSum) * 100.0

	// 로지스틱 함수 적용 (높은 점수일수록 올리기 어려움)
	return ApplyLogisticFunction(weightedAverage)
}

// ApplyLogisticFunction은 로지스틱 함수를 적용하여 점수를 조정합니다.
// f(x) = L / (1 + e^(-k*(x-x0)))
func ApplyLogisticFunction(score float64) float64 {
	// 로지스틱 함수 적용
	return L / (1 + math.Exp(-k*(score-x0)))
}
