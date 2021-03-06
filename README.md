# GitHub-tag-action

Inspired by [anothrNick/github-tag-action](https://github.com/anothrNick/github-tag-action)

Increase build number when PR is merged into develop branch. 
For example, let's assume that latest tag is `0.1.0-0`. After running this action, tag name will be set `0.1.0-1`(annotated tag) on head commit of specific(mostly `develop`) branch.


# Usage

in .github/workflows/main.yml

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
        REPO_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        WITHOUT_V: true  
```

> GitHub automatically creates a GITHUB_TOKEN secret to use in your workflow. You can use the GITHUB_TOKEN to authenticate in a workflow run.

- https://docs.github.com/en/actions/configuring-and-managing-workflows/authenticating-with-the-github_token


# Note

This action use annotated tag instead of lightweight tag. Because `man git-tag` says:

> Annotated tags are meant for release while lightweight tags are meant for private or temporary object labels.

- https://stackoverflow.com/a/35059291/4108346
