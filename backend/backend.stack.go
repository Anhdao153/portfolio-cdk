package backend

import (
	"cdk/config"

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

func NewBackendStack(scope constructs.Construct, id string, props *BackendStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)
	// Tạo một Lambda function

	lambda := awslambda.NewDockerImageFunction(stack, jsii.String("NEST-JS"), &awslambda.DockerImageFunctionProps{
		Code: awslambda.DockerImageCode_FromImageAsset(jsii.String(config.SOURCE_BE), &awslambda.AssetImageCodeProps{
			AssetName: jsii.String("portfolioLambda"),
		}),
		// Vpc: props.Vpc,
	})
	api := apiGateWayV2.NewHttpApi(stack, jsii.String("NEST-API"), &apiGateWayV2.HttpApiProps{
		CorsPreflight: &apiGateWayV2.CorsPreflightOptions{
			AllowOrigins:     &[]*string{jsii.String("http://localhost:3000"), jsii.String("https://profile.nguyenducthien.net")},
			AllowMethods:     &[]apiGateWayV2.CorsHttpMethod{apiGateWayV2.CorsHttpMethod_ANY},
			AllowHeaders:     &[]*string{jsii.String("*")},
			AllowCredentials: jsii.Bool(true),
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
