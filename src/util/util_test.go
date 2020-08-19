package util

import "testing"


func TestUtilValidateInputOk(t *testing.T) {
    nodes := 5
    volumeSize := 100
    subnets := []string{"sub1", "sub2", "sub3", "sub4", "sub4"}
    securityGroups := []string{"sg1"}
    instanceTypes := []string{"type1", "type1", "type1", "type1", "type1"}
    result := ValidateInputs(nodes, volumeSize, subnets, securityGroups, instanceTypes)
    if result != nil {
        t.Errorf("TestUtilValidateInput failed")
    }
}

func TestUtilValidateInputNotOkOne(t *testing.T) {
    nodes := 5
    volumeSize := 100
    subnets := []string{"sub1", "sub2", "sub3", "sub4"}
    securityGroups := []string{"sg1"}
    instanceTypes := []string{"type1", "type1", "type1", "type1", "type1"}
    result := ValidateInputs(nodes, volumeSize, subnets, securityGroups, instanceTypes)
    if result == nil {
        t.Errorf("TestUtilValidateInput failed")
    }
}

func TestUtilValidateInputNotOkTwo(t *testing.T) {
    nodes := 5
    volumeSize := 100
    subnets := []string{"", "sub2", "sub3", "sub4", "sub5"}
    securityGroups := []string{"sg1"}
    instanceTypes := []string{"type1", "type1", "type1", "type1", "type1"}
    result := ValidateInputs(nodes, volumeSize, subnets, securityGroups, instanceTypes)
    if result == nil {
        t.Errorf("TestUtilValidateInput failed")
    }
}
