# RDS CLi's

## To list all the RDS dbinstances,run the following command:

```bash
awsx-rds --zone <zone> --acccessKey <acccessKey> --secretKey <secretKey> --crossAccountRoleArn <crossAccountRoleArn> --externalId <externalId>
```

## To retrieve the cost details of a specific RDS dbInstancecmd, run the following command:

```bash
awsx-rds getCostData -t <table> --zone <zone> --accessKey <accessKey> --secretKey <secretKey> --crossAccountRoleArn <crossAccountRoleArn> --external <externalId>  --granularity <granularity> --startDate <startDate> --endDate <endDate>
```

## To retrieve the cost Spikes of a specific RDS dbInstancecmd, run the following command:

```bash
awsx-rds GetCostSpike -t <table> --zone <zone> --accessKey <accessKey> --secretKey <secretKey> --crossAccountRoleArn <crossAccountRoleArn> --external <externalId>  --granularity <granularity> --startDate <startDate> --endDate <endDate>
```
