locals {
  super-node = [for k, v in var.cluster_nodes : cidrhost(var.cluster_subnet, v.ip_offset) if v.role == "super"]
  master_nodes = [for k, v in var.cluster_nodes : cidrhost(var.cluster_subnet, v.ip_offset) if v.role == "master"]
  worker_nodes = [for k, v in var.cluster_nodes : cidrhost(var.cluster_subnet, v.ip_offset) if v.role == "worker"]
}
