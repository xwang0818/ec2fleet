package main

import "api"
import "flag"
import "fmt"
import "strings"
import "log"
import "os"

const SPOT_PERCENTAGE = 80
const volumeSizeDefault = 3
const amiIdDefault = "ubuntu-18.04"
const instanceTypeDefault = "t3.micro"

func validateArgs(nodes int, volumeSize int, subnets []string, securityGroups, instanceTypes []string) bool {
    if nodes == 0 {
        log.Fatal("Number of nodes can not be zero.")
        return false
    }
    if volumeSize == 0 {
        log.Fatal("Volume size can not be zero.")
        return false
    }
    if len(subnets) == 0  {
        log.Fatal("Must specify subnets.")
        return false
    }
    if len(securityGroups) == 0 {
        log.Fatal("Must specify securityGroups.")
        return false
    }
    if  len(subnets) != nodes || len(securityGroups) != nodes || len(instanceTypes) != nodes {
        log.Fatal("Number of subnets, securityGroups, instanceTypes must equal to number of nodes.")
        return false
    }
    return true
}

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
    if !validateArgs(nodes, volumeSize, subnets, securityGroups, instanceTypes) {
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
    fmt.Println(SPOT_PERCENTAGE)
    requestBody := api.GetCreateFleetRequestTemplate(nodes,
                                                    volumeSize,
                                                    amiId,
                                                    subnets,
                                                    securityGroups,
                                                    instanceTypes)
    api.CreateFleet(requestBody)
}
