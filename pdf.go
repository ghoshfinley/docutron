package docutron

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// WritePDF writes a pdf to pdf/title.pdf using wkhtmltopdf.
func WritePDF(inv Invoice) {
	// Create new PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	fileName := fmt.Sprintf("html/%s.html", inv.Title)

	// Set global options
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	//pdfg.Grayscale.Set(true)

	// Read HTML file into memory
	b, err := os.ReadFile(fileName)
	check(err)
	pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader(b)))

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	// Write buffer contents to file on disk
	outFile := fmt.Sprintf("pdf/%s.pdf", inv.Title)
	err = pdfg.WriteFile(outFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("wrote: %s\n", outFile)
}

/*
	WrotePDFChrome writes a PDF to disk using Chrome via chromedp.

This got me better results than wkhtmltopdf for tables split across pages.

I'd rather not have this dependency though.
*/
func WritePDFChrome(req UserRequest, inv Invoice) {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	fileName := fmt.Sprintf("%s/html/%s.html", req.Project, inv.Title)

	absPath, err := filepath.Abs(fileName)
	check(err)

	fileURL := fmt.Sprintf("file://%s", absPath)

	// capture pdf
	var buf []byte
	if err := chromedp.Run(ctx, printToPDF(fileURL, &buf)); err != nil {
		log.Fatal(err)
	}

	outFile := fmt.Sprintf("%s/pdf/%s.pdf", req.Project, inv.Title)

	if err := ioutil.WriteFile(outFile, buf, perms); err != nil {
		log.Fatal(err)
	}
	log.Printf("wrote: %s\n", outFile)

}

// print a specific pdf page.
// basd on the example code on https://github.com/chromedp/examples/blob/master/pdf/main.go
func printToPDF(urlstr string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(false).WithPreferCSSPageSize(true).WithMarginTop(1).Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}
