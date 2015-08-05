/*
Package response is a helper for the Core.
It provides syntactic sugar that lets you easily write responses on *core.Context.

Usage

Full example with status, string, bytes, JSON and view responses:

	package main

	import (
		"net/http"
		"strings"

		"github.com/volatile/core"
		"github.com/volatile/response"
		"github.com/volatile/route"
	)

	func main() {
		// We set functions for views templates.
		response.ViewsFuncs(response.FuncMap{
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

		// View response
		route.Get("^/(?P<name>[A-Za-z]+)$", func(c *core.Context, params map[string]string) {
			response.View(c, "hello", params)
		})

		core.Run()
	}

	type Car struct {
		ID    int    `json:"id"`
		Brand string `json:"brand"`
		Model string `json:"model"`
	}

Views

The views templates are recursively parsed from the "views" directory, just before running the server.
The file names and extensions and the directory organization doesn't matter. All the files are parsed.

Built-in functions

This package gives some functions out-of-the-box:

- "html" uses a raw string without escaping it
- "nl2br" replaces "\n" by "<br>"

Custom functions

To use functions in your views, use response.ViewsFuncs(response.FuncMap{}), like with the `html/template` standard package.
*/
package response
