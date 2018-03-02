package createrepo

import (
	"os/exec"
	"path/filepath"

	"github.com/golang/glog"
)

type createRepo struct {
}

func (cr *createRepo) Update(repo string) (err error) {
	args := []string{"--update", "--database", repo}
	glog.Infof("executing command: %s, %s", "createrepo", args)
	cmd := exec.Command("createrepo", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		glog.Errorf("output: %s", string(output))
	}
	return
}

func (cr *createRepo) SignRPM(rpmfile string) (err error) {
	args := []string{rpmfile}
	glog.Infof("executing command: %s, %s", "rpmautosign", args)
	cmd := exec.Command("rpmautosign", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		glog.Errorf("output: %s", string(output))
	}
	return
}

func (cr *createRepo) SignRepo(repo string) (err error) {
	repomdxml := filepath.Join(repo, "repodata/repomd.xml")
	args := []string{"--yes", "--detach-sign", "--armor", repomdxml}
	glog.Infof("executing command: %s, %s", "gpg", args)
	cmd := exec.Command("gpg", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		glog.Errorf("output: %s", string(output))
	}
	return
}

func NewCreateRepo() *createRepo {
	return &createRepo{}
}
