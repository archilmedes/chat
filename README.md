# chat

## Setup

```
cp config/config.go.template config/config.go
```

Fill in `config/config.go` with your mysql database credentials.

## Pre-commit hook

You can set up a pre-commit hook and a post-commit hook so that the docs are generated before the code is committed, and a post-commit hook that
adds the docs to your commit.

```
chmod +x scripts/pre-commit
chmod +x scripts/post-commit
ln -s -f ../../scripts/pre-commit .git/hooks/pre-commit
ln -s -f ../../scripts/post-commit .git/hooks/post-commit
```
