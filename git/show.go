package git

import (
	"log"
	"net/http"
	"os/exec"
	"strings"
)

const (
	errNotFound = "this is not the blob you're looking for"
	wellKnown   = ".well-known/"
)

func Show(path, defaultBranch, commitFormat string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		slash := strings.Index(r.URL.Path, "/")
		id := r.URL.Path[slash+1:]
		f := commitFormat
		switch {
		case id == "":
			id = defaultBranch
		case len(wellKnown) < len(id) && id[:len(wellKnown)] == wellKnown:
			f = "%B"
			id = id[1:]
		case id[0] == '-':
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errNotFound))
			return
		// TODO: send an error for files created by `git init --bare`
		default:
		}
		cmd := exec.Command("git", "-C", path, "show", "--format="+f, id)
		cmd.Stdout = w
		if err := cmd.Run(); err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errNotFound))
			log.Println(err)
		}
	}
}

func CheckRepository(path string) error {
	return exec.Command("git", "-C", path, "branch").Run()
}
