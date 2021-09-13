package commands

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-version"
)

/*
	Twilight of the thundergod
*/

type githubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
}

func OutputVersion(ver string) error {
	resp, err := http.Get("https://api.github.com/repos/rogerwelin/cassowary/releases/latest")
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
		fmt.Printf("Cfnctl version v%s\n\nYour version of Cfnctl is out of date. The latest version is v%s\n", v1, v2)
	} else {
		fmt.Printf("Cfnctl version v%s\n", v1)
	}

	return nil
}
