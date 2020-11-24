package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/tullo/invoice-mvp/domain"
)

// Interfaces

// InvoicePresenter ...
type InvoicePresenter interface {
	Present(i interface{})
}

// Implementations

// CacheableActivities decorates activities with last modified date.
type CacheableActivities struct {
	Activities   []byte
	LastModified time.Time
}

// ActivitiesPresenter implements the presenter interface.
type ActivitiesPresenter struct {
}

// NewActivitiesPresenter instantiates an activities presenter.
func NewActivitiesPresenter() ActivitiesPresenter {
	return ActivitiesPresenter{}
}

// DefaultPresenter ...
type DefaultPresenter struct {
}

// NewDefaultPresenter ...
func NewDefaultPresenter() DefaultPresenter {
	return DefaultPresenter{}
}

// JSONInvoicePresenter ...
type JSONInvoicePresenter struct {
	writer http.ResponseWriter
}

// NewJSONInvoicePresenter ...
func NewJSONInvoicePresenter(w http.ResponseWriter) JSONInvoicePresenter {
	return JSONInvoicePresenter{writer: w}
}

// PDFInvoicePresenter ...
type PDFInvoicePresenter struct {
	writer  http.ResponseWriter
	request *http.Request
}

// NewPDFInvoicePresenter ...
func NewPDFInvoicePresenter(w http.ResponseWriter, r *http.Request) PDFInvoicePresenter {
	return PDFInvoicePresenter{writer: w, request: r}
}

// Presentations

// Present knows how to present the activities list combined with the last
// modified date.
func (ActivitiesPresenter) Present(i interface{}) CacheableActivities {
	lm := time.Unix(0, 0)
	as := i.([]domain.Activity)
	for _, a := range as {
		if a.Updated.After(lm) {
			lm = a.Updated
		}
	}
	b, _ := json.Marshal(i)
	return CacheableActivities{Activities: b, LastModified: lm}
}

// Present ...
func (p DefaultPresenter) Present(i interface{}) {
}

// Present ...
func (p JSONInvoicePresenter) Present(i interface{}) {
	if b, err := json.Marshal(i); err == nil {
		p.writer.Header().Set("Content-Type", "application/json")
		p.writer.Write(b)
	}
}

// Present ...
func (p PDFInvoicePresenter) Present(i interface{}) {
	modTime := time.Now()
	invoice := i.(domain.Invoice)
	content := bytes.NewReader(invoice.ToPDF())
	http.ServeContent(p.writer, p.request, "invoice.pdf", modTime, content)
}
