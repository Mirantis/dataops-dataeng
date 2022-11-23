package contributors

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

// curl -su $username:$token https://api.github.com/repos/$ORG/$REPO/traffic/clones
// github get contributors --org $ORG --repo $REPO
func GetContributors() error {
	req, err := http.NewRequest("GET", os.ExpandEnv("https://api.github.com/repos/$ORG/$REPO/contributors"), nil)
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