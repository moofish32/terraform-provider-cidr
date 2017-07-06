---
layout: "cidr"
page_title: "Template: cidr_network"
sidebar_current: "docs-cidr-datasource-network"
description: |-
  Creates CIDR values which describe a network based on subnets
---

# cidr_network

Describes a network with subnets 

## Example Usage

```hcl
variable	"cidr_block" { default = "192.168.0.0/21" } 

data "aws_availability_zones" "available" {}

data "cidr_network" "cake" {
	cidr_block = "${var.cidr_block}"
	subnet {
		mask = 24
		name = "private_az1" 
	}
	subnet {
		mask = 24
		name = "private_az2" 
	}
	subnet {
		mask = 24
		name = "private_az3" 
	}
	subnet {
		mask = 25
		name = "public_az1"
	} 
	subnet {
		mask = 25
		name = "public_az2"
	} 
	subnet {
		mask = 25
		name = "public_az3"
	} 
	subnet {
		mask = 28
		name = "elb_az1" 
	}
	subnet {
		mask = 28
		name = "elb_az2" 
	}
	subnet {
		mask = 28
		name = "elb_az3" 
	}
}

resource "aws_vpc" "main" {
	cidr_block           = "${var.cidr_block}"
	enable_dns_hostnames = true
	enable_dns_support   = true
}
resource "aws_subnet" "public_subnets" {
	vpc_id            = "${aws_vpc.main.id}"
	count              = "${length(data.aws_availability_zones.available.names)}"
	cidr_block        = "${data.cidr_network.cake.subnet_cidrs["public_az${count.index + 1}"]}"
	availability_zone = "${data.aws_availability_zones.available.names[count.index]}"
}
resource "aws_subnet" "private_subnets" {
	vpc_id            = "${aws_vpc.main.id}"
	count              = "${length(data.aws_availability_zones.available.names)}"
	cidr_block        = "${data.cidr_network.cake.subnet_cidrs["private_az${count.index + 1}"]}"
	availability_zone = "${data.aws_availability_zones.available.names[count.index]}"
}
resource "aws_subnet" "elb_subnets" {
	vpc_id            = "${aws_vpc.main.id}"
	count              = "${length(data.aws_availability_zones.available.names)}"
	cidr_block        = "${data.cidr_network.cake.subnet_cidrs["elb_az${count.index + 1}"]}"
	availability_zone = "${data.aws_availability_zones.available.names[count.index]}"
}
```

The order of the subnets creation will impact the actual CIDR values. The order
is decided by order of appearance in the file.

```
private_az1 = 10.0.0.0/24
private_az2 = 10.0.1.0/24
private_az3 = 10.0.2.0/24
public_az1  = 10.0.3.0/25
public_az2  = 10.0.3.128/25
public_az3  = 10.0.4.0/25
elb_az1     = 10.0.4.128/28
elb_az2     = 10.0.4.144/28
elb_az3     = 10.0.4.160/28
```

## Argument Reference

The following arguments are supported:

* `cidr_block` - (Required) The CIDR block for the entire network (aka supernet)

* `subnet` - (Required) The list of subnets to create. Each subnet has a name
    and mask
  * `name` - (Required) The name of the subnet. The subnets will be created in
      the order of the name sorted ascending.
  * `mask` - (Required) The desired subnet mask for the subnet

## Attributes Reference

The following attributes are exported:

* `subnet_cidr` - A map with keys of the subnet name and the value is the CIDR
    notation network (e.g. 10.0.0.0/24)

