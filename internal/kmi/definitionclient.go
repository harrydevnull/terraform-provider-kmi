package kmi

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

func (client *KMIRestClient) CreateDefinition(collectionName string, definitionName string, kmigenerator kmigenerator) error {
	idenityengineurl := fmt.Sprintf("%s/Col=%s/Def=%s", client.Host, collectionName, definitionName)

	out, err := kmigenerator.RequestPayload()
	if err != nil {
		return err
	}

	_, err = client.httpclient.Post(idenityengineurl, "application/xml", bytes.NewBuffer(out))
	if err != nil {
		return err
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
	xml.Unmarshal(responseData, &responseDetails)
	return &responseDetails, nil
}
