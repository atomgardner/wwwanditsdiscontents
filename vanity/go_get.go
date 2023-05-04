package vanity

import (
	"html/template"
	"net/http"
)

var tmplGoGet *template.Template

func init() {
	tmpl := `<html>
<head>
	<meta name="go-import" content="{{ .Name }}  git {{ .Remote }}">
	<meta http-equiv="refresh" content="0;URL='{{ .Remote }}'">
</head>
</html>`

	tmplGoGet = template.Must(template.New("go_get").Parse(tmpl))
}

func GoGet(name, remote string, h http.Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.URL.Query()["go-get"]; !ok {
			h.ServeHTTP(w, r)
			return
		}
		tmplGoGet.Execute(w, struct {
			Name, Remote string
		}{
			Name:   name,
			Remote: remote,
		})
	}
}
