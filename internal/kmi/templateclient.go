package kmi

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

func (client *KMIRestClient) CreateTemplateOrSign(cacollectionName string, cadefinitionName string, templateName string, options Template) error {
	idenityengineurl := fmt.Sprintf("%s/template/Col=%s/Def=%s/Tmpl=%s", client.Host, cacollectionName, cadefinitionName, templateName)
	fmt.Println(idenityengineurl)
	out, err := xml.MarshalIndent(options, "", "")
	fmt.Printf("CreateTemplate payload %v\n", string(out))
	if err != nil {
		return err
	}

	resp, err := client.httpclient.Post(idenityengineurl, "application/xml", bytes.NewBuffer(out))

	if err != nil {
		fmt.Printf("error while calling CreateTemplate api posting  %s\n", err.Error())
		return err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// go write error handling code for 200
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error while calling CreateTemplate api  %s and payload is %v", resp.Status, string(b))
	}

	return nil
}

func (client *KMIRestClient) GetTemplate(cacollectionName string, cadefinitionName string, templateName string) (*Template, error) {
	idenityengineurl := fmt.Sprintf("%s/template/Col=%s/Def=%s/Tmpl=%s", client.Host, cacollectionName, cadefinitionName, templateName)

	resp, err := client.httpclient.Get(idenityengineurl)

	if err != nil {
		fmt.Printf("error while calling GetTemplate api posting  %s\n", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	// go write error handling code for 200
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error while calling GetTemplate api  %s and payload is %v", resp.Status, resp)
	}

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responseDetails Template
	err = xml.Unmarshal(responseData, &responseDetails)
	if err != nil {
		return nil, err
	}
	return &responseDetails, nil

}

func (client *KMIRestClient) DeleteTemplate(cacollectionName string, cadefinitionName string, templateName string) error {

	idenityengineurl := fmt.Sprintf("%s/template/Col=%s/Def=%s/Tmpl=%s", client.Host, cacollectionName, cadefinitionName, templateName)

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
