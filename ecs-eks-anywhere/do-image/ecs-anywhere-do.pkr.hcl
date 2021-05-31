variable "do_token" {
  type = string
  sensitive = true
}

source "digitalocean" "ecs-anywhere" {
  ssh_username = "root"
  api_token = var.do_token
  image = "ubuntu-20-04-x64"
  region = "nyc1"
  size = "s-1vcpu-1gb"
  snapshot_name = "ecs-anywhere"
}

build {
  sources = [
    "source.digitalocean.ecs-anywhere"
  ]

  provisioner "shell" {
    script = "user_data.sh"
    pause_before = "10s"
  }

  provisioner "ansible" {
    playbook_file = "playbook.yml"
    pause_before = "10s"
  }

  provisioner "shell" {
    script = "post_install.sh"
    pause_before = "10s"
  }
}
