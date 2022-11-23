package soql

import (
	"encoding/json"
	"fmt"

	"github.com/Mirantis/dataeng/dataengctl/pkg/client"
)

type Options struct {
	Output string
	Client client.DataClientInterface
}

func (o *Options) SOQLQuery(query string) error {
	client, err := o.Client.SalesForceClient()
	if err != nil {
		return err
	}

	result, err := client.Query(query)
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
