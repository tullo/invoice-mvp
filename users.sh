#!/bin/bash
# curl --no-progress-meter -H "Authorization: $API_KEY" \
#   -H "X-FusionAuth-TenantId: ${TENANT_ID}" \
#   http://localhost:9011/api/user/search?ids=83bf3724-adbd-462a-8c4a-f9e4f74f47b9 | jq

# admin@example.com
curl --no-progress-meter -H "Authorization: $API_KEY" \
  -H "X-FusionAuth-TenantId: ${TENANT_ID}" \
  -H 'Content-Type: application/json' \
 	-d @user-01.json http://localhost:9011/api/user/${USER_01_ID} | jq > users.json

# user@example.com
curl --no-progress-meter -H "Authorization: $API_KEY" \
  -H "X-FusionAuth-TenantId: ${TENANT_ID}" \
  -H 'Content-Type: application/json' \
 	-d @user-02.json http://localhost:9011/api/user/${USER_02_ID} | jq >> users.json

# test@example.com
curl --no-progress-meter -H "Authorization: $API_KEY" \
  -H "X-FusionAuth-TenantId: ${TENANT_ID}" \
  -H 'Content-Type: application/json' \
 	-d @user-03.json http://localhost:9011/api/user/${USER_03_ID} | jq >> users.json
