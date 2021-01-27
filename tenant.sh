#!/bin/bash
curl --no-progress-meter -H "Authorization: $API_KEY" \
-H'Content-Type: application/json' \
-d @tenant.json http://localhost:9011/api/tenant/${TENANT_ID} | jq > tenant-01.json
