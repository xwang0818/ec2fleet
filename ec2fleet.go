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
const amiIdDefault = "ami-0bbe28eb2173f6167" // ubuntu-18.04
const instanceTypeDefault = "t3.micro"

func main () {
    // Flags
    // mandatory
    nodesPtr          := flag.Int("nodes", 0, "Number of Nodes\n(Require)")
    subnetsPtr        := flag.String("subnets", "", "Subnets\neg. --subnets=sub1,sub2,...\n(Require)")
    securityGroupsPtr := flag.String("securityGroups", "", "Security groups\neg. --securityGroups=sg1,sg2,...\n(Require)")
    // optionalmultiAttachVolumeSize
    instanceTypesPtr  := flag.String("instanceTypes", "", "Instance Types\n(Optional)")
    volumeSizePtr     := flag.Int("volumeSize", 0, "Multi attach volume size\n(Optional)")
    amiIdPtr          := flag.String("amiId", "", "Amazon Machine Image\n(Optional)")
    configPtr         := flag.String("configFile", "", "JSON config file\neg. --configFile=etc/config.json\n(Optional)")
    envPtr            := flag.Bool("env", false, "Use environment variables")
    flag.Parse()

    var nodes int
    var volumeSize int
    var amiId string
    var subnets, securityGroups, instanceTypes []string

    // These zone names are obtained from cli `aws ec2 describe-availability-zones`
    // TODO: this can be dynamically retrieved from API `func (*EC2) DescribeAvailabilityZones`
    var availabilityZones = []string{ "us-east-1a", "us-east-1b" }
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
                                                            int64(volumeSize),
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
    util.CreateFleet(createFleetInput)

    // clean up launch template
    util.DeleteLaunchTemplate(launchTemplateId)

    os.Exit(0)
}
