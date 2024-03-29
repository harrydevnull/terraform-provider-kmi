package kmi

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

func (client *KMIRestClient) CreateDefinition(collectionName string, definitionName string, definition KMIDefinition) error {
	idenityengineurl := fmt.Sprintf("%s/definition/Col=%s/Def=%s", client.Host, collectionName, definitionName)
	fmt.Println(idenityengineurl)
	out, err := xml.MarshalIndent(definition, " ", "  ")
	if err != nil {
		return err
	}

	resp, err := client.httpclient.Post(idenityengineurl, "application/xml", bytes.NewBuffer(out))

	if err != nil {
		fmt.Printf("error while calling CreateDefinition api posting  %s\n", err.Error())
		return err
	}
	defer resp.Body.Close()
	// go write error handling code for 200
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error while calling CreateDefinition api  %s and payload is %v", resp.Status, resp)
	}

	return nil
}

func (client *KMIRestClient) CreateBlockSecret(collectionName string, definitionName string, opaque BlockSecret) error {
	idenityengineurl := fmt.Sprintf("%s/secret/Col=%s/Def=%s/Idx=AUTOINDEX", client.Host, collectionName, definitionName)
	fmt.Println(idenityengineurl)

	out, err := xml.MarshalIndent(opaque, " ", "  ")
	if err != nil {
		return err
	}

	resp, err := client.httpclient.Post(idenityengineurl, "application/xml", bytes.NewBuffer(out))

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error while calling CreateDefinition CreateOpaqueSecret api  %s and payload is %v", resp.Status, resp)
	}
	return nil
}

func (client *KMIRestClient) DeleteDefinition(collectionName string, definitionName string) error {

	idenityengineurl := fmt.Sprintf("%s/definition/Col=%s/Def=$%s", client.Host, collectionName, definitionName)
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

func (client *KMIRestClient) GetDefinition(collectionName string, definitionName string) (*KMIDefinitionResponse, error) {
	idenityengineurl := fmt.Sprintf("%s/definition/Col=%s/Def=%s", client.Host, collectionName, definitionName)

	response, err := client.httpclient.Get(idenityengineurl)
	if err != nil {
		return nil, err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var responseDetails KMIDefinitionResponse
	err = xml.Unmarshal(responseData, &responseDetails)
	if err != nil {
		return nil, err
	}
	return &responseDetails, nil
}
