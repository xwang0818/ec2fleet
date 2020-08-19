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
import "time"
import "log"
import "os"


func ValidateInputs(nodes, volumeSize int, subnets, securityGroups, instanceTypes []string) error {
    if nodes <= 0 {
        return errors.New("Number of nodes is invalid.")
    }
    if volumeSize < 4 || volumeSize > 16384 {
        return errors.New("Invalid volume size, must be between 4-16384 Gib inclusively.")
    }
    for _, sub := range subnets {
        if sub == "" {
            return errors.New("Subnet can not be empty.")
        }
    }
    for _, sg := range securityGroups {
        if sg == "" {
            return errors.New("Security group can not be empty.")
        }
    }
    for _, it := range instanceTypes {
        if it == "" {
            return errors.New("Instance type can not be empty.")
        }
    }
    if  len(subnets) != nodes || len(instanceTypes) != nodes {
        return errors.New("Number of subnets and instanceTypes must equal to number of nodes.")
    }
    return nil
}

func GetCreateLaunchTemplateInput(templateName string,
                                  amiId string,
                                  instanceTypeDefault string,
                                  securityGroups []string) *ec2.CreateLaunchTemplateInput {
    secGroups := []*string{}
    for i := range securityGroups {
        secGroups = append(secGroups, &securityGroups[i])
    }
    input := &ec2.CreateLaunchTemplateInput{
        LaunchTemplateData: &ec2.RequestLaunchTemplateData {
            ImageId:        aws.String(amiId),
            InstanceType:   aws.String(instanceTypeDefault),
            SecurityGroupIds: secGroups,
            Placement: &ec2.LaunchTemplatePlacementRequest {
                AvailabilityZone: aws.String("us-east-1a"),
            },
        },
        LaunchTemplateName: aws.String(templateName),
    }
    return input
}

func CreateLaunchTemplate(input *ec2.CreateLaunchTemplateInput) *ec2.CreateLaunchTemplateOutput {
    svc := ec2.New(session.New())
    responseBody, err := svc.CreateLaunchTemplate(input)
    if err != nil {
        log.Println("Create Launch Template error:")
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            case "DryRunOperation":
                log.Println("Create Launch Template DryRun succeeded.")
            default:
                log.Println("Create Launch Template status code: ", aerr.Code())
                log.Fatal(aerr.Error())
            }
        } else {
            log.Fatal(err.Error())
        }
        os.Exit(1)
    }
    log.Println("Launch Template created successfully:\n", responseBody)
    return responseBody
}

func DeleteLaunchTemplate(templateId string) *ec2.DeleteLaunchTemplateOutput {
    svc := ec2.New(session.New())
    input := &ec2.DeleteLaunchTemplateInput {
        LaunchTemplateId: aws.String(templateId),
    }
    responseBody, err := svc.DeleteLaunchTemplate(input)
    if err != nil {
        log.Println("Delete Launch Template error:")
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            default:
                log.Println("Delete Launch Template status code: ", aerr.Code())
                log.Fatal(aerr.Error())
            }
        } else {
            log.Fatal(err.Error())
        }
        os.Exit(1)
    }
    log.Println("Launch template", templateId, "was delete successfully.")
    return responseBody
}

func GetCreateFleetRequestInput(nodes int64,
                                launchTemplateId string,
                                subnets []string,
                                instanceTypes []string,
                                availabilityZones []string,
                                onDemandPercentage int64) *ec2.CreateFleetInput {
    onDemand := onDemandPercentage*nodes/100
    spot := nodes - onDemand
    overrides := []*ec2.FleetLaunchTemplateOverridesRequest {}
    size := int(nodes)
    for i := 0; i < size; i++ {
        az := availabilityZones[0]
        if i >= size/2 {
            az = availabilityZones[1]
        }
        overrides = append(overrides, &ec2.FleetLaunchTemplateOverridesRequest {
            AvailabilityZone: aws.String(az),
            InstanceType: aws.String(instanceTypes[i]),
            SubnetId: aws.String(subnets[i]),
        })
    }

    input := &ec2.CreateFleetInput {
        // TODO: add DryRun option for testing
        // DryRun: aws.Bool(true),
        LaunchTemplateConfigs: []*ec2.FleetLaunchTemplateConfigRequest {
            {
                LaunchTemplateSpecification: &ec2.FleetLaunchTemplateSpecificationRequest {
                    LaunchTemplateId: aws.String(launchTemplateId),
                    Version: aws.String("1"),
                },
                Overrides: overrides,
            },
        },
        SpotOptions: &ec2.SpotOptionsRequest {
            AllocationStrategy: aws.String("diversified"),
        },
        Type: aws.String("instant"),
        TargetCapacitySpecification: &ec2.TargetCapacitySpecificationRequest {
            OnDemandTargetCapacity: aws.Int64(onDemand),
            SpotTargetCapacity: aws.Int64(spot),
            TotalTargetCapacity: aws.Int64(nodes),
            DefaultTargetCapacityType: aws.String("spot"),
        },
    }
    return input
}

