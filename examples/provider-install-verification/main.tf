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
  host      = "Hostname of the KMI instance"
  api_key   = "path to api key"
  api_crt   = "path to api crt"
  akamai_ca = "Path to CA"

}



locals {
  account_name  = "PIM_TEST"
  workload_name = "instance_validator"
  clusters_ids = [
    "<CLUSTER_ID>",
    

  ]
  kubernetes_service_account   = "kmi-sa"
  kubernetes_namespace         = "app"
  reader_groupname             = "PIM_READERS-SHARED"
  adder_groupname              = "PIM_TEST_admins"
  modifier_groupname           = "PIM_TEST_admins"
  collection_name              = "testcollection-HAREE"
  definition_name              = "testdefinition-HAREE"
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


output "output_engine" {
  value     = resource.kmi_engine.identityengine
  sensitive = true
}


resource "kmi_group" "sec_group" {
  account_name = local.account_name
  group_name   = local.reader_groupname
}

output "group_output" {
  value = kmi_group.sec_group
}


locals {
  members = flatten([for cluster in local.clusters_ids : {
    name = format("workload:%s:%s:%s", local.account_name, data.linode_lke_cluster.lke_cluster[cluster].label, local.workload_name)

  }])

}

resource "kmi_group_membership" "group_membership" {

  group_name = local.reader_groupname
  members    = local.members
}



resource "kmi_collections" "collection" {
  depends_on   = [kmi_group.sec_group]
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
  depends_on      = [kmi_collections.collection]
  name            = local.definition_name
  collection_name = local.collection_name
  opaque = jsonencode({
    "username" = "bob"
    "password" = "pass123"
  })


}

resource "kmi_definitions" "ssl_defn" {
  name            = local.ssl_cert_definition_name
  collection_name = local.collection_name
  ssl_cert = {
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



# resource "kmi_definitions" "az_defn" {
#   name            = local.azure_sp_definition_name
#   collection_name = local.collection_name
#   azure_sp = {
#     "auto_generate" = true
#   }
#   depends_on = [kmi_collections.collection]
# }
