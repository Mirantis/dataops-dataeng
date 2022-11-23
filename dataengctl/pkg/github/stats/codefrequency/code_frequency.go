package code_frequency

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Mirantis/dataeng/dataengctl/pkg/client"
)

// Options Struct
type Options struct {
	Output string
	Client client.DataClientInterface
	Org	string
	Repo string
}

// curl -u $username:$token https://api.github.com/repos/$ORG/$REPO/stats/code_frequency
// github get stats --type codefrequency --org $ORG --repo $REPO
func GetRepoStatsCodeFrequencu() error {
	req, err := http.NewRequest("GET", os.ExpandEnv("https://api.github.com/repos/$ORG/$REPO/stats/code_frequency"), nil)
	if err != nil {
		// handle err
		fmt.Printf("ERROR: %s", err)
		return err
	}
	req.SetBasicAuth(os.ExpandEnv("$username"), os.ExpandEnv("$token"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
		fmt.Printf("ERROR: %s", err)
		return err
	}
	defer resp.Body.Close()
	fmt.Printf("ERROR: %s", err)
	return err
}
