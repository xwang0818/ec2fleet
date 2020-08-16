package api

import "github.com/aws/aws-sdk-go/service/ec2"
import "github.com/aws/aws-sdk-go/aws/session"
import "github.com/aws/aws-sdk-go/aws/awserr"
import "github.com/aws/aws-sdk-go/aws"
import "fmt"

func CreateFleet() {
    svc := ec2.New(session.New())
    input := &ec2.CreateFleetInput {
        DryRun: aws.Bool(true),
        LaunchTemplateConfigs: []*ec2.FleetLaunchTemplateConfigRequest {
        },
        TargetCapacitySpecification: &ec2.TargetCapacitySpecificationRequest {
            TotalTargetCapacity: aws.Int64(5),
        },

    }

    result, err := svc.CreateFleet(input)
    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            default:
                fmt.Println(aerr.Error())
            }
        } else {
            // Print the error, cast err to awserr.Error to get the Code and
            // Message from an error.
            fmt.Println(err.Error())
        }
        return
    }
    fmt.Println(result)
}
