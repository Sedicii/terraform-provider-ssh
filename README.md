# SSH Terraform Provider

This Terraform provides allows to use ssh tunnels. 

### Maintainers

This provider plugin is maintained by [Sedicii](https://sedicii.com/).

### Requirements

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)

### Installation

```bash
curl https://raw.githubusercontent.com/Sedicii/terraform-provider-ssh/master/scripts/install-ssh-tf-pluging.sh | bash
```

### Usage

```
provider "ssh" {
  version = "~> 0.1.0"
}

data "ssh_tunnel" "bastion" {
    bastion_host = "${var.bastion_host}"
    bastion_port = "${var.bastion_port}"
    bastion_user = "${var.bastion_user}"
    bastion_private_key = "${var.bastion_private_key}"
//  bastion_password = "${var.bastion_password}"
    bastion_host_key = "${var.bastion_host_key}"
    remote_host = "${var.remote_host}"
    remote_port = "${var.remote_port}"
}
```

**For a more detailed example look at the example directory !!**

### Building The Provider

Clone repository to: `$GOPATH/src/github.com/sedicii/terraform-provider-ssh`

```sh
$ mkdir -p $GOPATH/src/github.com/sedicii; cd $GOPATH/src/github.com/sedicii
$ git clone git@github.com:sedicii/terraform-provider-ssh
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/sedicii/terraform-provider-ssh
$ make build
```

### Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.8+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-ssh
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```
