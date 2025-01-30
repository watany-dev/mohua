package display

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrinterNoResources(t *testing.T) {
	tests := []struct {
		name     string
		useJSON  bool
		region   string
		expected string
	}{
		{
			name:    "Table format no resources",
			useJSON: false,
			region:  "ap-northeast-1",
			expected: "No SageMaker resources found in region ap-northeast-1",
		},
		{
			name:    "JSON format no resources",
			useJSON: true,
			region:  "ap-northeast-1",
			expected: `{
  "resources": [],
  "metadata": {
    "region": "ap-northeast-1",
    "message": "No resources found"
  }
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			printer := &Printer{
				useJSON: tt.useJSON,
				output:  &buf,
			}

			printer.PrintNoResources(tt.region)

			output := strings.TrimSpace(buf.String())
			if tt.useJSON {
				// Verify JSON structure
				var result map[string]interface{}
				err := json.Unmarshal([]byte(output), &result)
				assert.NoError(t, err)
				
				resources, ok := result["resources"].([]interface{})
				assert.True(t, ok)
				assert.Empty(t, resources)

				metadata, ok := result["metadata"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, tt.region, metadata["region"])
				assert.Equal(t, "No resources found", metadata["message"])
			} else {
				assert.Contains(t, output, tt.expected)
			}
		})
	}
}

func TestPrinterResourceOutput(t *testing.T) {
	tests := []struct {
		name     string
		useJSON  bool
		resource ResourceInfo
		expected string
	}{
		{
			name:    "Table format single resource",
			useJSON: false,
			resource: ResourceInfo{
				ResourceType:  "Endpoint",
				Name:         "test-endpoint",
				Status:       "InService",
				InstanceType: "ml.t3.medium",
				RunningTime:  "1h",
			},
			expected: `Type            Name                           Status       Instance        Running Time   
------------------------------------------------------------------------------------------------------------------------
Endpoint        test-endpoint                  InService    ml.t3.medium    1h             
------------------------------------------------------------------------------------------------------------------------`,
		},
		{
			name:    "JSON format single resource",
			useJSON: true,
			resource: ResourceInfo{
				ResourceType:  "Notebook",
				Name:         "test-notebook",
				Status:       "InService",
				InstanceType: "ml.t3.medium",
				RunningTime:  "2h",
			},
			expected: `[
  {"resourceType":"Notebook","name":"test-notebook","status":"InService","instanceType":"ml.t3.medium","runningTime":"2h"}
]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			printer := &Printer{
				useJSON: tt.useJSON,
				output:  &buf,
				isFirstResource: true,
			}

			printer.PrintHeader()
			printer.PrintResource(tt.resource)
			printer.PrintFooter()

			output := strings.TrimSpace(buf.String())
			assert.Equal(t, tt.expected, output)
		})
	}
}
