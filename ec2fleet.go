/* Copyright (C) Xiang Wang - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Xiang Wang <xwang1314@gmail.com>, August 2020
 */

package main

import "util"
import "flag"
import "strings"
import "os"
import "log"
import "fmt"


const SPOT_PERCENTAGE = 80
const volumeSizeDefault = 3
const amiIdDefault = "ubuntu-18.04"
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

    var nodes, volumeSize int
    var amiId string
    var subnets, securityGroups, instanceTypes []string

    volumeSize = volumeSizeDefault
    amiId = amiIdDefault

    if *configPtr != "" {
        fmt.Println("Using JSON config file", *configPtr)
    } else if *envPtr {
        fmt.Println("Using environment variables")
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

    ///*
    fmt.Println("nodes:", nodes)
    fmt.Println("volumeSize:", volumeSize)
    fmt.Println("amiId:", amiId)
    fmt.Println("subnets:", subnets)
    fmt.Println("securityGroups:", securityGroups)
    fmt.Println("instanceTypes:", instanceTypes)
    //*/

    requestBody := util.GetCreateFleetRequestTemplate(nodes,
                                                    volumeSize,
                                                    amiId,
                                                    subnets,
                                                    securityGroups,
                                                    instanceTypes,
                                                    SPOT_PERCENTAGE)
    response := util.CreateFleet(requestBody)
    fmt.Println(response)
}
