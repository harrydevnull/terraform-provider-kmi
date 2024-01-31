package kmi

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
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
	b, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return err
	}
	// go write error handling code for 200
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error while calling CreateGroup api  %s and payload is %v", resp.Status, string(b))
	}
	return nil
}

func (client *KMIRestClient) CreateGroupMembership(groupName string, child string) error {
	idenityengineurl := fmt.Sprintf("%s/group_membership/Parent=%s/Child=%s", client.Host, groupName, child)

	data := []byte(`<group_membership/>`)
	resp, err := client.httpclient.Post(idenityengineurl, "application/xml", bytes.NewBuffer(data))
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
		return fmt.Errorf("error while calling CreateGroupMembership api  %s and payload is %v", resp.Status, string(b))
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

	var kmiGroup KMIGroup
	err = xml.Unmarshal(responseData, &kmiGroup)
	if err != nil {
		return nil, err
	}

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

func (client *KMIRestClient) DeleteGroupMembership(groupName string, child string) error {
	idenityengineurl := fmt.Sprintf("%s/group_membership/Parent=%s/Child=%s", client.Host, groupName, child)
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
