#!/bin/bash
# admin@example.com
curl --no-progress-meter -H "Authorization: $API_KEY" \
  -H "X-FusionAuth-TenantId: ${TENANT_ID}" \
  -H 'Content-Type: application/json' \
  -d @reg-01.json http://localhost:9011/api/user/registration/${USER_01_ID} | jq

# user@example.com
curl --no-progress-meter -H "Authorization: $API_KEY" \
  -H "X-FusionAuth-TenantId: ${TENANT_ID}" \
  -H 'Content-Type: application/json' \
  -d @reg-02.json http://localhost:9011/api/user/registration/${USER_02_ID} | jq

# test@example.com
curl --no-progress-meter -H "Authorization: $API_KEY" \
  -H "X-FusionAuth-TenantId: ${TENANT_ID}" \
  -H 'Content-Type: application/json' \
  -d @reg-03.json http://localhost:9011/api/user/registration/${USER_03_ID} | jq