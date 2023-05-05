# WAF CLi's

## To list all the WAF webAclId ,run the following command:

```bash
awsx-waf --zone <zone> --acccessKey <acccessKey> --secretKey <secretKey> --crossAccountRoleArn <crossAccountRoleArn> --externalId <externalId>
```

## To retrieve the cost and usage details of a specific WAF wafcmd run the following command:

```bash
awsx-waf getCostData -t <table> --zone <zone> --accessKey <accessKey> --secretKey <secretKey> --crossAccountRoleArn <crossAccountRoleArn> --external <externalId>
```

## To retrieve the cost Spikes of a specific WAF wafcmd, run the following command:

```bash
awsx-waf GetCostSpike -t <table> --zone <zone> --accessKey <accessKey> --secretKey <secretKey> --crossAccountRoleArn <crossAccountRoleArn> --external <externalId>  --granularity <granularity> --startDate <startDate> --endDate <endDate>
```
