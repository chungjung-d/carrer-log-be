# Career Log 백엔드 API

이 프로젝트는 Career Log 애플리케이션의 백엔드 API를 제공합니다. Go 언어와 Fiber 웹 프레임워크를 사용하여 구현되었으며, PostgreSQL 데이터베이스를 사용합니다.

## 기술 스택

- **언어**: Go
- **웹 프레임워크**: Fiber
- **ORM**: GORM
- **데이터베이스**: PostgreSQL
- **인증**: JWT
- **컨테이너화**: Docker

## 주요 기능

- 사용자 인증 (회원가입, 로그인)
- JWT 기반 인증 시스템
- 구조화된 에러 처리
- 환경 설정 관리

## 프로젝트 구조

```
career-log-be/
├── config/                 # 설정 관련 코드
│   ├── database/           # 데이터베이스 설정
│   └── env/                # 환경 변수 관리
│       └── provider/       # 환경 변수 제공자
├── errors/                 # 에러 처리 시스템
├── middleware/             # 미들웨어
├── models/                 # 데이터 모델
├── routes/                 # 라우터
│   └── v1/                 # API v1
│       ├── auth/           # 인증 관련 라우트
│       └── ...             # 기타 라우트
├── services/               # 비즈니스 로직
│   ├── auth/               # 인증 관련 서비스
│   └── ...                 # 기타 서비스
├── utils/                  # 유틸리티 함수
│   └── jwt/                # JWT 관련 유틸리티
├── volumes/                # 도커 볼륨 (gitignore)
├── docker-compose.yml      # 도커 컴포즈 설정
├── main.go                 # 애플리케이션 진입점
└── README.md               # 프로젝트 문서
```

## 핵심 컴포넌트

### 1. 환경 변수 관리

환경 변수 관리를 위한 추상화된 인터페이스를 제공합니다. 이를 통해 다양한 환경(개발, 테스트, 프로덕션)에서 설정을 쉽게 관리할 수 있습니다.

```go
// config/env/provider/provider.go
type EnvProvider interface {
    GetString(key string) string
    GetInt(key string) int
    GetBool(key string) bool
}
```

### 2. 데이터베이스 연결

GORM을 사용하여 PostgreSQL 데이터베이스에 연결합니다.

```go
// config/database/postgres.go
func NewDatabase(config Config) (*gorm.DB, error) {
    // 데이터베이스 연결 설정
}
```

### 3. 사용자 인증

JWT 기반 인증 시스템을 구현하여 안전한 API 접근을 제공합니다.

```go
// utils/jwt/jwt.go
type JWTUtils struct {
    config *env.JWTConfig
}

func (j *JWTUtils) GenerateToken(userID, email string) (string, error) {
    // JWT 토큰 생성
}

func (j *JWTUtils) ValidateToken(tokenString string) (*UserClaims, error) {
    // JWT 토큰 검증
}
```

### 4. 미들웨어

다양한 미들웨어를 통해 요청 처리 파이프라인을 구성합니다.

- **데이터베이스 미들웨어**: 요청 컨텍스트에 DB 인스턴스 제공
- **JWT 미들웨어**: 요청 컨텍스트에 JWT 유틸리티 제공
- **인증 미들웨어**: JWT 토큰 검증 및 사용자 정보 설정
- **에러 핸들러**: 구조화된 에러 응답 제공

### 5. 에러 처리 시스템

일관된 에러 처리를 위한 시스템을 제공합니다.

```go
// errors/error.go
type AppError struct {
    Type      ErrorType
    Code      ErrorCode
    Message   string
    Details   any
    DebugInfo string
    Err       error
}
```

에러 타입에 따라 적절한 HTTP 상태 코드와 응답 형식을 자동으로 생성합니다.

### 6. ID 생성 시스템

사용자 및 리소스를 위한 커스텀 ID 생성 시스템을 제공합니다.

```go
// utils/id_generator.go
func GenerateID(prefix string) string {
    // PREFIX_TIMESTAMP_RANDOMSEED 형식의 ID 생성
}
```

## API 엔드포인트

API 명세에 대한 자세한 내용은 [API 문서](docs/API.md)를 참조하세요.

## 설치 및 실행

### 필수 조건

- Go 1.16 이상
- Docker 및 Docker Compose

### 설치

1. 저장소 클론:
```bash
git clone https://github.com/your-username/career-log-be.git
cd career-log-be
```

2. 의존성 설치:
```bash
go mod download
```

### 실행

1. 데이터베이스 실행:
```bash
docker-compose up -d
```

2. 애플리케이션 실행:
```bash
go run main.go
```

서버는 기본적으로 `http://localhost:3000`에서 실행됩니다.

## 개발 가이드

### 새 API 엔드포인트 추가

1. 모델 정의 (`models/` 디렉토리)
2. 서비스 로직 구현 (`services/` 디렉토리)
3. 라우터 설정 (`routes/` 디렉토리)

### 인증이 필요한 엔드포인트 추가

인증이 필요한 엔드포인트는 `AuthMiddleware`를 사용하여 보호할 수 있습니다:

```go
// 보호된 라우트 그룹 예시
protected := router.Group("/protected")
protected.Use(middleware.AuthMiddleware())

// 이 엔드포인트는 인증이 필요합니다
protected.Get("/resource", handleGetResource())
```

### 에러 처리

에러 처리 시스템을 사용하여 일관된 에러 응답을 제공하세요:

```go
// 에러 반환 예시
if err != nil {
    return appErrors.NewInternalError(
        appErrors.ErrorCodeDatabaseError,
        "Failed to query database",
        err,
    )
}
```

### 핸들러 구조

프로젝트는 도메인 중심의 구조를 따르며, 라우터와 핸들러를 명확하게 분리합니다:

#### 라우터-서비스 구조
```go
routes/
└── v1/
    ├── routes.go          # 메인 라우터 설정
    ├── auth/
    │   └── routers.go     # 인증 관련 라우터 그룹
    └── user/
        └── routers.go     # 사용자 관련 라우터 그룹

services/
└── user/
    ├── create_job_satisfaction_importance.go  # 직무 만족도 중요도 생성 핸들러
    └── create_user_profile.go                 # 사용자 프로필 생성 핸들러
```

- 각 도메인별로 라우터를 그룹화하여 `routes/v1/{도메인}/routers.go` 파일에서 관리
- 각 라우터에 연결되는 핸들러는 `services/{도메인}` 디렉토리에서 1핸들러 1파일 원칙으로 관리
- 이를 통해 코드의 관심사 분리와 유지보수성 향상

## 라이센스

이 프로젝트는 MIT 라이센스 하에 배포됩니다. 