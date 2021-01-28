#!/bin/bash
cd $(dirname $(dirname $(readlink -f $0)))
curl --no-progress-meter -H "Authorization: $API_KEY" \
-H'Content-Type: application/json' \
-d @tmp/tenant.json http://localhost:9011/api/tenant/${TENANT_ID} | jq > tmp/tenant-01.json
