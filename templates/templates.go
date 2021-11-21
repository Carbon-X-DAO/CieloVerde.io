package templates

import (
	_ "embed"
	"html/template"
	"log"
)

var (
	//go:embed listing.html
	listingSource string
	//go:embed notfound.html
	notFoundSource string
	//go:embed error.html
	errorSource string
)

var (
	Listing  *template.Template
	NotFound *template.Template
	Error    *template.Template
)

func init() {
	Listing = parse("listing", listingSource)
	NotFound = parse("notfound", notFoundSource)
	Error = parse("error", errorSource)
}

func parse(name, text string) *template.Template {
	tmpl, err := template.New(name).Parse(text)

	if err != nil {
		log.Fatal(err)
	}

	return tmpl
}
