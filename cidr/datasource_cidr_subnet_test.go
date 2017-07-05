package cidr

import (
	"fmt"
	"testing"

	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestCidrSubnet(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			r.TestStep{
				Config: subnets,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.cidr_subnet.public", "subnet_cidrs.#", "1"),
					r.TestCheckResourceAttr("data.cidr_subnet.public", "max_subnet",
						"10.0.3.0/25"),
					r.TestCheckResourceAttr("data.cidr_subnet.private", "subnet_cidrs.#", "3"),
					r.TestCheckResourceAttr("data.cidr_subnet.private", "max_subnet",
						"10.0.2.0/24"),
					r.TestCheckResourceAttr("data.cidr_subnet.elb", "max_subnet",
						"10.0.3.128/28"),
				),
			},
		},
	},
	)
}
func TestAccCIDRSubnet(t *testing.T) {
	r.Test(t, r.TestCase{
		Providers: testAccProviders,
		Steps: []r.TestStep{
			{
				Config: vpcSubnet,
				Check: r.ComposeTestCheckFunc(
					outputCheck,
				),
			},
		},
	})
}

func outputCheck(s *terraform.State) error {
	answers := [][]string{
		[]string{"private_subnet1", "192.168.0.0/24"},
		[]string{"private_subnet2", "192.168.1.0/24"},
		[]string{"private_subnet3", "192.168.2.0/24"},
		[]string{"public_subnet1", "192.168.3.0/25"},
		[]string{"public_subnet2", "192.168.3.128/25"},
		[]string{"public_subnet3", "192.168.4.0/25"},
		// the order won't change but this one proves sorting is not in your control
		[]string{"elb_subnet2", "192.168.4.128/28"},
		[]string{"elb_subnet3", "192.168.4.144/28"},
		[]string{"elb_subnet1", "192.168.4.160/28"},
	}
	for _, ans := range answers {
		got := s.RootModule().Outputs[ans[0]]
		if ans[1] != got.Value {
			fmt.Printf("Outputs %v\n", s.RootModule().Outputs)
			return fmt.Errorf("Output expected %s, got %s\n", ans[1], got.Value)
		}
	}
	return nil
}

const subnets = `
data "cidr_subnet" "private" {
	cidr_block = "10.0.0.0/21"
	subnet_mask = 24 
	subnet_count = 3
}

data "cidr_subnet" "public" {
	cidr_block = "10.0.0.0/21"
	subnet_mask = 25
	start_after = "${data.cidr_subnet.private.max_subnet}"
}

data "cidr_subnet" "elb" {
	cidr_block = "10.0.0.0/21"
	subnet_mask = 28
	start_after = "${data.cidr_subnet.public.max_subnet}"
}
`
const vpcSubnet = `
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
output "public_subnet1"  { value = "${aws_subnet.public_subnets.0.cidr_block}" }
output "public_subnet2"  { value = "${aws_subnet.public_subnets.1.cidr_block}" }
output "public_subnet3"  { value = "${aws_subnet.public_subnets.2.cidr_block}" }
output "private_subnet1"  { value = "${aws_subnet.private_subnets.0.cidr_block}" }
output "private_subnet2"  { value = "${aws_subnet.private_subnets.1.cidr_block}" }
output "private_subnet3"  { value = "${aws_subnet.private_subnets.2.cidr_block}" }
output "elb_subnet1"  { value = "${aws_subnet.elb_subnets.0.cidr_block}" }
output "elb_subnet2"  { value = "${aws_subnet.elb_subnets.1.cidr_block}" }
output "elb_subnet3"  { value = "${aws_subnet.elb_subnets.2.cidr_block}" }
`
