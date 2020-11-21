package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/tullo/invoice-mvp/domain"
)

// Embedded wraps only bookings for now.
type Embedded struct {
	Bookings []domain.Booking `json:"bookings,omitempty"`
}

// Link models a HAL link.
type Link struct {
	Href string `json:"href"`
}

// HALInvoice decorates an invoice with HAL conform _link elements.
type HALInvoice struct {
	domain.Invoice
	Links    map[domain.Operation]Link `json:"_links"`              // _links
	Embedded *Embedded                 `json:"_embedded,omitempty"` // _embedded
}

// NewHALInvoice instantiates a HAL invoice.
func NewHALInvoice(i domain.Invoice) HALInvoice {
	var links = make(map[domain.Operation]Link)
	links["self"] = Link{fmt.Sprintf("/invoice/%d", i.ID)}
	for _, os := range i.Operations() {
		if link, err := translate(os, i); err == nil {
			links[os] = link
		} else {
			log.Print(err)
		}
	}
	return HALInvoice{Invoice: i, Links: links}
}

func translate(o domain.Operation, i domain.Invoice) (Link, error) {
	switch o {
	case "book":
		return Link{fmt.Sprintf("/book/%d", i.ID)}, nil
	case "bookings":
		return Link{fmt.Sprintf("/invoice/%d/bookings", i.ID)}, nil
	case "charge":
		return Link{fmt.Sprintf("/charge/%d", i.ID)}, nil
	case "cancel":
		return Link{fmt.Sprintf("/invoice/%d", i.ID)}, nil
	case "payment":
		return Link{fmt.Sprintf("/payment/%d", i.ID)}, nil
	case "archive":
		return Link{fmt.Sprintf("/invoice/%d", i.ID)}, nil
	default:
		return Link{}, fmt.Errorf("No translation found for operation %s", o)
	}
}

// HALInvoicePresenter structure.
type HALInvoicePresenter struct {
	writer http.ResponseWriter
}

// NewHALInvoicePresenter instantiates a HAL invoice presenter.
func NewHALInvoicePresenter(w http.ResponseWriter) HALInvoicePresenter {
	return HALInvoicePresenter{writer: w}
}

// Present knows how to present a HAL invoice.
func (p HALInvoicePresenter) Present(i interface{}) {
	inv := i.(HALInvoice)
	if len(inv.Bookings) > 0 {
		var e Embedded
		e.Bookings = inv.Bookings
		inv.Embedded = &e
	}

	if b, err := json.Marshal(inv); err == nil {
		p.writer.Write(b)
	}
}
