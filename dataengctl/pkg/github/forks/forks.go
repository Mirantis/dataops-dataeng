package forks

import (
	"net/http"
	"os"

	"github.com/Mirantis/dataeng/dataengctl/pkg/client"
)

// Options Struct
type Options struct {
	Output string
	Client client.DataClientInterface
}

// curl -u $username:$token https://api.github.com/repos/$ORG/$REPO/forks
// github get forks --org $ORG --repo $REPO
func ListRepoForks() error {
	req, err := http.NewRequest("GET", os.ExpandEnv("https://api.github.com/repos/$ORG/$REPO/forks"), nil)
	if err != nil {
		// handle err
		return err
	}
	req.SetBasicAuth(os.ExpandEnv("$username"), os.ExpandEnv("$token"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
		return err
	}
	defer resp.Body.Close()
	return err
}