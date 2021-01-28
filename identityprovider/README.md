# Identity Provider

## HMAC signed tokens

For testing JWT based authentication and authorisation `go run main.go` can be used to generate a token with with sample claims:

```json
{
    "name": "Go Invoicer",
    "admin": true,
    "sub": "f8c39a31-9ced-4761-8a33-b9c628a67510"
}
```

together with `HMACKeyFunc` [../rest/auth.go](rest/auth.go) that verifies the token.

## Identity Provider Service

This MVP uses FusionAuth to provide an external service for user identity handling and token signing.

To launch the service follow these steps:

1. `make identityprovider-up`
1. Go to http://localhost:9011/ and complete the setup-wizard.
1. Create API key named "Invoice MVP" at http://localhost:9011/admin/api/
1. Copy the API key id.
1. Edit the makefile and replace the exported API_KEY with key from step 4.
1. Bootstrap Identity Provider config by running through the makefile targets in the "Auth Bootstrapping" section.
1. Finally: run `make test` to execute unit & integration tests.
