# ec2fleet

`ec2fleet` is a command line tool that can be used to provision AWS fleets.

## Instructions

#### 1. Clone and build this repo
NOTE: First time building this package may take longer because `go` needs to fetch ec2 dependencies
```
git clone https://github.com/xwang0818/ec2fleet.git
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
./ec2fleet -nodes=2 -volumeSize=4 -subnets=subnet-15288a34,subnet-d68bfc9b -securityGroups=sg-0e6218c9c2826b9dd -instanceTypes=t2.micro,t2.micro
```

### Using environment variables
Modify etc/env.config to include all the inputs
```
source etc/env.config
./ec2fleet -env
```

### Using JSON config file
Modify etc/config.json to include all the inputs
```
./ec2fleet -configFile=etc/config.json
```
