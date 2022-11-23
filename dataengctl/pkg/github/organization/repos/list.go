package repos

import (
	"fmt"
	"github.com/Mirantis/dataeng/dataengctl/pkg/client"
	"io"
	"net/http"
)

// List Repos Defined on Organization

// curl -s \
// -H "Accept: application/vnd.github.v3+json" \
// https://api.github.com/orgs/Mirantis/repos

// github list repos --organization {ORG_NAME}
// github list repos --organization Mirantis

// Options Struct
type Options struct {
	Output string
	Client client.DataClientInterface
	Org	string
	Repo string
}

// github get repos --org $ORG --repo $REPO
func (lo *Options) ListRepos(w io.Writer) error {
	req, err := http.NewRequest("GET", "https://api.github.com/orgs/Mirantis/repos", nil)
	if err != nil {
		// Error Handle
		fmt.Printf("ERROR: %s", err)
		return err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// Error Handle
		fmt.Printf("ERROR: %s", err)
		return err
	}
	defer resp.Body.Close()
	return nil
}
