
provider "ssh" {}

data "ssh_tunnel" "bastion" {
  bastion_host = "${var.bastion_host}"
  bastion_port = "${var.bastion_port}"
  bastion_user = "${var.bastion_user}"
  bastion_private_key = "${var.bastion_private_key}"
  remote_host = "${var.db_host}"
  remote_port = "${var.db_port}"
}

provider "postgresql" {
  host = "${data.ssh_tunnel.bastion.local_host}"
  port = "${data.ssh_tunnel.bastion.local_port}"
  username = "${var.db_username}"
  password = "${var.db_password}"
  sslmode = "${var.db_sslmode}"
}

resource "postgresql_database" "database" {
  name = "${var.database_name}"
}
