package flows_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/SimpleAjax/Xephyr/internal/dto"
	. "github.com/SimpleAjax/Xephyr/tests/e2e"
	"github.com/SimpleAjax/Xephyr/tests/e2e/helpers"
)

var _ = Describe("Task Assignment Flow", func() {
	var client *helpers.APIClient

	BeforeEach(func() {
		client = helpers.NewAPIClient(Config.BaseURL, Config.APIToken, Config.OrganizationID)
		client.WithContext(TestCtx)
	})

	Describe("Given a project manager needs to assign a high-priority task", func() {
		var taskID string = "task-ec-4"

		Context("When the PM requests assignment suggestions for the task", func() {
			It("should return a ranked list of candidate assignees with compatibility scores", func() {
				// Arrange
				reqPath := "/assignments/suggestions"

				// Act
				resp, err := client.Get(reqPath, 
					helpers.WithQueryParam("taskId", taskID),
					helpers.WithQueryParam("limit", "3"),
				)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.AssignmentSuggestionsResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Success).To(BeTrue())

				// Verify response structure
				Expect(result.Data.TaskID).To(Equal(taskID))
				Expect(result.Data.TaskTitle).ToNot(BeEmpty())
				Expect(result.Data.RequiredSkills).ToNot(BeEmpty())
				Expect(result.Data.Candidates).To(HaveLen(3))

				// Verify top candidate structure
				topCandidate := result.Data.Candidates[0]
				Expect(topCandidate.Rank).To(Equal(1))
				Expect(topCandidate.Score).To(BeNumerically(">", 0))
				Expect(topCandidate.Score).To(BeNumerically("<=", 100))
				Expect(topCandidate.Person.ID).ToNot(BeEmpty())
				Expect(topCandidate.Person.Name).ToNot(BeEmpty())

				// Verify score breakdown
				Expect(topCandidate.Breakdown.SkillMatch).To(BeNumerically(">=", 0))
				Expect(topCandidate.Breakdown.Availability).To(BeNumerically(">=", 0))
				Expect(topCandidate.Breakdown.WorkloadBalance).To(BeNumerically(">=", 0))
				Expect(topCandidate.Breakdown.PastPerformance).To(BeNumerically(">=", 0))

				// Verify AI explanation is provided
				Expect(topCandidate.AIExplanation).ToNot(BeEmpty())
			})

			It("should include skill match details for each candidate", func() {
				// Act
				resp, err := client.Get("/assignments/suggestions",
					helpers.WithQueryParam("taskId", taskID),
				)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.AssignmentSuggestionsResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())

				// Verify skill match details
				for _, candidate := range result.Data.Candidates {
					Expect(candidate.SkillMatchDetails).ToNot(BeEmpty())
					for _, skill := range candidate.SkillMatchDetails {
						Expect(skill.SkillID).ToNot(BeEmpty())
						Expect(skill.HasSkill).To(BeTrue())
						Expect(skill.Proficiency).To(BeNumerically(">=", 1))
						Expect(skill.Proficiency).To(BeNumerically("<=", 5))
					}
				}
			})

			It("should include workload and availability context for each candidate", func() {
				// Act
				resp, err := client.Get("/assignments/suggestions",
					helpers.WithQueryParam("taskId", taskID),
				)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.AssignmentSuggestionsResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())

				for _, candidate := range result.Data.Candidates {
					Expect(candidate.ContextSwitchAnalysis.ActiveProjects).To(BeNumerically(">=", 0))
					Expect(candidate.ContextSwitchAnalysis.CurrentWorkload).To(BeNumerically(">=", 0))
					Expect(candidate.ContextSwitchAnalysis.RiskLevel).ToNot(BeEmpty())
				}
			})

			It("should flag candidates with overallocation warnings", func() {
				// Act
				resp, err := client.Get("/assignments/suggestions",
					helpers.WithQueryParam("taskId", taskID),
				)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.AssignmentSuggestionsResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())

				// At least one candidate should have warnings (based on test data)
				var foundWarning bool
				for _, candidate := range result.Data.Candidates {
					if len(candidate.Warnings) > 0 {
						foundWarning = true
						for _, warning := range candidate.Warnings {
							Expect(warning).ToNot(BeEmpty())
						}
					}
				}
				// Based on fixtures, Alex is overallocated
				Expect(foundWarning).To(BeTrue(), "Expected at least one candidate with overallocation warning")
			})
		})

		Context("When the PM assigns the task to the top candidate", func() {
			It("should successfully assign the task and update workload", func() {
				// Arrange
				assignReq := dto.AssignTaskRequest{
					PersonID:       "user-mike",
					Note:           "Assigned based on AI suggestion - 92% match",
					SkipSuggestion: false,
				}

				// Act
				resp, err := client.Post("/assignments/tasks/"+taskID+"/assign", assignReq)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.AssignTaskResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Success).To(BeTrue())

				// Verify assignment details
				Expect(result.Data.TaskID).To(Equal(taskID))
				Expect(result.Data.AssignedTo.PersonID).To(Equal("user-mike"))
				Expect(result.Data.AssignedTo.Name).ToNot(BeEmpty())
				Expect(result.Data.Assignment.AssignedAt).ToNot(BeZero())
				Expect(result.Data.Assignment.AssignedBy).ToNot(BeEmpty())

				// Verify impact
				Expect(result.Data.Impact.WorkloadUpdated).To(BeTrue())
				Expect(result.Data.Impact.NewAllocation).To(BeNumerically(">", 0))
				Expect(result.Data.Impact.NotificationsSent).ToNot(BeEmpty())
			})

			It("should generate nudges if the assignment creates overallocation", func() {
				// Arrange - Assign to a user who would be overallocated
				assignReq := dto.AssignTaskRequest{
					PersonID: "user-emma", // Emma is already at 125% allocation
					Note:     "Testing overload detection",
				}

				// Act
				resp, err := client.Post("/assignments/tasks/"+taskID+"/assign", assignReq)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.AssignTaskResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Success).To(BeTrue())

				// Verify overload nudge was generated
				Expect(result.Data.Impact.NudgesGenerated).ToNot(BeEmpty())
			})
		})

		Context("When the PM uses auto-assign with best match strategy", func() {
			It("should automatically assign to the best matching candidate", func() {
				// Arrange
				autoAssignReq := dto.AutoAssignTaskRequest{
					Strategy: "best_match",
					Constraints: dto.AutoAssignConstraints{
						MaxAllocation:       100,
						RequiredProficiency: 3,
					},
				}

				// Act
				resp, err := client.Post("/assignments/tasks/"+taskID+"/auto-assign", autoAssignReq)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.AssignTaskResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Success).To(BeTrue())

				// Verify auto-assignment
				Expect(result.Data.TaskID).To(Equal(taskID))
				Expect(result.Data.AssignedTo.PersonID).ToNot(BeEmpty())
				Expect(result.Data.AssignedTo.Name).ToNot(BeEmpty())
			})
		})

		Context("When checking person-task compatibility", func() {
			It("should return detailed compatibility analysis", func() {
				// Act
				resp, err := client.Get("/assignments/compatibility",
					helpers.WithQueryParam("taskId", taskID),
					helpers.WithQueryParam("personId", "user-mike"),
				)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.AssignmentCompatibilityResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Success).To(BeTrue())

				Expect(result.Data.TaskID).To(Equal(taskID))
				Expect(result.Data.PersonID).To(Equal("user-mike"))
				Expect(result.Data.Score).To(BeNumerically(">=", 0))
				Expect(result.Data.Score).To(BeNumerically("<=", 100))
			})
		})
	})

	Describe("Given a project manager needs to redistribute workload", func() {
		Context("When performing bulk reassignments", func() {
			It("should successfully reassign multiple tasks", func() {
				// Arrange
				bulkReq := dto.BulkReassignRequest{
					Reassignments: []dto.ReassignmentItem{
						{
							TaskID:       "task-web-3",
							FromPersonID: "user-emma",
							ToPersonID:   "user-rachel",
						},
					},
					Reason: "Redistribute workload - Emma overallocated",
				}

				// Act
				resp, err := client.Post("/assignments/bulk-reassign", bulkReq)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.BulkReassignResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Success).To(BeTrue())

				Expect(result.Data.Processed).To(Equal(1))
				Expect(result.Data.Succeeded).To(Equal(1))
			})
		})
	})
})


