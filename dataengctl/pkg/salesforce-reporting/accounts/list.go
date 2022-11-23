package accounts

import (
	"encoding/json"
	"fmt"

	"github.com/Mirantis/dataeng/dataengctl/pkg/client"
)

type Options struct {
	Output string
	Client client.DataClientInterface
}

func (o *Options) ListAccounts() error {
	q := "SELECT+name+from+Account"
	client, err := o.Client.SalesForceClient()
	if err != nil {
		return err
	}

	result, err := client.Query(q)
	if err != nil {
		return err
	}
	for _, record := range result.Records {
		b, err := json.Marshal(record)
		if err != nil {
			// handle the error
			return err
		}

		fmt.Println(string(b))
	}
	return err
}
