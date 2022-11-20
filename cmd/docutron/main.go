package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/minimalistsoftware/docutron"
)

var newFlag bool
var calcFlag bool
var htmlFlag bool
var pdfFlag bool
var nameFlag string
var initFlag string

func init() {
	flag.BoolVar(&newFlag, "new", false, "Create new Invoice json file")
	flag.BoolVar(&htmlFlag, "html", false, "Output HTML file in html/name.html")
	flag.BoolVar(&pdfFlag, "pdf", false, "Output PDF file in pdf/name.pdf")
	flag.StringVar(&nameFlag, "name", "", "Input filename to operate on eg. json/INV1.json")
	flag.StringVar(&initFlag, "init", "", "Create new project directory")
	flag.Parse()
}

func main() {
	if !newFlag && !calcFlag && !htmlFlag && initFlag == "" && nameFlag == "" {
		flag.Usage()
		return
	}

	if initFlag != "" {
		docutron.InitProject(initFlag)
		return
	}

	config := docutron.ReadConfig()

	if newFlag {
		num := docutron.NextNumber()
		docutron.NewJSONFile(fmt.Sprintf("%s%d", config.Invoice.Prefix, num))
		return
	}

	if htmlFlag && nameFlag == "" {
		log.Fatalf("-name flag must be set to use -html\n")
	}
	if htmlFlag {
		inv := docutron.UnmarshalJSONFile(nameFlag)
		inv = docutron.CalculateTotals(inv)
		docutron.WriteHTML(inv, nameFlag, config.Invoice.Template)
	}
	if pdfFlag && nameFlag == "" {
		log.Fatalf("-name flag must be set to use -pdf\n")
	}
	if pdfFlag {
		inv := docutron.UnmarshalJSONFile(nameFlag)
		inv = docutron.CalculateTotals(inv)
		docutron.WriteHTML(inv, nameFlag, config.Invoice.Template)
		//docutron.WritePDF(inv)
		docutron.WritePDFChrome(inv)
	}

}
