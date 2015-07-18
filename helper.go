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

// FuncMap is the type of the map defining the mapping from names to functions.
// Each function must have either a single return value, or two return values of which the second has type error.
// In that case, if the second (error) argument evaluates to non-nil during execution, execution terminates and Execute returns that error.
// FuncMap has the same base type as FuncMap in "text/template", copied here so clients need not import "text/template".
type FuncMap map[string]interface{}

// ViewsFuncs adds a function that will be available to all templates.
func ViewsFuncs(funcMap FuncMap) {
	views.Funcs(template.FuncMap(funcMap))
}

// Status responds with the given status code.
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

// JSON set the correct header and responds with the marshalled content.
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

// View pass the data to the template from namea and responds with it.
func View(c *core.Context, name string, data interface{}) {
	err := views.ExecuteTemplate(c.ResponseWriter, name, data)
	if err != nil {
		log.Println(err)
		http.Error(c.ResponseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
