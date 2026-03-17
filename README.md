# 개요

디스코드 음성 채널 참여 시간을 자동으로 기록하고, 이를 웹 대시보드에서 시각화하여 보여주는 오픈소스 도구입니다. 스터디 그룹이나 커뮤니티에서 멤버들의 활동량을 측정하고 동기를 부여하기 위해 제작되었습니다.

## 주요 기능

### 디스코드 봇
- **음성 채널 트래킹**: 사용자가 음성 채널에 입장하고 퇴장하는 시간을 실시간으로 기록합니다.
- **입/퇴장 알림**: 설정된 채널로 사용자의 입장 및 퇴장 시간을 전송합니다.
- **슬래시 커맨드 지원**:
  - `/info`: 현재 서버 또는 특정 사용자의 활동 내역 확인 링크를 제공합니다.
  - `/set_channel`: 봇의 알림 메시지가 전송될 메인 채널을 설정합니다.

### 웹 대시보드
- **활동 통계**: 유저별 총 접속 시간을 확인할 수 있습니다.
- **잔디 그래프 (Contribution Graph)**: 깃허브 스타일의 잔디 그래프를 통해 일자별 참여도를 시각적으로 보여줍니다.

## 기술 스택

### Backend
- **Language**: Go (1.25.6+)
- **Framework**: Gin (Web Server)
- **Database**: MySQL (sqlx 사용)
- **Communication**: Discord API (WebSocket & REST)

### Frontend
- **Library**: React 19 (TypeScript)
- **Build Tool**: Vite
- **Styling**: Sass (SCSS)
- **State/Routing**: React Router, Axios, Dayjs

### Infrastructure
- **Deployment**: Docker (Multi-stage build)
- **CI/CD**: GitHub Actions

## 설정 및 설치

### 전제 조건
- Go 1.25.6 이상
- Node.js 24 이상
- MySQL 데이터베이스
- 디스코드 애플리케이션 (봇 토큰 및 클라이언트 정보)

### 환경 변수 설정 (`.env`)
루트 디렉토리에 `.env` 파일을 생성하고 아래 내용을 입력합니다.

```env
# Database
DB_URL=username:password@tcp(127.0.0.1:3306)/database_name?parseTime=true

# Discord
DISCORD_TOKEN=your_discord_bot_token
APPLICATION_ID=your_application_id
CLIENT_SECRET=your_client_secret

# General
HOST_URL=http://your-domain.com
TZ=Asia/Seoul
```

### 실행 방법

#### 1. 로컬에서 직접 실행
**Frontend 빌드:**
```bash
cd front
npm install
npm run build
```

**Backend 실행:**
```bash
go mod download
go run cmd/bot/main.go
```

#### 2. Docker를 사용하여 실행
```bash
docker build -t sgbot -f deploy/dockerfile .
docker run -p 8080:8080 --env-file .env sgbot
```

## 프로젝트 구조
- `cmd/bot/`: 애플리케이션 엔트리 포인트 (main.go)
- `internal/bot/`: 디스코드 봇 로직 (이벤트 핸들러, 커맨드 등)
- `internal/web/`: Gin 기반 웹 서버 및 API 엔드포인트
- `internal/repository/`: 데이터베이스 스키마 및 CRUD 로직
- `internal/discord/`: 디스코드 API 연동 모듈
- `front/`: React 기반 프론트엔드 소스 코드
- `deploy/`: 배포 관련 설정 (Dockerfile)
