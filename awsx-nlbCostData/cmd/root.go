package cmd

import (
	"log"
	"os"

	"github.com/Appkube-awsx/awsx-elbv2/authenticator"
	"github.com/Appkube-awsx/awsx-elbv2/client"
	"github.com/Appkube-awsx/awsx-elbv2/cmd/loadbalancercmd"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/spf13/cobra"
)

var AwsxLoadBalancerCmd = &cobra.Command{
	Use:   "get List of loadbalancers",
	Short: "get List of loadbalancers command gets resource counts",
	Long:  `get List of loadbalancers command gets resource counts details of an AWS account`,

	Run: func(cmd *cobra.Command, args []string) {

		log.Println("Command List of loadbalancers started")
		vaultUrl := cmd.PersistentFlags().Lookup("vaultUrl").Value.String()
		accountNo := cmd.PersistentFlags().Lookup("accountId").Value.String()
		region := cmd.PersistentFlags().Lookup("zone").Value.String()
		acKey := cmd.PersistentFlags().Lookup("accessKey").Value.String()
		secKey := cmd.PersistentFlags().Lookup("secretKey").Value.String()
		crossAccountRoleArn := cmd.PersistentFlags().Lookup("crossAccountRoleArn").Value.String()
		externalId := cmd.PersistentFlags().Lookup("externalId").Value.String()

		authFlag := authenticator.AuthenticateData(vaultUrl, accountNo, region, acKey, secKey, crossAccountRoleArn, externalId)

		if authFlag {
			getListloadbalncers(region, crossAccountRoleArn, acKey, secKey, externalId)
		}
	},
}

// json.Unmarshal
func getListloadbalncers(region string, crossAccountRoleArn string, accessKey string, secretKey string,  externalId string) (*elbv2.DescribeLoadBalancersOutput, error) {
	log.Println("getting load balncer list summary")

	listlbClient := client.GetClient(region, crossAccountRoleArn, accessKey, secretKey, externalId)
	lbRequest := &elbv2.DescribeLoadBalancersInput{}
	lbResponse, err := listlbClient.DescribeLoadBalancers(lbRequest)
	if err != nil {
		log.Fatalln("Error:in getting  load balncer list", err)
	}
	log.Println(lbResponse)
	return lbResponse, err
}

func Execute() {
	err := AwsxLoadBalancerCmd.Execute()
	if err != nil {
		log.Fatal("There was some error while executing the CLI: ", err)
		os.Exit(1)
	}
}

func init() {

	AwsxLoadBalancerCmd.AddCommand(loadbalancercmd.GetCostDataCmd)
	AwsxLoadBalancerCmd.AddCommand(loadbalancercmd.GetCostSpikeCmd)

	AwsxLoadBalancerCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxLoadBalancerCmd.PersistentFlags().String("accountId", "", "aws account number")
	AwsxLoadBalancerCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxLoadBalancerCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxLoadBalancerCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxLoadBalancerCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws crossAccountRoleArn is required")
	AwsxLoadBalancerCmd.PersistentFlags().String("externalId", "", "aws external id auth")

}
