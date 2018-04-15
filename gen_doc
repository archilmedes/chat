#!/bin/bash
set -u

build_result=$(go build)
build_rc=$?
if [ $build_rc -ne 0 ] ; then
	echo "git pre-commit check failed: build failed."
	exit 1
fi

PKG=chat
DOC_DIR=docs

# Run a godoc server which we will scrape
godoc -http=localhost:6060 & >/dev/null 2>&1
DOC_PID=$!

# Wait for the server to init
sleep 2

# Scrape the godocs for the CSS / JS
wget -r -m -k -E -p -erobots=off --include-directories="/pkg,/lib" --exclude-directories="*" "http://localhost:6060/pkg/$PKG/" >/dev/null 2>&1

# Stop the godoc server
kill -9 $DOC_PID >/dev/null 2>&1

# Delete the old directory and put the docs in place
rm -rf $DOC_DIR >/dev/null 2>&1
mv localhost\:6060 $DOC_DIR
