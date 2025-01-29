# SageMaker Monitor

## Overview

SageMaker Monitor is a lightweight and efficient CLI tool for monitoring AWS SageMaker compute resources and performing real-time cost analysis. This tool tracks the basic status and associated costs of the following SageMaker resources:

- Endpoints
- Notebook Instances
- Studio Applications

## Key Features

- 🔍 Basic resource state monitoring
- 💰 Simplified cost analysis
  - Current cumulative costs
  - Hourly costs
  - Monthly projected costs
- 📊 Flexible output formats (Table/JSON)

## Important Notes

This tool is a lightweight implementation focused on basic monitoring and estimated cost calculation for SageMaker resources:
- Endpoint instance types are simplified and do not retrieve detailed configuration information
- Notebook instance volume size information is omitted
- For more detailed information, please check the AWS Console

## Prerequisites

- Go 1.16 or higher
- AWS CLI configured
- AWS IAM access permissions

## Installation

### Method 1: Build from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/mohua.git
cd mohua

# Install dependencies
go mod tidy

# Build
go build -o mohua

# Optional: Install
go install
```

### Method 2: Download Binary

Download the latest binary from the [Releases](https://github.com/yourusername/mohua/releases) page.

## Usage

### Basic Usage

```bash
# Display in table format
./mohua --region us-east-1

# Output in JSON format
./mohua --region us-east-1 --json
```

### Command Line Options

- `--region, -r`: Specify AWS region (required)
- `--json, -j`: Output in JSON format

## Environment Configuration

AWS credentials can be configured using one of the following methods:

1. AWS CLI configuration
```bash
aws configure
```

2. Environment variables
```bash
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key
export AWS_DEFAULT_REGION=us-east-1
```

3. IAM role (EC2 or ECS)

## Output Examples

### Table Format
```
Type            Name                          Status     Instance       Running Time  Hourly($)  Current($)  Projected($)
Endpoint        my-ml-endpoint                InService  unknown        72h 15m       $1.24      $89.54      $912.80
Notebook        dev-notebook                  Running    ml.t3.medium   168h 30m      $0.11      $18.54      $81.40

Total Current Cost: $108.08    Projected Monthly Cost: $994.20
```

### JSON Format
```json
[
  {
    "resourceType": "Endpoint",
    "name": "my-ml-endpoint",
    "status": "InService",
    "instanceType": "unknown",
    "runningTime": "72h 15m",
    "hourlyCost": 1.24,
    "currentCost": 89.54,
    "projectedMonthlyCost": 912.80
  },
  ...
]
```

## Troubleshooting

- AWS authentication error: Check IAM policies and permissions
- Region specification error: Use the correct region name
- Unexpected results: Verify AWS SDK version

## Contributing

Pull requests and feature suggestions are welcome. See `CONTRIBUTING.md` for details.

## License

This project is published under the [MIT License](LICENSE).

## Disclaimer

This tool is for informational purposes only. Always verify billing information through the AWS Console. Cost calculations are approximate and may differ from actual billing amounts.
