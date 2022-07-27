# Changelog
All notable changes to this project will be documented in this file.
See updating [Changelog example here](https://keepachangelog.com/en/1.0.0/).

## 0.1.10 (unreleased)
### Updated

* Dependencies for CI and Testing

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
