package frontend

import (
	"cdk/config"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscertificatemanager"
	cloudFront "github.com/aws/aws-cdk-go/awscdk/v2/awscloudfront"
	origins "github.com/aws/aws-cdk-go/awscdk/v2/awscloudfrontorigins"
	s3 "github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	s3Deploy "github.com/aws/aws-cdk-go/awscdk/v2/awss3deployment"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type FrontendStackProps struct {
	awscdk.StackProps
}

func NewFrontendStacks(scope constructs.Construct, id string, props *FrontendStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)
	frontendBucket := s3.NewBucket(stack, jsii.String("portfolioBucket"), &s3.BucketProps{
		AutoDeleteObjects: jsii.Bool(true),
		BlockPublicAccess: s3.BlockPublicAccess_BLOCK_ALL(),
		RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
	})
	frontEndDistributions := cloudFront.NewDistribution(stack, jsii.String("portfolioDistribution"), &cloudFront.DistributionProps{
		DefaultBehavior: &cloudFront.BehaviorOptions{
			Origin:               origins.NewS3Origin(frontendBucket, &origins.S3OriginProps{}),
			ViewerProtocolPolicy: cloudFront.ViewerProtocolPolicy_REDIRECT_TO_HTTPS,
			// CachePolicy:          cloudFront.CachePolicy_CACHING_DISABLED(),
			// AllowedMethods:       cloudFront.AllowedMethods_ALLOW_ALL(),
			// OriginRequestPolicy:  cloudFront.OriginRequestPolicy_ALL_VIEWER(),
		},

		DefaultRootObject: jsii.String("index.html"),
		ErrorResponses: &[]*cloudFront.ErrorResponse{
			{
				HttpStatus:         jsii.Number(404),
				ResponseHttpStatus: jsii.Number(200),
				Ttl:                awscdk.Duration_Seconds(jsii.Number(10)),
				ResponsePagePath:   jsii.String("/index.html"),
			},
			{
				HttpStatus:         jsii.Number(403),
				ResponseHttpStatus: jsii.Number(200),
				Ttl:                awscdk.Duration_Seconds(jsii.Number(10)),
				ResponsePagePath:   jsii.String("/index.html"),
			},
		},
		DomainNames: &[]*string{
			jsii.String(config.WEB_DOMAIN),
		},
		Certificate: awscertificatemanager.Certificate_FromCertificateArn(stack, jsii.String("webPortfolioCertificate"), &config.ARN_SSL_DOMAIN),
	})
	s3Deploy.NewBucketDeployment(stack, jsii.String("portfolioBucketDeployment"), &s3Deploy.BucketDeploymentProps{
		DistributionPaths: &[]*string{
			jsii.String("/*"),
		},
		Distribution: frontEndDistributions,
		Sources: &[]s3Deploy.ISource{
			s3Deploy.Source_Asset(jsii.String(config.SOURCE_FE+"/build"), nil),
		},
		DestinationBucket: frontendBucket,
	})

	awscdk.NewCfnOutput(stack, jsii.String("CloudFrontURL"), &awscdk.CfnOutputProps{
		Value:       frontEndDistributions.DomainName(),
		Description: jsii.String("The distribution portfolio URL"),
		ExportName:  jsii.String("PortfolioCloudFrontURL"),
	})
	awscdk.NewCfnOutput(stack, jsii.String("BucketURL"), &awscdk.CfnOutputProps{
		Value:       frontendBucket.BucketDomainName(),
		Description: jsii.String("The bucket portfolio s3 URL"),
		ExportName:  jsii.String("PortfolioBucketURL"),
	})
	return stack
}
