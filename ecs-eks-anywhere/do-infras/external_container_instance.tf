resource "digitalocean_droplet" "ecs-1" {
  image = 85934373
  name = "ecs-1"
  region = "nyc1"
  size = "s-1vcpu-1gb"
  private_networking = false
  ssh_keys = [
    digitalocean_ssh_key.default.fingerprint
  ]

  connection {
    host = self.ipv4_address
    user = "root"
    type = "ssh"
    private_key = file(var.pvt_key)
    timeout = "2m"
  }

  provisioner "remote-exec" {
    inline = [
      "mkdir ~/ecs-anywhere",
    ]
  }

  provisioner "file" {
    source      = "ssm-activation.json"
    destination = "~/ecs-anywhere/ssm-activation.json"
  }

  provisioner "file" {
    source      = "user-data.sh"
    destination = "~/ecs-anywhere/user-data.sh"
  }

  provisioner "remote-exec" {
    inline = [
      "cd ~/ecs-anywhere && chmod +x user-data.sh && ./user-data.sh",
    ]
  }
}

resource "digitalocean_ssh_key" "default" {
  name = "default"
  public_key = file("~/.ssh/id_rsa.pub")
}