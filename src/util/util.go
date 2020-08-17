/* Copyright (C) Xiang Wang - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Xiang Wang <xwang1314@gmail.com>, August 2020
 */

package util

import "github.com/aws/aws-sdk-go/service/ec2"
import "github.com/aws/aws-sdk-go/aws/session"
import "github.com/aws/aws-sdk-go/aws/awserr"
import "github.com/aws/aws-sdk-go/aws"
import "errors"
import "log"


func ValidateArgs(nodes int, volumeSize int, subnets []string, securityGroups, instanceTypes []string) error {
    if nodes == 0 {
        return errors.New("Number of nodes can not be zero.")
    }
    if volumeSize == 0 {
        return errors.New("Volume size can not be zero.")
    }
    if len(subnets) == 0  {
        return errors.New("Must specify subnets.")
    }
    if len(securityGroups) == 0 {
        return errors.New("Must specify securityGroups.")
    }
    if  len(subnets) != nodes || len(securityGroups) != nodes || len(instanceTypes) != nodes {
        return errors.New("Number of subnets, securityGroups, instanceTypes must equal to number of nodes.")
    }
    return nil
}

func GetCreateFleetRequestTemplate(nodes int,
                                    volumeSize int,
                                    amiId string,
                                    subnets []string,
                                    securityGroups,
                                    instanceTypes []string,
                                    spot int) *ec2.CreateFleetInput {

    template := &ec2.CreateFleetInput {
        // DryRun is set because of testing
        DryRun: aws.Bool(true),
        LaunchTemplateConfigs: []*ec2.FleetLaunchTemplateConfigRequest {
            {
                LaunchTemplateSpecification: &ec2.FleetLaunchTemplateSpecificationRequest {
                    LaunchTemplateId: aws.String("lt-0e8c754449b27161c"),
                    Version: aws.String("1"),
                },
                Overrides: []*ec2.FleetLaunchTemplateOverridesRequest {
                    {
                        AvailabilityZone: aws.String("us-east-1"),
                        InstanceType: aws.String("t1.micro"),
                        SubnetId: aws.String("ididid"),
                    },
                },
            },
        },
        TargetCapacitySpecification: &ec2.TargetCapacitySpecificationRequest {
            TotalTargetCapacity: aws.Int64(2),
            DefaultTargetCapacityType: aws.String("spot"),
        },
    }
    return template
}

func CreateFleet(requestBody *ec2.CreateFleetInput) (*ec2.CreateFleetOutput) {
    svc := ec2.New(session.New())
    result, err := svc.CreateFleet(requestBody)
    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            default:
                log.Fatal(aerr.Error())
            }
        } else {
            // Print the error, cast err to awserr.Error to get the Code and
            // Message from an error.
            log.Fatal(err.Error())
        }
    }
    return result
}
