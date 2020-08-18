# ec2fleet

`ec2fleet` is a command line tool that can be used to provision AWS fleets.

## Instructions

#### 1. Clone and build this repo
NOTE: First time building this package may take longer because `go` needs to fetch ec2 dependencies
```
git clone git@github.com:xwang0818/ec2fleet.git
cd ec2fleet/
make
```

### 2. Modify etc/aws.configs to include your AWS credentials
```
source etc/aws.config
```

### 3. Use the help page to learn how to use the CLI
```
cd build/
.ec2fleet -help
# eg.
.ec2fleet -nodes=5 -volumeSize=4 -subnets=subnet1,subnet2,subnet3,subnet4,subnet5 -securityGroups=sg1,sg2
```
