package response

import (
	"html/template"
	"strings"
)

func viewsFuncHTML(s string) template.HTML {
	return template.HTML(s)
}

func viewsFuncNL2BR(s string) template.HTML {
	return template.HTML(strings.Replace(s, "\n", "<br>", -1))
}
