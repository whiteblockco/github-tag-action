# github-tag-action

Inspired by [anothrNick/github-tag-action](https://github.com/anothrNick/github-tag-action)

Increase build number when PR is merged into develop branch.

# Usage

```yaml
name: Increase build number
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
      uses: whiteblockco/github-tag-action@master
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```
