module github.com/taiidani/terraform-provider-jenkins

go 1.16

require (
	github.com/aws/aws-sdk-go v1.30.12 // indirect
	github.com/bndr/gojenkins v1.1.1
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.10.1
)

replace github.com/bndr/gojenkins v1.1.1 => github.com/yarlson/gojenkins v1.1.104
