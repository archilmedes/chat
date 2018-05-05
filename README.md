# chat

## First-time Setup

```
cp config/config.go.template config/config.go
```

You also need npm. We use localtunnel to create a reverse proxy to forward requests on a websocket without a central server. Install with:

```
npm install -g localtunnel
```

Fill in `config/config.go` with your mysql database credentials.

## Running the application

```
go build
./chat
```

###  Development

#### Hooks

You can set up a pre-commit hook so that the application `gofmt`s your code before committing.
You can also set up a pre-push hook so that the application builds and tests your code before pushing.
Note that the scripts must be executable for the hooks to run.

```
chmod +x scripts/pre-commit
ln -s -f ../../scripts/pre-commit .git/hooks/pre-commit
chmod +x scripts/pre-push
ln -s -f ../../scripts/pre-push .git/hooks/pre-push
```

#### Generate Documentation

There is a script that generates documentation with `godoc` for this package.
It can be run with the following script:

```
chmod +x scripts/generate_docs
./scripts/generate_docs
```

It starts up a godoc server and downloads the `chat` package html, css, and js relevant to the package.
