/*
Package response is a helper for the core (https://godoc.org/github.com/volatile/core).
It provides syntactic sugar that lets you easily write responses on the context.

Templates

Templates are recursively parsed from the "templates" directory, just before running the server.
Filenames (including extensions) and directory organization doesn't matter: all the files are parsed.

Global functions are predefined:

	html	uses a raw string without escaping it
	nl2br	replaces new lines by "<br>"

TemplatesFuncs sets custom templates functions:

	response.TemplatesFuncs(response.FuncMap{
		"toUpper": strings.ToUpper,
		"toLower": strings.ToLower,
	})

TemplatesData sets custom templates functions:

	response.TemplatesData(response.DataMap{
		"toUpper": strings.ToUpper,
		"toLower": strings.ToLower,
	})

Template responds with a template:

	response.Template(c, "hello", response.DataMap{
		"foo": "bar",
	})

The context is always part of the data under key "c".
*/
package response
