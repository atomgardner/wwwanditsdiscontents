// SPDX-License-Identifier: CC0-1.0
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"atomgardner.com/wwwanditsdiscontents/git"
	"atomgardner.com/wwwanditsdiscontents/multi"
	"atomgardner.com/wwwanditsdiscontents/robots"
	"atomgardner.com/wwwanditsdiscontents/vanity"
	"atomgardner.com/wwwanditsdiscontents/with"

	"golang.org/x/crypto/acme/autocert"
)

var (
	repo          = flag.String("repo", ".", "path to git repo")
	defaultBranch = flag.String("branch", "main", "default branch")
	commitFormat  = flag.String("format", "%h/%cI%n%s%n%n%b", "default format")

	port           = flag.String("port", "http", "listen for non-tls connections on this port")
	tlsPort        = flag.String("tls-port", "https", "listen for tls connections on this port")
	disallowRobots = flag.Bool("disallow-robots", true, "essentially whether to disable crawlers")
	randoFavicon   = flag.Bool("favicon", true, "use a random favicon")

	// TODO: find a way to store config in git refs.
	hosts  multi.String
	assets multi.String
)

func init() {
	flag.Var(&hosts, "auto-cert-host", "hostname for Let's Encrypt! Automatic Certificate Management")
	flag.Var(&assets, "asset-dir", "path to a directory served by an http file server.")
	flag.Parse()

	err := exec.Command("git", "-C", *repo, "branch").Run()
	if err != nil {
		log.Printf("`%s` is not a git repo\n", *repo)
		os.Exit(1)
	}
}

func main() {
	mux := http.NewServeMux()
	show := git.Show(*repo, *defaultBranch, *commitFormat)
	patterns := map[string]http.Handler{"/": http.HandlerFunc(show)}

	if *randoFavicon {
		patterns["/favicon.ico"] = http.HandlerFunc(vanity.Favicon)
	}
	if *disallowRobots {
		patterns["/robots.txt"] = http.HandlerFunc(robots.Disallow)
		patterns["/.well-known/robots.txt"] = http.HandlerFunc(robots.Disallow)
	}
	for _, dir := range assets {
		slash := strings.LastIndex(dir, "/")
		if slash == -1 || slash+1 == len(dir) {
			log.Printf("skipping asset dir `%s`: unusable final slash", dir)
			continue
		}
		patterns[dir[slash+1:]+"/"] = http.FileServer(http.Dir(dir))
	}
	for pattern, handler := range patterns {
		mux.Handle(pattern, with.Feedback(handler))
	}

	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(hosts...),
		Cache:      autocert.DirCache("certs"),
	}

	go http.ListenAndServe(":"+*port, m.HTTPHandler(with.Feedback(http.HandlerFunc(show))))
	tlsServer := &http.Server{Handler: mux, Addr: ":" + *tlsPort, TLSConfig: m.TLSConfig()}
	log.Fatal(tlsServer.ListenAndServeTLS("", ""))
}
