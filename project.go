package docutron

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path"
)

const perms = 0755

//go:embed templates/*
var templates embed.FS

// InitProject creates a new directory and the skeleton config files.
func InitProject(dirName string) {
	dir := path.Clean(dirName)

	// If the directory already exists, stop init
	if _, err := os.Stat(dir); err == nil {
		fmt.Println("This project already has been setup. Nothing to add.")
		return
	}

	// Otherwise create the directory
	err := os.MkdirAll(dir, perms)
	check(err)

	subDirs := []string{"html", "json", "templates", "pdf"}
	for _, d := range subDirs {
		dPath := path.Join(dir, d)
		err = os.MkdirAll(dPath, perms)
		check(err)
	}
	confPath := path.Join(dir, "config.json")
	WriteConfig(confPath)
	WriteTemplates(dir)
}

// WriteTemplates writes templates from the embedded FS to the new project templates/ directory.
func WriteTemplates(dir string) {
	files, _ := templates.ReadDir("templates")

	for _, f := range files {
		embedPath := path.Join("templates", f.Name())
		if f.IsDir() {
			continue
		}
		// newdir/templates/filename.html
		fpath := path.Join(dir, "templates", f.Name())
		b, err := templates.ReadFile(embedPath)
		check(err)
		err = os.WriteFile(fpath, b, perms)
		log.Printf("wrote %s", fpath)
		check(err)

	}
}
