terraform {
  required_providers {
    kmi = {
      source = "registry.terraform.io/akamai/kmi"
    }
    linode = {
      source  = "linode/linode"
      version = "2.12.0"
    }

  }
}

provider "linode" {
  token = "84395df1bee3c8e6ddf78033908e3e48c332169967e0b265248f7b79877dcb2d"
}

# data "linode_lke_cluster" "cluster" {
#   id = 143738
# }




provider "kmi" {
  host      = "https://kdc-api.shared.qa.akamai.com:11838"
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
#     serviceaccount = "kmi-sa"
#     namespace      = "app"
#     region         = data.linode_lke_cluster.cluster.region
#   }]

# }



# output "kmi_engine_create" {
#   value = kmi_engine.identityengine
#   sensitive = true
# }


# # resource "kmi_group" "group_name" {

# #   account_name  = "PIM_TEST"
# #   group_name    = "kmiGrouptest"
# #   engine        = "pi-dev-usiad-l1-2023-1127-121839"
# #   workload_name = "instance_validator"
# # }

# resource "kmi_collections" "collection" {
#   account_name = "PIM_TEST"
#   adders       = "PIM_ADMIN"
#   modifiers    = "PIM_ADMIN"
#   readers      = "PIM_READERS"
#   name         = "testCollection"

# }


# resource "kmi_definitions" "collectiondefinition" {
#   name            = "testdefinition"
#   collection_name = "testCollection"
#   opaque          = ""

# }


locals {
  account_name  = "PIM_TEST"
  workload_name = "instance_validator"
  clusters_ids = [
    "103067",
    "140548",
    "143042",
    "143738",
    "143752",
    "143888",
    "143892",
    "143893",
    "144139",
    "144142",
    "144149",
    "145885",
    "147806",
    "149191",
    "149192"
  ]
}


data "linode_lke_cluster" "lke_cluster" {
  for_each = toset(local.clusters_ids)
  id       = each.value
}

resource "kmi_engine" "identityengine" {
  for_each = toset(local.clusters_ids)

  engine       = data.linode_lke_cluster.lke_cluster[each.value].label
  account_name = local.account_name
  api_endpoint = yamldecode(base64decode(data.linode_lke_cluster.lke_cluster[each.value].kubeconfig)).clusters[0].cluster.server
  cas_base64   = yamldecode(base64decode(data.linode_lke_cluster.lke_cluster[each.value].kubeconfig)).clusters[0].cluster.certificate-authority-data
  workloads = [{
    name           = local.workload_name
    serviceaccount = "kmi-sa"
    namespace      = "app"
    region         = data.linode_lke_cluster.lke_cluster[each.value].region
  }]

}


output "kmi_engine_create" {
  value     = kmi_engine.identityengine
  sensitive = true
}