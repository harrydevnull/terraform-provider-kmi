resource "kmi_engine" "identityengine" {

  engine       = data.linode_lke_cluster.cluster.label
  account_name = "ACCOUNT_NAME"
  api_endpoint = yamldecode(base64decode(data.linode_lke_cluster.cluster.kubeconfig)).clusters[0].cluster.server
  cas_base64   = yamldecode(base64decode(data.linode_lke_cluster.cluster.kubeconfig)).clusters[0].cluster.certificate-authority-data
  workloads = [{
    name           = "instance_validator"
    serviceaccount = "SA1"
    namespace      = "app"
    region         = data.linode_lke_cluster.cluster.region
  }]

}