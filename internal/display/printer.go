package display

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// ResourceInfo represents the information to be displayed for each resource
type ResourceInfo struct {
	ResourceType  string    `json:"resourceType"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	InstanceType string    `json:"instanceType"`
	RunningTime  string    `json:"runningTime"`
	HourlyCost   float64   `json:"hourlyCost"`
	CurrentCost  float64   `json:"currentCost"`
	ProjectedCost float64  `json:"projectedMonthlyCost"`
	StorageCost  float64   `json:"storageCost,omitempty"`
	TotalCost    float64   `json:"totalCost"`
}

// Printer handles the formatting and display of resource information
type Printer struct {
	useJSON bool
	resources []ResourceInfo
}

// NewPrinter creates a new printer instance
func NewPrinter(useJSON bool) *Printer {
	return &Printer{
		useJSON: useJSON,
	}
}

// AddResource adds a resource to be printed
func (p *Printer) AddResource(info ResourceInfo) {
	p.resources = append(p.resources, info)
}

// Print outputs all resources in the specified format
func (p *Printer) Print() {
	if p.useJSON {
		p.printJSON()
	} else {
		p.printTable()
	}
}

// printJSON outputs the resources in JSON format
func (p *Printer) printJSON() {
	data, err := json.MarshalIndent(p.resources, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}
	fmt.Println(string(data))
}

// printTable outputs the resources in a formatted table
func (p *Printer) printTable() {
	if len(p.resources) == 0 {
		color.Yellow("No active SageMaker resources found")
		return
	}

	// Print header
	headerFmt := color.New(color.FgGreen, color.Bold).SprintfFunc()
	fmt.Printf("%s\n", headerFmt(
		"%-15s %-30s %-12s %-15s %-15s %-12s %-12s %-15s",
		"Type", "Name", "Status", "Instance", "Running Time", "Hourly($)", "Current($)", "Projected($)",
	))
	fmt.Println(strings.Repeat("-", 120))

	// Print resources
	var totalCurrentCost, totalProjectedCost float64
	statusColor := map[string]func(a ...interface{}) string{
		"InService":  color.New(color.FgGreen).SprintFunc(),
		"Running":    color.New(color.FgGreen).SprintFunc(),
		"Stopped":    color.New(color.FgYellow).SprintFunc(),
		"Failed":     color.New(color.FgRed).SprintFunc(),
		"Deleting":   color.New(color.FgRed).SprintFunc(),
	}

	for _, r := range p.resources {
		status := r.Status
		if colorFunc, ok := statusColor[status]; ok {
			status = colorFunc(status)
		}

		fmt.Printf("%-15s %-30s %-12s %-15s %-15s $%-11.2f $%-11.2f $%-14.2f\n",
			r.ResourceType,
			truncateString(r.Name, 29),
			status,
			r.InstanceType,
			r.RunningTime,
			r.HourlyCost,
			r.CurrentCost,
			r.ProjectedCost,
		)

		totalCurrentCost += r.CurrentCost
		totalProjectedCost += r.ProjectedCost
	}

	// Print summary
	fmt.Println(strings.Repeat("-", 120))
	summaryFmt := color.New(color.FgCyan, color.Bold).SprintfFunc()
	fmt.Printf("%s\n", summaryFmt("Total Current Cost: $%.2f    Projected Monthly Cost: $%.2f\n",
		totalCurrentCost, totalProjectedCost))

	// Print warnings for high-cost resources
	p.printWarnings()
}

// printWarnings outputs warnings for resources with high costs
func (p *Printer) printWarnings() {
	warningFmt := color.New(color.FgYellow).SprintfFunc()
	
	for _, r := range p.resources {
		if r.ProjectedCost > 1000 {
			fmt.Printf("%s\n", warningFmt("WARNING: %s '%s' has a high projected monthly cost: $%.2f",
				r.ResourceType, r.Name, r.ProjectedCost))
		}
		
		if r.CurrentCost > 100 {
			fmt.Printf("%s\n", warningFmt("WARNING: %s '%s' has accumulated a significant cost: $%.2f",
				r.ResourceType, r.Name, r.CurrentCost))
		}
	}
}

// Helper function to truncate long strings
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
