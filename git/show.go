package git

import (
	"log"
	"net/http"
	"os/exec"
	"strings"
)

const errNotFound = "this is not the blob you're looking for"

func Show(repo, defaultBranch, commitFormat string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		slash := strings.Index(r.URL.Path, "/")
		id := r.URL.Path[slash+1:]
		switch {
		case id == "":
			id = defaultBranch
		case id[0] == '-':
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errNotFound))
			return
		// TODO: send an error for files created by `git init --bare`
		default:
		}
		cmd := exec.Command("git", "-C", repo, "show", "--format="+commitFormat, id)
		cmd.Stdout = w
		if err := cmd.Run(); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errNotFound))
		}
	}
}
