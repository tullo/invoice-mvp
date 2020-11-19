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
