# docutron

Minimalist invoice generation.

Docutron allows you to easily generate invoices from plain text JSON files. 

It's deliberately simple and should be easy to customise to suit your needs.

## Why JSON files?

Because the world has enough YAML (and the Go standard library handles it nicely).


## Build and install

```
git clone git@github.com:minimalistsoftware/docutron.git
cd docutron/cmd/docutron
go install .
```

## Dependencies

```
Go
Chrome/Chromium (optional, for PDF generation)
```

## Usage

### Make a directory to work in, with required config and directories.

```
docutron -init mybusiness/
```



Create a new invoice
```
cd mybusiness
docutron -new
2022/11/18 17:40:58 wrote json/INV1.json
```

Generate an HTML version
```
docutron -name json/INV1.json -html
```

Generate a PDF version using Chrome/Chromium
```
docutron -name json/INV1.json -pdf
```


### Customing the template

The invoice template is in templates/invoice.html