func CreateFleet(requestBody *ec2.CreateFleetInput) (*ec2.CreateFleetOutput, error) {
    svc := ec2.New(session.New())
    responseBody, err := svc.CreateFleet(requestBody)
    if err != nil {
        log.Println("Create Fleet error:")
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            case "DryRunOperation":
                log.Println("Create Fleet DryRun succeeded.")
            default:
                log.Println("Create Fleet status code: ", aerr.Code())
                log.Fatal(aerr.Error())
            }
        } else {
            log.Fatal(err.Error())
        }
    }
    log.Println("EC2 fleet created successfully:", responseBody)
    return responseBody, err
}

func CreateVolume(vSize int64, aZone string) *ec2.Volume {
    svc := ec2.New(session.New())
    input := &ec2.CreateVolumeInput {
        Size:               aws.Int64(vSize),
        Iops:               aws.Int64(200),
        VolumeType:         aws.String("io1"),
        AvailabilityZone:   aws.String(aZone),
        MultiAttachEnabled: aws.Bool(true),
    }
    responseBody, err := svc.CreateVolume(input)
    if err != nil {
        log.Println("Create volume error:")
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            case "DryRunOperation":
                log.Println("Create volume DryRun succeeded.")
            default:
                log.Println("Create volume status code: ", aerr.Code())
                log.Fatal(aerr.Error())
            }
        } else {
            log.Fatal(err.Error())
        }
        os.Exit(1)
    }
    log.Println("Created volume in", aZone," successfully.")
    return responseBody
}

func AttachVolume(instanceId, volumeId string) *ec2.VolumeAttachment {
    svc := ec2.New(session.New())
    input := &ec2.AttachVolumeInput {
        Device:     aws.String("/dev/sdf"),
        InstanceId: aws.String(instanceId),
        VolumeId:   aws.String(volumeId),
    }
    // Check for instance status for 180 seconds or 3 mins
    for i := 0; i < 6; i++ {
        if GetInstanceStatus(instanceId) == "running" {
            break
        }
        log.Println("Checking instance status before attaching volume. Sleep 30 seconds...")
        time.Sleep(30 * time.Second)
    }
    responseBody, err := svc.AttachVolume(input)
    if err != nil {
        log.Println("Attach volume error:")
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            default:
                log.Println("Attach volume status code: ", aerr.Code())
                log.Fatal(aerr.Error())
            }
        } else {
            log.Fatal(err.Error())
        }
        os.Exit(1)
    }
    log.Println("Volume attached successfully.")
    return responseBody
}

func GetInstanceStatus(instanceId string) string{
    svc := ec2.New(session.New())
    input := &ec2.DescribeInstanceStatusInput{
        InstanceIds: []*string{
            aws.String(instanceId),
        },
    }
    log.Println("GetInstanceStatus for instance ID:", instanceId)
    responseBody, err := svc.DescribeInstanceStatus(input)
    if err != nil {
        log.Println("GetInstanceStatus error:")
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            default:
                log.Println("GetInstanceStatus status code: ", aerr.Code())
                log.Fatal(aerr.Error())
            }
        } else {
            log.Fatal(err.Error())
        }
        return ""
    }
    log.Println("GetInstanceStatus successfully.", responseBody)
    if len(responseBody.InstanceStatuses) > 0 {
        return *responseBody.InstanceStatuses[0].InstanceState.Name
    }else {
        return ""
    }
}
