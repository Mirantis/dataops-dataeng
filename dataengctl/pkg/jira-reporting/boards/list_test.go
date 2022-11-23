package boards

import (
	//"bytes"
	"bytes"
	"testing"

	"github.com/Mirantis/dataeng/dataengctl/pkg/client"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	testCases := []struct {
		client         *client.DataClient
		writer         *bytes.Buffer
		errorString    string
		expectedOutput string
	}{
		{
			client: &client.DataClient{
				ConfigFilePath: "doesn't exit",
			},
			writer:      nil,
			errorString: "no such file or directory",
		},
		{
			client: &client.DataClient{
				ConfigFilePath: "doesn't exit",
			},
			writer:      bytes.NewBuffer([]byte{}),
			errorString: "",
		},
	}

	for _, tt := range testCases {
		options := &Options{
			Config: tt.client,
		}
		err := options.List(tt.writer)

		if tt.errorString != "" {
			assert.Contains(t, err.Error(), tt.errorString)
		} else {
			assert.NoError(t, err)
			assert.Contains(t, tt.writer.String(), tt.expectedOutput)
		}
	}

}
