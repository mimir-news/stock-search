#!/bin/bash

TEST_ID=$(openssl rand -hex 5)
GIT_COMMIT=$(git rev-parse HEAD)
SHORT_COMMIT="${GIT_COMMIT:0:7}"

# Service metadata
APPV_FILE="../appv.json"
SVC_NAME=$(jq '.name' -r $APPV_FILE)
SVC_VERSION=$(jq '.version' -r $APPV_FILE)
SVC_REGISTRY=$(jq '.registry' -r $APPV_FILE)
SVC_IMAGE="$SVC_REGISTRY/$SVC_NAME:$SVC_VERSION"
SVC_CONTAINER_NAME="$SVC_NAME-$SVC_VERSION-$SHORT_COMMIT-$TEST_ID"

echo "Testing $SVC_NAME v: $SVC_VERSION commit: $SHORT_COMMIT. Test ID: $TEST_ID"
echo ""

# Database metadata
DB_IMAGE='postgres:11.1-alpine'
DB_CONTAINER_NAME="$SVC_NAME-db-$TEST_ID"

# Database setup
echo "Starting database: $DB_CONTAINER_NAME"
docker run -d --rm --name $DB_CONTAINER_NAME -p 5432:5432 \
   -e POSTGRES_PASSWORD=password $DB_IMAGE

echo "Sleeping for 5 seconds to make database ready"
sleep 5

echo 'Setup up database and user'
docker exec -i $DB_CONTAINER_NAME psql -U postgres < conf/db_setup.sql
docker exec -i $DB_CONTAINER_NAME psql -U streamlistner streamlistner < conf/schema.sql
docker exec -i $DB_CONTAINER_NAME psql -U postgres streamlistner < conf/db_user_setup.sql

echo 'Database ready'
