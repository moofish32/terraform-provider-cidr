---
layout: "cidr"
page_title: "CIDR: subnet"
sidebar_current: "docs-cidr-datasource-subnet"
description: |-
  Used to create subnet configurations of sequential networks
---

# cidr_subnet

Calculates the desired subnets from the CIDR block and given offset value

## Example Usage

```hcl
variable	"cidr_block" { default = "192.168.0.0/21" } 

data "aws_availability_zones" "available" {}

data "cidr_subnet" "private" {
	cidr_block = "${var.cidr_block}"
	subnet_mask = 24 
	subnet_count = "${length(data.aws_availability_zones.available.names)}"
}

data "cidr_subnet" "public" {
	cidr_block = "${var.cidr_block}"
	subnet_mask = 25
	subnet_count = "${length(data.aws_availability_zones.available.names)}"
	start_after = "${data.cidr_subnet.private.max_subnet}"
}

data "cidr_subnet" "elb" {
	cidr_block = "${var.cidr_block}"
	subnet_mask = 28
	subnet_count = "${length(data.aws_availability_zones.available.names)}"
	start_after = "${data.cidr_subnet.public.max_subnet}"
}

resource "aws_vpc" "main" {
	cidr_block           = "${var.cidr_block}"
	enable_dns_hostnames = true
	enable_dns_support   = true
}

resource "aws_subnet" "public_subnets" {
	vpc_id            = "${aws_vpc.main.id}"
	count              = "${length(data.aws_availability_zones.available.names)}"
	cidr_block        = "${data.cidr_subnet.public.subnet_cidrs[count.index]}"
	availability_zone = "${data.aws_availability_zones.available.names[count.index]}"
}
resource "aws_subnet" "private_subnets" {
	vpc_id            = "${aws_vpc.main.id}"
	count              = "${length(data.aws_availability_zones.available.names)}"
	cidr_block        = "${data.cidr_subnet.private.subnet_cidrs[count.index]}"
	availability_zone = "${data.aws_availability_zones.available.names[count.index]}"
}
resource "aws_subnet" "elb_subnets" {
	vpc_id            = "${aws_vpc.main.id}"
	count              = "${length(data.aws_availability_zones.available.names)}"
	cidr_block        = "${data.cidr_subnet.elb.subnet_cidrs[count.index]}"
	availability_zone = "${data.aws_availability_zones.available.names[count.index]}"
}
```

The above configuration will create the following subnets: 

```
		private_subnet1 = 192.168.0.0/24
		private_subnet2 = 192.168.1.0/24
		private_subnet3 = 192.168.2.0/24
		public_subnet1  = 192.168.3.0/25
		public_subnet2  = 192.168.3.128/25
		public_subnet3  = 192.168.4.0/25
		elb_subnet2     = 192.168.4.128/28 #Note the order here
		elb_subnet3     = 192.168.4.144/28
		elb_subnet1     = 192.168.4.160/28
```

## Argument Reference

The following arguments are supported:

* `cidr_block` - (Required) The network to create subnets within 

* `subnet_mask` - (Required) The desired subnet mask for each subent

* `start_after` - (Optional) A CIDR network within the CIDR block to offset the
    start of subnet creation from (e.g. 10.0.2.0/24 would start all subnets in
    10.0.3.0 and above)

* `subnet_count` - (Option) The number of sequential subnets to create

## Attributes Reference

The following attributes are exported:

* `max_subnet` - The highest IP based subnet generated from this configuration

* `subnet_cidrs` - The set of sequential subnets created from this
    configuration. While the set will be sequential and the order will be
    consistent between runs, the order may not be ascending. For example, if the
    configuration called for 3 /24 networks you will get set will contain three
    /24 networks, but the first array value may not be lowest IP range.
