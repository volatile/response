package response

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/volatile/core"
	"github.com/volatile/core/httputil"
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
			panic("response: " + err.Error())
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
	StringStatus(c, http.StatusOK, s)
}

// StringStatus responds with the given string and status code.
func StringStatus(c *core.Context, code int, s string) {
	httputil.SetDetectedContentType(c.ResponseWriter, []byte(s))
	c.ResponseWriter.WriteHeader(code)
	c.ResponseWriter.Write([]byte(s))
}

// Bytes responds with the given slice of byte.
func Bytes(c *core.Context, b []byte) {
	BytesStatus(c, http.StatusOK, b)
}

// BytesStatus responds with the given slice of byte and status code.
func BytesStatus(c *core.Context, code int, b []byte) {
	httputil.SetDetectedContentType(c.ResponseWriter, b)
	c.ResponseWriter.WriteHeader(code)
	c.ResponseWriter.Write(b)
}

// JSON set the correct header and responds with the marshalled content.
func JSON(c *core.Context, v interface{}) {
	JSONStatus(c, http.StatusOK, v)
}

// JSONStatus set the correct header and responds with the marshalled content and status code.
func JSONStatus(c *core.Context, code int, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	c.ResponseWriter.Header().Set("Content-Type", "application/json")
	c.ResponseWriter.WriteHeader(code)
	c.ResponseWriter.Write(b)
}

// View pass the data to the template associated to name, and responds with it.
func View(c *core.Context, name string, data map[string]interface{}) {
	if views == nil {
		panic(`response: views can't be used without a "views" directory`)
	}

	if data == nil {
		data = make(map[string]interface{})
	}
	data["c"] = c

	c.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.ExecuteTemplate(c.ResponseWriter, name, data); err != nil {
		panic(err)
	}
}
