# EKS CLi's

## To list all the EKS cluster, run the following command:

```bash
awsx-eks --zone <zone> --acccessKey <acccessKey> --secretKey <secretKey> --crossAccountRoleArn <crossAccountRoleArn> --externalId <externalId>
```

# To retrieve the cost details of a specific EKS clustercmd, run the following command:

```bash
awsx-eks getCostData -t <table> --zone <zone> --accessKey <accessKey> --secretKey <secretKey> --crossAccountRoleArn <crossAccountRoleArn> --external <externalId>  --granularity <granularity> --startDate <startDate> --endDate <endDate>
```

## To retrieve the cost Spikes details of a specific EKS clustercmd, run the following command:

```bash
awsx-eks GetCostSpike -t <table> --zone <zone> --accessKey <accessKey> --secretKey <secretKey> --crossAccountRoleArn <crossAccountRoleArn> --external <externalId>  --granularity <granularity> --startDate <startDate> --endDate <endDate>
```
