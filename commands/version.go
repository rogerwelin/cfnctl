package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-version"
)

var latestVersion = "https://api.github.com/repos/rogerwelin/cassowary/releases/latest"

type githubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
}

// OutputVersion queries github releases and will output whenever there are a newer release available
func OutputVersion(ver string, writer io.Writer) error {
	resp, err := http.Get(latestVersion)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	v := &githubRelease{}
	json.NewDecoder(resp.Body).Decode(v)

	v1, err := version.NewVersion(ver)
	if err != nil {
		return err
	}
	v2, err := version.NewVersion(v.TagName)
	if err != nil {
		return err
	}

	if v1.LessThan(v2) {
		fmt.Fprintf(writer, "Cfnctl version v%s\n\nYour version of Cfnctl is out of date. The latest version is v%s\n", v1, v2)
	} else {
		fmt.Fprintf(writer, "Cfnctl version v%s\n", v1)
	}

	return nil
}
