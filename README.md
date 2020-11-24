# Invoice MVP

A minimum viable product for invoicing.

## Level-3: HyperMedia/HATEOAS

- 9th iteration [HTTP Response Caching](https://github.com/tullo/invoice-mvp/tree/a73ff21913bc54feea14d87101eaf1ca319b1900) (the current mvp implementation)
- 8th iteration [User Authorisation](https://github.com/tullo/invoice-mvp/tree/326343f99c518ea4b6cbfd706f790cc296fe03aa)
- 7th iteration [Transport Layer Encryption and Server Authn](https://github.com/tullo/invoice-mvp/tree/0a5848995a5056fc2c62b1fba3572c4b8629faab)
- 6th iteration [Authentication: Basic, Digest, JWT](https://github.com/tullo/invoice-mvp/tree/db1c1f28fb9b5270cf6bb6ee6dc4c12f17a93303)
- 5th iteration [HAL invoice representation](https://github.com/tullo/invoice-mvp/tree/e0a6377b1b76a0cbb73af9440d76145d048d499d)
- 4th iteration [Use Case unit tests and HTTP tests](https://github.com/tullo/invoice-mvp/blob/6f41b13fd163fe7c15ad6164b01b20a6f21438cf/usecase/update_invoice_test.go)
- 3rd iteration [Use Cases with Ports & Adapters Architecture](https://github.com/tullo/invoice-mvp/tree/2d20583dadb7c70d34a6f72102c8d72074d0642f)
- 2nd iteration [CRUD-Services](https://github.com/tullo/invoice-mvp/tree/1194c3a456374718ce2a5e6e9bf1102c063ee2f4)
- [Richardson Maturity Model](https://devopedia.org/richardson-maturity-model) (RMM)

## HAL Invoice

HAL invoice representation with allowed actions depending on current invoice state.

```json
{
  "id": 1,
  "month": 9,
  "year": 2020,
  "status": "open",
  "customerId": 1,
  "_links": {
    "book": {
      "href": "/book/1"
    },
    "bookings": {
      "href": "/invoice/1/bookings"
    },
    "cancel": {
      "href": "/invoice/1"
    },
    "charge": {
      "href": "/charge/1"
    },
    "self": {
      "href": "/invoice/1"
    }
  }
}
```
