package cmd

import (
	"log"
	"os"

	"github.com/Appkube-awsx/awsx-rdsCostData/authenticator"
	"github.com/Appkube-awsx/awsx-rdsCostData/client"
	"github.com/Appkube-awsx/awsx-rdsCostData/cmd/dbInstancecmd"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/spf13/cobra"
)

var AwsxdbInstanceMetadataCmd = &cobra.Command{
	Use:   "DescribeRDSMetdataDetails",
	Short: "DescribeRDSMetdataDetails command gets resource counts",
	Long:  `DescribeRDSMetdataDetails command gets resource counts details of an AWS account`,

	Run: func(cmd *cobra.Command, args []string) {

		log.Println("Command rds cost Data started")
		vaultUrl := cmd.PersistentFlags().Lookup("vaultUrl").Value.String()
		accountNo := cmd.PersistentFlags().Lookup("accountId").Value.String()
		region := cmd.PersistentFlags().Lookup("zone").Value.String()
		acKey := cmd.PersistentFlags().Lookup("accessKey").Value.String()
		secKey := cmd.PersistentFlags().Lookup("secretKey").Value.String()
		crossAccountRoleArn := cmd.PersistentFlags().Lookup("crossAccountRoleArn").Value.String()
		externalId := cmd.PersistentFlags().Lookup("externalId").Value.String()

		authFlag := authenticator.AuthenticateData(vaultUrl, accountNo, region, acKey, secKey, crossAccountRoleArn, externalId)

		if authFlag {
			getRDSResourceList(region, crossAccountRoleArn, acKey, secKey, externalId)
		}

	},
}

func getRDSResourceList(region string, crossAccountRoleArn string, accessKey string, secretKey string, externalId string) (*rds.DescribeDBInstancesOutput, error) {
	log.Println(" aws describe rds instance metadata count summary")
	dbclient := client.GetClient(region, crossAccountRoleArn, accessKey, secretKey, externalId)
	dbRequest := rds.DescribeDBInstancesInput{}
	dbResponse, err := dbclient.DescribeDBInstances(&dbRequest)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	log.Println(dbResponse)

	
	return dbResponse, err
}

func Execute() {
	err := AwsxdbInstanceMetadataCmd.Execute()
	if err != nil {
		log.Fatal("There was some error while executing the CLI: ", err)
		os.Exit(1)
	}
}

func init() {
	
	AwsxdbInstanceMetadataCmd.AddCommand(dbInstancecmd.GetCostDataCmd)
	AwsxdbInstanceMetadataCmd.AddCommand(dbInstancecmd.GetCostSpikeCmd)

	AwsxdbInstanceMetadataCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxdbInstanceMetadataCmd.PersistentFlags().String("accountId", "", "aws account number")
	AwsxdbInstanceMetadataCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxdbInstanceMetadataCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxdbInstanceMetadataCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxdbInstanceMetadataCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws crossAccountRoleArn Required")
	AwsxdbInstanceMetadataCmd.PersistentFlags().String("externalId", "", "aws externalId Required")

}
