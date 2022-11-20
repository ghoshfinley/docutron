package docutron

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	GSTPercent int
	Invoice    struct {
		NumOffset   int    // padd this to the numbers generated
		Template    string // path to template file
		Prefix      string // prefix to the numeric filename. eg. MIN to generate MIN46
		CompanyName string // TODO this should move out of the Invoice struct
		CompanyURL  string
	}
}

func ReadConfig() Config {
	var c Config

	b, err := os.ReadFile("config.json")
	check(err)

	err = json.Unmarshal(b, &c)
	check(err)

	config = c

	return c
}

// WriteConfig initialises a new config JSON file
func WriteConfig(name string) {
	var c Config
	c.GSTPercent = 10
	c.Invoice.Template = "templates/invoice.html"
	c.Invoice.Prefix = "INV"
	b, err := json.MarshalIndent(c, "", " ")
	check(err)
	err = os.WriteFile(name, b, perms)
	check(err)
	log.Printf("wrote %s", name)
}
