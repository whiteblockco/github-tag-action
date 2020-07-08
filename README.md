# github-tag-action

Inspired by [anothrNick/github-tag-action](https://github.com/anothrNick/github-tag-action)

Increase build number when PR is merged into develop branch. 
For example, let's assume that latest tag is `0.1.0-0`. After running this action, tag name will be set `0.1.0-1`(annotated tag) on head commit of specific(mostly `develop`) branch.


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
        REPO_TOKEN: ${{ secrets.REPO_TOKEN }}
```

REPO_TOKEN: For reading repository content, `REPO_TOKEN` is a deploy key of repository that run the github action. 
You should create on `github.com/USERNAME/REPO/settings/keys` and set secret (https://github.com/USERNAME/REPO/settings/secrets/new)

# Note

This action use annotated tag instead of lightweight tag. Because `man git-tag` says:

> Annotated tags are meant for release while lightweight tags are meant for private or temporary object labels.

- https://stackoverflow.com/a/35059291/4108346

If buildNumber is 0 this action increase patch part and set build number to 1
ex) `0.1.1 -> 0.1.2-1`, `1.2.3->1.2.4-1`  
