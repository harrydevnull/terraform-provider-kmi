locals {
  account_name  = "PIM_TEST"
  workload_name = "instance_validator"
  clusters_ids = [
    "CLUSTER_ID",
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


resource "kmi_engine" "identityengine" {

  engine       = data.linode_lke_cluster.cluster.label
  account_name = local.account_name
  api_endpoint = yamldecode(base64decode(data.linode_lke_cluster.cluster.kubeconfig)).clusters[0].cluster.server
  cas_base64   = yamldecode(base64decode(data.linode_lke_cluster.cluster.kubeconfig)).clusters[0].cluster.certificate-authority-data
  workloads = [{
    name           = local.workload_name
    serviceaccount = local.kubernetes_service_account
    namespace      = local.kubernetes_namespace
    region         = data.linode_lke_cluster.cluster.region
  }]

}



resource "kmi_group" "group_name" {
  account_name = local.account_name
  group_name   = local.reader_groupname
}

output "group_output" {
  value = kmi_group.group_name
}
resource "kmi_collections" "collection" {
  account_name = local.account_name
  adders       = local.adder_groupname
  modifiers    = local.modifier_groupname
  readers      = local.reader_groupname
  name         = local.collection_name
}

output "collection_output" {
  value = kmi_collections.collection
}

resource "kmi_definitions" "defn" {
  readers         = local.reader_groupname
  adders          = local.adder_groupname
  modifiers       = local.modifier_groupname
  name            = local.definition_name
  collection_name = local.collection_name
  opaque = jsonencode({
    "username" = "bob"
    "password" = "pass123"
  })

  depends_on = [kmi_collections.collection]
}

resource "kmi_definitions" "ssl_defn" {
  adders          = local.adder_groupname
  name            = local.ssl_cert_definition_name
  collection_name = local.collection_name
  ssl_cert = {
    "auto_generate" = true
  }
  depends_on = [kmi_collections.collection]
}
resource "kmi_definitions" "az_defn" {
  modifiers       = local.modifier_groupname
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
