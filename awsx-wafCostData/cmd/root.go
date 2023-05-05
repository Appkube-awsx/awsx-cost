package cmd

import (
	"log"
	"os"

	"github.com/Appkube-awsx/awsx-wafCostdata/authenticator"
	"github.com/Appkube-awsx/awsx-wafCostdata/client"
	"github.com/Appkube-awsx/awsx-wafCostdata/cmd/wafcmd"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/spf13/cobra"
)

var AwsxWafListCmd = &cobra.Command{
	Use:   "getWafListDetails",
	Short: "getWafListDetails command gets resource counts",
	Long:  `getWafListDetails command gets resource counts details of an AWS account`,

	Run: func(cmd *cobra.Command, args []string) {

		log.Println("Command get waf list started")
		vaultUrl := cmd.PersistentFlags().Lookup("vaultUrl").Value.String()
		accountNo := cmd.PersistentFlags().Lookup("accountId").Value.String()
		region := cmd.PersistentFlags().Lookup("zone").Value.String()
		acKey := cmd.PersistentFlags().Lookup("accessKey").Value.String()
		secKey := cmd.PersistentFlags().Lookup("secretKey").Value.String()
		crossAccountRoleArn := cmd.PersistentFlags().Lookup("crossAccountRoleArn").Value.String()
		externalId := cmd.PersistentFlags().Lookup("externalId").Value.String()

		authFlag := authenticator.AuthenticateData(vaultUrl, accountNo, region, acKey, secKey, crossAccountRoleArn, externalId)

		if authFlag {
			getWebAclList(region, crossAccountRoleArn, acKey, secKey, externalId)
		}

	},
}

func getWebAclList(region string, crossAccountRoleArn string, accessKey string, secretKey string, externalId string) (*waf.ListWebACLsOutput, error) {
	log.Println(" aws waf list details count summary")
	dbclient := client.GetClient(region, crossAccountRoleArn, accessKey, secretKey, externalId)
	dbRequest := waf.ListWebACLsInput{}
	dbclusterResponse, err := dbclient.ListWebACLs(&dbRequest)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	log.Println(dbclusterResponse)
	return dbclusterResponse, err
}

func Execute() {
	err := AwsxWafListCmd.Execute()
	if err != nil {
		log.Fatal("There was some error while executing the CLI: ", err)
		os.Exit(1)
	}
}

func init() {
	AwsxWafListCmd.AddCommand(wafcmd.GetCostDataCmd)
	AwsxWafListCmd.AddCommand(wafcmd.GetCostSpikeCmd)

	AwsxWafListCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxWafListCmd.PersistentFlags().String("accountId", "", "aws account number")
	AwsxWafListCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxWafListCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxWafListCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxWafListCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws crossAccountRoleArn Required")
	AwsxWafListCmd.PersistentFlags().String("externalId", "", "aws externalId Required")

}
