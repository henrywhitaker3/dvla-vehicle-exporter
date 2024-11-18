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
| `dvla_vehicle_details_tax_expiry_seconds` | gauge | The number of seconds until tax expiry |
| `dvla_vehicle_details_co2_emissions` | gauge | CO2 expiry |
| `dvla_vehicle_details_engine_capactiy` | gauge | Engine capacity |
| `dvla_vehicle_details_year_of_manufacture` | gauge | The year of vehicle manufacture |
| `dvla_vehicle_details_tax_status` | gauge | Whether the vehicle is taxed (1 = taxed, 0 = not taxed) |
| `dvla_vehicle_details_mot_status` | gauge | Whether the vehicle is MOT'd (1 = taxed, 0 = not taxed) |
