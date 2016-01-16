package response

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/volatile/core"
	"github.com/volatile/core/httputil"
)

const templatesDir = "templates"

var (
	// ErrNoTemplatesDir is used when a template feature is used without having the templates directory.
	ErrNoTemplatesDir = fmt.Errorf("response: templates can't be used without a %q directory", templatesDir)

	templates *template.Template
)

func init() {
	if _, err := os.Stat(templatesDir); err != nil {
		return
	}

	templates = template.New(templatesDir)

	// Built-in templates funcs
	templates.Funcs(template.FuncMap{
		"html":  templatesFuncHTML,
		"nl2br": templatesFuncNL2BR,
	})

	core.BeforeRun(func() {
		if err := filepath.Walk(templatesDir, templatesWalk); err != nil {
			panic("response: " + err.Error())
		}
	})
}

// walk is the path/filepath.WalkFunc used to walk templatesDir in order to initialize templates.
// It will try to parse all files it encounters and recurse into subdirectories.
func templatesWalk(path string, f os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if f.IsDir() {
		return nil
	}

	_, err = templates.ParseFiles(path)
	return err
}

// FuncMap is the type of the map defining the mapping from names to functions.
// Each function must have either a single return value, or two return values of which the second has type error.
// In that case, if the second (error) argument evaluates to non-nil during execution, execution terminates and Execute returns that error.
// FuncMap has the same base type as FuncMap in "text/template", copied here so clients need not import "text/template".
type FuncMap map[string]interface{}

// TemplatesFuncs adds a function that will be available to all templates.
func TemplatesFuncs(funcMap FuncMap) {
	if templates == nil {
		panic(ErrNoTemplatesDir)
	}
	templates.Funcs(template.FuncMap(funcMap))
}

// Status responds with the status code.
func Status(c *core.Context, code int) {
	http.Error(c.ResponseWriter, http.StatusText(code), code)
}

// String responds with the string s.
func String(c *core.Context, s string) {
	StringStatus(c, http.StatusOK, s)
}

// StringStatus responds with the status code and the string s.
func StringStatus(c *core.Context, code int, s string) {
	httputil.SetDetectedContentType(c.ResponseWriter, []byte(s))
	c.ResponseWriter.WriteHeader(code)
	c.ResponseWriter.Write([]byte(s))
}

// Bytes responds with the slice of bytes b.
func Bytes(c *core.Context, b []byte) {
	BytesStatus(c, http.StatusOK, b)
}

// BytesStatus responds with the status code and the slice of bytes b.
func BytesStatus(c *core.Context, code int, b []byte) {
	httputil.SetDetectedContentType(c.ResponseWriter, b)
	c.ResponseWriter.WriteHeader(code)
	c.ResponseWriter.Write(b)
}

// JSON responds with the JSON marshalled v.
func JSON(c *core.Context, v interface{}) {
	JSONStatus(c, http.StatusOK, v)
}

// JSONStatus responds with the status code and the JSON marshalled v.
func JSONStatus(c *core.Context, code int, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	c.ResponseWriter.Header().Set("Content-Type", "application/json")
	c.ResponseWriter.WriteHeader(code)
	c.ResponseWriter.Write(b)
}

// Template responds with the template associated to name.
func Template(c *core.Context, name string, data map[string]interface{}) {
	TemplateStatus(c, http.StatusOK, name, data)
}

// TemplateStatus responds with the status code and the template associated to name.
func TemplateStatus(c *core.Context, code int, name string, data map[string]interface{}) {
	if templates == nil {
		panic(ErrNoTemplatesDir)
	}

	if data == nil {
		data = make(map[string]interface{})
	}
	data["c"] = c

	c.ResponseWriter.Header().Set("Content-Type", "text/html; charsets=utf-8")
	c.ResponseWriter.WriteHeader(code)
	if err := ExecuteTemplate(c.ResponseWriter, name, data); err != nil {
		panic(err)
	}
}

// ExecuteTemplate works like the standard html/template.Template.ExecuteTemplate function.
func ExecuteTemplate(wr io.Writer, name string, data map[string]interface{}) error {
	return templates.ExecuteTemplate(wr, name, data)
}
