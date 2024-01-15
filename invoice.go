package docutron

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"time"

	"golang.org/x/text/message"
)

type InvoiceConfig struct {
	GSTPercent int `json:"GSTPercent"`
	Invoice    struct {
		NumOffset   int    `json:"NumOffset"`   // padd this to the numbers generated
		Template    string `json:"Template"`    // path to template file
		Prefix      string `json:"Prefix"`      // prefix to the numeric filename. eg. MIN to generate MIN46
		CompanyName string `json:"CompanyName"` // TODO this should move out of the Invoice struct
		CompanyURL  string `json:"CompanyURL"`
	}
}

type UserRequest struct {
	Project string        `json:"project"`
	Config  InvoiceConfig `json:"config"`
	Items   []LineItem    `json:"items"`
}

type Invoice struct {
	Title          string     `json:"title"`
	CustomerName   string     `json:"customer_name"`
	Date           time.Time  `json:"date"`
	Items          []LineItem `json:"items"`
	CurrencyCode   string     `json:"currency_code"`
	CurrencySymbol string     `json:"currency_symbol"`
	Subtotal       int        `json:"subtotal"`    // cents
	GSTApplies     bool       `json:"gst_applies"` // do we calculate GST
	GST            int        `json:"gst"`         // cents
	Total          int        `json:"total" `      // cents
}

type LineItem struct {
	Position         int     `json:"position"`
	Quantity         int     `json:"quantity"`
	Description      string  `json:"description"`
	UnitPriceDollars float64 `json:"unit_price_dollars"` // UnitPrice in dollars
	UnitPrice        int     // UnitPrice in cents
	TotalPrice       int     // TotalPrice in cents
	GST              int     // GST in cents
}

var config Config

// WriteHTML outputs html/name.html.
func WriteHTML(req UserRequest, inv Invoice, templatePath string) {

	funcMap := template.FuncMap{
		"currency": CentsToString,
		"date":     FormatDate,
	}

	fileName := fmt.Sprintf("%s/html/%s.html", req.Project, inv.Title)

	fh, err := os.Create(fileName)
	check(err)
	defer fh.Close()

	tmpl, err := template.New("invoice").Funcs(funcMap).ParseFiles(templatePath)
	check(err)

	err = tmpl.Execute(fh, inv)
	check(err)
	log.Printf("wrote: %s\n", fileName)
}

func check(err error) {
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}

func CalculateTotals(GSTPercent int, inv Invoice) Invoice {
	var newItems []LineItem

	var subtotal int

	inv.GST = 0
	for _, item := range inv.Items {
		item.UnitPrice = int(item.UnitPriceDollars * 100)
		item.TotalPrice = item.Quantity * item.UnitPrice
		//Item.GSTAmount = item.Quantity / GST
		newItems = append(newItems, item)
		subtotal += item.TotalPrice
	}
	inv.Subtotal = subtotal
	inv.Items = newItems
	if inv.GSTApplies && GSTPercent != 0 {
		inv.GST = subtotal / GSTPercent
	}
	inv.Total = inv.Subtotal + inv.GST
	return inv
}

// NewJSONFile creates a new file json/name.json.
func NewJSONFile(req UserRequest, name string) Invoice {
	fpath := fmt.Sprintf("%s/json/%s.json", req.Project, name)
	fh, err := os.Create(fpath)
	check(err)
	defer fh.Close()

	var inv Invoice

	// Read invoice template from file.
	tb, err := os.ReadFile("templates/invoice.json")
	check(err)
	err = json.Unmarshal(tb, &inv)
	check(err)

	// Add the new title and current date.
	inv.Title = name
	inv.Date = time.Now()
	inv.Items = req.Items

	inv = CalculateTotals(req.Config.GSTPercent, inv)

	b, err := json.MarshalIndent(inv, "", " ")
	check(err)

	_, err = fh.Write(b)
	check(err)
	log.Printf("wrote %s", fpath)

	// Return final Invoice for use
	return inv
}

// UnmarshalJSONFile opens a file in json/name.json.
func UnmarshalJSONFile(name string) Invoice {
	b, err := os.ReadFile(name)
	check(err)
	var inv Invoice
	err = json.Unmarshal(b, &inv)
	check(err)
	return inv
}

// MarshalJSONFile writes inv to name.json as JSON.
func MarshalJSONFile(inv Invoice, name string) {
	b, err := json.MarshalIndent(inv, "", " ")
	check(err)
	err = os.WriteFile(name, b, perms)
	check(err)
	log.Printf("wrote %s", name)
}

// NextNumber determines the next filename number based on number of files in json/ directory + 1.
// TODO make this aware of different document types.
func NextNumber(req UserRequest) int {
	offset := req.Config.Invoice.NumOffset //Read offset from Config file

	jsonDir := fmt.Sprintf("%s/json/", req.Project)
	entries, err := os.ReadDir(jsonDir)
	check(err)
	nextNum := offset + len(entries) + 1
	return nextNum
}

// CentsToString converts integer cents to a human-readable currency string
func CentsToString(c int) string {
	n := float64(c) / 100
	p := message.NewPrinter(message.MatchLanguage("en"))
	return p.Sprintf("%.2f", n)
}

// Format Date in a readable way for the invoice.
func FormatDate(t time.Time) string {
	return t.Format("2 January 2006")
}
