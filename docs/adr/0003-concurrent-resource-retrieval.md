# ADR-0003: Concurrent Resource Retrieval Strategy

## Status

Accepted

## Context

Challenges in SageMaker Resource Monitoring:
- Retrieving multiple resource types (Endpoints, Notebooks, Studio Apps)
- Potential latency from individual API calls
- Need for efficient resource information collection

Existing Problems:
- Sequential resource retrieval is time-consuming
- A single resource retrieval error could block the entire process
- Optimization of response time is required

## Decisions

1. Concurrent Processing Using Go's Goroutines and Channels
   - Independent goroutine for each resource type
     - Endpoints
     - Notebook Instances
     - Studio Applications

2. Concurrent Processing Implementation Strategy
   - Use of `sync.WaitGroup` for goroutine synchronization
   - Dedicated channels for each resource type
   - Separated error and result handling

3. Error Handling Design
   - Individual resource retrieval errors do not stop the entire process
   - Preserve the first error while continuing resource retrieval
   - Special handling for retryable errors

4. Resource Information Integration
   - Unified type using common `ResourceInfo` struct
   - Support for flexible output formatting

## Consequences

Benefits:
- Significant reduction in response time
- Independent resource retrieval processes
- High parallel processing efficiency
- Robust error handling

Drawbacks:
- Slightly more complex implementation
- Minimal increase in memory usage

Anticipated Impacts:
- Accelerated resource information retrieval
- Improved system responsiveness
- Enhanced scalability

## Specific Implementation Details

```go
func runMonitor() error {
    // Channel creation
    endpointsChan := make(chan ResourceResult, 1)
    notebooksChan := make(chan ResourceResult, 1)
    appsChan := make(chan ResourceResult, 1)

    // Synchronization with WaitGroup
    var wg sync.WaitGroup
    wg.Add(3)

    // Concurrent Endpoints retrieval
    go func() {
        defer wg.Done()
        endpoints, err := client.ListEndpoints(ctx)
        endpointsChan <- ResourceResult{Resources: endpoints, Error: err}
    }()

    // Concurrent Notebooks retrieval
    go func() {
        defer wg.Done()
        notebooks, err := client.ListNotebooks(ctx)
        notebooksChan <- ResourceResult{Resources: notebooks, Error: err}
    }()

    // Concurrent Studio Apps retrieval
    go func() {
        defer wg.Done()
        apps, err := client.ListStudioApps(ctx)
        appsChan <- ResourceResult{Resources: apps, Error: err}
    }()

    // Channel closure
    go func() {
        wg.Wait()
        close(endpointsChan)
        close(notebooksChan)
        close(appsChan)
    }()

    // Resource information processing
    // ...
}
```

## Performance Considerations

- Theoretical maximum response time: Time of the slowest resource retrieval
- Expected speedup of approximately 3x compared to sequential processing
- Consideration of network latency and API limitations

## References

- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)
- [Go Concurrency Patterns Presentation](https://talks.golang.org/2012/concurrency.slide)
