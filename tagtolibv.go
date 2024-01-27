package tagtolibv

import (
	"os/exec"
	"strings"

	"github.com/BurntSushi/toml"
)

func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	out, err := cmd.CombinedOutput()

	if err != nil {
		return "", err
	}

	return strings.Trim(string(out), "\n"), nil
}

func GetTags(lim int) ([]string, error) {
	cmd := exec.Command("git", "tag", "-l", "v*", "--sort=-creatordate")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	tags := strings.Split(string(out), "\n")
	if len(tags) > 0 {
		tags = tags[:len(tags)-1]
		if lim < len(tags) {
			tags = tags[:lim]
		}
	}

	return tags, nil
}

type tomlConfig struct {
	Versions map[string]string
}

func GetLibVersion(tag string, lib string) (string, error) {
	GitCheckout(tag)

	tomlFile := "gradle/libs.versions.toml"
	var config tomlConfig
	_, err := toml.DecodeFile(tomlFile, &config)
	if err != nil {
		return "", err
	}

	version := config.Versions[lib]

	return version, nil
}

func GitCheckout(target string) error {
	cmd := exec.Command("git", "checkout", target)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
