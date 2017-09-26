package aptly

import (
	"os/exec"
	"strings"
)

type aptly struct {
}

func (a *aptly) RepoCreate(repo string) (err error) {
	args := []string{"repo", "create", repo}
	cmd := exec.Command("aptly", args...)
	_, err = cmd.CombinedOutput()
	return
}

func (a *aptly) RepoAdd(repo, debfile string) (err error) {
	args := []string{"repo", "add", repo, debfile}
	cmd := exec.Command("aptly", args...)
	_, err = cmd.CombinedOutput()
	return
}

func (a *aptly) RepoList() (repos []string, err error) {
	repos = make([]string, 0)
	args := []string{"repo", "list", "-raw"}
	cmd := exec.Command("aptly", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(output), "\n") {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) < 1 {
			continue
		}
		repos = append(repos, fields[0])
	}
	return
}

func (a *aptly) PublishUpdate(distribution string) (err error) {
	args := []string{"publish", "update", distribution}
	cmd := exec.Command("aptly", args...)
	_, err = cmd.CombinedOutput()
	return
}

func (a *aptly) PublishRepo(repo, distribution string) (err error) {
	args := []string{"publish", "repo", "-distribution", distribution, repo}
	cmd := exec.Command("aptly", args...)
	_, err = cmd.CombinedOutput()
	return
}

func (a *aptly) PublishList(distribution string) (publishes []string, err error) {
	publishes = make([]string, 0)
	args := []string{"publish", "list", "-raw"}
	cmd := exec.Command("aptly", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(output), "\n") {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) < 2 {
			continue
		}
		publishes = append(publishes, fields[1])
	}
	return
}

func NewAptly() *aptly {
	return &aptly{}
}
