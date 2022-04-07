#!/bin/sh

# pg_isready -d <database-name> -h <host> -p <port> -U <user>
until pg_isready -d "$POSTGRES_DB" -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER"; do
	>&2 echo "DB is unavailable - sleeping"
	sleep 1
done

>&2 echo "DB is up - executing command"
exec "$@"