#!/bin/sh
status=0
for file in $(git diff --cached --name-only | grep -e '\.go$'); do
    badfile=$(gofmt -l $file)
    if test -n "$badfile" ; then
        echo "git pre-commit check failed: file needs gofmt: $badfile"
        status=1
    fi
done
exit $status
