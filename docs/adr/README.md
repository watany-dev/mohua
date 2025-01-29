# Architecture Decision Records (ADRs)

## Overview

This directory contains the key architectural design decisions for the SageMaker Monitor project. Each ADR documents specific design challenges, considered options, and the final decisions with their rationale.

## ADR List

### [ADR-0001: Project Architecture](0001-project-architecture.md)
- Adoption of Go and AWS SDK v2
- CLI tool basic design
- Component structure

### [ADR-0002: Error Handling and Retry Strategy](0002-error-handling-and-retry-strategy.md)
- Exponential backoff and jitter implementation
- Retry possibility determination
- Improved fault tolerance

### [ADR-0003: Concurrent Resource Retrieval](0003-concurrent-resource-retrieval.md)
- Goroutines and channels for concurrent processing
- Independent resource retrieval
- Efficient error handling

### [ADR-0004: Output Formatting Design](0004-output-formatting.md)
- Table and JSON format support
- Color-coded output
- Flexible display options

## Purpose of ADRs

- Ensure transparency of design decisions
- Knowledge transfer to future developers
- Track technical evolution of the project

## Reading ADRs

Each ADR follows this structure:
- **Status**: Current state (Proposed/Accepted/Deprecated)
- **Context**: Background necessitating the decision
- **Decisions**: Specific technical choices
- **Consequences**: Benefits, drawbacks, anticipated impacts
- **References**: Related resources and documentation

## Contributing

To add a new ADR:
1. Copy [0000-adr-template.md](0000-adr-template.md)
2. Create a new file with an incremented number
3. Provide clear and concise explanations
4. Include specific implementation details where possible

## License

This documentation follows the project's main license.
