package vanity

import (
	"fmt"
	"math/rand"
	"net/http"
)

const defaultFaviconSVG = `<svg width="32" height="32"
	xmlns="http://www.w3.org/2000/svg">
	<rect width="32" height="32" style="fill:rgba(%d,%d,%d,%d)" />
</svg>`

func Favicon(w http.ResponseWriter, _ *http.Request) {
	r, g, b, a := rand.Int31n(256), rand.Int31n(256), rand.Int31n(256), rand.Int31n(256)
	w.Header().Add("Content-Type", "image/svg+xml")
	fmt.Fprintf(w, defaultFaviconSVG, r, g, b, a)
}
