#!/bin/bash
# curl --no-progress-meter -H "Authorization: $API_KEY" \
# "X-FusionAuth-TenantId: ${TENANT_ID}" \
# http://localhost:9011/api/application/26d006ae-58ef-40d8-862a-6a1db34754df | jq

# Invoice MVP =================================================================
curl --no-progress-meter -H "Authorization: $API_KEY" \
    -H "X-FusionAuth-TenantId: ${TENANT_ID}" \
    -H 'Content-Type: application/json' \
 	-d  @app-01.json http://localhost:9011/api/application/${INVOICE_APP_ID} | jq > apps.json

# Test APP ====================================================================
curl --no-progress-meter -H "Authorization: $API_KEY" \
    -H "X-FusionAuth-TenantId: ${TENANT_ID}" \
    -H 'Content-Type: application/json' \
 	-d  @app-02.json http://localhost:9011/api/application/${TEST_APP_ID} | jq >> apps.json
