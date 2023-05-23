package cmd

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/Appkube-awsx/awsx-costSpikes/authenticator"
	"github.com/Appkube-awsx/awsx-costSpikes/client"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/spf13/cobra"
)

var AwsxCostSpikeCmd = &cobra.Command{
	Use:   "get Cost Spikes Data Details",
	Short: "get Cost Spikes Data Details command gets resource counts",
	Long:  `get Cost Spikes Data Details command gets resource counts details of an AWS account`,

	Run: func(cmd *cobra.Command, args []string) {

		log.Println("Command get Cost Spikes Data started")
		vaultUrl := cmd.PersistentFlags().Lookup("vaultUrl").Value.String()
		accountNo := cmd.PersistentFlags().Lookup("accountId").Value.String()
		region := cmd.PersistentFlags().Lookup("zone").Value.String()
		acKey := cmd.PersistentFlags().Lookup("accessKey").Value.String()
		secKey := cmd.PersistentFlags().Lookup("secretKey").Value.String()
		crossAccountRoleArn := cmd.PersistentFlags().Lookup("crossAccountRoleArn").Value.String()
		externalId := cmd.PersistentFlags().Lookup("externalId").Value.String()
		serviceName := cmd.PersistentFlags().Lookup("serviceName").Value.String()
		granularity := cmd.PersistentFlags().Lookup("granularity").Value.String()
		startDate := cmd.PersistentFlags().Lookup("startDate").Value.String()
		endDate := cmd.PersistentFlags().Lookup("endDate").Value.String()

		authFlag := authenticator.AuthenticateData(vaultUrl, accountNo, region, acKey, secKey, crossAccountRoleArn, externalId, serviceName)

		if authFlag {
			wrapperCostSpike(region, crossAccountRoleArn, acKey, secKey, externalId, serviceName, granularity, startDate, endDate)
		}
	},
}

// json.Unmarshal
// Wrapper function to get cost, spike percentage and print them.
func wrapperCostSpike(region string, crossAccountRoleArn string, acKey string, secKey string, externalId string, serviceName string, granularity string, startDate string, endDate string) (string, error) {
	costClient := client.GetCostClient(region, crossAccountRoleArn, acKey, secKey, externalId, serviceName)
	fmt.Println("Cost Spike for :" + serviceName)

	switch granularity {
	case "DAILY":
		// Call CostSpikes function for the date period
		layout := "2006-01-02" // layout must be the same as start date format
		startDate, err := time.Parse(layout, startDate)
		if err != nil {
			fmt.Println("Error parsing start date:", err)
			return "", err
		}
		endDate, err := time.Parse(layout, endDate)
		if err != nil {
			fmt.Println("Error parsing end date:", err)
			return "", err
		}

		for d := startDate; d.Before(endDate.AddDate(0, 0, 1)); d = d.AddDate(0, 0, 1) {
			prevDate := d.AddDate(0, 0, -1)
			// fmt.Printf("%s (%s)\n", d.Format("2006-01-02"), prevDate.Format("2006-01-02"))
			CostSpikes(region, crossAccountRoleArn, acKey, secKey, externalId, serviceName, granularity, prevDate.Format("2006-01-02"), d.Format("2006-01-02"), costClient)
		}
		return "", nil

	case "MONTHLY":
		// Call CostSpikes function for the month period
		layout := "2006-01-02" // layout must be the same as start date format
		startDate, err := time.Parse(layout, startDate)
		if err != nil {
			fmt.Println("Error parsing start date:", err)
			return "", err
		}
		endDate, err := time.Parse(layout, endDate)
		if err != nil {
			fmt.Println("Error parsing end date:", err)
			return "", err
		}

		for d := startDate; d.Before(endDate.AddDate(0, 1, 0)); d = d.AddDate(0, 1, 0) {
			if d.Equal(endDate) {
				break
			}
			prevDate := d.AddDate(0, -1, 0)
			// fmt.Printf("%s (%s)\n", d.Format("2006-01-02"), prevDate.Format("2006-01-02"))
			CostSpikes(region, crossAccountRoleArn, acKey, secKey, externalId, serviceName, granularity, prevDate.Format("2006-01-02"), d.Format("2006-01-02"), costClient)
		}
		return "", nil

	case "HOURLY":
		// Call CostSpikes function for the hour period
		layout := "2006-01-02T15:04:05Z" // layout must be the same as start date format
		startDateTime, err := time.Parse(layout, startDate)
		if err != nil {
			fmt.Println("Error parsing start date:", err)
			return "", err
		}
		endDateTime, err := time.Parse(layout, endDate)
		if err != nil {
			fmt.Println("Error parsing end date:", err)
			return "", err
		}

		for d := startDateTime; d.Before(endDateTime); d = d.Add(time.Hour) {
			prevHour := d.Add(-time.Hour)
			// fmt.Println(prevHour.Format(layout), d.Format(layout))
			CostSpikes(region, crossAccountRoleArn, acKey, secKey, externalId, granularity, serviceName, prevHour.Format("2006-01-02T15:04:05Z"), d.Format("2006-01-02T15:04:05Z"), costClient)
		}

		return "", nil

	default:
		return "", fmt.Errorf("unsupported granularity: %s", granularity)
	}

}

