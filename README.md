# Mohua

A CLI tool for real-time monitoring of AWS SageMaker resources and cost analysis.

## Key Features

- 🔍 SageMaker Resource Monitoring
  - Check status of Endpoints, Notebook Instances, and Studio Applications
  - Fast resource information retrieval through parallel processing
- 💰 Cost Analysis
  - Current cumulative costs
  - Hourly costs
  - Monthly projected costs
- 📊 Flexible Output Formats
  - Color-coded table view (default)
  - JSON output

## Prerequisites

- Go 1.16 or higher
- AWS CLI configured
- AWS IAM access permissions

## Installation

```bash
# Clone the repository
git clone https://github.com/watany-dev/mohua.git
cd mohua

# Build and install
make
```

## Usage

Basic usage examples:

```bash
# Display in table format
./mohua

# Output in JSON format
./mohua --region us-east-1 --json
```

### Command Line Options

- `--region, -r`: Specify AWS region
- `--json, -j`: Output in JSON format

## Output Example

```text
Type            Name               Status     Instance      Running Time
Endpoint        ml-endpoint        InService  ml.t3.medium  72h 15m 
Notebook        dev-notebook       Running    ml.t3.medium  168h 30m
```

## Development

### Testing

The project includes both unit tests and integration tests:

```bash
# Run unit tests only (default)
make test

# Run integration tests (requires AWS credentials)
make test-integ

# Run all tests (both unit and integration)
make test-all
```

#### Test Organization
- Unit tests: Tests that don't require AWS credentials
- Integration tests: Tests that interact with AWS services
  - Requires valid AWS credentials
  - Automatically skipped if credentials are not available
  - Use build tags to separate from unit tests

## Additional Information

- Refer to [ADR](docs/adr/) for architecture and design details
- This tool provides estimated cost calculations. Verify exact billing in the AWS Console

## License

[MIT License](LICENSE)
