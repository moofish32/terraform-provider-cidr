Terraform Provider
==================
[![Build Status](https://travis-ci.com/moofish32/terraform-provider-cidr.svg?branch=master)](https://travis-ci.com/moofish32/terraform-provider-cidr)
[![Test Coverage](https://api.codeclimate.com/v1/badges/78c89d99d59df51cc9b1/test_coverage)](https://codeclimate.com/github/moofish32/terraform-provider-cidr/test_coverage)
[![Maintainability](https://api.codeclimate.com/v1/badges/78c89d99d59df51cc9b1/maintainability)](https://codeclimate.com/github/moofish32/terraform-provider-cidr/maintainability)
- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.11 or higher (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/hashicorp/terraform-provider-$PROVIDER_NAME`

```sh
$ mkdir -p $GOPATH/src/github.com/hashicorp; cd $GOPATH/src/github.com/hashicorp
$ git clone git@github.com:hashicorp/terraform-provider-$PROVIDER_NAME
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/hashicorp/terraform-provider-$PROVIDER_NAME
$ make build
```

Using the provider
----------------------
```hcl
data "cidr_network" "order" {
	cidr_block = "10.0.0.0/21"
	subnet {
		mask = 28
		name = "private_az1" 
	}
	subnet {
		mask = 24
		name = "private_az2" 
	}
	subnet {
		mask = 28
		name = "elb_az1" 
	}
	subnet {
		mask = 27
		name = "elb_az2" 
	}
}

// then later use the outputs 
```

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.8+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-$PROVIDER_NAME
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
