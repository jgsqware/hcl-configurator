# HCL Generator for Cloud Run Services

A Go application that generates HCL configuration files for Terraform with parameterizable Cloud Run service configurations.

## Features

- Generate HCL files with `locals.cloudrun_services` structure
- Support for various service features:
  - Firebase authentication
  - Cloud Run invoker permissions
  - Storage bucket access (creator, reader, writer)
  - Firestore access
  - MySQL and PostgreSQL database access with granular grants
- JSON-based configuration for easy parameterization
- CLI interface for flexibility

## Usage

### Build the application:
```bash
go build -o envhcl
```

### Run with configuration file:
```bash
./envhcl -config example_config.json -output terraform.hcl
```

### Command line options:
- `-config`: Path to JSON configuration file (required)
- `-output`: Output HCL file path (default: "terraform.hcl")

## Configuration Format

The application uses JSON configuration files to define services and their features. Example:

```json
{
  "services": {
    "service-name": {
      "features": {
        "firebaseauth": "viewer",
        "cloudrun_invoker": true,
        "bucket_creator": ["bucket1"],
        "firestore_access": true,
        "bucket_reader": ["bucket1", "bucket2"],
        "bucket_writer": ["bucket3"],
        "mysql_access": true,
        "postgres_access": true
      },
      "mysql_grants": {
        "read": ["table1", "table2"]
      },
      "postgres_grants": {
        "read": ["table1"],
        "read_write": ["table2"],
        "read_write_delete": ["table3"]
      }
    }
  }
}
```

## Supported Features

- `firebaseauth`: String value for Firebase auth role
- `cloudrun_invoker`: Boolean for Cloud Run invoker permissions
- `bucket_creator`: Array of bucket names for creator access
- `bucket_reader`: Array of bucket names for read access
- `bucket_writer`: Array of bucket names for write access
- `firestore_access`: Boolean for Firestore access
- `mysql_access`: Boolean for MySQL access
- `postgres_access`: Boolean for PostgreSQL access

## Database Grants

### MySQL Grants
- `read`: Array of table names for read access

### PostgreSQL Grants
- `read`: Array of table names for read access
- `read_write`: Array of table names for read/write access
- `read_write_delete`: Array of table names for full access

## Output

The application generates HCL files with the following structure:

```hcl
locals {
  cloudrun_services = {
    service-name = {
      features = {
        // service features
      }
      mysql_grants = {
        // mysql grants if applicable
      }
      postgres_grants = {
        // postgres grants if applicable
      }
    }
  }
}
```