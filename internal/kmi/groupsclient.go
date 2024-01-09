package kmi

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

func (client *KMIRestClient) CreateGroup(account string, groupName string) error {
	idenityengineurl := fmt.Sprintf("%s/group/Acct=%s/Name=%s", client.Host, account, groupName)

	group := GroupRequest{
		Account: account,
		Type:    "union",
	}
	out, err := xml.MarshalIndent(group, " ", "  ")
	if err != nil {
		return err
	}

	resp, err := client.httpclient.Post(idenityengineurl, "application/xml", bytes.NewBuffer(out))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *KMIRestClient) CreateGroupMembership(account string, groupName string, engineName string, workloadName string) error {
	idenityengineurl := fmt.Sprintf("%s/group_membership/Parent=%s/Child=workload:%s:%s:%s", client.Host, groupName, account, engineName, workloadName)

	data := []byte(`<group_membership/>`)
	_, err := client.httpclient.Post(idenityengineurl, "application/xml", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	return nil
}

func (client *KMIRestClient) GetGroup(groupName string) (*KMIGroup, error) {
	idenityengineurl := fmt.Sprintf("%s/group/Name=%s", client.Host, groupName)
	response, err := client.httpclient.Get(idenityengineurl)
	if err != nil {
		return nil, err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(responseData))

	var kmiGroup KMIGroup
	xml.Unmarshal(responseData, &kmiGroup)

	return &kmiGroup, nil
}

func (client *KMIRestClient) DeleteGroup(groupName string) error {

	idenityengineurl := fmt.Sprintf("%s/group/Name=%s", client.Host, groupName)
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

type GroupRequest struct {
	XMLName xml.Name `xml:"group"`
	Text    string   `xml:",chardata"`
	Type    string   `xml:"type,attr"`
	Account string   `xml:"account,attr"`
}
