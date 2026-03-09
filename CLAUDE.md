## 프로젝트 개요
10만 건 동시 업데이트를 처리하는 GCP 기반 실시간 랭킹 시스템

## 기술 스택
- Language: Go 1.22+
- Module: `gorankd` (go.mod)
- Communication: gRPC + Protocol Buffers
- Cache: Redis Cluster (ZADD/ZREVRANGE)
- DB: Cloud Spanner (Emulator로 로컬 개발)
- Infra: Docker, GKE
- 주요 의존성: `google.golang.org/grpc`

## API (gRPC)
Proto 파일: `api/proto/ranking.proto` (package: `ranking.v1`)
- `UpdateScore` - 플레이어 점수 업데이트
- `GetRank` - 플레이어 순위 조회
- `GetTopN` - 상위 N명 조회
- `GetPlayerScore` - 플레이어 점수 조회

## 디렉토리 구조
```
gorankd/
├── cmd/server/main.go         # 진입점 (slog, gRPC 서버, graceful shutdown)
├── internal/
│   ├── ranking/               # 핵심 랭킹 비즈니스 로직
│   │   ├── interface.go       # Service 인터페이스, PlayerRank 타입
│   │   └── service.go         # 구현체 (cache + store 의존)
│   ├── cache/                 # Redis 클라이언트 래퍼
│   │   ├── interface.go       # Cache 인터페이스, RankEntry 타입
│   │   └── redis.go           # 구현체 (stub)
│   ├── store/                 # Spanner 클라이언트 래퍼
│   │   ├── interface.go       # Store 인터페이스, PlayerScore 타입
│   │   └── spanner.go         # 구현체 (stub)
│   └── server/                # gRPC 서버 핸들러
│       └── grpc.go            # GRPCServer (ranking.Service 의존)
├── api/proto/                 # .proto 파일
│   ├── ranking.proto
│   └── gen/                   # protoc 생성 코드 (go_package: gorankd/api/proto/gen/rankingpb)
├── pkg/                       # 외부 공개 가능한 공용 유틸
├── deploy/
│   ├── docker/                # Dockerfile (멀티스테이지), docker-compose.yml
│   └── k8s/                   # GKE 매니페스트
├── scripts/proto-gen.sh       # protoc 실행 스크립트
├── configs/config.dev.yaml    # 로컬 개발 설정
└── Makefile                   # proto-gen, build, run, test, lint
```

## 아키텍처 레이어
```
gRPC Handler (internal/server)
  → ranking.Service (internal/ranking)
    → cache.Cache (internal/cache)   # Redis - 실시간 순위 연산
    → store.Store (internal/store)   # Spanner - 영속성
```

## 코드 컨벤션

### 일반
- 에러 핸들링: errors.As/Is 사용, panic 금지
- 동시성: Channel 우선, Mutex는 최소화
- 객체 재사용: sync.Pool 적용 대상 명시할 것
- 패키지명: 단수형 소문자 (ranking, cache, store)
- 인터페이스 정의: `{package}/interface.go` 파일에 집중

### gRPC
- 모든 핸들러에서 context deadline 하위로 전파 필수
- 로깅/메트릭은 Interceptor에서 처리 (핸들러 내부 금지)
- 에러 코드: `google.golang.org/grpc/codes` 표준 사용

### Redis
- 키 형식: `rank:{namespace}:{id}` (콜론 구분)
- 배치 ZADD는 반드시 Pipeline으로 묶기
- 모든 키에 명시적 TTL 설정 필수

### Spanner
- Read-only / Read-write 트랜잭션 명확히 구분
- DML보다 Mutation 방식 우선 (throughput 유리)
- 핫스팟 방지: 키에 UUID v4 또는 해시 prefix 사용

### 로깅
- `slog` 패키지만 사용 (Go 1.22 표준, 구조화 로깅)
- 프로덕션: JSON 포맷 / 개발: text 포맷
- request_id를 context에 항상 포함

## 성능 & 동시성 주의사항
- **Goroutine 누수**: goroutine 생성 시 종료 조건 반드시 명시
- **메모리 할당**: hot path에서 []byte 재사용, sync.Pool 우선 적용
- **backpressure**: 채널 버퍼 크기 명시, 버퍼 포화 시 drop/reject 정책을 주석으로 명시
- **타임아웃 계층**: gRPC deadline → Redis timeout → Spanner timeout 순으로 설정
- **벤치마크**: 핵심 랭킹 연산은 `*_bench_test.go` 필수 작성

## 로컬 개발 환경
- Spanner Emulator는 실제 Spanner와 동작 차이 있음 (commit timestamp 등)
- Redis Cluster 로컬 환경: docker-compose 6노드 + init 컨테이너 구성 사용
- proto 파일 변경 후 반드시 `make proto-gen` 실행 후 커밋
- 로컬 실행: `make run` (포트 50051)
- docker-compose: `docker compose -f deploy/docker/docker-compose.yml up`

## Makefile 타겟
- `make proto-gen` - protoc으로 Go 코드 생성
- `make build` - `bin/gorankd` 바이너리 빌드
- `make run` - 로컬 실행
- `make test` - 테스트 실행
- `make lint` - golangci-lint 실행

## 현재 상태 (구현 진행도)
- [x] 프로젝트 스캐폴딩 (디렉토리, go.mod, Makefile, Docker)
- [x] Proto 파일 정의 (ranking.proto)
- [x] 인터페이스 정의 (ranking, cache, store)
- [x] 랭킹 서비스 구현체 (cache/store 연동 로직)
- [x] gRPC 서버 진입점 (main.go, graceful shutdown)
- [ ] proto-gen 실행 및 gRPC 핸들러 등록
- [ ] Redis 클라이언트 구현 (cache/redis.go)
- [ ] Spanner 클라이언트 구현 (store/spanner.go)
- [ ] gRPC Interceptor (로깅, 메트릭)
- [ ] 벤치마크 테스트
- [ ] 설정 파일 로딩 (configs/config.dev.yaml → 구조체 바인딩)