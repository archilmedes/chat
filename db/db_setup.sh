#!/bin/sh
if [ "$2" = "" ]; then
    mysql -u "$1" < "db_setup.sql"
else
    mysql -u "$1" -p"$2" < "db_setup.sql"
fi;