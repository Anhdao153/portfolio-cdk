package vpc

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsssm"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type VpcStackProps struct {
	awscdk.StackProps
}

func NewStackVPC(scope constructs.Construct, id string, props *VpcStackProps) awsec2.Vpc {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)
	cirdMask := float64(24)
	vpc := awsec2.NewVpc(stack, jsii.String("PortfolioVpc"),
		&awsec2.VpcProps{
			Cidr:   jsii.String("10.0.0.0/16"),
			MaxAzs: jsii.Number(2),
			SubnetConfiguration: &[]*awsec2.SubnetConfiguration{
				{
					CidrMask:   &cirdMask,
					Name:       jsii.String("IngressSubnet"),
					SubnetType: awsec2.SubnetType_PUBLIC,
				},
				{
					CidrMask:   &cirdMask,
					Name:       jsii.String("RdsSubnet"),
					SubnetType: awsec2.SubnetType_PRIVATE_ISOLATED,
				},
				{
					CidrMask:   &cirdMask,
					Name:       jsii.String("Application"),
					SubnetType: awsec2.SubnetType_PRIVATE_WITH_NAT,
				},
			},
			NatGateways: jsii.Number(1),
		},
	)

	awsssm.NewStringParameter(stack, jsii.String("PortfolioVpcParams"),
		&awsssm.StringParameterProps{
			Description:   jsii.String("portfolio VPC"),
			ParameterName: jsii.String("/network/portfolio"),
			StringValue:   vpc.VpcId(),
		},
	)

	return vpc
}