// Function to do the cost comparison.
func CostSpikes(region string, crossAccountRoleArn string, accessKey string, secretKey string, externalId string, serviceName string, granularity string, startDateTime string, endDateTime string, costClient *costexplorer.CostExplorer) (string, error) {
	// Get cost data for latest time period
	startCostData, err := ServiceCostDetails(region, crossAccountRoleArn, accessKey, secretKey, externalId, serviceName, granularity, startDateTime, endDateTime, costClient)
	if err != nil {
		log.Fatalln("Error: in getting cost data for start date", err)
		return "", err
	}

	var endCost float64

	// Get cost data for previous time period
	switch granularity {
	case "MONTHLY":
		// Get start date and end date for previous time period
		previousStartDateTime, previousEndDateTime, err := generateDatesForMonthlyGranularity(startDateTime, endDateTime)
		if err != nil {
			log.Fatalln("Error: in getting previous time period date", err)
			return "", err
		}
		endCostData, err := ServiceCostDetails(region, crossAccountRoleArn, accessKey, secretKey, externalId, serviceName, granularity, previousStartDateTime, previousEndDateTime, costClient)
		if err != nil {
			log.Fatalln("Error: in getting cost data for end date", err)
			return "", err
		}
		endCost = convertCostDataToFloat(endCostData)

	default:
		endCostData, err := ServiceCostDetails(region, crossAccountRoleArn, accessKey, secretKey, externalId, serviceName, granularity, endDateTime, endDateTime, costClient)
		if err != nil {
			log.Fatalln("Error: in getting cost data for end date", err)
			return "", err
		}
		endCost = convertCostDataToFloat(endCostData)
	}

	// Convert cost data to float and positive
	startCost := convertCostDataToFloat(startCostData)
	// endCost = convertCostDataToFloat(endCostData)

	// Calculate cost difference
	costDifference := endCost - startCost

	// Calculate cost difference percentage
	costDifferencePercentage := (costDifference / startCost) * 100

	// |2/12/2023| 5.95|+3%| --- Print format
	if costDifferencePercentage >= 0 {
		output := fmt.Sprintf("|%s| %f | +%f%% |", endDateTime, endCost, costDifferencePercentage)
		fmt.Println(output)
		return output, nil
	}
	if costDifferencePercentage < 0 {
		output := fmt.Sprintf("|%s| %f | %f%% |", endDateTime, endCost, costDifferencePercentage)
		fmt.Println(output)
		return output, nil
	}

	return "", nil
}

// Function to get cost for a given service for given time period.
func ServiceCostDetails(region string, crossAccountRoleArn string, accessKey string, secretKey string, externalId string, serviceName string, granularity string, startDateTime string, endingDateTime string, costClient *costexplorer.CostExplorer) (string, error) {
	// costClient := client.GetCostClient(region, crossAccountRoleArn, accessKey, secretKey, externalId)

	// Get endDateTime from startDateTime for DAILY/WEEKLY/HOURLY
	endDateTime, err := generateEndDateTime(granularity, startDateTime)
	if err != nil {
		log.Fatalln("Error: in generating end date time", err)
		return "", err
	}

	var start, end string
	switch granularity {
	case "DAILY":
		// fmt.Println("startDateTime: ", startDateTime, " endDateTime: ", endDateTime)
		start = startDateTime //"2023-03-01"
		end = endDateTime
	case "MONTHLY":
		// fmt.Println( startDateTime, endDateTime)
		start = startDateTime
		end = endDateTime // Use input dates to the function for monthly granularity
	case "HOURLY":
		start = startDateTime
		end = endDateTime
	default:
		return "", fmt.Errorf("unsupported granularity: %s", granularity)
	}
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
			Start: aws.String(start),
			End:   aws.String(end),
		},
		Metrics: []*string{
			aws.String("UNBLENDED_COST"),
			aws.String("BLENDED_COST"),
			aws.String("AMORTIZED_COST"),
		},
		GroupBy: []*costexplorer.GroupDefinition{
			{
				Type: aws.String("DIMENSION"),
				Key:  aws.String("SERVICE"),
			},
			{
				Type: aws.String("DIMENSION"),
				Key:  aws.String("RECORD_TYPE"),
			},
		},
		Granularity: aws.String(granularity),
		Filter:      filter,
	}

	costData, err := costClient.GetCostAndUsage(input)
	if err != nil {
		log.Fatalln("Error: in getting cost data", err)
	}

	// fmt.Println("Cost Data: ", costData)

	// Extract the blended cost from the response (change this to get the cost you want)
	blendedCost := float64(0)
	for _, result := range costData.ResultsByTime {
		for _, group := range result.Groups {
			if metrics := group.Metrics; metrics != nil {
				if blendedCostMetric, ok := metrics["BlendedCost"]; ok && blendedCostMetric != nil && blendedCostMetric.Amount != nil {
					if amount, err := strconv.ParseFloat(*blendedCostMetric.Amount, 64); err == nil {
						blendedCost += math.Abs(amount)
					}
				}
			}
		}
	}

	// log.Println(costData.ResultsByTime)
	return strconv.FormatFloat(blendedCost, 'f', -1, 64), err
}

