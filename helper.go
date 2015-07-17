package response

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/volatile/core"
)

var (
	views    *template.Template
	viewsDir = "views"
)

func init() {
	views = template.New("views")

	core.BeforeRun(func() {
		if _, err := os.Stat(viewsDir); err == nil {
			err = filepath.Walk(viewsDir, func(path string, f os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if !f.IsDir() {
					if _, err = views.ParseFiles(path); err != nil {
						return err
					}
				}

				return nil
			})
			if err != nil {
				panic(err)
			}
		}
	})
}

// AddViewFunc adds a function that will be available to all templates.
func AddViewFunc(name string, f interface{}) {
	views.Funcs(template.FuncMap{name: f})
}

// Status responds with the given status.
func Status(c *core.Context, v int) {
	c.ResponseWriter.WriteHeader(v)
}

// String responds with the given string.
func String(c *core.Context, s string) {
	c.ResponseWriter.Write([]byte(s))
}

// Bytes responds with the given slice of byte.
func Bytes(c *core.Context, b []byte) {
	c.ResponseWriter.Write(b)
}

// JSON responds with the marshalled content with Content-Type header.
func JSON(c *core.Context, v interface{}) {
	c.ResponseWriter.Header().Set("Content-Type", "application/json")

	var js []byte
	var err error

	if core.Production {
		js, err = json.Marshal(v)
	} else {
		js, err = json.MarshalIndent(v, "", "\t")
	}

	if err != nil {
		log.Println(err)
		http.Error(c.ResponseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	c.ResponseWriter.Write(js)
}

// View responds with the given template with data.
func View(c *core.Context, name string, data interface{}) {
	err := views.ExecuteTemplate(c.ResponseWriter, name, data)
	if err != nil {
		log.Println(err)
		http.Error(c.ResponseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
