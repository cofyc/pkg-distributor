package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cofyc/pkg-distributor/pkg/utils"
	"github.com/golang/glog"
	"github.com/gorilla/handlers"
)

var (
	optListen    string
	optDir       string
	optBasicAuth string
)

func init() {
	flag.StringVar(&optListen, "listen", "0.0.0.0:1973", "host and port to listen on (default: 0.0.0.0:1973)")
	flag.StringVar(&optDir, "dir", "", "repo directory")
	flag.StringVar(&optBasicAuth, "basic-auth", "", "basic auth info (e.g. user:pass)")
	flag.Set("logtostderr", "true")
}

// upload handles "/v1/upload", it receives a file and add it into repo.
func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "only POST is allowed", http.StatusBadRequest)
		return
	}
	r.ParseMultipartForm(32 << 20)
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	glog.Infof("header.Filename: %s, header.Header: %v, header.Size: %d", header.Filename, header.Header, header.Size)
	tmpfile := fmt.Sprintf("/tmp/%s", header.Filename)
	err = utils.Store(tmpfile, file)
	if err != nil {
		glog.Errorf("error: %v", err)
		return
	}
	for _, distro := range []string{"xenial"} {
		args := []string{"includedeb", distro, tmpfile}
		cmd := exec.Command("reprepro", args...)
		cmd.Dir = filepath.Join(optDir, "apt")
		cmd.Stdout = w
		cmd.Stderr = w
		err = cmd.Run()
		if err != nil {
			glog.Errorf("error: %v", err)
			return
		}
	}
	return
}

func main() {
	flag.Parse()

	basicAuthPairs := make(map[string]string)
	if optBasicAuth != "" {
		seps := strings.Split(optBasicAuth, ":")
		if len(seps) != 2 || len(seps[0]) <= 0 || len(seps[0]) <= 0 {
			glog.Fatalf("invalid basic auth option: %s", optBasicAuth)
		}
		basicAuthPairs[seps[0]] = seps[1]
	}

	var uploadHandler http.Handler
	uploadHandler = http.HandlerFunc(upload)
	if len(basicAuthPairs) > 0 {
		glog.Infof("basic auth is enabled")
		uploadHandler = utils.NewBasicAuthHandler("pkg-distributor", basicAuthPairs)(uploadHandler)
	}
	http.Handle("/v1/upload", uploadHandler)
	http.Handle("/", http.FileServer(http.Dir(optDir)))
	serveMux := handlers.LoggingHandler(os.Stderr, http.DefaultServeMux)
	glog.Infof("listen on %s", optListen)
	glog.Fatal(http.ListenAndServe(optListen, serveMux))
}
