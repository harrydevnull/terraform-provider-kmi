package kmi

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

func (client *KMIRestClient) SaveIdentityEngine(account string, engineName string, kmiEngine KMIEngine) error {

	idenityengineurl := fmt.Sprintf("%s/engine/Acct=%s/Eng=%s", client.Host, account, engineName)

	out, err := xml.MarshalIndent(kmiEngine, " ", "  ")
	if err != nil {
		return err
	}

	resp, err := client.httpclient.Post(idenityengineurl, "application/xml", bytes.NewBuffer(out))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// go write error handling code for 200
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error while calling SaveIdentityEngine api  %s and payload is %v", resp.Status, resp)
	}
	return nil
}

func SetCloudTypeIdentityEngine(cloud string) string {
	if cloud != "" {
		return cloud
	}
	return "linode"
}

func (client *KMIRestClient) GetIdentityEngine(account string, engineName string) (*IdentityEngine, error) {

	idenityengineurl := fmt.Sprintf("%s/engine/Acct=%s/Eng=%s", client.Host, account, engineName)
	response, err := client.httpclient.Get(idenityengineurl)
	if err != nil {
		return nil, err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var engine IdentityEngine
	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'engine' which we defined above
	err = xml.Unmarshal(responseData, &engine)
	if err != nil {
		return nil, err
	}

	// we iterate through every user within o

	return &engine, nil
}

func (client *KMIRestClient) DeleteIdentityEngine(account string, engineName string) error {

	idenityengineurl := fmt.Sprintf("%s/engine/Acct=%s/Eng=%s", client.Host, account, engineName)
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
