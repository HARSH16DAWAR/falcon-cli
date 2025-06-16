# Falcon CLI

A command-line interface for interacting with the CrowdStrike Falcon API.

## Features

- Interactive configuration setup

## Installation

```bash
go install github.com/HARSH16DAWAR/falcon-cli@latest
```

## Configuration

Before using the CLI, you need to configure your Falcon API credentials. Create a configuration file at `~/.falcon/config.yaml` with the following content:

```yaml
falcon:
  client_id: "your_client_id"
  client_secret: "your_client_secret"
  cloud_region: "us-1"  # or your preferred region
```

## Usage

### List Hosts

To list all hosts in your Falcon environment:

```bash
falcon-cli hosts
```

This will display:
- Total number of hosts found
- Query execution time
- Trace ID for the request
- Any errors that occurred during the request

Example output:
```
Found 150 hosts
Query time: 0.25 seconds
Trace ID: abc123def456
```

## Development

### Prerequisites

- Go 1.21 or later
- CrowdStrike Falcon API credentials

### Building from Source

1. Clone the repository:
```bash
git clone https://github.com/HARSH16DAWAR/falcon-cli.git
cd falcon-cli
```

2. Build the project:
```bash
go build
```

3. Run the CLI:
```bash
./falcon-cli
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Author

HARSH DAWAR - [harsh160102@gmail.com](mailto:harsh160102@gmail.com) 

## TODO 

[] add error handling to the client
[] add the ability for the user to save filters ( only work on hosts for now )
[] add the ability to export to CSV 
[] implement reusable filter system:
  - Create centralized filter management
  - Support predefined filters in config
  - Add filter validation and composition
  - Enable filter caching for performance
[] implement command composition:
  - Add machine-readable (JSON) output format for hosts command
  - Add human-readable (table) output format
  - Enable pipe-based command chaining
  - Implement describe-hosts command with both direct and pipe input support
[] add filter templates and versioning support
[] implement dynamic filter system:
  - Add template-based filters with parameter support
  - Add parameter validation for common patterns (AWS IDs, regions, etc.)
  - Add template categories (system, user, team)
  - Add template documentation and examples
  - Add template sharing and versioning
  - Add integration with other tools (AWS CLI, Terraform)
  - Add template chaining support
  - Add template export/import functionality
[] implement simplified output strategy:
  - Add clean table format for human-readable output (default)
  - Add JSON format for machine-readable output and debugging
  - Add CSV format for analytics and data processing
  - Ensure consistent output structure across all commands 