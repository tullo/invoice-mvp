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

1. `docker-compose up`
1. Go to http://localhost:9011/ and complete the setup-wizard
1. Create "Invoice MVP" app: http://localhost:9011/admin/application/
1. Generate a RSA key pair: http://localhost:9011/admin/key/
    1. Name: "Invoice MVP Keys"
    1. Issuer: "invoice.mvp"
1. Edit the "Invoice MVP" app. Click on *JWT* tab and choose the app specific key for "Access Token signing key" field.
1. *Save* to persist the change.
1. *View* to see settings like login url, clientId, etc.
