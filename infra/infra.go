package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type InfraStackProps struct {
	awscdk.StackProps
}

func UserServiceStack(scope constructs.Construct, id string, props *InfraStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	httpApi := awsapigatewayv2.NewHttpApi(stack, jsii.String("UserServiceAPI"), nil)

	getUserFn := awslambda.NewFunction(stack, jsii.String("UserServiceGetUser"), &awslambda.FunctionProps{
		FunctionName: jsii.String("UserServiceGetUser"),
		Runtime:      awslambda.Runtime_PROVIDED_AL2023(),
		Code:         awslambda.Code_FromAsset(jsii.String("../bin/user_service/routes/get_user"), nil),
		Handler:      jsii.String("bootstrap"),
	})

	getUserIntegration := awsapigatewayv2integrations.NewHttpLambdaIntegration(jsii.String("UserServiceGetUserIntegration"), getUserFn, nil)

	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path: jsii.String("/users/{user_id}"),
		Methods: &[]awsapigatewayv2.HttpMethod{
			awsapigatewayv2.HttpMethod_GET,
		},
		// TODO: Add the authorizer to the API
		Authorizer:  nil,
		Integration: getUserIntegration,
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	UserServiceStack(app, "InfraStack", &InfraStackProps{
		awscdk.StackProps{
			Env:         env(),
			Description: jsii.String("User Service"),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	return nil
}
