# ADR-0004: Output Formatting Design

## Status

Accepted

## Context

Requirements for output display in SageMaker resource monitoring tool:
- Human-readable format
- Machine-parsable format
- Support for different use cases
- Clear visualization of resource information

Existing Problems:
- Single output format lacks flexibility
- Balance between readability and machine-readability
- Need to address different user requirements

## Decisions

1. Output Format Diversity
   - Table Format (Default)
     - Easy for human reading
     - Enhanced visibility through color coding
   - JSON Format (Optional)
     - Optimal for machine analysis
     - Facilitates integration with other tools and scripts

2. Table Output Characteristics
   - Color Coding
     - Green: Running resources
     - Yellow: Paused/Warning state
     - Red: Error/Stopped state
   - Concise Information Display
     - Resource type
     - Name (long names truncated)
     - Status
     - Instance type
     - Running time

3. JSON Output Characteristics
   - Structured detailed information
   - Preservation of all resource information
   - Compliance with RFC 8259 JSON format

4. Output Control
   - Switching via command-line flags
     - `--json` or `-j`
   - Default is table format

## Consequences

Benefits:
- High readability
- Flexible output options
- Support for different use cases
- Visual information transmission

Drawbacks:
- Increased implementation complexity
- Slight performance overhead

Anticipated Impacts:
- Improved user experience
- Support for diverse usage scenarios
- Easier script integration

## Specific Implementation Details

```go
type Printer struct {
    useJSON bool
    output  io.Writer
    isFirstResource bool
}

func (p *Printer) PrintResource(info ResourceInfo) {
    if p.useJSON {
        p.printJSONResource(info)
    } else {
        p.printTableResource(info)
    }
}

func (p *Printer) printTableResource(info ResourceInfo) {
    statusColor := map[string]func(a ...interface{}) string{
        "InService": color.New(color.FgGreen).SprintFunc(),
        "Running":   color.New(color.FgGreen).SprintFunc(),
        "Stopped":   color.New(color.FgYellow).SprintFunc(),
        "Failed":    color.New(color.FgRed).SprintFunc(),
    }
    
    // Color-coded table output
    // ...
}

func (p *Printer) printJSONResource(info ResourceInfo) {
    // JSON format output
    // ...
}
```

## Output Examples

### Table Format
```
Type            Name                Status      Instance       Running Time
Endpoint        ml-endpoint         InService   ml.t3.medium   72h 15m
Notebook        dev-notebook        Running     ml.t3.large    168h 30m
```

### JSON Format
```json
[
  {
    "resourceType": "Endpoint",
    "name": "ml-endpoint",
    "status": "InService",
    "instanceType": "ml.t3.medium",
    "runningTime": "72h 15m"
  },
  ...
]
```

## References

- [Go Formatting Library](https://github.com/fatih/color)
- [JSON Specification (RFC 8259)](https://tools.ietf.org/html/rfc8259)
- [CLI Output Design Guidelines](https://clig.dev/)
