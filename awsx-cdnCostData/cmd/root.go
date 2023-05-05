package cmd

import (
	"log"
	"os"

	"github.com/Appkube-awsx/awsx-cloudfrontCostData/authenticator"
	"github.com/Appkube-awsx/awsx-cloudfrontCostData/client"
	"github.com/Appkube-awsx/awsx-cloudfrontCostData/cmd/cloudfrontcmd"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/spf13/cobra"
)

var awsxCloudFunctionCmd = &cobra.Command{
	Use:   "cloudFunctionListDetails",
	Short: "cloudFunctionListDetails command gets resource counts",
	Long:  `cloudFunctionListDetails command gets resource counts details of an AWS account`,

	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Command get cloudFunction List Details started")
		vaultUrl := cmd.PersistentFlags().Lookup("vaultUrl").Value.String()
		accountNo := cmd.PersistentFlags().Lookup("accountId").Value.String()
		region := cmd.PersistentFlags().Lookup("zone").Value.String()
		acKey := cmd.PersistentFlags().Lookup("accessKey").Value.String()
		secKey := cmd.PersistentFlags().Lookup("secretKey").Value.String()
		crossAccountRoleArn := cmd.PersistentFlags().Lookup("crossAccountRoleArn").Value.String()
		externalId := cmd.PersistentFlags().Lookup("externalId").Value.String()

		authFlag := authenticator.AuthenticateData(vaultUrl, accountNo, region, acKey, secKey, crossAccountRoleArn,  externalId)

		if authFlag {
			cloudFunctionList(region, crossAccountRoleArn, acKey, secKey,  externalId)
		}
	},
}

func cloudFunctionList(region string, crossAccountRoleArn string, accessKey string, secretKey string,  externalId string) (*cloudfront.ListFunctionsOutput, error) {
	log.Println("Getting aws cloudFunction Count summary")
	getClient := client.GetClient(region, crossAccountRoleArn, accessKey, secretKey, externalId)
	input := &cloudfront.ListFunctionsInput{}
	functionResponse, err := getClient.ListFunctions(input)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	log.Println(functionResponse)
	return functionResponse, err
}

func Execute() {
	err := awsxCloudFunctionCmd.Execute()
	if err != nil {
		log.Fatal("There was some error while executing the CLI: ", err)
		os.Exit(1)
	}
}

func init() {
	
	awsxCloudFunctionCmd.AddCommand(cloudfrontcmd.GetCostDataCmd)
	awsxCloudFunctionCmd.AddCommand(cloudfrontcmd.GetCostSpikeCmd)

	awsxCloudFunctionCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	awsxCloudFunctionCmd.PersistentFlags().String("accountId", "", "aws account number")
	awsxCloudFunctionCmd.PersistentFlags().String("zone", "", "aws region")
	awsxCloudFunctionCmd.PersistentFlags().String("accessKey", "", "aws access key")
	awsxCloudFunctionCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	awsxCloudFunctionCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws crossAccountRoleArn is required")
	awsxCloudFunctionCmd.PersistentFlags().String("externalId", "", "aws external id auth")

}
