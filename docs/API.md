# Career Log API 명세서

## 개요
이 문서는 Career Log 애플리케이션의 API 명세를 제공합니다.

## 기본 정보
- 기본 URL: `http://localhost:3000/api/v1`
- 모든 요청/응답은 `application/json` 형식을 사용합니다
- 인증이 필요한 엔드포인트는 요청 헤더에 `Authorization: Bearer {token}` 형식으로 JWT 토큰을 포함해야 합니다

## 공통 응답 형식

### 성공 응답
```json
{
  "success": true,
  "data": {
    // 응답 데이터
  }
}
```

### 에러 응답
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "에러 메시지",
    "details": "상세 에러 정보 (옵션)"
  }
}
```

## API 엔드포인트

### 인증 API

#### 회원가입
- **POST /auth/register**
  - 설명: 새로운 사용자 계정을 생성합니다
  - 요청 바디:
    ```json
    {
      "email": "user@example.com",
      "password": "password123",
      "name": "홍길동"
    }
    ```
  - 응답: 201 Created
    ```json
    {
      "success": true,
      "data": {
        "id": "user_id",
        "email": "user@example.com",
        "name": "홍길동",
        "createdAt": "2024-02-28T12:00:00Z"
      }
    }
    ```

#### 로그인
- **POST /auth/login**
  - 설명: 사용자 인증 및 JWT 토큰 발급
  - 요청 바디:
    ```json
    {
      "email": "user@example.com",
      "password": "password123"
    }
    ```
  - 응답: 200 OK
    ```json
    {
      "success": true,
      "data": {
        "token": "eyJhbGciOiJIUzI1NiIs...",
        "user": {
          "id": "user_id",
          "email": "user@example.com",
          "name": "홍길동"
        }
      }
    }
    ```

### 사용자 API

#### 프로필 생성
- **POST /user/profile**
  - 설명: 사용자 프로필 정보를 생성합니다
  - 인증: 필요
  - 요청 바디:
    ```json
    {
      "nickname": "개발왕",
      "bio": "열정적인 개발자입니다",
      "jobTitle": "백엔드 개발자",
      "company": "테크 컴퍼니",
      "yearsOfExperience": 3
    }
    ```
  - 응답: 201 Created
    ```json
    {
      "success": true,
      "data": {
        "id": "user_id",
        "nickname": "개발왕",
        "bio": "열정적인 개발자입니다",
        "jobTitle": "백엔드 개발자",
        "company": "테크 컴퍼니",
        "yearsOfExperience": 3,
        "createdAt": "2024-02-28T12:00:00Z",
        "updatedAt": "2024-02-28T12:00:00Z"
      }
    }
    ```

#### 직무 만족도 중요도 생성
- **POST /user/job-satisfaction-importance**
  - 설명: 사용자의 직무 만족도 요소별 중요도를 설정합니다
  - 인증: 필요
  - 요청 바디:
    ```json
    {
      "workload": 80,          // 업무량 중요도 (0-100)
      "compensation": 70,      // 보상 중요도 (0-100)
      "growth": 90,           // 성장 중요도 (0-100)
      "workEnvironment": 85,   // 근무환경 중요도 (0-100)
      "workRelationships": 75, // 직장 내 관계 중요도 (0-100)
      "workValues": 95        // 직장 가치관 중요도 (0-100)
    }
    ```
  - 응답: 201 Created
    ```json
    {
      "success": true,
      "data": {
        "id": "user_id",
        "workload": 80,
        "compensation": 70,
        "growth": 90,
        "workEnvironment": 85,
        "workRelationships": 75,
        "workValues": 95,
        "createdAt": "2024-02-28T12:00:00Z",
        "updatedAt": "2024-02-28T12:00:00Z"
      }
    }
    ```
  - 에러 응답:
    - 400 Bad Request: 이미 직무 만족도 중요도가 존재하는 경우
    - 422 Unprocessable Entity: 입력값 검증 실패

## 에러 코드

| 에러 코드 | 설명 |
|----------|------|
| INVALID_INPUT | 잘못된 입력값 |
| RESOURCE_EXISTS | 이미 리소스가 존재함 |
| RESOURCE_NOT_FOUND | 리소스를 찾을 수 없음 |
| UNAUTHORIZED | 인증되지 않은 요청 |
| DATABASE_ERROR | 데이터베이스 오류 |
| INTERNAL_ERROR | 내부 서버 오류 | 