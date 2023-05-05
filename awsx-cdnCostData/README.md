# CLoudfront CLi's

## To list all the Cloudfront function,run the following command:

```bash
awsx-cloudfront --zone <zone> --acccessKey <acccessKey> --secretKey <secretKey> --crossAccountRoleArn <crossAccountRoleArn> --externalId <externalId>
```

## To retrieve the cost and usage details of a specific CDN function in cloudfrontcmd run the following command:

````bash
awsx-cloudfront getCostData -t <table> --zone <zone> --accessKey <accessKey> --secretKey <secretKey> --crossAccountRoleArn <crossAccountRoleArn> --external <externalId>

## To retrieve the cost Spikes of a specific CDN function in cloudfrontcmd run the following command:

```bash
awsx-cloudfront GetCostSpike -t <table> --zone <zone> --accessKey <accessKey> --secretKey <secretKey> --crossAccountRoleArn <crossAccountRoleArn> --external <externalId>  --granularity <granularity> --startDate <startDate> --endDate <endDate>
````