// Function to generate endDateTime according to granularity
func generateEndDateTime(granularity string, startDateTime string) (string, error) {

	switch granularity {
	case "DAILY":
		layout := "2006-01-02" // layout must be the same as start date format
		startDate, err := time.Parse(layout, startDateTime)
		if err != nil {
			return "", err
		}
		// Add 1 day to the start date to get the end date
		endDate := startDate.AddDate(0, 0, 1)
		end := endDate.Format(layout)

		// fmt.Println("Start Date: ", startDate, "End Date: ", endDate)
		return end, nil

	case "MONTHLY":
		layout := "2006-01-02" // layout must be the same as start date format
		startDate, err := time.Parse(layout, startDateTime)
		if err != nil {
			return "", err
		}
		// Add 1 month to the start date to get the end date
		endDate := startDate.AddDate(0, 1, 0)
		end := endDate.Format(layout)

		// fmt.Println("Start Date: ", startDate, "End Date: ", endDate)
		return end, nil

	case "HOURLY":
		layout := "2006-01-02T15:04:05Z" // layout must be the same as start date format
		startDate, err := time.Parse(layout, startDateTime)
		if err != nil {
			return "", err
		}
		// Add 1 hour to the start date to get the end date
		endDate := startDate.Add(time.Hour)
		end := endDate.Format(layout)

		// fmt.Println("Start Date: ", startDate, "End Date: ", endDate)
		return end, nil

	default:
		return "", fmt.Errorf("unsupported granularity: %s", granularity)
	}
}

func generateDatesForMonthlyGranularity(startDateTime string, endingDateTime string) (string, string, error) {
	layout := "2006-01-02" // layout must be the same as start date format
	startDate, err := time.Parse(layout, startDateTime)
	if err != nil {
		return "", "", err
	}

	endDate, err := time.Parse(layout, endingDateTime)
	if err != nil {
		return "", "", err
	}

	// Calculate time period between start and end date
	timePeriod := endDate.Sub(startDate)

	// Calclulate starting date of previous time period by subtracting time period from start date
	previousStartDate := startDate.AddDate(0, 0, -int(timePeriod.Hours()/24))

	// subtract 1 day from start date to get the end date of previous time period
	previousEndDate := startDate.AddDate(0, 0, -1)

	// Convert the dates to string
	start := previousStartDate.Format(layout)
	end := previousEndDate.Format(layout)

	// fmt.Println("Start Date: ", startDate, "End Date: ", endDate)
	// fmt.Println("Start Date: ", start, "End Date: ", end)
	return start, end, nil
}

func convertCostDataToFloat(CostData string) float64 {
	// Convert the cost data to float
	cost, err := strconv.ParseFloat(CostData, 64)
	if err != nil {
		log.Fatalln("Error: in converting cost data to float", err)
	}

	// Convert the cost to positive if it is negative
	if cost < 0 {
		cost = cost * -1
	}

	return cost
}

func Execute() {
	err := AwsxCostSpikeCmd.Execute()
	if err != nil {
		log.Fatal("There was some error while executing the CLI: ", err)
		os.Exit(1)
	}
}

func init() {

	AwsxCostSpikeCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxCostSpikeCmd.PersistentFlags().String("accountId", "", "aws account number")
	AwsxCostSpikeCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxCostSpikeCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxCostSpikeCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxCostSpikeCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws crossAccountRoleArn is required")
	AwsxCostSpikeCmd.PersistentFlags().String("externalId", "", "aws external id ")
	AwsxCostSpikeCmd.PersistentFlags().String("serviceName", "", "aws serviceName auth")
	AwsxCostSpikeCmd.PersistentFlags().String("granularity", "", "aws granularity")
	AwsxCostSpikeCmd.PersistentFlags().String("startDate", "", "aws startDate ")
	AwsxCostSpikeCmd.PersistentFlags().String("endDate", "", "aws endDate auth")

}
