# Invoice MVP

A minimum viable product for invoicing.

## Level-2-Services

- 4th iteration [Use Case unit tests and HTTP tests](https://github.com/tullo/invoice-mvp/blob/6f41b13fd163fe7c15ad6164b01b20a6f21438cf/usecase/update_invoice_test.go) (the current mvp implementation)
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
