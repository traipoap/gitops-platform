data "local_file" "ssh_public_key" {
  filename = "./id_ed25519.pub"
}

resource "local_file" "ansible_inventory" {
  filename = "../ansible/inventory/hosts"
  content = templatefile("${path.module}/hosts.tpl", {
    super_ip  = local.super-node[0]
    master_ip = local.master_nodes[0]
    worker_ip = local.worker_nodes[0]
  })
}

resource "null_resource" "run_ansible_after_provisioning" {
  depends_on = [
    proxmox_virtual_environment_file.cloud_configs
  ]

  provisioner "local-exec" {
    command = <<-EOT
      echo 'Waiting for VMs to be reachable via SSH...'
      nc -zv "${local.super-node[0]}" 22
      nc -zv "${local.master_nodes[0]}" 22
      nc -zv "${local.worker_nodes[0]}" 22
    EOT
  }
}

# สร้างไฟล์ Cloud-init Snippets แบบ Dynamic ทั้ง Cluster
resource "proxmox_virtual_environment_file" "cloud_configs" {
  for_each = var.cluster_nodes

  content_type = "snippets"
  datastore_id = "local"
  node_name    = "local"

  source_raw {
    data      = <<-EOF
    #cloud-config
    hostname: ${each.value.hostname}
    timezone: Asia/Bangkok
    users:
      - default
      - name: traipoap
        groups: [sudo]
        shell: /bin/bash
        ssh_authorized_keys:
          - ${trimspace(data.local_file.ssh_public_key.content)}
        sudo: ALL=(ALL) N_PASSWD:ALL
    package_update: true
    packages:
      - qemu-guest-agent
    runcmd:
      - systemctl enable qemu-guest-agent
      - systemctl start qemu-guest-agent
      - echo "done" > /tmp/cloud-config.done
    EOF
    file_name = "${each.key}.yaml"
  }
}

resource "proxmox_virtual_environment_vm" "nodes" {
  for_each = var.cluster_nodes

  name      = each.key
  node_name = "local"
  vm_id     = each.value.vm_id

  clone {
    vm_id = 2000
    full  = false
  }

  agent { enabled = true }

  cpu {
    cores = each.value.cpu_cores
    type  = "host"
  }

  memory {
    dedicated = each.value.ram_mb
  }

  # --- THE DYNAMIC DISK LOOP ---
  dynamic "disk" {
    for_each = each.value.disks
    content {
      datastore_id = disk.value.datastore_id
      interface    = disk.value.interface
      size         = disk.value.size
    }
  }

  network_device {
    bridge = "vmbr0"
  }

  serial_device {}

  # --- Cloud-config for K3s Worker (agent only, joins the master) ---
  initialization {
    user_data_file_id = proxmox_virtual_environment_file.cloud_configs[each.key].id
    ip_config {
      ipv4 {
        address = "${cidrhost(var.cluster_subnet, each.value.ip_offset)}/24"
        gateway = "192.168.1.1"
      }
    }
    dns {
      servers = ["8.8.8.8", "8.8.4.4"]
      domain  = "."
    }
  }
  depends_on = [
    proxmox_virtual_environment_file.cloud_configs
  ]
}
