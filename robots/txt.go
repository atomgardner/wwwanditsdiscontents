package robots

import (
	"fmt"
	"net/http"
)

const robotsTXT = `User-agent: *
Disallow: /
`

func Disallow(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, robotsTXT)
}
