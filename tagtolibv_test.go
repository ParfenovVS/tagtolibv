package tagtolibv_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ParfenovVS/tagtolibv"
)

func TestGetCurrentBranch(t *testing.T) {
	defaultWd, _ := os.Getwd()
	defer os.Chdir(defaultWd)

	dir, err := createTempRepository(".temp_TestGetCurrentBranch")
	if err != nil {
		t.Fatalf("cannot create temp repository: %s", err.Error())
	}
	defer os.RemoveAll(dir)

	exp := "master"

	branch, err := tagtolibv.GetCurrentBranch()
	if err != nil {
		t.Errorf("failed with %q", err)
	} else if exp != branch {
		t.Errorf("expected %q, got %q instead.", exp, branch)
	}
}

func TestGetTags(t *testing.T) {
	defaultWd, _ := os.Getwd()
	defer os.Chdir(defaultWd)

	dir, err := createTempRepository(".temp_TestGetTags")
	if err != nil {
		t.Fatalf("cannot create temp repository: %s", err.Error())
	}
	defer os.RemoveAll(dir)

	expTags := []string{
		"v1.0.0",
		"v1.1.0-beta",
	}

	os.Create("temp")
	cmd := exec.Command("git", "add", "temp")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("cannot add to git: %s, %s", err.Error(), string(out))
	}
	cmd = exec.Command("git", "commit", "-m", "\"add temp\"")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("cannot commit: %s, %s", err.Error(), string(out))
	}

	for _, tag := range expTags {
		cmd := exec.Command("git", "tag", tag)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("cannot create tag %s: %s, %s", tag, err.Error(), string(out))
		}
	}

	tags, err := tagtolibv.GetTags()
	if err != nil {
		t.Fatalf("cannot get tags: %s", err.Error())
	}

	if !reflect.DeepEqual(expTags, tags) {
		t.Errorf("expected:\n")
		for _, tag := range expTags {
			t.Errorf("%q\n", tag)
		}
		t.Errorf("actual:\n")
		for _, tag := range tags {
			t.Errorf("%q\n", tag)
		}
	}
}

func TestGetLibVersion(t *testing.T) {
	defaultWd, _ := os.Getwd()
	defer os.Chdir(defaultWd)

	dir, err := createTempRepository(".temp_TestGetTags")
	if err != nil {
		t.Fatalf("cannot create temp repository: %s", err.Error())
	}
	defer os.RemoveAll(dir)

	tomlContent := []byte("[versions]\nlib = \"1.2.3\"")
	os.Mkdir("gradle", os.ModePerm)
	os.WriteFile("gradle/libs.versions.toml", tomlContent, 0644)

	cmd := exec.Command("git", "add", "gradle/libs.versions.toml")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("cannot add to git: %s, %s", err.Error(), string(out))
	}
	cmd = exec.Command("git", "commit", "-m", "\"add toml file\"")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("cannot commit: %s, %s", err.Error(), string(out))
	}

	cmd = exec.Command("git", "tag", "v1.0.0")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("cannot create tag v1.0.0: %s, %s", err.Error(), string(out))
	}

	version, err := tagtolibv.GetLibVersion("v1.0.0", "lib")
	if err != nil {
		t.Fatal(err)
	}

	if version != "1.2.3" {
		t.Errorf("expected 1.2.3, got %s instead.", version)
	}
}

func TestGitCheckout(t *testing.T) {
	defaultWd, _ := os.Getwd()
	defer os.Chdir(defaultWd)

	dir, err := createTempRepository(".temp_TestGitCheckout")
	if err != nil {
		t.Fatalf("cannot create temp repository: %s", err.Error())
	}
	defer os.RemoveAll(dir)

	initialBranch, _ := tagtolibv.GetCurrentBranch()

	os.Create("temp")
	cmd := exec.Command("git", "add", "temp")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("cannot add to git: %s, %s", err.Error(), string(out))
	}
	cmd = exec.Command("git", "commit", "-m", "\"add temp\"")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("cannot commit: %s, %s", err.Error(), string(out))
	}

	cmd = exec.Command("git", "checkout", "-b", "test_branch")
	if err := cmd.Run(); err != nil {
		t.Fatalf("cannot create test_branch: %s", err.Error())
	}
	os.Create("temp2")
	cmd = exec.Command("git", "add", "temp2")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("cannot add to git: %s, %s", err.Error(), string(out))
	}
	cmd = exec.Command("git", "commit", "-m", "\"add temp2\"")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("cannot commit: %s, %s", err.Error(), string(out))
	}
	cmd = exec.Command("git", "checkout", initialBranch)
	if err := cmd.Run(); err != nil {
		t.Fatalf("cannot checkout %s: %s", initialBranch, err.Error())
	}

	tagtolibv.GitCheckout("test_branch")

	actualBranch, _ := tagtolibv.GetCurrentBranch()

	if actualBranch != "test_branch" {
		t.Errorf("expected test_branch, got %s instead.", actualBranch)
	}
}

func createTempRepository(name string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot get work directory")
		return "", err
	}

	err = os.Mkdir(name, os.ModePerm)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot create temp directory")
		return "", err
	}

	repoPath := filepath.Join(wd, name)
	os.Chdir(repoPath)

	cmd := exec.Command("git", "init")
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "cannot create repository")
		return "", err
	}

	fmt.Printf("created repository: %s\n", repoPath)
	return repoPath, nil
}
