variable "cluster_subnet" {
  type        = string
  default     = "192.168.1.0/24"
}

variable "proxmox_username" {
  type    = string
  default = "root@pam"
}

variable "proxmox_password" {
  type    = string
  default = "Password in secrets.auto.tfvars file."
}

variable "proxmox_endpoint" {
  type    = string
  default = "https://proxmox.example.com"
}

variable "ssh_username" {
  type    = string
  default = "traipoap"
}

variable "timezone" {
  type    = string
  default = "Asia/Bangkok"
}

variable "datastore_id" {
  type    = string
  default = "local"
}

variable "node_name" {
  type    = string
  default = "local"
}

variable "clone_vm_id" {
  type    = number
  default = 2000
}

variable "network_bridge" {
  type    = string
  default = "vmbr0"
}

variable "gateway_ip" {
  type    = string
  default = "192.168.1.1"
}

variable "dns_servers" {
  type    = list(string)
  default = ["8.8.8.8", "8.8.4.4"]
}

variable "ssh_public_key" {
  type    = string
  default = "ssh-ed25519"
}

variable "ansible_inventory_path" {
  type    = string
  default = "../ansible/inventory/hosts"
}

variable "cluster_nodes" {
  type = map(object({
    vm_id        = number
    role         = string
    cpu_cores    = number
    ram_mb       = number
    ip_offset    = number
    hostname     = string
    disks        = list(object({
      datastore_id = string
      interface    = string
      size         = number
    }))
  }))
  default = {
    super-node   = {
      vm_id = 201, role = "super", cpu_cores = 1, ram_mb = 2048, ip_offset = 4, hostname = "super-node",
      disks = [
        { datastore_id = "st500", interface = "scsi0", size = 32 },
        { datastore_id = "data-st1000", interface = "scsi1", size = 100 }
      ]
    },
    k3s-master-1   = {
      vm_id = 202, role = "master", cpu_cores = 1, ram_mb = 4096, ip_offset = 11, hostname = "k3s-master",
      disks = [
        { datastore_id = "st500", interface = "scsi0", size = 32 }
      ]
    },
    k3s-worker-1 = {
      vm_id = 203, role = "worker", cpu_cores = 1, ram_mb = 8192, ip_offset = 21, hostname = "k3s-worker-1",
      disks = [
        { datastore_id = "st500", interface = "scsi0", size = 32 }
      ]
    }
  }
}
