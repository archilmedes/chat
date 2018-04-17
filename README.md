# chat

## First-time Setup

```
cp config/config.go.template config/config.go
```

Fill in `config/config.go` with your mysql database credentials.

###  Pre-commit hook

You can set up a pre-commit hook so that the application `gofmt`s your code before committing.
You can also set up a pre-push hook so that the application builds and tests your code before pushing.
Note that the scripts must be executable for the hooks to run.

```
chmod +x scripts/pre-commit
ln -s -f ../../scripts/pre-commit .git/hooks/pre-commit
chmod +x scripts/pre-push
ln -s -f ../../scripts/pre-push .git/hooks/pre-push
```

### Generate Documentation

There is a script that generates documentation with `godoc` for this package.
It can be run with the following script:

```
chmod +x scripts/generate_docs
./scripts/generate_docs
```

It starts up a godoc server and downloads the `chat` package html, css, and js relevant to the package.

### Database Migrations (Devs Only)

Every schema change, the database will need to be "migrated". This can be done by dropping and recreating the database with the following script.

```
sh db/db_setup.sh [mysql username] [mysql password]
```
