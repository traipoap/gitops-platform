provider "proxmox" {
  endpoint = "https://proxmox.example.com"
  username = "root@pam"
  password = "REDACTED_PROXMOX_PASSWORD"
  insecure = true

  ssh {
    agent    = true
    username = "root"
    password = "REDACTED_PROXMOX_PASSWORD"
  }
}
