terraform {
  required_providers {
    kmi = {
      source = "registry.terraform.io/akamai/kmi"
    }
    # linode = {
    #   source  = "linode/linode"
    #   version = "2.12.0"
    # }

  }
}

# provider "linode" {
#   token = "<>"
# }

# data "linode_lke_cluster" "cluster" {
#   id = 143738
# }




provider "kmi" {
  host      = "https://kdc-api.shared.qa.akamai.com"
  api_key   = "/Users/hachandr/Documents/work/kmi-k8/api_cert.key"
  api_crt   = "/Users/hachandr/Documents/work/kmi-k8/api_cert.crt"
  akamai_ca = "/Users/hachandr/Documents/work/kmi-k8/akamai_ca_list.pem"

}

# data "kmi_account" "example" {
#   account_name = "PIM_TEST"

# }
# output "kmioutput" {
#   value = data.kmi_account.example
# }

# resource "kmi_engine" "identityengine" {

#   engine       = data.linode_lke_cluster.cluster.label
#   account_name = "PIM_TEST"
#   api_endpoint = yamldecode(base64decode(data.linode_lke_cluster.cluster.kubeconfig)).clusters[0].cluster.server
#   cas_base64   = yamldecode(base64decode(data.linode_lke_cluster.cluster.kubeconfig)).clusters[0].cluster.certificate-authority-data
#   workloads = [{
#     name           = "instance_validator"
#     serviceaccount = "SA1"
#     namespace      = "app"
#     region         = data.linode_lke_cluster.cluster.region
#   }]

# }



# output "kmi_engine_create" {
#   value = kmi_engine.identityengine
#   sensitive = true
# }


# resource "kmi_group" "example" {

#   account_name = "PIM_TEST"
#   group_name   = "somevalue"
#   engine       = ""
# }

resource "kmi_collections" "collection" {
  account_name = "PIM_TEST"
  adders       = "PIM_ADMIN"
  modifiers    = "PIM_ADMIN"
  readers      = "PIM_READERS"
  name         = "testCollection"

}
