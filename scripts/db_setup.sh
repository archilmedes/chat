#!/bin/sh

CURR="$(pwd)"
cd "$(dirname "$0")" 

if [ -z "$2" ]; then
	mysql -u "$1" < "db_setup.sql"
else
	mysql -u "$1" -p"$2" < "db_setup.sql"
fi

cd "$CURR"
