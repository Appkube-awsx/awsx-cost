package cmd

import (
	"log"
	"os"

	"github.com/Appkube-awsx/awsx-glue/authenticator"
	"github.com/Appkube-awsx/awsx-glue/client"
	"github.com/Appkube-awsx/awsx-glue/cmd/gluecmd"
	"github.com/aws/aws-sdk-go/service/glue"
	"github.com/spf13/cobra"
)

var AwsxGlueCmd = &cobra.Command{
	Use:   "getGlueDetails",
	Short: "getGlueDetails command gets resource counts",
	Long:  `getGlueDetails command gets resource counts details of an AWS account`,

	Run: func(cmd *cobra.Command, args []string) {

		log.Println("Command get glue details started")
		vaultUrl := cmd.PersistentFlags().Lookup("vaultUrl").Value.String()
		accountNo := cmd.PersistentFlags().Lookup("accountId").Value.String()
		region := cmd.PersistentFlags().Lookup("zone").Value.String()
		acKey := cmd.PersistentFlags().Lookup("accessKey").Value.String()
		secKey := cmd.PersistentFlags().Lookup("secretKey").Value.String()
		crossAccountRoleArn := cmd.PersistentFlags().Lookup("crossAccountRoleArn").Value.String()
		externalId := cmd.PersistentFlags().Lookup("externalId").Value.String()

		authFlag := authenticator.AuthenticateData(vaultUrl, accountNo, region, acKey, secKey, crossAccountRoleArn, externalId)

		if authFlag {
			getListGlue(region, crossAccountRoleArn, acKey, secKey, externalId)
		}
	
	},
}

// json.Unmarshal
func getListGlue(region string, crossAccountRoleArn string, accessKey string, secretKey string, externalId string) (*glue.ListJobsOutput, error) {
	log.Println("getting glue job  list summary")

	listClient := client.GetClient(region, crossAccountRoleArn, accessKey, secretKey, externalId)
	listRequest := &glue.ListJobsInput{}
	listResponse, err := listClient.ListJobs(listRequest)
	if err != nil {
		log.Fatalln("Error:in getting glue list", err)
	}
	log.Println(listResponse)
	return listResponse, err
}


func Execute() {
	err := AwsxGlueCmd.Execute()
	if err != nil {
		log.Fatal("There was some error while executing the CLI: ", err)
		os.Exit(1)
	}
}

func init() {
	
	AwsxGlueCmd.AddCommand(gluecmd.GetCostDataCmd)
	AwsxGlueCmd.AddCommand(gluecmd.GetCostSpikeCmd)

	AwsxGlueCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxGlueCmd.PersistentFlags().String("accountId", "", "aws account number")
	AwsxGlueCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxGlueCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxGlueCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxGlueCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws crossAccountRoleArn is required")
	AwsxGlueCmd.PersistentFlags().String("externalId", "", "aws external id auth")
	AwsxGlueCmd.PersistentFlags().String("jobName", "", "aws jobName auth")


}
