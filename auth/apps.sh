#!/bin/bash
cd $(dirname $(dirname $(readlink -f $0)))

# Invoice MVP =================================================================
curl --no-progress-meter -H "Authorization: $API_KEY" \
    -H "X-FusionAuth-TenantId: ${TENANT_ID}" \
    -H 'Content-Type: application/json' \
 	-d  @tmp/app-01.json http://localhost:9011/api/application/${INVOICE_APP_ID} | jq > tmp/apps.json

# Test APP ====================================================================
curl --no-progress-meter -H "Authorization: $API_KEY" \
    -H "X-FusionAuth-TenantId: ${TENANT_ID}" \
    -H 'Content-Type: application/json' \
 	-d  @tmp/app-02.json http://localhost:9011/api/application/${TEST_APP_ID} | jq >> tmp/apps.json
