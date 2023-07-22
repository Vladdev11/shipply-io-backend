#!/bin/sh

set -e

ENV="dev"

if [[ $ENVIRONMENT == "production" ]]; then
 ENV="prod"
fi

echo ":: Fetching s3://${CONFIG_BUCKET}/config.yml for env ${ENVIRONMENT} (${ENV}) ..."

aws s3 cp s3://$CONFIG_BUCKET/config.yml ./config.$ENV.yml

echo ":: Starting backend app..."

exec "$@"