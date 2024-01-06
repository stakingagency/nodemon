package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type githubRelease struct {
	TagName    string `json:"tag_name"`
	Name       string `json:"name"`
	Draft      bool   `json:"draft"`
	PreRelease bool   `json:"prerelease"`
}

func GetGoVersion() (string, error) {
	cmd := exec.CommandContext(context.Background(), "go", "version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	parts := strings.Split(string(output), " ")
	for _, part := range parts {
		if !strings.Contains(part, ".") {
			continue
		}

		return strings.TrimPrefix(part, "go"), nil
	}

	return "", nil
}

func GetLatestVersion(repo string) (string, error) {
	bytes, err := GetHTTP(fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo), "")
	if err != nil {
		return "", err
	}

	latest := &githubRelease{}
	err = json.Unmarshal(bytes, &latest)
	if err != nil {
		return "", err
	}

	return latest.TagName, nil
}

func SelfUpdate(repoBinaryPath string) error {
	parts := strings.Split(repoBinaryPath, "/")
	repoWithToken := strings.Join(parts[:3], "/")
	githubRepo := strings.Join(parts[1:3], "/")
	binaryPath := strings.Join(parts[3:], "/")
	repo := repoWithToken
	if strings.Contains(repo, "@") && strings.Index(repo, "@") < strings.Index(repo, "/") {
		parts = strings.Split(repoWithToken, "@")
		repo = strings.Join(parts[1:], "@")
	}

	parts = strings.Split(os.Args[0], "/")
	appName := parts[len(parts)-1]
	appPath, err := os.Getwd()
	if err != nil {
		log.Error("get current path", "error", err)
		return err
	}

	goPath := strings.Split(os.Getenv("GOPATH"), ":")[0]

	// get latest tag
	latestVersion, err := GetLatestVersion(githubRepo)
	if err != nil {
		log.Error("get latest version", "error", err)
		return err
	}

	// clone repo
	fullRepoPath := goPath + "/" + repo
	if fullRepoPath+"/"+binaryPath == appPath {
		return errors.New(appName + " runs from the repo path")
	}

	parts = strings.Split(fullRepoPath, "/")
	repoPath := strings.Join(parts[:len(parts)-1], "/")
	os.MkdirAll(repoPath, os.ModePerm)

	defer os.Chdir(appPath)

	if exists, _ := pathExists(fullRepoPath); exists {
		err = os.Chdir(fullRepoPath)
		if err != nil {
			log.Error("chdir fullRepoPath", "error", err, "path", fullRepoPath)
			return err
		}

		err = exec.CommandContext(context.Background(), "git", "stash").Run()
		if err != nil {
			log.Error("run git stash", "error", err, "path", fullRepoPath)
			return err
		}

		err = exec.CommandContext(context.Background(), "git", "pull").Run()
		if err != nil {
			log.Error("run git pull", "error", err, "path", fullRepoPath)
			return err
		}
	} else {
		err = os.Chdir(repoPath)
		if err != nil {
			log.Error("chdir repoPath", "error", err, "path", repoPath)
			return err
		}

		err = exec.CommandContext(context.Background(), "git", "clone", "https://"+repoWithToken, "--branch="+latestVersion, "--single-branch", "--depth=1").Run()
		if err != nil {
			log.Error("run git clone", "error", err, "path", repoPath)
			return err
		}
	}

	// build app
	buildPath := fullRepoPath + "/" + binaryPath
	err = os.Chdir(buildPath)
	if err != nil {
		log.Error("chdir buildPath", "error", err, "path", buildPath)
		return err
	}

	tmpFilename := fmt.Sprintf("tmp%v", time.Now().Unix())
	err = exec.CommandContext(context.Background(), "go", "build", "-v", "-o", tmpFilename, "-ldflags", `-X main.appVersion=`+latestVersion).Run()
	if err != nil {
		log.Error("run go build", "error", err, "path", buildPath)
		return err
	}

	err = os.Rename(tmpFilename, appPath+"/"+appName)
	if err != nil {
		log.Error("rename file", "error", err, "tmpFile", tmpFilename)
		return err
	}

	err = os.Chdir(appPath)
	if err != nil {
		log.Error("chdir appPath", "error", err, "path", appPath)
		return err
	}

	os.Exit(0)

	return nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func UpdateOS() error {
	err := exec.CommandContext(context.Background(), "sudo", "apt", "update").Run()
	if err != nil {
		return err
	}

	err = runCommandWithInput("sudo", []string{"apt", "upgrade"}, "y")
	if err != nil {
		return err
	}

	err = runCommandWithInput("sudo", []string{"apt", "autoremove"}, "y")
	if err != nil {
		return err
	}

	return nil
}

func runCommandWithInput(command string, args []string, answer string) error {
	cmd := exec.CommandContext(context.Background(), command, args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	_, err = stdin.Write([]byte(answer + "\n"))
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

func RebootHost() error {
	return exec.CommandContext(context.Background(), "sudo", "shutdown", "-r", "now").Run()
}
