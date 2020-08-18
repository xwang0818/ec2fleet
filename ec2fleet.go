/* Copyright (C) Xiang Wang - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Xiang Wang <xwang1314@gmail.com>, August 2020
 */

package main

import "strings"
import "util"
import "flag"
import "log"
import "os"


const ON_DEMAND_PERCENTAGE = 20
const volumeSizeDefault = 3
const amiIdDefault = "ami-0bcc094591f354be2" // ubuntu-18.04
const instanceTypeDefault = "t3.micro"

func main () {
    // Flags
    // mandatory
    nodesPtr          := flag.Int("nodes", 0, "Number of Nodes\n(Require)\neg. -nodes=2")
    subnetsPtr        := flag.String("subnets", "", "Network IDs for each instance to attach to\n(Require)\neg. -subnets=sub1,sub2,...")
    securityGroupsPtr := flag.String("securityGroups", "", "Security group IDs that will be applied on all instances\n(Require)\neg. -securityGroups=sg1,sg2,...")
    // optionalmultiAttachVolumeSize
    instanceTypesPtr  := flag.String("instanceTypes", "", "Instance types\n(Optional) Default: t3.micro.\neg. -instanceTypes=t3.micro\nMulti-Attach volume can only be attached to instance types that are Nitro System\nhttps://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-types.html#ec2-nitro-instances")
    volumeSizePtr     := flag.Int("volumeSize", 0, "Multi-attach volume size\n(Optional) Default: 3\neg. -volumeSize=4\nMin: 4 GiB, Max: 16384 GiB")
    amiIdPtr          := flag.String("amiId", "", "Amazon Machine Image ID\n(Optional) Default: ami-0bbe28eb2173f6167 (ubuntu-18.04)\neg. -amiId=ami-0bbe28eb2173f6167")
    configPtr         := flag.String("configFile", "", "JSON config file\n(Optional) Default: empty\neg. -configFile=etc/config.json")
    envPtr            := flag.Bool("env", false, "Use environment variables\n(Optional) Default: false\neg. -env")
    flag.Parse()

    var nodes int
    var volumeSize int
    var amiId string
    var subnets, securityGroups, instanceTypes []string

    // These zone names are obtained from cli `aws ec2 describe-availability-zones`
    // According to this resource https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ebs-volumes-multi.html
    // Multi-attach volume is available only in us-east-1, us-west-2, eu-west-1, and ap-northeast-2 Regions
    // TODO: this can be dynamically retrieved from API `func (*EC2) DescribeAvailabilityZones`
    var availabilityZones = []string{ "us-east-1b", "us-east-1a" }
    volumeSize = volumeSizeDefault
    amiId = amiIdDefault

    if *configPtr != "" {
        log.Println("Using JSON config file", *configPtr)
        // TODO make support for json file
    } else if *envPtr {
        log.Println("Using environment variables")
        // TODO make support for env vars
    } else {
        nodes = *nodesPtr
        subnets = strings.Split(*subnetsPtr, ",")
        securityGroups = strings.Split(*securityGroupsPtr, ",")
        if *volumeSizePtr != 0 {
            volumeSize = *volumeSizePtr
        }
        if *amiIdPtr != "" {
            amiId = *amiIdPtr
        }
        if *instanceTypesPtr != "" {
            instanceTypes = strings.Split(*instanceTypesPtr, ",")
        } else {
            instanceTypes = make([]string, nodes)
            for i := range instanceTypes {
                instanceTypes[i] = instanceTypeDefault
            }
        }
    }
    err := util.ValidateArgs(nodes, volumeSize, subnets, securityGroups, instanceTypes)
    if  err != nil {
        log.Fatal(err)
        os.Exit(1)
    }

    launchTemplateInput := util.GetCreateLaunchTemplateInput("ec2fleet-template",
                                                            amiId,
                                                            instanceTypeDefault,
                                                            securityGroups)
    log.Println("Creating Launch Template with the following parameters:\n", launchTemplateInput)

    launchTemplateResponse := util.CreateLaunchTemplate(launchTemplateInput)
    launchTemplateId := *launchTemplateResponse.LaunchTemplate.LaunchTemplateId

    createFleetInput := util.GetCreateFleetRequestInput(int64(nodes),
                                                        launchTemplateId,
                                                        subnets,
                                                        instanceTypes,
                                                        availabilityZones,
                                                        ON_DEMAND_PERCENTAGE)
    log.Println("Creating EC2 Fleet with the following parameters:\n", createFleetInput)
    fleet, err := util.CreateFleet(createFleetInput)

    // clean up launch template
    util.DeleteLaunchTemplate(launchTemplateId)

    log.Println(fleet.Instances)

    if err == nil {
        responseOne := util.CreateVolume(int64(volumeSize), availabilityZones[0])
        volumeOne := *responseOne.VolumeId
        responseTwo := util.CreateVolume(int64(volumeSize), availabilityZones[1])
        volumeTwo := *responseTwo.VolumeId
        log.Println(volumeOne, volumeTwo)
    }
    os.Exit(0)
}
