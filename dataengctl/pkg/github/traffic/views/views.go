package views

import (
	"fmt"
	"net/http"
	"os"
)

// curl -u $username:$token https://api.github.com/repos/$ORG/$REPO/traffic/views
// github get traffic --type views --org $ORG --repo $REPO
func GetViews() error {
	req, err := http.NewRequest("GET", os.ExpandEnv("https://api.github.com/repos/$ORG/$REPO/traffic/views"), nil)
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
	return err
}