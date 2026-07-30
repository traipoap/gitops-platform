variable "cluster_subnet" {
  type        = string
  default     = "192.168.1.0/24"
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
    k3s-master   = {
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
