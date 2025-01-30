package display

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
)

// ResourceInfo represents the information to be displayed for each resource
type ResourceInfo struct {
	ResourceType  string `json:"resourceType"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	InstanceType string `json:"instanceType"`
	RunningTime  string `json:"runningTime"`
}

// Printer handles the formatting and display of resource information
type Printer struct {
	useJSON bool
	output  io.Writer
	isFirstResource bool
}

// NewPrinter creates a new printer instance
func NewPrinter(useJSON bool) *Printer {
	return &Printer{
		useJSON: useJSON,
		output:  os.Stdout,
		isFirstResource: true,
	}
}

// PrintHeader prepares the output for resource listing
func (p *Printer) PrintHeader() {
	if p.useJSON {
		fmt.Fprint(p.output, "[\n")
	} else {
		headerFmt := color.New(color.FgGreen, color.Bold).SprintfFunc()
		fmt.Fprintf(p.output, "%s\n", headerFmt(
			"%-15s %-30s %-12s %-15s %-15s",
			"Type", "Name", "Status", "Instance", "Running Time",
		))
		fmt.Fprintln(p.output, strings.Repeat("-", 120))
	}
}

// PrintResource outputs a single resource
func (p *Printer) PrintResource(info ResourceInfo) {
	if p.useJSON {
		p.printJSONResource(info)
	} else {
		p.printTableResource(info)
	}
}

// printJSONResource outputs a single resource in JSON format
func (p *Printer) printJSONResource(info ResourceInfo) {
	if !p.isFirstResource {
		fmt.Fprint(p.output, ",\n")
	}
	
	data, err := json.Marshal(info)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		return
	}
	fmt.Fprint(p.output, "  ", string(data))
	
	p.isFirstResource = false
}

// printTableResource outputs a single resource in table format
func (p *Printer) printTableResource(info ResourceInfo) {
	statusColor := map[string]func(a ...interface{}) string{
		"InService":  color.New(color.FgGreen).SprintFunc(),
		"Running":    color.New(color.FgGreen).SprintFunc(),
		"Stopped":    color.New(color.FgYellow).SprintFunc(),
		"Failed":     color.New(color.FgRed).SprintFunc(),
		"Deleting":   color.New(color.FgRed).SprintFunc(),
	}

	status := info.Status
	if colorFunc, ok := statusColor[status]; ok {
		status = colorFunc(status)
	}

	fmt.Fprintf(p.output, "%-15s %-30s %-12s %-15s %-15s\n",
		info.ResourceType,
		truncateString(info.Name, 29),
		status,
		info.InstanceType,
		info.RunningTime,
	)
}

// PrintFooter finalizes the output
func (p *Printer) PrintFooter() {
	if p.useJSON {
		fmt.Fprint(p.output, "\n]\n")
	} else {
		fmt.Fprintln(p.output, strings.Repeat("-", 120))
	}
}

// PrintNoResources handles the case when no resources are found
func (p *Printer) PrintNoResources(region string) {
	if p.useJSON {
		// Create a JSON object with metadata about no resources
		noResourcesJSON := struct {
			Resources []interface{} `json:"resources"`
			Metadata  struct {
				Region  string `json:"region"`
				Message string `json:"message"`
			} `json:"metadata"`
		}{
			Resources: []interface{}{},
			Metadata: struct {
				Region  string `json:"region"`
				Message string `json:"message"`
			}{
				Region:  region,
				Message: "No resources found",
			},
		}

		jsonData, err := json.MarshalIndent(noResourcesJSON, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
			return
		}
		fmt.Fprintln(p.output, string(jsonData))
	} else {
		// Use color for the no resources message in table format
		noResourceMsg := color.New(color.FgYellow).SprintfFunc()
		fmt.Fprintf(p.output, "%s\n", noResourceMsg("No SageMaker resources found in region %s", region))
	}
}

// Helper function to truncate long strings
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
