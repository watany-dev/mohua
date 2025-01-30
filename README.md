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

## Additional Information

- Refer to [ADR](docs/adr/) for architecture and design details
- This tool provides estimated cost calculations. Verify exact billing in the AWS Console

## License

[MIT License](LICENSE)
