package flows_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/SimpleAjax/Xephyr/internal/dto"
	. "github.com/SimpleAjax/Xephyr/tests/e2e"
	"github.com/SimpleAjax/Xephyr/tests/e2e/helpers"
)

var _ = Describe("Scenario Simulation Flow", func() {
	var client *helpers.APIClient
	var createdScenarioID string

	BeforeEach(func() {
		client = helpers.NewAPIClient(Config.BaseURL, Config.APIToken, Config.OrganizationID)
		client.WithContext(TestCtx)
	})

	AfterEach(func() {
		// Cleanup created scenario if exists
		if createdScenarioID != "" && !Config.SkipCleanup {
			_, _ = client.Delete("/scenarios/" + createdScenarioID)
			createdScenarioID = ""
		}
	})

	Describe("Given a project manager wants to simulate an employee leave scenario", func() {
		Context("When creating a new scenario", func() {
			It("should successfully create the scenario with draft status", func() {
				// Arrange
				createReq := dto.CreateScenarioRequest{
					Title:       "Emma takes 1-week vacation",
					Description: "Simulate impact of Emma taking vacation next week",
					ChangeType:  "employee_leave",
					ProposedChanges: dto.ProposedChanges{
						PersonID:         strPtr("user-emma"),
						LeaveStartDate:   strPtr("2026-02-24"),
						LeaveEndDate:     strPtr("2026-02-28"),
						CoverageStrategy: strPtr("reassign"),
					},
				}

				// Act
				resp, err := client.Post("/scenarios", createReq)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))

				var result helpers.Response[dto.ScenarioResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Success).To(BeTrue())

				// Verify scenario was created
				Expect(result.Data.ScenarioID).ToNot(BeEmpty())
				createdScenarioID = result.Data.ScenarioID

				Expect(result.Data.Title).To(Equal(createReq.Title))
				Expect(result.Data.ChangeType).To(Equal("employee_leave"))
				Expect(result.Data.Status).To(Equal("draft"))
				Expect(result.Data.SimulationStatus).To(Equal("pending"))

				// Verify proposed changes
				Expect(result.Data.ProposedChanges.PersonID).To(Equal(createReq.ProposedChanges.PersonID))
				Expect(result.Data.ProposedChanges.LeaveStartDate).To(Equal(createReq.ProposedChanges.LeaveStartDate))
				Expect(result.Data.ProposedChanges.LeaveEndDate).To(Equal(createReq.ProposedChanges.LeaveEndDate))
				Expect(result.Data.CreatedAt).ToNot(BeZero())
			})
		})

		Context("When listing existing scenarios", func() {
			BeforeEach(func() {
				// Create a test scenario first
				createReq := dto.CreateScenarioRequest{
					Title:      "Test Scenario for Listing",
					ChangeType: "employee_leave",
					ProposedChanges: dto.ProposedChanges{
						PersonID: strPtr("user-emma"),
					},
				}
				resp, err := client.Post("/scenarios", createReq)
				Expect(err).ToNot(HaveOccurred())

				var result helpers.Response[dto.ScenarioResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				createdScenarioID = result.Data.ScenarioID
			})

			It("should return scenarios with pagination", func() {
				// Act
				resp, err := client.Get("/scenarios",
					helpers.WithQueryParam("limit", "10"),
					helpers.WithQueryParam("offset", "0"),
				)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.ScenarioListResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Success).To(BeTrue())

				// Verify response structure
				Expect(result.Data.Scenarios).ToNot(BeNil())
				Expect(result.Data.Total).To(BeNumerically(">=", 1))
			})
		})

		Context("When running a full simulation", func() {
			BeforeEach(func() {
				// Create a test scenario
				createReq := dto.CreateScenarioRequest{
					Title:       "Emma Vacation Simulation",
					Description: "Full simulation test",
					ChangeType:  "employee_leave",
					ProposedChanges: dto.ProposedChanges{
						PersonID:         strPtr("user-emma"),
						LeaveStartDate:   strPtr("2026-02-24"),
						LeaveEndDate:     strPtr("2026-02-28"),
						CoverageStrategy: strPtr("reassign"),
					},
				}
				resp, err := client.Post("/scenarios", createReq)
				Expect(err).ToNot(HaveOccurred())

				var result helpers.Response[dto.ScenarioResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				createdScenarioID = result.Data.ScenarioID
			})

			It("should complete simulation with impact analysis", func() {
				// Arrange
				simulateReq := dto.SimulateScenarioRequest{
					Depth:                  "full",
					IncludeRecommendations: true,
				}

				// Act
				resp, err := client.Post("/scenarios/"+createdScenarioID+"/simulate", simulateReq)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.SimulateScenarioResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Success).To(BeTrue())

				// Verify simulation status
				Expect(result.Data.ScenarioID).To(Equal(createdScenarioID))
				Expect(result.Data.SimulationStatus).To(Equal("completed"))

				// Verify impact analysis
				Expect(result.Data.ImpactAnalysis.AffectedProjects).ToNot(BeEmpty())
				Expect(result.Data.ImpactAnalysis.AffectedTasks).ToNot(BeEmpty())

				// Verify affected project structure
				for _, project := range result.Data.ImpactAnalysis.AffectedProjects {
					Expect(project.ProjectID).ToNot(BeEmpty())
					Expect(project.Name).ToNot(BeEmpty())
					Expect(project.Impact).ToNot(BeEmpty())
					Expect(project.AffectedTasks).ToNot(BeNil())
				}

				// Verify timeline comparison
				Expect(result.Data.ImpactAnalysis.TimelineComparison.TotalDelayDays).To(BeNumerically(">=", 0))

				// Verify cost analysis
				Expect(result.Data.ImpactAnalysis.CostAnalysis.TotalCost).To(BeNumerically(">=", 0))
				Expect(result.Data.ImpactAnalysis.CostAnalysis.Confidence).To(BeNumerically(">=", 0))
				Expect(result.Data.ImpactAnalysis.CostAnalysis.Confidence).To(BeNumerically("<=", 1))

				// Verify AI recommendations
				Expect(result.Data.AIRecommendations).ToNot(BeEmpty())
				for _, rec := range result.Data.AIRecommendations {
					Expect(rec.Priority).To(BeNumerically(">=", 1))
					Expect(rec.Action).ToNot(BeEmpty())
					Expect(rec.Reasoning).ToNot(BeEmpty())
					Expect(rec.EstimatedImpact).ToNot(BeEmpty())
				}

				Expect(result.Data.CalculatedAt).ToNot(BeZero())
				Expect(result.Data.SimulationDuration).ToNot(BeEmpty())
			})

			It("should identify affected tasks with delay information", func() {
				// Arrange
				simulateReq := dto.SimulateScenarioRequest{
					Depth:                  "full",
					IncludeRecommendations: true,
				}

				// Act
				resp, err := client.Post("/scenarios/"+createdScenarioID+"/simulate", simulateReq)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.SimulateScenarioResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())

				// Verify affected tasks
				for _, task := range result.Data.ImpactAnalysis.AffectedTasks {
					Expect(task.TaskID).ToNot(BeEmpty())
					Expect(task.Title).ToNot(BeEmpty())
					Expect(task.DelayDays).To(BeNumerically(">=", 0))
					Expect(task.Reason).ToNot(BeEmpty())

					// Verify suggested reassignment if present
					if task.SuggestedReassignment != nil {
						Expect(task.SuggestedReassignment.ToPersonID).ToNot(BeEmpty())
						Expect(task.SuggestedReassignment.Compatibility).To(BeNumerically(">=", 0))
						Expect(task.SuggestedReassignment.Compatibility).To(BeNumerically("<=", 100))
					}
				}
			})

			It("should provide resource impact analysis", func() {
				// Arrange
				simulateReq := dto.SimulateScenarioRequest{
					Depth:                  "full",
					IncludeRecommendations: true,
				}

				// Act
				resp, err := client.Post("/scenarios/"+createdScenarioID+"/simulate", simulateReq)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.SimulateScenarioResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())

				// Verify resource impacts
				Expect(result.Data.ImpactAnalysis.ResourceImpacts).ToNot(BeEmpty())
				for _, impact := range result.Data.ImpactAnalysis.ResourceImpacts {
					Expect(impact.PersonID).ToNot(BeEmpty())
					Expect(impact.CurrentAllocation).To(BeNumerically(">=", 0))
					Expect(impact.NewAllocation).To(BeNumerically(">=", 0))
					Expect(impact.Risk).ToNot(BeEmpty())
				}
			})
		})

		Context("When applying a simulated scenario", func() {
			BeforeEach(func() {
				// Create and simulate a scenario
				createReq := dto.CreateScenarioRequest{
					Title:      "Apply Test Scenario",
					ChangeType: "employee_leave",
					ProposedChanges: dto.ProposedChanges{
						PersonID:         strPtr("user-emma"),
						LeaveStartDate:   strPtr("2026-02-24"),
						LeaveEndDate:     strPtr("2026-02-28"),
						CoverageStrategy: strPtr("reassign"),
					},
				}
				resp, err := client.Post("/scenarios", createReq)
				Expect(err).ToNot(HaveOccurred())

				var createResult helpers.Response[dto.ScenarioResponse]
				err = helpers.ParseResponse(resp, &createResult)
				Expect(err).ToNot(HaveOccurred())
				createdScenarioID = createResult.Data.ScenarioID

				// Run simulation first
				simulateReq := dto.SimulateScenarioRequest{
					Depth:                  "full",
					IncludeRecommendations: true,
				}
				_, err = client.Post("/scenarios/"+createdScenarioID+"/simulate", simulateReq)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should successfully apply the scenario with selected recommendations", func() {
				// Arrange
				applyReq := dto.ApplyScenarioRequest{
					ApplyRecommendations:    true,
					SelectedRecommendations: []int{1},
					NotifyStakeholders:      true,
				}

				// Act
				resp, err := client.Post("/scenarios/"+createdScenarioID+"/apply", applyReq)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.ApplyScenarioResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Success).To(BeTrue())

				// Verify scenario status
				Expect(result.Data.ScenarioID).To(Equal(createdScenarioID))
				Expect(result.Data.Status).To(Equal("applied"))
				Expect(result.Data.AppliedAt).ToNot(BeZero())

				// Verify changes were applied
				Expect(result.Data.Changes.TasksReassigned).ToNot(BeNil())
				Expect(result.Data.Changes.DatesAdjusted).ToNot(BeNil())
				Expect(result.Data.Changes.NotificationsSent).To(BeNumerically(">=", 0))

				// Verify follow-up actions
				Expect(result.Data.FollowUp.NudgesCreated).ToNot(BeNil())
			})
		})

		Context("When retrieving scenario details", func() {
			BeforeEach(func() {
				createReq := dto.CreateScenarioRequest{
					Title:       "Detail Test Scenario",
					Description: "Testing detail view",
					ChangeType:  "employee_leave",
					ProposedChanges: dto.ProposedChanges{
						PersonID: strPtr("user-emma"),
					},
				}
				resp, err := client.Post("/scenarios", createReq)
				Expect(err).ToNot(HaveOccurred())

				var result helpers.Response[dto.ScenarioResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				createdScenarioID = result.Data.ScenarioID
			})

			It("should return complete scenario details with history", func() {
				// Act
				resp, err := client.Get("/scenarios/" + createdScenarioID)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.ScenarioDetailResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Success).To(BeTrue())

				// Verify basic info
				Expect(result.Data.ScenarioID).To(Equal(createdScenarioID))
				Expect(result.Data.Title).ToNot(BeEmpty())
				Expect(result.Data.Description).ToNot(BeEmpty())
				Expect(result.Data.Status).ToNot(BeEmpty())

				// Verify history is present
				Expect(result.Data.History).ToNot(BeEmpty())
				Expect(result.Data.History[0].Action).To(Equal("created"))
				Expect(result.Data.History[0].Timestamp).ToNot(BeZero())
			})
		})

		Context("When rejecting a scenario", func() {
			BeforeEach(func() {
				createReq := dto.CreateScenarioRequest{
					Title:      "Reject Test Scenario",
					ChangeType: "employee_leave",
					ProposedChanges: dto.ProposedChanges{
						PersonID: strPtr("user-emma"),
					},
				}
				resp, err := client.Post("/scenarios", createReq)
				Expect(err).ToNot(HaveOccurred())

				var result helpers.Response[dto.ScenarioResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				createdScenarioID = result.Data.ScenarioID
			})

			It("should successfully reject the scenario", func() {
				// Act
				resp, err := client.Post("/scenarios/"+createdScenarioID+"/reject", nil)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.RejectScenarioResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Success).To(BeTrue())

				Expect(result.Data.ScenarioID).To(Equal(createdScenarioID))
				Expect(result.Data.Status).To(Equal("rejected"))
				Expect(result.Data.RejectedAt).ToNot(BeZero())
			})
		})

		Context("When modifying a scenario", func() {
			BeforeEach(func() {
				createReq := dto.CreateScenarioRequest{
					Title:      "Modify Test Scenario",
					ChangeType: "employee_leave",
					ProposedChanges: dto.ProposedChanges{
						PersonID: strPtr("user-emma"),
					},
				}
				resp, err := client.Post("/scenarios", createReq)
				Expect(err).ToNot(HaveOccurred())

				var result helpers.Response[dto.ScenarioResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				createdScenarioID = result.Data.ScenarioID
			})

			It("should successfully update the scenario details", func() {
				// Arrange
				newTitle := "Updated Scenario Title"
				modifyReq := dto.ModifyScenarioRequest{
					Title: &newTitle,
				}

				// Act
				resp, err := client.Patch("/scenarios/"+createdScenarioID+"/modify", modifyReq)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.ScenarioResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Success).To(BeTrue())

				Expect(result.Data.Title).To(Equal(newTitle))
			})
		})
	})
})

// Helper function
func strPtr(s string) *string {
	return &s
}
