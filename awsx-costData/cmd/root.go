package cmd

import (
	"log"
	"os"

	"github.com/Appkube-awsx/awsx-costData/authenticator"
	"github.com/Appkube-awsx/awsx-costData/client"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/spf13/cobra"
)

//  awsxCostAData Cmd represents the base Command
var AwsxCostDataCmd = &cobra.Command{
	Use:   "get Cost Data Details",
	Short: "get Cost Data Details command gets resource counts",
	Long:  `get Cost Data Details command gets resource counts details of an AWS account`,

	Run: func(cmd *cobra.Command, args []string) {
        // Required flag for cost Data
		log.Println("Command get Cost Data started")
		vaultUrl := cmd.PersistentFlags().Lookup("vaultUrl").Value.String()
		accountNo := cmd.PersistentFlags().Lookup("accountId").Value.String()
		region := cmd.PersistentFlags().Lookup("zone").Value.String()
		acKey := cmd.PersistentFlags().Lookup("accessKey").Value.String()
		secKey := cmd.PersistentFlags().Lookup("secretKey").Value.String()
		crossAccountRoleArn := cmd.PersistentFlags().Lookup("crossAccountRoleArn").Value.String()
		externalId := cmd.PersistentFlags().Lookup("externalId").Value.String()
		// Retrieve all service data through serviceName flag
		serviceName := cmd.PersistentFlags().Lookup("serviceName").Value.String()
		// Retrieve value of granularity flag
		granularity := cmd.PersistentFlags().Lookup("granularity").Value.String()
		// Retrieve value of start and end Date
		startDate := cmd.PersistentFlags().Lookup("startDate").Value.String()
		endDate := cmd.PersistentFlags().Lookup("endDate").Value.String()
        
		authFlag := authenticator.AuthenticateData(vaultUrl, accountNo, region, acKey, secKey, crossAccountRoleArn, externalId, serviceName)
         
		if authFlag {
			getClusterCostDetail(region, crossAccountRoleArn, acKey, secKey, externalId, serviceName, granularity, startDate, endDate)
		}
	},
}
//  function to get cost for all service for the given time period
// json.Unmarshal
func getClusterCostDetail(region string, crossAccountRoleArn string, accessKey string, secretKey string, externalId string, serviceName string, granularity string, startDate string, endDate string) (*costexplorer.GetCostAndUsageOutput, error) {
	log.Println("Getting cost data for all service")
	costClient := client.GetCostClient(region, crossAccountRoleArn, accessKey, secretKey, externalId, serviceName)
    
	sName := serviceName
	var filter = &costexplorer.Expression{
		And: []*costexplorer.Expression{
			{
				Dimensions: &costexplorer.DimensionValues{
					Key: aws.String("SERVICE"),
					Values: []*string{
						aws.String(sName),
					},
				},
			},
			{
				Dimensions: &costexplorer.DimensionValues{
					Key: aws.String("RECORD_TYPE"),
					Values: []*string{
						aws.String("Credit"),
					},
				},
			},
		},
	}
	// compare  both condition of serviceName
	if serviceName == "ALL" || serviceName == "" {
		sName = ""
		filter = &costexplorer.Expression{
			Or: []*costexplorer.Expression{
				{
					Dimensions: &costexplorer.DimensionValues{
						Key: aws.String("SERVICE"),
						Values: []*string{
							aws.String(sName),
						},
					},
				},
				{
					Dimensions: &costexplorer.DimensionValues{
						Key: aws.String("RECORD_TYPE"),
						Values: []*string{
							aws.String("Credit"),
						},
					},
				},
			},
		}
	}
	input := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &costexplorer.DateInterval{
			// Use input dates to the function for monthly granularity
			Start: aws.String(startDate),
			End:   aws.String(endDate),
		},

		Metrics: []*string{
			aws.String("UNBLENDED_COST"),
			aws.String("BLENDED_COST"),
			aws.String("AMORTIZED_COST"),
			aws.String("NET_AMORTIZED_COST"),
		},
		GroupBy: []*costexplorer.GroupDefinition{
			{
				Type: aws.String("DIMENSION"),
				Key:  aws.String("SERVICE"),
			},
			{
				Type: aws.String("DIMENSION"),
				Key:  aws.String("REGION"),
			},
		},
		Granularity: aws.String("MONTHLY"),
		Filter:      filter,
	}

	costData, err := costClient.GetCostAndUsage(input)
	if err != nil {
		log.Fatalln("Error: in getting cost data", err)
	}
	log.Println(costData)
	return costData, err
}

func Execute() {
	err := AwsxCostDataCmd.Execute()
	if err != nil {
		log.Fatal("There was some error while executing the CLI: ", err)
		os.Exit(1)
	}
}
//  initialization the above flags which we used 
func init() {

	AwsxCostDataCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxCostDataCmd.PersistentFlags().String("accountId", "", "aws account number")
	AwsxCostDataCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxCostDataCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxCostDataCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxCostDataCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws crossAccountRoleArn is required")
	AwsxCostDataCmd.PersistentFlags().String("externalId", "", "aws external id auth")
	AwsxCostDataCmd.PersistentFlags().String("serviceName", "", "aws serviceName auth")
	AwsxCostDataCmd.PersistentFlags().String("granularity", "", "aws granularity")
	AwsxCostDataCmd.PersistentFlags().String("startDate", "", "aws startDate")
	AwsxCostDataCmd.PersistentFlags().String("endDate", "", "aws endDate")

}
