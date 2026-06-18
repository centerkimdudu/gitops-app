# gitops-demo-app

Go로 작성된 GitOps 데모용 HTTP 서버. CI/CD 파이프라인의 소스가 되는 앱 레포입니다.

## 엔드포인트

| Path | 설명 |
|---|---|
| `GET /` | `Hello from GitOps Demo! ENV=<env> VERSION=<version>` |
| `GET /health` | liveness/readiness probe 용 (`ok` 반환) |
| `GET /version` | 현재 환경·버전 JSON 반환 |

## 로컬 실행

```bash
cd src
go run . 
# → http://localhost:8080
```

## 테스트

```bash
cd src
go test ./... -v -race
```

## CI/CD 흐름

```
git push (feature/**) → PR → main merge
    └─► GitHub Actions CI
            ├─ go test
            ├─ docker build & push (Docker Hub, sha 태그)
            └─ gitops-manifest PR 생성 (dev overlay image tag 업데이트)
```

prod 배포는 `Actions` → `Promote to Prod` 워크플로우를 수동 실행합니다.

## 관련 레포

- **gitops-manifest**: K8s 매니페스트 + ArgoCD Application
- **ARCHITECTURE.md**: 전체 시스템 아키텍처
- **TEST.md**: 테스트 절차서
