package main

import (
	"fmt"

	dockercreds "github.com/docker/docker-credential-helpers/credentials"
	awscreds "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/Jeff-Hanson/docker-credential-helpers/awsssm"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func main() {

	sess, err := session.NewSession(aws.NewConfig().
	  WithRegion("us-east-2").
	  WithCredentials(awscreds.NewSharedCredentials("", "credo-auth")))
	if err != nil {
		panic(fmt.Sprintf("Failed to get session: %v", err))
	}
	
	ssmClient := ssm.New(sess, aws.NewConfig().WithRegion("us-east-2"))
	dockercreds.Serve(awsssm.Awsssm{Svc: ssmClient})
}
