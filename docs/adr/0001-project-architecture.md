# ADR-0001: SageMaker Monitor Project Architecture

## Status

Accepted

## Context

Challenges in developing a CLI tool for monitoring AWS SageMaker resources and cost analysis:
- SageMaker resource management is complex and cost tracking is difficult
- Existing tools are either too detailed or lack sufficient information
- Need for a tool that allows developers and administrators to easily understand resource status

## Decisions

1. Implemented as a CLI application using Go
   - Reasons:
     - High performance through compiled language
     - Cross-platform compatibility
     - Easy single binary distribution

2. Adoption of AWS SDK v2
   - Reasons:
     - Support for latest AWS APIs
     - Better compatibility with modern Go language features
     - Improved type safety

3. Key Architectural Components:
   - `cmd`: Command-line interface
   - `internal/sagemaker`: Interaction with AWS SageMaker
   - `internal/display`: Output formatting
   - `internal/retry`: Error handling and retry strategy

4. Concurrent Processing
   - Use of goroutines and channels for resource retrieval
   - Simultaneous monitoring of multiple resources (Endpoints, Notebooks, Studio Apps)

5. Output Formats
   - Table format (default)
   - JSON format (optional)

## Consequences

Benefits:
- Lightweight and fast SageMaker resource monitoring tool
- Flexible output options
- Minimal dependencies
- High extensibility

Drawbacks:
- Omission of detailed configuration information
- Simplified resource information

Anticipated Impacts:
- Improved resource management efficiency for developers
- Increased visibility of AWS SageMaker usage costs
- Contribution to open-source community

## References

- [AWS SDK for Go V2 Documentation](https://aws.github.io/aws-sdk-go-v2/docs/)
- [Cobra CLI Framework](https://github.com/spf13/cobra)
- [Go Concurrency Guide](https://go.dev/doc/effective_go#concurrency)
