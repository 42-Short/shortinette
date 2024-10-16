#!/bin/sh

if ! command -v atlas &> /dev/null
then
    echo "Atlas is not installed or not in the PATH"
    echo "please install Atlas first - visit https://atlasgo.io for installation instructions"
    exit 1
fi

echo "applying schema changes..."
if atlas schema apply --url "sqlite3://shortinette.db" --to "file://schema.hcl"; then
    echo "schema changes applied successfully."
    
    echo "saving migration"
    if atlas schema inspect --url "sqlite3://shortinette.db" > "migrations/migration-$(date '+%Y-%m-%d-%T').hcl"; then
        echo "migration file created."
    else
        echo "error: could not create migration file"
        exit 1
    fi
else
    echo "error: could not apply schema changes"
    exit 1
fi
