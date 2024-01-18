package kmi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCollection(t *testing.T) {

	expected := "`<collection name=\"testcollection1\" source=\"restserv:user:hachandr_kmi_cert\" readers=\"PIM_READERS\" adders=\"PIM_TEST_admins\" modifiers=\"PIM_TEST_admins\" modified=\"355396416\" distributed=\"355397350\" distributed_date=\"2024-01-14 04:19:42\" account=\"PIM_TEST\"></collection>`"
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, expected)
	}))
	defer svr.Close()
	client := &KMIRestClient{
		Host:       svr.URL,
		httpclient: svr.Client(),
	}

	// Define the collection name for testing
	collectionName := "test-collection"

	// Call the GetCollection function
	collection, err := client.GetCollection(collectionName)

	// Check if there was an error
	assert.NoError(t, err, "Expected no error")

	// Check if the collection is not nil
	assert.NotNil(t, collection, "Expected collection to not be nil")

	// Add additional assertions to validate the collection object
	// For example, you can check if the collection name matches the expected value

	// assert.Equal(t, collectionName, collection.Name, "Expected collection name to match")

	// Add more assertions as needed
}

func TestGetGroup(t *testing.T) {

	expected := "`<group name=\"PIM_ADMIN\" type=\"union\" source=\"reqserv:user:hachandr_kmi_cert\" account=\"PIM_TEST\"><adders>superusers</adders><modifiers>superusers</modifiers></group>`"
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, expected)
	}))
	defer svr.Close()
	client := &KMIRestClient{
		Host:       svr.URL,
		httpclient: svr.Client(),
	}

	// Define the collection name for testing
	group_name := "PIM_ADMIN"

	group, err := client.GetGroup(group_name)

	// Check if there was an error
	assert.NoError(t, err, "Expected no error")

	// Check if the collection is not nil
	assert.NotNil(t, group, "Expected Group to not be nil")

	// Add additional assertions to validate the collection object
	// For example, you can check if the collection name matches the expected value

	// assert.Equal(t, collectionName, collection.Name, "Expected collection name to match")

	// Add more assertions as needed
}
