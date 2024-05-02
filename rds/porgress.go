package rds

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	rds "github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type RdsPostgreSQLClusterStackProps struct {
	awscdk.StackProps
	awsec2.Vpc
}

func NewRdsMySqlClusterStack(scope constructs.Construct, id string, props *RdsPostgreSQLClusterStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, &id, &sprops)

	rds.NewDatabaseInstance(stack, jsii.String("hello"), &rds.DatabaseInstanceProps{
		InstanceType: awsec2.InstanceType_Of(awsec2.InstanceClass_T3, awsec2.InstanceSize_MICRO),
		Engine: rds.DatabaseInstanceEngine_Postgres(&rds.PostgresInstanceEngineProps{
			Version: rds.PostgresEngineVersion_VER_16(),
		}),
		DatabaseName: jsii.String("postgres"),
		Credentials:  rds.Credentials_FromGeneratedSecret(jsii.String("postgres"), &rds.CredentialsBaseOptions{}),
		Vpc:          props.Vpc,
		VpcSubnets: &awsec2.SubnetSelection{
			SubnetType: awsec2.SubnetType_PUBLIC,
		},
		MultiAz: jsii.Bool(false),
	})
	return stack
}
