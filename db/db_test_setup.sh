#!/bin/sh

if [ -z "$2" ]; then
	mysql -u "$1" < "db_test_setup.sql"
else
	mysql -u "$1" -p"$2" < "db_test_setup.sql"
fi
