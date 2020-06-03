# github-tag-action

dev 스테이지 배포 자동화를 위한 [GitHub Action](https://github.com/features/actions)

k8s의 [namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/) 를 사용해 prod 환경과 테스트 및 개발 환경을 분리함. 같은 클러스터 다른 namespace

전체적인 흐름은 다음과 같음

_@는 Tag, => 는 Merge_

1. master 보다 하나 높은 버전으로 develop 브랜치 태그 생성. 
예: `master@0.1.0` 이면 `develop@0.1.1`
1. develop@0.1.0 를 base 삼아 `feature/브랜치이름-JIRA이슈번호` 브랜치 생성
1. `feature/branch_name-JIR-123` => develop 로 merge 
1. **PR이 develop 브랜치에 merge 될 때 build_number 자동으로 increase** -> 이 작업을 해주는 게 본 repository (github-tag-action)
1. develop@0.1.1-1 을 docuwallet-dev namespace 에 배포 

develop 브랜치는 다음과 같은 포맷으로 태깅함
`{major}.{minor}.{patch}-{build_number}`
예: `0.1.1-1`

# 주석

- 2\) 과정의 feature 브랜치 이름은 변경될 여지가 있음 (20.06.03 현재 확정되지 않은 상태)
- docuwallet-dev 는 원칙적으론 `docuwallet-dev-node-pool` 이라는 노드 풀에 배포 되어야 함. 
 
# Usage

name: Bump version
on:
  push:
    branches:
      - develop
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@master
      with:
        fetch-depth: '0'
    - name: Bump version and push tag
      uses: whiteblockco/github-tag-action@1.0.0
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
