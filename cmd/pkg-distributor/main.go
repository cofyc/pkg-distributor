package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/golang/glog"
	"github.com/gorilla/handlers"
)

var (
	optListen string
	optDir    string
)

func init() {
	flag.StringVar(&optListen, "listen", "0.0.0.0:1973", "host and port to listen on (default: 0.0.0.0:1973)")
	flag.StringVar(&optDir, "dir", "", "repo directory")
	flag.Set("logtostderr", "true")
}

// store stores string read from `body` into file at `dstfile`.
func store(dstfile string, body io.ReadCloser) error {
	dstfile = filepath.Clean(dstfile)
	tmpfile := filepath.Join(filepath.Dir(dstfile), fmt.Sprintf(".%s.tmp", filepath.Base(dstfile)))
	// Check file exists or not.
	if _, err := os.Stat(tmpfile); !os.IsNotExist(err) {
		return fmt.Errorf("%s does exist", tmpfile)
	}
	if _, err := os.Stat(dstfile); !os.IsNotExist(err) {
		return fmt.Errorf("%s does exist", dstfile)
	}
	// Open.
	f, err := os.OpenFile(tmpfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	// Copy.
	if _, err := io.Copy(f, body); err != nil {
		os.Remove(tmpfile)
		return err
	}
	// Sync.
	if err := f.Sync(); err != nil {
		os.Remove(tmpfile)
		return err
	}
	// Rename.
	if err := os.Rename(tmpfile, dstfile); err != nil {
		os.Remove(tmpfile)
		return err
	}
	return nil
}

// upload handles "/v1/upload", it receives a file and add it into repo.
func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "only POST is allowed\n")
		return
	}
	r.ParseMultipartForm(32 << 20)
	file, header, err := r.FormFile("file")
	if err != nil {
		glog.Errorf("error: %v", err)
		return
	}
	glog.Infof("header.Filename: %s, header.Header: %v, header.Size: %d", header.Filename, header.Header, header.Size)
	tmpfile := fmt.Sprintf("/tmp/%s", header.Filename)
	err = store(tmpfile, file)
	if err != nil {
		glog.Errorf("error: %v", err)
		return
	}
	for _, distro := range []string{"trusty", "xenial"} {
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

	http.HandleFunc("/v1/upload", upload)
	http.Handle("/", http.FileServer(http.Dir(optDir)))
	serveMux := handlers.LoggingHandler(os.Stderr, http.DefaultServeMux)

	glog.Infof("listen on %s", optListen)
	glog.Fatal(http.ListenAndServe(optListen, serveMux))
}
