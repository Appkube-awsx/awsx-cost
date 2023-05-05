package cmd

import (
	"log"
	"os"

	"github.com/Appkube-awsx/awsx-eksCostData/authenticator"
	"github.com/Appkube-awsx/awsx-eksCostData/client"
	"github.com/Appkube-awsx/awsx-eksCostData/cmd/clustercmd"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/spf13/cobra"
)

var AwsxClusterMetadataCmd = &cobra.Command{
	Use:   "getListClusterMetaDataDetails",
	Short: "getListClusterMetaDataDetails command gets resource counts",
	Long:  `getListClusterMetaDataDetails command gets resource counts details of an AWS account`,

	Run: func(cmd *cobra.Command, args []string) {

		log.Println("Command getElementDetails started")
		vaultUrl := cmd.PersistentFlags().Lookup("vaultUrl").Value.String()
		accountNo := cmd.PersistentFlags().Lookup("accountId").Value.String()
		region := cmd.PersistentFlags().Lookup("zone").Value.String()
		acKey := cmd.PersistentFlags().Lookup("accessKey").Value.String()
		secKey := cmd.PersistentFlags().Lookup("secretKey").Value.String()
		crossAccountRoleArn := cmd.PersistentFlags().Lookup("crossAccountRoleArn").Value.String()
		externalId := cmd.PersistentFlags().Lookup("externalId").Value.String()

		authFlag := authenticator.AuthenticateData(vaultUrl, accountNo, region, acKey, secKey, crossAccountRoleArn, externalId)

		if authFlag {
			getListCluster(region, crossAccountRoleArn, acKey, secKey, externalId)
		}
	},
}

// type Tags struct {
// 	Environment string `json:"Environment"`
// }

// type Cluster struct {
// 	Tags Tags
// }

// type Response struct {
// 	Cluster Cluster
// }

// json.Unmarshal
func getListCluster(region string, crossAccountRoleArn string, accessKey string, secretKey string, externalId string) (*eks.ListClustersOutput, error) {
	log.Println("getting cluster metadata list summary")

	listClusterClient := client.GetClient(region, crossAccountRoleArn, accessKey, secretKey, externalId)
	listClusterRequest := &eks.ListClustersInput{}
	listClusterResponse, err := listClusterClient.ListClusters(listClusterRequest)
	if err != nil {
		log.Fatalln("Error:in getting  cluster list", err)
	}
	log.Println(listClusterResponse)
	return listClusterResponse, err
}

// for _, clusterName := range listClusterResponse.Clusters {
// 		clusterDetails := getClusterDetails(region, accessKey, secretKey, *clusterName)

// 		// clusterDetails.Cluster.Tags.Environment
// 		var responseObject Response
// 		jsonedResponse, _ := json.Marshal(clusterDetails)
// 		json.Unmarshal([]byte(string(jsonedResponse)), &responseObject)
// 		if env == responseObject.Cluster.Tags.Environment {
// 			log.Println(responseObject.Cluster.Tags)
// 		} else {
// 			log.Println("No cluster present for Env:", env)

// 		}
// 	}

func Execute() {
	err := AwsxClusterMetadataCmd.Execute()
	if err != nil {
		log.Fatal("There was some error while executing the CLI: ", err)
		os.Exit(1)
	}
}

func init() {
	
	AwsxClusterMetadataCmd.AddCommand(clustercmd.GetCostDataCmd)
	AwsxClusterMetadataCmd.AddCommand(clustercmd.GetCostSpikeCmd)

	AwsxClusterMetadataCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxClusterMetadataCmd.PersistentFlags().String("accountId", "", "aws account number")
	AwsxClusterMetadataCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxClusterMetadataCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxClusterMetadataCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxClusterMetadataCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws crossAccountRoleArn is required")
	AwsxClusterMetadataCmd.PersistentFlags().String("externalId", "", "aws external id auth")

}
