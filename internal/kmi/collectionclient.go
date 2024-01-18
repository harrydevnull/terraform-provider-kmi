package kmi

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (client *KMIRestClient) CreateCollection(account string, collectionName string, collection CollectionRequest) error {
	idenityengineurl := fmt.Sprintf("%s/collection/Acct=%s/Col=%s", client.Host, account, collectionName)

	out, err := xml.MarshalIndent(collection, " ", "  ")
	if err != nil {
		return err
	}
	tflog.Info(context.Background(), "CreateCollection payload %v\n"+string(out))
	resp, err := client.httpclient.Post(idenityengineurl, "application/xml", bytes.NewBuffer(out))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return err
	}
	// go write error handling code for 200
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error while calling CreateCollection api  %s and payload is %v", resp.Status, string(b))
	}
	return nil
}

func (client *KMIRestClient) DeleteCollection(collectionName string) error {

	idenityengineurl := fmt.Sprintf("%s/collection/Col=%s", client.Host, collectionName)
	req, err := http.NewRequest("DELETE", idenityengineurl, nil)
	if err != nil {
		return err
	}
	resp, err := client.httpclient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *KMIRestClient) GetCollection(collectionName string) (*Collection, error) {
	idenityengineurl := fmt.Sprintf("%s/collection/Col=%s", client.Host, collectionName)

	response, err := client.httpclient.Get(idenityengineurl)
	if err != nil {
		return nil, err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var collectionDetails Collection
	err = xml.Unmarshal(responseData, &collectionDetails)
	if err != nil {
		return nil, err
	}
	return &collectionDetails, nil
}
