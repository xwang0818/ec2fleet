package api

import "github.com/aws/aws-sdk-go/service/ec2"
import "github.com/aws/aws-sdk-go/aws/session"
import "github.com/aws/aws-sdk-go/aws/awserr"
import "github.com/aws/aws-sdk-go/aws"
import "fmt"


func GetCreateFleetRequestTemplate(nodes int,
                                    volumeSize int,
                                    amiId string,
                                    subnets []string,
                                    securityGroups,
                                    instanceTypes []string) *ec2.CreateFleetInput {
    
    result := &ec2.CreateFleetInput {
        // DryRun is set because of testing
        DryRun: aws.Bool(true),
        LaunchTemplateConfigs: []*ec2.FleetLaunchTemplateConfigRequest {
            {
                LaunchTemplateSpecification: &ec2.FleetLaunchTemplateSpecificationRequest {
                    LaunchTemplateId: aws.String("lt-0e8c754449b27161c"),
                    Version: aws.String("1"),
                },
            },
        },
        TargetCapacitySpecification: &ec2.TargetCapacitySpecificationRequest {
            TotalTargetCapacity: aws.Int64(2),
            DefaultTargetCapacityType: aws.String("spot"),
        },
    }
    return result
}

func CreateFleet(requestBody *ec2.CreateFleetInput) {
    svc := ec2.New(session.New())
    result, err := svc.CreateFleet(requestBody)
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
    fmt.Println("Response: ", result)
}
