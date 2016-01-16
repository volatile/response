/*
Package response is a helper for the Core (https://github.com/volatile/core).
It provides syntactic sugar that lets you easily write responses on the context.

Installation

In the terminal:

	$ go get github.com/volatile/response

Usage

Example:

	package main

	import (
		"net/http"
		"strings"

		"github.com/volatile/core"
		"github.com/volatile/response"
		"github.com/volatile/route"
	)

	func main() {
		// We set functions for templates.
		response.TemplatesFuncs(response.FuncMap{
			"toUpper": strings.ToUpper,
			"toLower": strings.ToLower,
		})

		// Status response
		route.Get("^/status$", func(c *core.Context) {
			response.Status(c, http.StatusForbidden)
		})

		// String response
		route.Get("^/string$", func(c *core.Context) {
			response.String(c, "Hello, World!")
		})

		// Bytes response
		route.Get("^/bytes$", func(c *core.Context) {
			response.Bytes(c, []byte("Hello, World!"))
		})

		// JSON response
		route.Get("^/json$", func(c *core.Context) {
			response.JSON(c, &Car{
				ID:    1,
				Brand: "Bentley",
				Model: "Continental GT",
			})
		})

		// Templates response
		route.Get("^/(?P<name>[A-Za-z]+)$", func(c *core.Context, params map[string]string) {
			response.Templates(c, "hello", params)
		})

		core.Run()
	}

	type Car struct {
		ID    int    `json:"id"`
		Brand string `json:"brand"`
		Model string `json:"model"`
	}

Templates

The templates templates are recursively parsed from the "templates" directory, just before running the server.
Filenames (including extensions) and directory organization doesn't matter. All the files are parsed.

Built-in functions

This package gives some functions out-of-the-box:

● "html" uses a raw string without escaping it

● "nl2br" replaces "\n" by "<br>"

Custom functions

To use functions in your templates, use TemplatesFuncs, like with the `html/template` standard package.
*/
package response
