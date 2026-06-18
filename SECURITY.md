# SECURITY.md

## 보안 정책

### 지원 버전

| 버전 | 지원 여부 |
|---|---|
| main 브랜치 최신 | ✅ 지원 |
| 이전 배포본 | ❌ 미지원 (이미지 태그 업데이트로 대응) |

---

## 취약점 신고

이 레포지토리에서 보안 취약점을 발견한 경우 **절대 공개 이슈로 등록하지 마세요.**

아래 방법으로 비공개 신고해 주세요:

1. **GitHub Private Security Reporting** (권장)
   - `Security` 탭 → `Report a vulnerability`
2. **이메일**: 프로젝트 Approver에게 직접 연락

신고 시 포함할 내용:
- 취약점 유형 (CVE 분류 등)
- 재현 단계
- 영향 범위 (어떤 환경, 어떤 데이터)
- 심각도 평가

---

## 보안 설계 원칙

### 1. 최소 권한 (Least Privilege)

| 역할 | dev | qa | prod |
|---|---|---|---|
| Developer | sync ✅ | ❌ | ❌ |
| QA Engineer | sync ✅ | sync ✅ | ❌ |
| Approver | sync ✅ | sync ✅ | sync ✅ |

### 2. 3중 prod 배포 게이트

```
Gate 1: GitHub Environment "production" → Approver 수동 승인
Gate 2: gitops-manifest main Branch Protection → PR + CODEOWNERS 리뷰
Gate 3: ArgoCD app-prod → automated 없음 → 수동 Sync만 허용
```

### 3. 불변 이미지 (Immutable Image Tags)

- 모든 이미지는 `git-sha` 7자리 태그로 식별
- `latest` 태그는 참조용으로만 사용, prod 배포에는 sha 태그 사용
- 한번 push된 이미지는 덮어쓰지 않음

### 4. 감사 추적 (Audit Trail)

모든 배포 이벤트는 다음 경로에서 추적 가능:

| 이벤트 | 추적 위치 |
|---|---|
| 코드 변경 | `gitops-app` git log |
| 이미지 빌드 | GitHub Actions 워크플로우 실행 기록 |
| manifest 변경 | `gitops-manifest` git log + PR 기록 |
| 클러스터 배포 | ArgoCD UI → Application → History |
| prod 승인 | GitHub Environment → Deployments |

### 5. 컨테이너 보안

```yaml
securityContext:
  readOnlyRootFilesystem: true    # 파일시스템 쓰기 차단
  allowPrivilegeEscalation: false # 권한 상승 차단
  runAsNonRoot: false             # (scratch 이미지 — 추후 nonroot 적용)
  seccompProfile:
    type: RuntimeDefault          # 시스템 콜 제한
```

### 6. Secret 관리 정책

| Secret | 저장 위치 | 순환 주기 |
|---|---|---|
| DOCKERHUB_TOKEN | GitHub Secrets | 90일 |
| MANIFEST_REPO_TOKEN | GitHub Secrets (Fine-grained PAT) | 90일 |
| ArgoCD admin 비밀번호 | 초기 설정 후 즉시 변경 | 즉시 |

**절대 하지 말 것:**
- Secret을 코드나 커밋에 포함
- `latest` 태그로 prod 배포
- Branch Protection 규칙 bypass (`--no-verify` 등)

---

## 의존성 취약점 스캔

Go 의존성 스캔:
```bash
# 외부 의존성 없음 (stdlib만 사용)
go mod tidy
govulncheck ./...
```

컨테이너 이미지 스캔 (선택):
```bash
trivy image <dockerhub-id>/gitops-demo-app:<tag>
```
