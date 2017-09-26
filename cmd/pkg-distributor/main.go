package main

import (
	"flag"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cofyc/pkg-distributor/pkg/aptly"
	"github.com/cofyc/pkg-distributor/pkg/utils"
	"github.com/golang/glog"
	"github.com/gorilla/handlers"
)

var (
	optListen    string
	optBasicAuth string
	repo         string = "stable"
	dataDir      string = "/data"
	publicDir    string = "/data/public"
	filesDir     string = "/data/files"
)

func init() {
	flag.StringVar(&optListen, "listen", "0.0.0.0:1973", "host and port to listen on (default: 0.0.0.0:1973)")
	flag.StringVar(&optBasicAuth, "basic-auth", "", "basic auth info (e.g. user:pass)")
	flag.Set("logtostderr", "true")
}

func inArray(v string, ss []string) bool {
	for _, s := range ss {
		if s == v {
			return true
		}
	}
	return false
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
	tmpfile := filepath.Join(filesDir, header.Filename)
	err = utils.Store(tmpfile, file, true)
	if err != nil {
		glog.Errorf("error: %v", err)
		return
	}
	aptly := aptly.NewAptly()
	// create repo if does not exist
	repos, err := aptly.RepoList()
	if err != nil {
		glog.Errorf("failed to list repos: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !inArray(repo, repos) {
		err = aptly.RepoCreate(repo)
		if err != nil {
			glog.Errorf("failed to create repo %s: %v", repo, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// add deb into repo
	err = aptly.RepoAdd(repo, tmpfile)
	if err != nil {
		glog.Errorf("failed to add %s into repo %s: %v", tmpfile, repo, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// publish
	for _, distro := range []string{"xenial"} {
		publishes, err := aptly.PublishList(distro)
		if err != nil {
			glog.Errorf("failed to list publishes %s: %v", distro, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if inArray(distro, publishes) {
			err = aptly.PublishUpdate(distro)
		} else {
			err = aptly.PublishRepo(repo, distro)
		}
		if err != nil {
			glog.Errorf("failed to publish repo %s: %v", repo, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	return
}

func main() {
	flag.Parse()

	if os.Getenv("DATA_DIR") != "" {
		dataDir = os.Getenv("DATA_DIR")
		publicDir = filepath.Join(dataDir, "public")
		filesDir = filepath.Join(dataDir, "files")
	}

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
	http.Handle("/", http.FileServer(http.Dir(publicDir)))
	serveMux := handlers.LoggingHandler(os.Stderr, http.DefaultServeMux)
	glog.Infof("listen on %s", optListen)
	glog.Fatal(http.ListenAndServe(optListen, serveMux))
}
