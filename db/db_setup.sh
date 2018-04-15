#!/bin/bash

if [ -z "$2" ]; then
	mysql -u "$1" < "db/db_setup.sql"
else
	mysql -u "$1" -p "$2" < "db/db_setup.sql"
fi
