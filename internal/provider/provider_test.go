// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

const (
// providerConfig is a shared configuration to combine with the actual
// test configuration so the KMI client is properly configured.
// It is also possible to use the KMI_ environment variables instead,
// such as updating the Makefile and running the testing through that tool.
// 	providerConfig = `
// 		provider "kmi" {
// 			host      = "http://localhost:19090"
// 			api_key   = "api_cert.key"
// 			api_crt   = "api_cert.crt"
// 			akamai_ca = "akamai_ca_list.pem"

//			}
//	  `
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
// var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
// 	"kmi": providerserver.NewProtocol6WithError(New("test")()),
// }

// func testAccPreCheck(t *testing.T) {
// 	// You can add code here to run prior to any test case execution, for example assertions
// 	// about the appropriate environment variables being set are common to see in a pre-check
// 	// function.
// }
