provider "kmi" {
  host      = "https://kdc-api.shared.qa.akamai.com"
  api_key   = "<path to api_cert.key>"
  api_crt   = "<path to api_cert.crt>"
  akamai_ca = "<path to akamai_ca_list.pem>"

}