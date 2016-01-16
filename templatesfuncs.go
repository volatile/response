package response

import (
	"html/template"
	"strings"
)

func templatesFuncHTML(s string) template.HTML {
	return template.HTML(s)
}

func templatesFuncNL2BR(s string) template.HTML {
	return template.HTML(strings.Replace(s, "\n", "<br>", -1))
}
