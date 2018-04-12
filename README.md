<<<<<<< HEAD
# chat

## Setup

```
cp config/config.go.template config/config.go
```

Fill in `config/config.go` with your mysql database credentials.
=======
# Setup

```
go build
```

# Pre-commit hook

You can set up a pre-commit hook so that the docs are generated before the code is committed, and a post-commit hook that
adds the docs to your commit.

```
ln -s gen_docs.sh .git/hooks/pre-commit
ln -s add_doc_to_git.sh .git/hooks/post-commit
```
>>>>>>> Generate the docs and add instructions to pre-commit hook
