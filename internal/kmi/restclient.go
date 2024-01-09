package kmi

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"

	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
)

type KMIRestClient struct {
	Host       string
	ApiKey     string
	ApiCrt     string
	AkamaiCA   string
	httpclient *http.Client
}

func NewKMIRestClient(host string, apiKey string, apiCrt string, akamaiCA string) (*KMIRestClient, error) {

	cert, err := tls.LoadX509KeyPair(apiCrt, apiKey)
	if err != nil {
		return nil, err
	}

	// Load CA cert
	caCert, err := os.ReadFile(akamaiCA)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	return &KMIRestClient{Host: host, ApiKey: apiKey, ApiCrt: apiCrt, AkamaiCA: akamaiCA, httpclient: client}, nil
}

func (client *KMIRestClient) SaveIdentityEngine(account string, engineName string, kmiEngine KMIEngine) error {

	idenityengineurl := fmt.Sprintf("%s:11838/engine/Acct=%s/Eng=%s", client.Host, account, engineName)

	out, err := xml.MarshalIndent(kmiEngine, " ", "  ")
	if err != nil {
		return err
	}

	_, err = client.httpclient.Post(idenityengineurl, "application/xml", bytes.NewBuffer(out))
	if err != nil {
		return err
	}
	return nil
}

func (client *KMIRestClient) GetWorkloadDetails(account string, engineName string, workloadName string) (*Workload, error) {

	idenityengineurl := fmt.Sprintf("%s:11838/workload/Acct=%s/Eng=%s/Proj=%s", client.Host, account, engineName, workloadName)
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

func (client *KMIRestClient) GetIdentityEngine(account string, engineName string) (*IdentityEngine, error) {

	idenityengineurl := fmt.Sprintf("%s:11838/engine/Acct=%s/Eng=%s", client.Host, account, engineName)
	response, err := client.httpclient.Get(idenityengineurl)
	if err != nil {
		return nil, err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(responseData))

	var engine IdentityEngine
	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'engine' which we defined above
	xml.Unmarshal(responseData, &engine)

	// we iterate through every user within o

	return &engine, nil
}

func (client *KMIRestClient) DeleteIdentityEngine(account string, engineName string) error {

	idenityengineurl := fmt.Sprintf("%s:11838/engine/Acct=%s/Eng=%s", client.Host, account, engineName)
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

func (client *KMIRestClient) DeleteWorkload(account string, engineName string, workloadName string) error {

	idenityengineurl := fmt.Sprintf("%s:11838/workload/Acct=%s/Eng=%s/Proj=%s", client.Host, account, engineName, workloadName)
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

func (client *KMIRestClient) GetAccountDetails(account string) (*Account, error) {

	idenityengineurl := fmt.Sprintf("%s:11838/account/Acct=%s/children", client.Host, account)
	response, err := client.httpclient.Get(idenityengineurl)
	if err != nil {
		return nil, err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var accountEngines Account
	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'engine' which we defined above
	xml.Unmarshal(responseData, &accountEngines)

	// we iterate through every user within o

	return &accountEngines, nil
}

func (client *KMIRestClient) CreateCollection(account string, collectionName string, collection CollectionRequest) error {
	idenityengineurl := fmt.Sprintf("%s:11838/collection/Acct=%s/Col=%s", client.Host, account, collectionName)

	out, err := xml.MarshalIndent(collection, " ", "  ")
	if err != nil {
		return err
	}

	_, err = client.httpclient.Post(idenityengineurl, "application/xml", bytes.NewBuffer(out))
	if err != nil {
		return err
	}
	return nil
}

func (client *KMIRestClient) DeleteCollection(collectionName string) error {

	idenityengineurl := fmt.Sprintf("%s:11838/collection/Col=%s", client.Host, collectionName)
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
	idenityengineurl := fmt.Sprintf("%s:11838/collection/Col=%s", client.Host, collectionName)

	response, err := client.httpclient.Get(idenityengineurl)
	if err != nil {
		return nil, err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var collectionDetails Collection
	xml.Unmarshal(responseData, &collectionDetails)
	return &collectionDetails, nil
}
