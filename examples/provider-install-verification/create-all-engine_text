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
  host      = "https://kdc-api.shared.qa.akamai.com"
  api_key   = "/Users/igada/.config/kmi-auth/api_cert.key"
  api_crt   = "/Users/igada/.config/kmi-auth/api_cert.crt"
  akamai_ca = "/Users/igada/.config/kmi-auth/akamai_ca_list.pem"

}

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