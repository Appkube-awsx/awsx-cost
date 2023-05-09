- [What is awsx-costData](#awsx-getCostData)
- [How to build / Test](#how-to-build--test)
- [what it does ](#what-it-does)
- [command input](#command-input)
- [command output](#command-output)
- [How to run ](#how-to-run)

# awsx-getCostData

This is a plugin subcommand for awsx cli ( https://github.com/Appkube-awsx/awsx#awsx ) cli.

For details about awsx commands and how its used in Appkube platform , please refer to the diagram below:

![alt text](https://d1.awsstatic.com/aws-cloud-financial-managment/aws-cost-explorer-diagram.36df791eafa0210c0c5b0ccdad499e77e1d303f4.png)

This plugin subcommand will implement the Apis' related to costData for all services , primarily the following API's:

- CostData

In cost data you can receive reports that break down your costs by the hour, day, or month, by product or product resource, or by tags that you define yourself.

This cli collect data from metric / logs / traces of the costData for all services and produce the cost data in a form that Appkube Platform expects.

This CLI , interacts with other Appkube services like Appkube vault , Appkube cloud CMDB so that it can talk with cloud services as well as filter and sort the information in terms of product/services, so that Appkube platform gets the data that it expects from the cli.

# How to write plugin subcommand

Please refer to the instruction -
https://github.com/Appkube-awsx/awsx#how-to-write-a-plugin-subcommand

It has detailed instruction on how to write a subcommand plugin , build / test / debug / publish and integrate into the main commmand.

# How to build / Test

            go run main.go
                - Program will print Calling aws-costData on console

            Another way of testing is by running go install command
            go install
            - go install command creates an exe with the name of the module (e.g. awsx-getCostData) and save it in the GOPATH
            - Now we can execute this command on command prompt as below

            awsx-getCostData --zone=us-east-1 --accessKey=xxxxxxxxxx --secretKey=xxxxxxxxxx --crossAccountRoleArn=xxxxxxxxxx  --externalId=xxxxxxxxxx

# what it does

This command implement the following functionalities -
getCostData - It will get the cost data for specific service and for all service also.

# command input

1. --valutURL = URL location of vault - that stores credentials to call API
2. --acountId = The AWS account id.
3. --zone = AWS region
4. --accessKey = Access key for the AWS account
5. --secretKey = Secret Key for the Aws Account
6. --crossAccountRoleArn = Cross Acount Rols Arn for the account.
7. --external Id = The AWS External id.
8. --granularity = DAILY,WEEKLY,MONTHLY.
9. --startDate = Required according
10. --endDate = Required according

# command output

{[
Keys: ["AWS WAF","us-east-1"],
Metrics: {
AmortizedCost: {
Amount: "-6.93950583",
Unit: "USD"
},
BlendedCost: {
Amount: "-6.93950583",
Unit: "USD"
},
UnblendedCost: {
Amount: "-6.93950583",
Unit: "USD"
},
NetAmortizedCost: {
Amount: "-6.93950583",
Unit: "USD"
}
}
]},

# How to run

From main awsx command , it is called as follows:

```bash
awsx getCostData --zone=us-east-1 --accessKey=<> --secretKey=<> --crossAccountRoleArn=<>  --externalId=<> GetCostSpike --granularity=DAILY --startDate=2023-03-01 --endDate=2023-03-10 --serviceName="ALL"
```

If you build it locally , you can simply run it as standalone command as

```bash
  go run main.go  getCostData--zone=us-east-1 --accessKey=<> --secretKey=<> --crossAccountRoleArn=<>  --externalId=<> GetCostSpike --granularity=DAILY --startDate=2023-03-01 --endDate=2023-03-10 --serviceName="ALL"
```
