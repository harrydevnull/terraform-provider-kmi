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

func (client *KMIRestClient) GetWorkloadDetails(account string, engineName string, workloadName string) (*Workload, error) {

	idenityengineurl := fmt.Sprintf("%s/workload/Acct=%s/Eng=%s/Proj=%s", client.Host, account, engineName, workloadName)
	response, err := client.httpclient.Get(idenityengineurl)
	if err != nil {
		return nil, err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var workload Workload
	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'engine' which we defined above
	err = xml.Unmarshal(responseData, &workload)
	if err != nil {
		return nil, err
	}

	// we iterate through every user within o

	return &workload, nil
}

func (client *KMIRestClient) CreateWorkloadDetails(account string, engineName string, workload Workload) (*Workload, error) {
	idenityengineurl := fmt.Sprintf("%s/workload/Acct=%s/Eng=%s/Proj=%s", client.Host, account, engineName, workload.Projection)
	tflog.Info(context.Background(), "CreateWorkloadDetails requesturl %v\n"+idenityengineurl)

	out, err := xml.MarshalIndent(workload, " ", "  ")
	if err != nil {
		return nil, err
	}
	tflog.Info(context.Background(), "CreateWorkloadDetails payload %v\n"+string(out))
	resp, err := client.httpclient.Post(idenityengineurl, "application/xml", bytes.NewBuffer(out))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, err
	}
	// go write error handling code for 200
	if resp.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("error while calling CreateWorkloadDetails api  %s and payload is %v", resp.Status, string(b))
	}
	return nil, err

}

func (client *KMIRestClient) DeleteWorkload(account string, engineName string, workloadName string) error {

	idenityengineurl := fmt.Sprintf("%s/workload/Acct=%s/Eng=%s/Proj=%s", client.Host, account, engineName, workloadName)
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
