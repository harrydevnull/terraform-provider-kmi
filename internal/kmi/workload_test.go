package kmi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateWorkload(t *testing.T) {

	// expected := `<workload projection="test"><region>test_region</region><linode_label>*</linode_label></workload>`
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprint(w, expected)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer svr.Close()
	client := &KMIRestClient{
		Host:       svr.URL,
		httpclient: svr.Client(),
	}

	// Define the collection name for testing

	_, err := client.CreateWorkloadDetails("test", "test", Workload{
		Projection: "test",
		Region: struct {
			Text   string "xml:\",chardata\""
			Source string "xml:\"source,attr,omitempty\""
		}{
			Text: "test_region",
		},
		LinodeLabel: &LinodeLabel{
			Text: "*",
		},
	})

	// Check if there was an error
	assert.NoError(t, err, "Expected no error")

}
