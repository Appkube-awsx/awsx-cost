package loadbalancercmd

import (
	"log"

	"github.com/Appkube-awsx/awsx-elbv2/authenticator"
	"github.com/Appkube-awsx/awsx-elbv2/client"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/spf13/cobra"
)

// getConfigDataCmd represents the getConfigData command
var GetCostDataCmd = &cobra.Command{
	Use:   "getCostData",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		vaultUrl := cmd.Parent().PersistentFlags().Lookup("vaultUrl").Value.String()
		accountNo := cmd.Parent().PersistentFlags().Lookup("accountId").Value.String()
		region := cmd.Parent().PersistentFlags().Lookup("zone").Value.String()
		acKey := cmd.Parent().PersistentFlags().Lookup("accessKey").Value.String()
		secKey := cmd.Parent().PersistentFlags().Lookup("secretKey").Value.String()
		crossAccountRoleArn := cmd.Parent().PersistentFlags().Lookup("crossAccountRoleArn").Value.String()
		externalId := cmd.Parent().PersistentFlags().Lookup("externalId").Value.String()

		granularity, err := cmd.Flags().GetString("granularity")
		startDate, err := cmd.Flags().GetString("startDate")
		endDate, err := cmd.Flags().GetString("endDate")

		if err != nil {
			log.Fatalln("Error: in getting granularity flag value", err)
		}

		authFlag := authenticator.AuthenticateData(vaultUrl, accountNo, region, acKey, secKey, crossAccountRoleArn, externalId)

		if authFlag {
			getClusterCostDetail(region, crossAccountRoleArn, acKey, secKey, externalId, granularity, startDate, endDate)
		}
	},
}

func getClusterCostDetail(region string, crossAccountRoleArn string, accessKey string, secretKey string, externalId string, granularity string, startDate string, endDate string) (*costexplorer.GetCostAndUsageOutput, error) {
	log.Println("Getting cost data")
	costClient := client.GetCostClient(region, crossAccountRoleArn, accessKey, secretKey, externalId)

	input := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &costexplorer.DateInterval{
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
		Granularity: aws.String("DAILY"),
		Filter: &costexplorer.Expression{
			And: []*costexplorer.Expression{
				{
					Dimensions: &costexplorer.DimensionValues{
						Key: aws.String("SERVICE"),
						Values: []*string{
							aws.String("Amazon Elastic Load Balancing"),
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
		},
	}

	costData, err := costClient.GetCostAndUsage(input)
	if err != nil {
		log.Fatalln("Error: in getting cost data", err)
	}
	log.Println(costData)
	return costData, err
}

func init() {
	GetCostDataCmd.Flags().StringP("granularity", "t", "", "granularity name")

	if err := GetCostDataCmd.MarkFlagRequired("granularity"); err != nil {
		log.Println(err)
	}

	GetCostDataCmd.Flags().StringP("startDate", "u", "", "startDate name")

	if err := GetCostDataCmd.MarkFlagRequired("startDate"); err != nil {
		log.Println(err)
	}

	GetCostDataCmd.Flags().StringP("endDate", "v", "", "endDate name")

	if err := GetCostDataCmd.MarkFlagRequired("endDate"); err != nil {
		log.Println(err)
	}

}
