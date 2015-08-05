<p align="center"><img src="http://volatile.whitedevops.com/images/repositories/response/logo.png" alt="Volatile Response" title="Volatile Response"><br><br></p>

Volatile Response is a helper for the [Core](https://github.com/volatile/core).  
It provides syntactic sugar that lets you easily write responses on `*core.Context`.

## Installation

```Shell
$ go get github.com/volatile/response
```

## Usage

```Go
package main

import (
	"net/http"
	"strings"

	"github.com/volatile/core"
	"github.com/volatile/response"
	"github.com/volatile/route"
)

func main() {
	// Functions for views templates
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
```

In `views/hello.gohtml`:

```HTML
{{define "hello"}}
	<!DOCTYPE html>
	<html>
		<head>
			<title>Hello</title>
		</head>
		<body>
			Hello, {{toUpper .name}}!
		</body>
	</html>
{{end}}
```

[![GoDoc](https://godoc.org/github.com/volatile/response?status.svg)](https://godoc.org/github.com/volatile/response)

### Views

The views templates are recursively parsed from the `views` directory, just before running the server.  
The file names and extensions and the directory organization doesn't matter. All the files are parsed.

#### Built-in functions

This package gives some functions out-of-the-box:

- `html` uses a raw string without escaping it
- `nl2br` replaces `\n` by `<br>\n`

#### Custom functions

To use functions in your views, use `response.ViewsFuncs(response.FuncMap{})`, like with the `html/template` standard package.
