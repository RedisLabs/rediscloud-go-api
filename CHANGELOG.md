# Changelog
All notable changes to this project will be documented in this file.
See updating [Changelog example here](https://keepachangelog.com/en/1.0.0/).

## 0.43.0 (11th November 2025)

### Added:
* Added `PersistentStorageEncryptionType` and `DeletionGracePeriod` fields to `Subscription` model for customer-managed key support

## 0.42.0 (10th November 2025)

### Added:

* Added `ID` field to `Region` struct in account model
* Added Redis version support for Essentials databases: `RedisVersion` field and `UpgradeRedisVersion()` method

## 0.41.0 (3rd November 2025)

### Added:
* Added `AWSAccountID` field to `CloudDetail` struct in subscriptions model for AWS account identification
* Added `govulncheck` to CI pipeline for automated vulnerability detection
* New `.github/workflows/vulnerability.yml` workflow for parallel vulnerability checking
* Enabled Go module caching in GitHub Actions for improved build performance

### Updated:
* Updated Go toolchain to 1.25.3 to address stdlib vulnerabilities (GO-2025-4007, GO-2025-3751, GO-2025-3750, GO-2025-3749, GO-2025-3563)

### Tests:
* Added AWS account ID to subscription test fixtures (`TestSubscription_List`, `TestSubscription_Get`, `TestSubscription_Get_PublicEndpointAccess`)

## 0.40.0 (31st October 2025)

### Added:
* Added global values fields to `ActiveActiveDatabase` read response: `GlobalDataPersistence`, `GlobalSourceIP`, `GlobalPassword`, `GlobalAlerts`, and `GlobalEnableDefaultUser`
* New `TestAADatabase_Get` test for ActiveActive Get method coverage
* Updated `TestAADatabase_List` test to include global values in response

### Updated:
* Bump golang.org/x/tools to 0.38.0
* Bump github.com/avast/retry-go/v4 to 4.7.0

### Fixed:
* Fixed an issue with endpoint scripts not unmarshalling correctly.

## 0.39.0 (21st October 2025)

### Added:
* Added `AutoMinorVersionUpgrade` to Pro and Active Active models, for upgrade and create endpoints.

## 0.38.0 (15th October 2025)

### Added:
* Added `PublicEndpointAccess` field to Subscription create/update API
* New tests to check this field is used correctly

## 0.37.0 (8th October 2025)

### Added:
* Added `GlobalEnableDefaultUser` field to UpdateActiveActiveDatabase struct
* New tests for ActiveActiveCreate and ActiveActiveUpdate methods

## v0.36.5

### Changed:
* Several fields in the PrivateLink connections datatype now have data types aligned with the API.
* Fixed an issue with a malformed URL for the PrivateLink endpoint scripts and added tests

## 0.36.4

### Changed:
* Fixing the EnableDefaultUser field to be part of region properties

## 0.36.3

### Added:
* The `regionId` property is now supported on `GET` ActiveActiveRegion requests.

### Changed:
* Changing `regionId` back to an int.
* Bumping dependency versions

## v0.36.2

### Added:
* Adding model and service for new PrivateLink endpoints

### Changed:
* Modified `delete` in API client so that it takes a `requestBody` parameter.
* Updating Testify to v1.11.1

## 0.35.0

### Added:
* Adding field `security` to Active Active Database model
* Adding field `enableDefaultUser` to active active database update request

## 0.34.1

### Added:
* Adding field `RedisVersion` to Redis Active Active subscription create

## 0.34.0

### Added:
* Adding field `RedisVersion` to Redis Active Active subscription model (get and list)


## 0.33.1

### Fixed:
* Added pending status for database upgrade

## 0.33.0

### Added:
* Adding Upgrade Redis database endpoint
* Adding Get Redis versions for subscription endpoint

## 0.32.0

### Added

* Adding redisVersion support on the create, get and list endpoints for pro databases.

## 0.31.0

### Added

* Adding an API call and endpoint for updating customer-managed encryption keys (CMKs) to an existing subscription.
* Adding a new status to support the `encryption_key_pending` status of a subscription.

## 0.30.0

### Added

* Adding in support for `persistentStorageEncryptionType`, to support customer-managed encryption keys (CMKs) across pro and active-active subscription creation

## 0.29.0

### Added

* Adding a new status for the database whilst it is pending: `dynamic-endpoints-creation-pending`

## 0.28.0

### Added

* Adding an endpoint to get a TLS certificate for a database.

## 0.27.0

### Updated

* Splitting Fixed Subscription structs into response and request structs, as JSON fields differ between the two uses.

## 0.26.0

### Fixed

* Fixes to the AA regions endpoint: removing the unfilled RegionId field and correcting the JSON for the DeploymentCidr field.

## 0.25.0

### Added

* Added payment method to the fixed subscription API. This should allow payments via the marketplace as well as credit cards.

## 0.24.0

### Added
* 
* Add endpoint for Active Active regions.

## 0.23.0

### Added

* Add Query Performance Factor property to subscription and database models

## 0.22.0

### Added

* Add endpoints to the new Private Service Connect API for Pro and Active-Active subscriptions.

## 0.21.0

### Added

* Handling API rate limits to wait or retry calls depending on the current window or remaining limits.   

### Fixed

* Removed extra API call when retrieving latest backup and import Statuses.

## 0.8.0

### Added

* A `resp_version` property for Regions (creation only)

## 0.7.0

### Added

* A `redis_version` property Subscriptions
* A `resp_version` property for ActiveActive Databases (creation only)

## 0.6.0

### Added

* A `status` property for ACL Users
* A `resp_version` property for Databases
* An `enable_default_user` property for Databases

## 0.5.4

### Updated

* Fixed `alert` properties on Databases and Regions to allow empty lists

## 0.5.3

### Updated

* Fixed some json serialization rules for ACL entities

## 0.5.2

### Added

* A `status` property for ACL Rules
* A `status` property for ACL Roles

## 0.5.1

### Added

* Support for ACL APIs: `Users`, `Roles` and `Rules`

## 0.1.10 - 0.4.1

* ...

## 0.1.10
### Updated

* Dependencies for CI and Testing

### Added
* Add support for port number and backup to databases
* Add support for multiple CIDR ranges when peering

## 0.1.9 (July 5 2022)
### Added 

* Adds wrap404Error for database service GET method
* Add proxy-policy-change and -draft to pending state constants

## 0.1.8 (May 12 2022)

### Added

* Adds region attribute to VPCPeering
* Adds paymentMethod field to CreateSubscription and Subscription structs

## 0.1.7 (April 25 2022)

### Removed

* Removed the PersistentStorageEncryption (deprecated) field from the API calls.

## 0.1.6 (January 14 2022)

### Added

* Include "Content-Type: application/json" in the header of all requests to the API.

## 0.1.5 (December 9 2021)

### Added

* enableTls field for the database APIs.

## 0.1.4

### Changed
* Changed the UpdateDatabase struct to allow replicaOf to be set as empty array

## 0.1.3

### Changed
* Expanded VPC Peering with fields for GCP and AWS

## 0.1.2

### Changed
* Expanded VPC Peering with fields for GCP

## 0.1.1

### Added

### Changed
* Changed maximum number of retries when waiting for a task to finish so that it is compatible with 32bit environments.

### Removed

## 0.1.0

### Added
* List payment methods through the Accounts API 
* Cloud accounts API
* Subscription API
* Database API

### Changed

### Removed
