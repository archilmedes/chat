#!/bin/sh
mysql -u "$1" -p"$2" < "db_setup.sql"