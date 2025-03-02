# Career Log API 문서

## 기본 정보
- 기본 URL: `http://localhost:3000`
- API 버전: v1
- 인증: JWT 토큰 기반 인증
- 요청/응답 형식: JSON

## 인증 (Authentication)

### 회원가입
- **엔드포인트**: `POST /api/v1/auth/register`
- **설명**: 새로운 사용자 계정을 생성합니다.
- **인증 필요**: 아니오
- **요청 본문**:
  ```json
  {
    "email": "string",
    "password": "string"
  }
  ```
- **응답**:
  ```json
  {
    "status": "success",
    "message": "User registered successfully"
  }
  ```

### 로그인
- **엔드포인트**: `POST /api/v1/auth/login`
- **설명**: 사용자 인증 및 JWT 토큰 발급
- **인증 필요**: 아니오
- **요청 본문**:
  ```json
  {
    "email": "string",
    "password": "string"
  }
  ```
- **응답**:
  ```json
  {
    "status": "success",
    "data": {
      "token": "JWT_TOKEN_STRING"
    }
  }
  ```

## 사용자 프로필 (User Profile)

### 프로필 생성
- **엔드포인트**: `POST /api/v1/user/profile`
- **설명**: 사용자 프로필 정보를 생성합니다.
- **인증 필요**: 예
- **헤더**:
  ```
  Authorization: Bearer {JWT_TOKEN}
  ```
- **요청 본문**:
  ```json
  {
    "name": "string",
    "jobTitle": "string",
    "company": "string"
  }
  ```
- **응답**:
  ```json
  {
    "status": "success",
    "data": {
      "profile": {
        "id": "string",
        "name": "string",
        "jobTitle": "string",
        "company": "string",
        "createdAt": "timestamp",
        "updatedAt": "timestamp"
      }
    }
  }
  ```

## 직무 만족도 (Job Satisfaction)

### 직무 만족도 중요도 설정
- **엔드포인트**: `POST /api/v1/job-satisfaction/importance`
- **설명**: 사용자의 직무 만족도 요소별 중요도를 설정합니다.
- **인증 필요**: 예
- **헤더**:
  ```
  Authorization: Bearer {JWT_TOKEN}
  ```
- **요청 본문**:
  ```json
  {
    "factors": [
      {
        "factor": "string",
        "importance": number
      }
    ]
  }
  ```
- **응답**:
  ```json
  {
    "status": "success",
    "message": "Job satisfaction importance created successfully"
  }
  ```

### 직무 만족도 초기화
- **엔드포인트**: `POST /api/v1/job-satisfaction/init`
- **설명**: 사용자의 직무 만족도 데이터를 초기화합니다.
- **인증 필요**: 예
- **헤더**:
  ```
  Authorization: Bearer {JWT_TOKEN}
  ```
- **응답**:
  ```json
  {
    "status": "success",
    "message": "Job satisfaction initialized successfully"
  }
  ```

## 에러 응답
모든 API 엔드포인트는 에러 발생 시 다음과 같은 형식으로 응답합니다:
```json
{
  "status": "error",
  "message": "에러 메시지",
  "error": {
    "code": "에러_코드",
    "details": "상세 에러 정보"
  }
}
```

## 인증 헤더
보호된 엔드포인트에 접근할 때는 반드시 다음 형식의 인증 헤더를 포함해야 합니다:
```
Authorization: Bearer {JWT_TOKEN}
```
여기서 `{JWT_TOKEN}`은 로그인 시 발급받은 JWT 토큰입니다. 