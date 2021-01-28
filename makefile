SHELL = /bin/bash -o pipefail
# make -n
# make -np 2>&1 | less

# Generate an api key via FA UI and paste here.
export API_KEY=ecNmeuOEhG2rZgggPZkPoes3sNK4G522bIMAbTEOll8w-sWzGI6rSbPe
export SIGNING_KEY_ID=032f0a62-e111-4144-9f37-ca55112383ec
export TENANT_ID=15b07c80-66ca-1a4b-f313-f519f5bf37cd
export INVOICE_APP_ID=1e9dbfe9-8f90-4cc4-90c1-5b6a4ca029c9
export TEST_APP_ID=f2455f9f-c1b6-44d2-b70f-a8abf90c70dc
export USER_01_ID=f73e4176-915c-423d-a092-ae97f0c8de88
export USER_02_ID=d68df2ec-d79a-4fbd-b290-22ae3a91532b
export USER_03_ID=206563d6-1636-4fc2-9bd8-0a68cbdf6ea3

test: bootstrap_dependencies
	go test -v -count=1 ./...

go-deps-reset:
	@git checkout -- go.mod
	@go mod tidy

# -d flag ...download the source code needed to build ...
# -t flag ...consider modules needed to build tests ...
# -u flag ...use newer minor or patch releases when available 
go-deps-upgrade:
	@go get -d -t -u -v ./...
	@go mod tidy

go-mod-tidy:
	@go mod tidy

identityprovider-down:
	cd identityprovider/; docker-compose down

identityprovider-up:
	cd identityprovider/; docker-compose up -d --remove-orphans

bootstrap_dependencies: identityprovider-up
	cd identityprovider/; docker-compose run --rm bootstrap_dependencies

# =============================================================================
# Auth Bootstrapping ==========================================================
# =============================================================================

# 1. Create signing key.
auth-signing-key:
	@./auth/key.sh
	@echo "OK"

# 2. Use default tenant as template.
auth-tenant-template: export DEFAULT_TENANT_ID=d3b6781e-1356-2f7c-ea1a-04549bb42dcf
auth-tenant-template:
	@curl --no-progress-meter -H "Authorization: ${API_KEY}" \
		http://localhost:9011/api/tenant/${DEFAULT_TENANT_ID} | jq > tmp/default-tenant.json
	@echo "OK"

# 3. Create tenant.
# Prepare default-tenant.json with substitutions for:
# NAME, ISSUER, SIGNING_KEY_ID
# See example from auth/tenant-substitutions.json
auth-tenant: export NAME=Invoice MVP
auth-tenant: export ISSUER=invoice.mvp
auth-tenant:
	@envsubst < tmp/default-tenant.json > tmp/tenant.json
	@./auth/tenant.sh
	@echo "OK"

# 4. Create apps with tenant-ref and signing key ref.
auth-apps:
	@envsubst < auth/app-invoice.json > tmp/app-01.json
	@envsubst < auth/app-test.json > tmp/app-02.json
	@./auth/apps.sh
	@echo "OK"

# 5. Create users with tenant-ref.
auth-users:
	@envsubst < auth/user-admin.json > tmp/user-01.json
	@envsubst < auth/user-user.json > tmp/user-02.json
	@envsubst < auth/user-test.json > tmp/user-03.json
	@./auth/users.sh
	@echo "OK"

# 6. Create user regs with tenant-ref.
auth-users-registrations:
	@envsubst < auth/reg-admin.json > tmp/reg-01.json
	@envsubst < auth/reg-user.json > tmp/reg-02.json
	@envsubst < auth/reg-test.json > tmp/reg-03.json
	@./auth/registrations.sh
	@echo "OK"

# 7. Create .env config file.
env: export IID=$(shell cat tmp/apps.json | jq -c '.[] | select( .id == "${INVOICE_APP_ID}" ) | .oauthConfiguration.clientId')
env: export ISECRET=$(shell cat tmp/apps.json | jq -c '.[] | select( .id == "${INVOICE_APP_ID}" ) | .oauthConfiguration.clientSecret')
env: export IUSER=$(shell cat auth/user-admin.json | jq '.user.email')
env: export IPASS=$(shell cat auth/user-admin.json | jq '.user.password')
env: export TID=$(shell cat tmp/apps.json | jq -c '.[] | select( .id == "${TEST_APP_ID}" ) | .oauthConfiguration.clientId')
env: export TSECRET=$(shell cat tmp/apps.json | jq -c '.[] | select( .id == "${TEST_APP_ID}" ) | .oauthConfiguration.clientSecret')
env: export TUSER=$(shell cat auth/user-test.json | jq '.user.email')
env: export TPASS=$(shell cat auth/user-test.json | jq '.user.password')
env:
	@echo "# Invoice App" > env
	@echo "BASE_DIR=$(PWD)" >> env
	@echo "MVP_USERNAME=${IUSER}" >> env
	@echo "MVP_PASSWORD=${IPASS}" >> env
	@echo "USER_ID=${USER_01_ID}" >> env
	@echo "# OAUTH2" >> env
	@echo "AUTH_REALM=invoice.mvp" >> env
	@echo "IDP_ISSUER=invoice.mvp" >> env
	@echo "CLIENT_ID=${IID}" >> env
	@echo "CLIENT_SECRET=${ISECRET}" >> env
	@echo "TENANT_ID=${TENANT_ID}" >> env
	@echo "GRANT_TYPE=authorization_code" >> env
	@echo "TOKEN_URI=http://localhost:9011/oauth2/token" >> env
	@echo "REDIRECT_URI=https://127.0.0.1:8443/auth/token" >> env
	@echo "# Test App" >> env
	@echo "TEST_LOGIN=${TUSER}" >> env
	@echo "TEST_PASSWD=${TPASS}" >> env
	@echo "TEST_CLIENT_ID=${TID}" >> env
	@echo "TEST_CLIENT_SECRET=${TSECRET}" >> env
	@mv env .env
