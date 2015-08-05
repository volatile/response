package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/volatile/core"
	"github.com/volatile/core/log"
)

const viewsDir = "views"

var views *template.Template

func init() {
	if _, err := os.Stat(viewsDir); err != nil {
		return
	}

	views = template.New("views")

	// Built-in views funcs
	views.Funcs(template.FuncMap{
		"html":  viewsFuncHTML,
		"nl2br": viewsFuncNL2BR,
	})

	core.BeforeRun(func() {
		if err := filepath.Walk(viewsDir, viewsWalk); err != nil {
			panic(fmt.Errorf("response: %v", err))
		}
	})
}

// walk is the path/filepath.WalkFunc used to walk viewsDir in order to initialize views.
// It will try to parse all files it encounters and recurse into subdirectories.
func viewsWalk(path string, f os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if f.IsDir() {
		return nil
	}

	_, err = views.ParseFiles(path)

	return err
}

// FuncMap is the type of the map defining the mapping from names to functions.
// Each function must have either a single return value, or two return values of which the second has type error.
// In that case, if the second (error) argument evaluates to non-nil during execution, execution terminates and Execute returns that error.
// FuncMap has the same base type as FuncMap in "text/template", copied here so clients need not import "text/template".
type FuncMap map[string]interface{}

// ViewsFuncs adds a function that will be available to all templates.
func ViewsFuncs(funcMap FuncMap) {
	if views == nil {
		panic(`response: views can't be used without a "views" directory`)
	}
	views.Funcs(template.FuncMap(funcMap))
}

// Status responds with the given status code.
func Status(c *core.Context, v int) {
	http.Error(c.ResponseWriter, http.StatusText(v), v)
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
	if err := json.NewEncoder(c.ResponseWriter).Encode(v); err != nil {
		log.Stack(err)
		http.Error(c.ResponseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// View pass the data to the template associated to name, and responds with it.
func View(c *core.Context, name string, data map[string]interface{}) {
	if views == nil {
		log.Stack(errors.New(`views can't be used without a "views" directory`))
	}
	data["c"] = c
	if err := views.ExecuteTemplate(c.ResponseWriter, name, data); err != nil {
		log.Stack(err)
		http.Error(c.ResponseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
