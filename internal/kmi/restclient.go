package kmi

import (
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

func (client *KMIRestClient) GetAccountDetails(account string) (*Account, error) {

	idenityengineurl := fmt.Sprintf("%s/account/Acct=%s/children", client.Host, account)
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
