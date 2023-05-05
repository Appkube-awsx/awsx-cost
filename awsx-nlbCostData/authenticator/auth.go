package authenticator

import (
	"log"

	"github.com/Appkube-awsx/awsx-elbv2/vault"
)

func AuthenticateData(vaultUrl string, accountNo string, region string, acKey string, secKey string, crossAccountRoleArn string,  externalId string) bool {

	if vaultUrl != "" && accountNo != "" {
		if region == "" {
			log.Fatalln("Zone not provided. Program exit")
			return false
		}
		log.Println("Getting account details")
		data, err := vault.GetAccountDetails(vaultUrl, accountNo)
		if err != nil {
			log.Println("Error in calling the account details api. \n", err)
			return false
		}
		if data.AccessKey == "" || data.SecretKey == "" || data.CrossAccountRoleArn == "" {
			log.Println("Account details not found.")
			return false
		}
		return true

	} else if region != "" && acKey != "" && secKey != "" && crossAccountRoleArn != "" && externalId != "" {
		return true
	} else {
		log.Fatal("AWS credentials like accesskey/secretkey/region/crossAccountRoleArn/externalId/granularity/startdate/endDate not provided. Program exit")
		return false
	}
}
