package reporters

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/onsi/ginkgo/v2/types"
)

// GherkinSpec represents a test case in Given-When-Then format
type GherkinSpec struct {
	Feature     string            `json:"feature"`
	Given       []string          `json:"given"`
	When        string            `json:"when"`
	Then        string            `json:"then"`
	Status      string            `json:"status"`
	Error       string            `json:"error,omitempty"`
	Duration    time.Duration     `json:"duration"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// GherkinReporter outputs tests in Given-When-Then format
type GherkinReporter struct {
	specs       []GherkinSpec
	currentSpec *GherkinSpec
	suiteName   string
}

// NewGherkinReporter creates a new Gherkin-style reporter
func NewGherkinReporter() *GherkinReporter {
	return &GherkinReporter{
		specs: make([]GherkinSpec, 0),
	}
}

// SuiteWillBegin is called before the test suite starts
func (r *GherkinReporter) SuiteWillBegin(config types.SuiteConfig, summary types.SuiteSummary) error {
	r.suiteName = summary.SuiteDescription
	fmt.Printf("\n🧪 Test Suite: %s\n", r.suiteName)
	fmt.Println(strings.Repeat("=", 60))
	return nil
}

// WillRun is called before each spec runs
func (r *GherkinReporter) WillRun(spec types.SpecReport) error {
	r.currentSpec = r.parseToGherkin(spec)
	return nil
}

// DidRun is called after each spec completes
func (r *GherkinReporter) DidRun(spec types.SpecReport) error {
	if r.currentSpec == nil {
		return nil
	}

	// Update with results
	r.currentSpec.Status = spec.State.String()
	r.currentSpec.Duration = spec.RunTime

	if spec.State.Is(types.SpecStateFailed) || spec.State.Is(types.SpecStatePanicked) {
		r.currentSpec.Error = spec.Failure.Message
	}

	r.specs = append(r.specs, *r.currentSpec)
	r.printSpec(*r.currentSpec)
	r.currentSpec = nil

	return nil
}

// SuiteDidEnd is called after the test suite completes
func (r *GherkinReporter) SuiteDidEnd(summary types.SuiteSummary) error {
	fmt.Println(strings.Repeat("=", 60))
	
	// Print summary
	passed := 0
	failed := 0
	skipped := 0
	
	for _, spec := range r.specs {
		switch spec.Status {
		case "passed":
			passed++
		case "failed", "panicked":
			failed++
		case "skipped", "pending":
			skipped++
		}
	}

	fmt.Printf("\n📊 Test Results:\n")
	fmt.Printf("   ✅ Passed:  %d\n", passed)
	fmt.Printf("   ❌ Failed:  %d\n", failed)
	fmt.Printf("   ⏭️  Skipped: %d\n", skipped)
	fmt.Printf("   📈 Total:   %d\n", len(r.specs))
	fmt.Println()

	return nil
}

// SpecSuiteWillBegin implements deprecated interface
func (r *GherkinReporter) SpecSuiteWillBegin(config types.SuiteConfig, summary types.SuiteSummary) {}

// SpecSuiteDidEnd implements deprecated interface  
func (r *GherkinReporter) SpecSuiteDidEnd(summary types.SuiteSummary) {}

// SpecWillRun implements deprecated interface
func (r *GherkinReporter) SpecWillRun(spec types.SpecReport) {}

// SpecDidComplete implements deprecated interface
func (r *GherkinReporter) SpecDidComplete(spec types.SpecReport) {}

// FailureState is called when a spec fails
func (r *GherkinReporter) FailureState(failure types.Failure) {}

func (r *GherkinReporter) parseToGherkin(spec types.SpecReport) *GherkinSpec {
	containers := spec.ContainerHierarchyTexts
	
	gherkin := &GherkinSpec{
		Labels: make(map[string]string),
	}

	// Extract feature from suite or container
	if len(containers) > 0 {
		gherkin.Feature = containers[0]
	} else {
		gherkin.Feature = r.suiteName
	}

	// Parse containers into Given-When-Then
	if len(containers) > 1 {
		// All containers except last are "Given" context
		for i := 0; i < len(containers)-1; i++ {
			context := containers[i+1]
			// Clean up "Given" or "And" prefixes if present
			context = strings.TrimPrefix(context, "Given ")
			context = strings.TrimPrefix(context, "And ")
			gherkin.Given = append(gherkin.Given, context)
		}
	}

	// Last container is usually "When"
	if len(containers) > 0 {
		when := containers[len(containers)-1]
		when = strings.TrimPrefix(when, "When ")
		gherkin.When = when
	}

	// The It block is "Then"
	then := spec.LeafNodeText
	then = strings.TrimPrefix(then, "should ")
	gherkin.Then = then

	return gherkin
}

func (r *GherkinReporter) printSpec(spec GherkinSpec) {
	// Print with appropriate emoji based on status
	statusEmoji := "✅"
	if spec.Status == "failed" || spec.Status == "panicked" {
		statusEmoji = "❌"
	} else if spec.Status == "skipped" || spec.Status == "pending" {
		statusEmoji = "⏭️"
	}

	fmt.Printf("\n%s Scenario: %s\n", statusEmoji, spec.Then)
	
	for _, given := range spec.Given {
		fmt.Printf("   Given %s\n", given)
	}
	
	if spec.When != "" {
		fmt.Printf("   When %s\n", spec.When)
	}
	
	fmt.Printf("   Then %s", spec.Then)
	
	if spec.Status == "failed" && spec.Error != "" {
		// Print abbreviated error
		errorPreview := spec.Error
		if len(errorPreview) > 100 {
			errorPreview = errorPreview[:100] + "..."
		}
		fmt.Printf("\n   💥 Error: %s", errorPreview)
	}
	
	fmt.Printf(" (%v)\n", spec.Duration.Round(time.Millisecond))
}

// GenerateJSONOutput outputs all specs as JSON for external consumption
func (r *GherkinReporter) GenerateJSONOutput() ([]byte, error) {
	return json.MarshalIndent(r.specs, "", "  ")
}

// GetSpecs returns all captured specs
func (r *GherkinReporter) GetSpecs() []GherkinSpec {
	return r.specs
}

// PrintSummary prints a human-readable summary
func (r *GherkinReporter) PrintSummary() {
	passed := 0
	failed := 0
	
	for _, spec := range r.specs {
		switch spec.Status {
		case "passed":
			passed++
		case "failed", "panicked":
			failed++
		}
	}

	fmt.Printf("\n📋 Summary: %d passed, %d failed\n", passed, failed)
}
