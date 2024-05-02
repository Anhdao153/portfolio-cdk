package backend

import (
	"cdk/config"
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	apiGateWayV2 "github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type BackendStackProps struct {
	awscdk.StackProps
	awsec2.Vpc
}
type ICommandHooks interface {
	AfterBundling(inputDir, outputDir string) []string
	BeforeBundling(inputDir, outputDir string) []string
	BeforeInstall(inputDir, outputDir string) []string
}

// CommandHooks là một cấu trúc thực hiện ICommandHooks interface
type CommandHooks struct{}

// AfterBundling triển khai phương thức AfterBundling của ICommandHooks interface
func (m *CommandHooks) AfterBundling(*string, *string) *[]*string {
	return &[]*string{}
}

// BeforeBundling triển khai phương thức BeforeBundling của ICommandHooks interface
func (m *CommandHooks) BeforeBundling(inputDir, outputDir *string) *[]*string {
	command := fmt.Sprintf("cd %s && npm ci", *inputDir)
	return &[]*string{&command}
}

// BeforeInstall triển khai phương thức BeforeInstall của ICommandHooks interface
func (m *CommandHooks) BeforeInstall(inputDir, outputDir *string) *[]*string {
	return &[]*string{}
}

func NewBackendStack(scope constructs.Construct, id string, props *BackendStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)
	// Tạo một Lambda function

	// layer := awslambda.NewLayerVersion(stack, jsii.String("NEST-Layer"), &awslambda.LayerVersionProps{
	// 	Description: jsii.String("layer nestJS include node Module"),
	// 	Code:        awslambda.AssetCode_FromAsset(jsii.String(config.SOURCE_BE+"/nodejs"), nil),
	// 	CompatibleRuntimes: &[]awslambda.Runtime{
	// 		awslambda.Runtime_NODEJS_20_X(),
	// 	},
	// })

	lambda := awslambda.NewDockerImageFunction(stack, jsii.String("NEST-JS"), &awslambda.DockerImageFunctionProps{
		Code: awslambda.DockerImageCode_FromImageAsset(jsii.String(config.SOURCE_BE), &awslambda.AssetImageCodeProps{
			AssetName: jsii.String("portfolioLambda"),
		}),
		// Vpc: props.Vpc,
	})
	api := apiGateWayV2.NewHttpApi(stack, jsii.String("NEST-API"), &apiGateWayV2.HttpApiProps{
		CorsPreflight: &apiGateWayV2.CorsPreflightOptions{
			AllowOrigins: &[]*string{jsii.String("*")}, //
			AllowMethods: &[]apiGateWayV2.CorsHttpMethod{apiGateWayV2.CorsHttpMethod_ANY},
		},
	})

	// Add a lambda proxy lambdaApi.
	lambdaApi := awsapigatewayv2integrations.NewHttpLambdaIntegration(jsii.String("lamdba-proxy"), lambda, &awsapigatewayv2integrations.HttpLambdaIntegrationProps{})

	// Add a route to api.
	api.AddRoutes(&apiGateWayV2.AddRoutesOptions{
		Integration: lambdaApi,
		Path:        jsii.String("/{proxy+}"),
		Methods: &[]apiGateWayV2.HttpMethod{
			apiGateWayV2.HttpMethod_ANY,
		},
	})

	return stack
}
