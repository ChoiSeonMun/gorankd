## 프로젝트 개요
10만 건 동시 업데이트를 처리하는 GCP 기반 실시간 랭킹 시스템

## 기술 스택
- Language: Go 1.22+
- Communication: gRPC + Protocol Buffers
- Cache: Redis Cluster (ZADD/ZREVRANGE)
- DB: Cloud Spanner (Emulator로 로컬 개발)
- Infra: Docker, GKE

## 디렉토리 구조
```
gorankd/
├── cmd/server/          # main.go 진입점
├── internal/
│   ├── ranking/         # 핵심 랭킹 비즈니스 로직
│   ├── cache/           # Redis 클라이언트 래퍼
│   ├── store/           # Spanner 클라이언트 래퍼
│   └── server/          # gRPC 서버 핸들러
├── api/proto/           # .proto 파일 + gen/ (생성 코드)
├── pkg/                 # 외부 공개 가능한 공용 유틸
├── deploy/
│   ├── docker/          # Dockerfile, docker-compose
│   └── k8s/             # GKE 매니페스트
├── scripts/             # 빌드, proto-gen, 마이그레이션
└── configs/             # 환경별 설정 파일
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
- Redis Cluster 로컬 환경: docker-compose 7노드 구성 사용
- proto 파일 변경 후 반드시 `make proto-gen` 실행 후 커밋