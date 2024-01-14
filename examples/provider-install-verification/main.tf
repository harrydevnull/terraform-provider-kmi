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
  config_profile = "dev"
}


provider "kmi" {
  host      = "https://kdc-api.shared.qa.akamai.com:11838"
  api_key   = "/Users/hachandr/Documents/work/kmi-k8/api_cert.key"
  api_crt   = "/Users/hachandr/Documents/work/kmi-k8/api_cert.crt"
  akamai_ca = "/Users/hachandr/Documents/work/kmi-k8/akamai_ca_list.pem"

}



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

# data "kmi_account" "example" {
#   account_name = "PIM_TEST"

# }
# output "kmioutput" {
#   value = data.kmi_account.example
# }

# output "kmi_engine_create" {
#   value = kmi_engine.identityengine
#   sensitive = true
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
  kubernetes_service_account   = "kmi-sa"
  kubernetes_namespace         = "app"
  reader_groupname             = "PIM_READERS"
  adder_groupname              = "PIM_TEST_admins"
  modifier_groupname           = "PIM_TEST_admins"
  collection_name              = "testcollection2"
  definition_name              = "testdefinition"
  ssl_cert_definition_name     = "ssl_cert_definition_name"
  azure_sp_definition_name     = "azurespdefinitionname"
  symetric_key_definition_name = "symetrickeydefinitionname"
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
    serviceaccount = local.kubernetes_service_account
    namespace      = local.kubernetes_namespace
    region         = data.linode_lke_cluster.lke_cluster[each.value].region
  }]

}



resource "kmi_group" "group_name" {

  for_each = toset(local.clusters_ids)

  account_name  = local.account_name
  group_name    = local.reader_groupname
  engine        = data.linode_lke_cluster.lke_cluster[each.value].label
  workload_name = local.workload_name
}

output "group_output" {
  value = kmi_group.group_name
}
resource "kmi_collections" "collection" {
  account_name = local.account_name
  adders       = local.adder_groupname
  modifiers    = local.modifier_groupname
  readers      = local.reader_groupname // for the time being entering the adder group name as the reader group name as UserId hasn't been added to that group
  name         = local.collection_name
}

output "collection_output" {
  value = kmi_collections.collection
}

resource "kmi_definitions" "defn" {
  name            = local.definition_name
  collection_name = local.collection_name
  opaque = jsonencode({
    "username" = "bob"
    "password" = "pass123"
  })

  depends_on = [kmi_collections.collection]
}

resource "kmi_definitions" "ssl_defn" {
  name            = local.ssl_cert_definition_name
  collection_name = local.collection_name
  ssl_cert = {
    "auto_generate" = true
  }
  depends_on = [kmi_collections.collection]
}
resource "kmi_definitions" "az_defn" {
  name            = local.azure_sp_definition_name
  collection_name = local.collection_name
  azure_sp = {
    "auto_generate" = true
  }
  depends_on = [kmi_collections.collection]
}

resource "kmi_definitions" "symetric_defn" {
  name            = local.symetric_key_definition_name
  collection_name = local.collection_name
  symmetric_key = {
    "auto_generate"  = true
    "key_size_bytes" = 16
    "expire_period"  = "3 months"
    "refresh_period" = "1 month"

  }
  depends_on = [kmi_collections.collection]
}

output "definitions_output" {
  value = kmi_definitions.symetric_defn
}

# data "linode_lke_cluster" "lke_cluster" {
#   for_each = toset(local.clusters_ids)
#   id       = each.value
# }

# resource "kmi_engine" "identityengine" {
#   for_each = toset(local.clusters_ids)

#   engine       = data.linode_lke_cluster.lke_cluster[each.value].label
#   account_name = local.account_name
#   api_endpoint = yamldecode(base64decode(data.linode_lke_cluster.lke_cluster[each.value].kubeconfig)).clusters[0].cluster.server
#   cas_base64   = yamldecode(base64decode(data.linode_lke_cluster.lke_cluster[each.value].kubeconfig)).clusters[0].cluster.certificate-authority-data
#   workloads = [{
#     name           = local.workload_name
#     serviceaccount = "kmi-sa"
#     namespace      = "app"
#     region         = data.linode_lke_cluster.lke_cluster[each.value].region
#   }]

# }


# output "kmi_engine_create" {
#   value     = kmi_engine.identityengine
#   sensitive = true
# }
