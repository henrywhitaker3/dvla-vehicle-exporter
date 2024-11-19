# DVLA Vehicle Exporter

A prometheus exporter to extract vehicle details from DVLA APIs

## Setup

### VES API

You need to request and API key fro the Vehicle Enquiry Service API [here](https://developer-portal.driver-vehicle-licensing.api.gov.uk/apis/vehicle-enquiry-service/vehicle-enquiry-service-description.html#vehicle-enquiry-service-ves-api-guide), then set the config value `vesApiKey`.

### Config

```yaml
logLevel: error
interval: 1h
vesApiKey: XXXXXXX
vehicles:
  - AB12CDE
```

## Metrics

| Name | Type | Description |
| --- | --- | --- |
| `dvla_vehicle_details_collection_errors_count` | counter | The number of errors encountered when collecting vehicle details |
| `dvla_vehicle_details` | counter | Always set to 1, containes other vehicle details as labels (e.g. co2 emissions) |
| `dvla_vehicle_details_tax_expiry_seconds` | gauge | The number of seconds until tax expiry |
| `dvla_vehicle_details_mot_expiry_seconds` | gauge | The number of seconds until MOT expiry |
| `dvla_vehicle_details_tax_status` | gauge | Whether the vehicle is taxed (1 = taxed, 0 = not taxed) |
| `dvla_vehicle_details_mot_status` | gauge | Whether the vehicle is MOT'd (1 = MOT'd, 0 = not MOT'd) |
