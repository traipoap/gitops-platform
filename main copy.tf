data "local_file" "ssh_public_key" {
  filename = "./id_ed25519.pub"
}

locals {
  cluster_subnet = "192.168.1.0/24"
  super_node_ip = cidrhost(local.cluster_subnet, 4)   # 192.168.1.4
  master_node_ip = cidrhost(local.cluster_subnet, 11) # 192.168.1.11
  worker_node_ip = cidrhost(local.cluster_subnet, 21) # 192.168.1.21
}

resource "local_file" "ansible_inventory" {
  filename = "../ansible/inventory/hosts"
  content = templatefile("${path.module}/hosts.tpl", {
    super_ip = local.super_node_ip
    master_ip = local.master_node_ip
    worker_ip = local.worker_node_ip
  })
}

# -2. Orchestrator: Trigger Ansible Playbook after all VMs are provisioned
resource "null_resource" "run_ansible_after_provisioning" {
  depends_on = [
    proxmox_virtual_environment_vm.super_node,
    proxmox_virtual_environment_vm.k3s_master,
    proxmox_virtual_environment_vm.k3s_worker
  ]

  provisioner "local-exec" {
    command = <<-EOT
      echo 'Waiting for VMs to be reachable via SSH...'
      nc -zv ${local.super_node_ip} 22
      nc -zv ${local.master_node_ip} 22
      nc -zv ${local.worker_node_ip} 22
    EOT
  }
}

# -1. Create file Cloud-config keep in to Proxmox Node (Storage: local)
resource "proxmox_virtual_environment_file" "user_data_cloud_config_for_super_node" {
  content_type = "snippets"
  datastore_id = "local"
  node_name    = "local"

  source_raw {
    data      = <<-EOF
    #cloud-config
    hostname: super-node
    timezone: Asia/Bangkok
    users:
      - default
      - name: traipoap
        groups:
          - sudo
        shell: /bin/bash
        ssh_authorized_keys:
          - ${trimspace(data.local_file.ssh_public_key.content)}
        sudo: ALL=(ALL) NOPASSWD:ALL
    package_update: true
    packages:
      - qemu-guest-agent
    runcmd:
      - systemctl enable qemu-guest-agent
      - systemctl start qemu-guest-agent
      - echo "done" > /tmp/cloud-config.done
    EOF
    file_name = "super-node.yaml"
  }
}

# 0. Create file Cloud-config keep in to Proxmox Node (Storage: local)
resource "proxmox_virtual_environment_file" "user_data_cloud_config_for_master" {
  content_type = "snippets"
  datastore_id = "local"
  node_name    = "local"

  source_raw {
    data      = <<-EOF
    #cloud-config
    hostname: k3s-master
    timezone: Asia/Bangkok
    users:
      - default
      - name: traipoap
        groups:
          - sudo
        shell: /bin/bash
        ssh_authorized_keys:
          - ${trimspace(data.local_file.ssh_public_key.content)}
        sudo: ALL=(ALL) NOPASSWD:ALL
    package_update: true
    packages:
      - qemu-guest-agent
    runcmd:
      - systemctl enable qemu-guest-agent
      - systemctl start qemu-guest-agent
      - echo "done" > /tmp/cloud-config.done
    EOF
    file_name = "k3s-master.yaml"
  }
}

# 1. Create file Cloud-config keep in to Proxmox Node (Storage: local)
resource "proxmox_virtual_environment_file" "user_data_cloud_config_for_worker" {
  content_type = "snippets"
  datastore_id = "local"
  node_name    = "local"

  source_raw {
    data      = <<-EOF
    #cloud-config
    hostname: k3s-worker
    timezone: Asia/Bangkok
    users:
      - default
      - name: traipoap
        groups:
          - sudo
        shell: /bin/bash
        ssh_authorized_keys:
          - ${trimspace(data.local_file.ssh_public_key.content)}
        sudo: ALL=(ALL) NOPASSWD:ALL
    package_update: true
    packages:
      - qemu-guest-agent
    runcmd:
      - systemctl enable qemu-guest-agent
      - systemctl start qemu-guest-agent
      - echo "done" > /tmp/cloud-config.done
    EOF
    file_name = "k3s-worker.yaml"
  }
}

# 2. Create VM from template
resource "proxmox_virtual_environment_vm" "super_node" {
  name      = "super-node"
  node_name = "local"
  vm_id     = 201

  clone {
    vm_id = 2000
    full  = false
  }

  agent {
    enabled = true
  }

  cpu {
    cores = 1
    type = "host"
  }

  memory {
    dedicated = 2048
  }

  # For root filesystem
  disk {
    datastore_id = "st500"
    interface    = "scsi0"
    size         = 32

  }

  # For s3 and nfs storage
  disk {
    datastore_id = "data-st1000"
    interface    = "scsi1"
    size         = 100
  }

  network_device {
    bridge = "vmbr0"
  }

  serial_device {}

  initialization {
    user_data_file_id = proxmox_virtual_environment_file.user_data_cloud_config_for_super_node.id

    ip_config {
      ipv4 {
        address = "${local.super_node_ip}/24"
        gateway = "192.168.1.1"
      }
    }
    dns {
      servers = ["8.8.8.8", "8.8.4.4"]
      domain  = "."
    }
  }
}

# 3. Create VM from template
resource "proxmox_virtual_environment_vm" "k3s_master" {
  name      = "k3s-master"
  node_name = "local"
  vm_id     = 202

  clone {
    vm_id = 2000
    full  = false
  }

  agent {
    enabled = true
  }

  cpu {
    cores = 1
    type = "host"
  }

  memory {
    dedicated = 4096
  }

  disk {
    datastore_id = "st500"
    interface    = "scsi0"
    size         = 32

  }

  network_device {
    bridge = "vmbr0"
  }

  serial_device {}

  initialization {
    user_data_file_id = proxmox_virtual_environment_file.user_data_cloud_config_for_master.id

    ip_config {
      ipv4 {
        address = "${local.master_node_ip}/24"
        gateway = "192.168.1.1"
      }
    }
    dns {
      servers = ["8.8.8.8", "8.8.4.4"]
      domain  = "."
    }
  }
}

# 4. Create VM from template
resource "proxmox_virtual_environment_vm" "k3s_worker" {
  name      = "k3s-worker"
  node_name = "local"
  vm_id     = 203

  clone {
    vm_id = 2000
    full  = false
  }

  agent {
    enabled = true
  }

  cpu {
    cores = 1
    type = "host"
  }

  memory {
    dedicated = 8192
  }

  disk {
    datastore_id = "st500"
    interface    = "scsi0"
    size         = 32

  }

  network_device {
    bridge = "vmbr0"
  }

  serial_device {}

  # --- Cloud-config for K3s Worker (agent only, joins the master) ---
  initialization {
    user_data_file_id = proxmox_virtual_environment_file.user_data_cloud_config_for_worker.id

    ip_config {
      ipv4 {
        address = "${local.worker_node_ip}/24"
        gateway = "192.168.1.1"
      }
    }
    dns {
      servers = ["8.8.8.8", "8.8.4.4"]
      domain  = "."
    }
  }
}
