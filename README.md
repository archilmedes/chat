# chat

## Setup

```
cp config/config.go.template config/config.go
```

Fill in `config/config.go` with your mysql database credentials.

## Pre-commit hook

You can set up a pre-commit hook so that the docs are generated before the code is committed, and a post-commit hook that
adds the docs to your commit.

```
ln -s -f ../../formatting .git/hooks/pre-commit
```
