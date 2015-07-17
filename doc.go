/*
Package response is a helper for the Core.
It provides syntactic sugar that let's you easily write responses on *core.Context.

Usage

Here is the full example status, string, bytes, JSON and view responses:

	package main

	import (
		"net/http"
		"strings"

		"github.com/volatile/core"
		"github.com/volatile/response"
		"github.com/volatile/route"
	)

	func main() {
		response.AddViewFunc("toUpper", strings.ToUpper)

		// Status
		route.Get("^/status$", func(c *core.Context) {
			response.Status(c, http.StatusForbidden)
		})

		// String
		route.Get("^/string$", func(c *core.Context) {
			response.String(c, "Hello, World!")
		})

		// Bytes
		route.Get("^/bytes$", func(c *core.Context) {
			response.Bytes(c, []byte("Hello, World!"))
		})

		// JSON
		route.Get("^/json$", func(c *core.Context) {
			response.JSON(c, &Car{
				ID:    1,
				Brand: "Bentley",
				Model: "Continental GT",
			})
		})

		// View
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
*/
package response
