#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd ${DIR}

if [ -z "$2" ]; then
	mysql -u "$1" < "db_setup.sql"
else
	mysql -u "$1" -p "$2" < "db_setup.sql"
fi
