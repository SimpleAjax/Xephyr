package flows_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/SimpleAjax/Xephyr/internal/dto"
	. "github.com/SimpleAjax/Xephyr/tests/e2e"
	"github.com/SimpleAjax/Xephyr/tests/e2e/helpers"
)

// CriticalPathsSmokeTest validates the most critical user journeys
// These tests ensure core functionality is always working
var _ = Describe("Critical Paths Smoke Test", func() {
	var client *helpers.APIClient
	var projectID string = "proj-ecommerce"
	var taskID string = "task-ec-4"

	BeforeEach(func() {
		client = helpers.NewAPIClient(Config.BaseURL, Config.APIToken, Config.OrganizationID)
		client.WithContext(TestCtx)
	})

	Describe("Smoke Test: Priority API", func() {
		It("should retrieve task priority", func() {
			resp, err := client.Get("/priorities/tasks/" + taskID)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var result helpers.Response[dto.TaskPriorityResponse]
			err = helpers.ParseResponse(resp, &result)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.Success).To(BeTrue())
			Expect(result.Data.TaskID).To(Equal(taskID))
			Expect(result.Data.UniversalPriorityScore).To(BeNumerically(">=", 0))
			Expect(result.Data.UniversalPriorityScore).To(BeNumerically("<=", 100))
		})

		It("should get bulk task priorities", func() {
			req := dto.BulkPriorityRequest{
				TaskIds:          []string{"task-ec-1", "task-ec-2", "task-ec-3"},
				IncludeBreakdown: true,
			}

			resp, err := client.Post("/priorities/tasks/bulk", req)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var result helpers.Response[dto.BulkPriorityResponse]
			err = helpers.ParseResponse(resp, &result)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.Success).To(BeTrue())
			Expect(len(result.Data.Priorities)).To(Equal(3))
		})

		It("should get project task ranking", func() {
			resp, err := client.Get("/priorities/projects/"+projectID+"/ranking",
				helpers.WithQueryParam("limit", "10"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var result helpers.Response[dto.ProjectTaskRankingResponse]
			err = helpers.ParseResponse(resp, &result)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.Success).To(BeTrue())
			Expect(result.Data.ProjectID).To(Equal(projectID))
			Expect(len(result.Data.Rankings)).To(BeNumerically(">", 0))
		})
	})

	Describe("Smoke Test: Health API", func() {
		It("should retrieve portfolio health", func() {
			resp, err := client.Get("/health/portfolio")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var result helpers.Response[dto.PortfolioHealthResponse]
			err = helpers.ParseResponse(resp, &result)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.Success).To(BeTrue())
			Expect(result.Data.PortfolioHealthScore).To(BeNumerically(">=", 0))
		})

		It("should retrieve project health", func() {
			resp, err := client.Get("/health/projects/" + projectID)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var result helpers.Response[dto.ProjectHealthResponse]
			err = helpers.ParseResponse(resp, &result)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.Success).To(BeTrue())
			Expect(result.Data.ProjectID).To(Equal(projectID))
		})
	})

	Describe("Smoke Test: Assignment API", func() {
		It("should get assignment suggestions", func() {
			resp, err := client.Get("/assignments/suggestions",
				helpers.WithQueryParam("taskId", taskID),
				helpers.WithQueryParam("limit", "3"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var result helpers.Response[dto.AssignmentSuggestionsResponse]
			err = helpers.ParseResponse(resp, &result)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.Success).To(BeTrue())
			Expect(result.Data.TaskID).To(Equal(taskID))
		})

		It("should check assignment compatibility", func() {
			resp, err := client.Get("/assignments/compatibility",
				helpers.WithQueryParam("taskId", taskID),
				helpers.WithQueryParam("personId", "user-mike"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var result helpers.Response[dto.AssignmentCompatibilityResponse]
			err = helpers.ParseResponse(resp, &result)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.Success).To(BeTrue())
		})
	})

	Describe("Smoke Test: Dependency API", func() {
		It("should get task dependencies", func() {
			resp, err := client.Get("/dependencies/tasks/" + taskID)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var result helpers.Response[dto.TaskDependenciesResponse]
			err = helpers.ParseResponse(resp, &result)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.Success).To(BeTrue())
			Expect(result.Data.TaskID).To(Equal(taskID))
		})

		It("should get critical path", func() {
			resp, err := client.Get("/dependencies/critical-path/" + projectID)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var result helpers.Response[dto.CriticalPathResponse]
			err = helpers.ParseResponse(resp, &result)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.Success).To(BeTrue())
			Expect(result.Data.ProjectID).To(Equal(projectID))
		})
	})

	Describe("Smoke Test: Nudge API", func() {
		It("should list nudges", func() {
			resp, err := client.Get("/nudges",
				helpers.WithQueryParam("limit", "10"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var result helpers.Response[dto.NudgeListResponse]
			err = helpers.ParseResponse(resp, &result)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.Success).To(BeTrue())
		})

		It("should get nudge statistics", func() {
			resp, err := client.Get("/nudges/stats",
				helpers.WithQueryParam("period", "30d"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var result helpers.Response[dto.NudgeStatsResponse]
			err = helpers.ParseResponse(resp, &result)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.Success).To(BeTrue())
		})
	})

	Describe("Smoke Test: Scenario API", func() {
		It("should list scenarios", func() {
			resp, err := client.Get("/scenarios",
				helpers.WithQueryParam("limit", "10"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var result helpers.Response[dto.ScenarioListResponse]
			err = helpers.ParseResponse(resp, &result)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.Success).To(BeTrue())
		})
	})

	Describe("Smoke Test: Workload API", func() {
		It("should get team workload", func() {
			resp, err := client.Get("/workload/team")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("should get individual workload", func() {
			resp, err := client.Get("/workload/people/user-mike")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})
	})

	Describe("Smoke Test: Complete Critical User Journey", func() {
		It("should execute the PM decision-making workflow", func() {
			// Step 1: Check portfolio health
			resp, err := client.Get("/health/portfolio")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var healthResult helpers.Response[dto.PortfolioHealthResponse]
			err = helpers.ParseResponse(resp, &healthResult)
			Expect(err).ToNot(HaveOccurred())

			// Step 2: Get priority ranking for at-risk project
			if healthResult.Data.Summary.AtRisk > 0 {
				resp, err = client.Get("/priorities/projects/"+projectID+"/ranking",
					helpers.WithQueryParam("limit", "5"),
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			}

			// Step 3: Check for any nudges
			resp, err = client.Get("/nudges",
				helpers.WithQueryParam("status", "unread"),
				helpers.WithQueryParam("limit", "5"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var nudgeResult helpers.Response[dto.NudgeListResponse]
			err = helpers.ParseResponse(resp, &nudgeResult)
			Expect(err).ToNot(HaveOccurred())

			// Step 4: Get assignment suggestions for top priority task
			resp, err = client.Get("/assignments/suggestions",
				helpers.WithQueryParam("taskId", taskID),
				helpers.WithQueryParam("limit", "3"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var assignResult helpers.Response[dto.AssignmentSuggestionsResponse]
			err = helpers.ParseResponse(resp, &assignResult)
			Expect(err).ToNot(HaveOccurred())

			// Verify we have candidates
			Expect(len(assignResult.Data.Candidates)).To(BeNumerically(">", 0))
		})
	})
})
