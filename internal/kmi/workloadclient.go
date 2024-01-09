package kmi

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
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
	xml.Unmarshal(responseData, &workload)

	// we iterate through every user within o

	return &workload, nil
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
