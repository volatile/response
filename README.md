<p align="center"><img src="https://cloud.githubusercontent.com/assets/9503891/8758179/67dfe466-2ce1-11e5-9539-2740ca0179bf.png" alt="Volatile Response" title="Volatile Response"><br><br></p>

Volatile Response is a helper for the [Core](https://github.com/volatile/core).  
It provides syntactic sugar that let's you easily write responses on `*core.Context`.

## Installation

```Shell
$ go get -u github.com/volatile/response
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
