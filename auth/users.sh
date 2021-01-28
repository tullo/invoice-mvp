#!/bin/bash
cd $(dirname $(dirname $(readlink -f $0)))

# admin@example.com
curl --no-progress-meter -H "Authorization: $API_KEY" \
  -H "X-FusionAuth-TenantId: ${TENANT_ID}" \
  -H 'Content-Type: application/json' \
 	-d @tmp/user-01.json http://localhost:9011/api/user/${USER_01_ID} | jq > tmp/users.json

# user@example.com
curl --no-progress-meter -H "Authorization: $API_KEY" \
  -H "X-FusionAuth-TenantId: ${TENANT_ID}" \
  -H 'Content-Type: application/json' \
 	-d @tmp/user-02.json http://localhost:9011/api/user/${USER_02_ID} | jq >> tmp/users.json

# test@example.com
curl --no-progress-meter -H "Authorization: $API_KEY" \
  -H "X-FusionAuth-TenantId: ${TENANT_ID}" \
  -H 'Content-Type: application/json' \
 	-d @tmp/user-03.json http://localhost:9011/api/user/${USER_03_ID} | jq >> tmp/users.json
