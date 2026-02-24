package flows_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/SimpleAjax/Xephyr/internal/dto"
	. "github.com/SimpleAjax/Xephyr/tests/e2e"
	"github.com/SimpleAjax/Xephyr/tests/e2e/helpers"
)

var _ = Describe("Health Monitoring Flow", func() {
	var client *helpers.APIClient
	var projectID string = "proj-ecommerce"

	BeforeEach(func() {
		client = helpers.NewAPIClient(Config.BaseURL, Config.APIToken, Config.OrganizationID)
		client.WithContext(TestCtx)
	})

	Describe("Given a project manager monitors project health", func() {
		Context("When retrieving portfolio health overview", func() {
			It("should return complete portfolio health with summary", func() {
				// Act
				resp, err := client.Get("/health/portfolio")

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.PortfolioHealthResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Success).To(BeTrue())

				// Verify portfolio health score
				Expect(result.Data.PortfolioHealthScore).To(BeNumerically(">=", 0))
				Expect(result.Data.PortfolioHealthScore).To(BeNumerically("<=", 100))
				Expect(result.Data.Status).ToNot(BeEmpty())
				Expect(result.Data.CalculatedAt).ToNot(BeZero())

				// Verify summary
				Expect(result.Data.Summary.TotalProjects).To(BeNumerically(">=", 0))
				Expect(result.Data.Summary.Healthy).To(BeNumerically(">=", 0))
				Expect(result.Data.Summary.Caution).To(BeNumerically(">=", 0))
				Expect(result.Data.Summary.AtRisk).To(BeNumerically(">=", 0))
				Expect(result.Data.Summary.Critical).To(BeNumerically(">=", 0))

				// Verify that counts add up
				totalCount := result.Data.Summary.Healthy +
					result.Data.Summary.Caution +
					result.Data.Summary.AtRisk +
					result.Data.Summary.Critical
				Expect(totalCount).To(Equal(result.Data.Summary.TotalProjects))

				// Verify projects list
				Expect(result.Data.Projects).ToNot(BeEmpty())
				for _, project := range result.Data.Projects {
					Expect(project.ProjectID).ToNot(BeEmpty())
					Expect(project.Name).ToNot(BeEmpty())
					Expect(project.HealthScore).To(BeNumerically(">=", 0))
					Expect(project.HealthScore).To(BeNumerically("<=", 100))
					Expect(project.Status).ToNot(BeEmpty())
					Expect(project.Trend).ToNot(BeEmpty())
				}
			})

			It("should identify at-risk projects in the portfolio", func() {
				// Act
				resp, err := client.Get("/health/portfolio")

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.PortfolioHealthResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())

				// Look for at-risk projects based on test data
				var atRiskFound bool
				for _, project := range result.Data.Projects {
					if project.Status == "at_risk" || project.HealthScore < 50 {
						atRiskFound = true
						Expect(project.Name).ToNot(BeEmpty())
						Expect(project.HealthScore).To(BeNumerically("<", 50))
					}
				}

				// Based on fixtures, proj-mobile should be at risk
				Expect(atRiskFound).To(BeTrue(), "Expected at least one at-risk project")
			})
		})

		Context("When retrieving detailed project health", func() {
			It("should return complete health breakdown with all components", func() {
				// Act
				resp, err := client.Get("/health/projects/"+projectID,
					helpers.WithQueryParam("includeBreakdown", "true"),
				)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.ProjectHealthResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Success).To(BeTrue())

				// Verify project info
				Expect(result.Data.ProjectID).To(Equal(projectID))
				Expect(result.Data.ProjectName).ToNot(BeEmpty())
				Expect(result.Data.HealthScore).To(BeNumerically(">=", 0))
				Expect(result.Data.HealthScore).To(BeNumerically("<=", 100))
				Expect(result.Data.Status).ToNot(BeEmpty())
				Expect(result.Data.CalculatedAt).ToNot(BeZero())

				// Verify health breakdown
				Expect(result.Data.Breakdown.ScheduleHealth).To(BeNumerically(">=", 0))
				Expect(result.Data.Breakdown.ScheduleHealth).To(BeNumerically("<=", 100))
				Expect(result.Data.Breakdown.CompletionHealth).To(BeNumerically(">=", 0))
				Expect(result.Data.Breakdown.CompletionHealth).To(BeNumerically("<=", 100))
				Expect(result.Data.Breakdown.DependencyHealth).To(BeNumerically(">=", 0))
				Expect(result.Data.Breakdown.DependencyHealth).To(BeNumerically("<=", 100))
				Expect(result.Data.Breakdown.ResourceHealth).To(BeNumerically(">=", 0))
				Expect(result.Data.Breakdown.ResourceHealth).To(BeNumerically("<=", 100))
				Expect(result.Data.Breakdown.CriticalPathHealth).To(BeNumerically(">=", 0))
				Expect(result.Data.Breakdown.CriticalPathHealth).To(BeNumerically("<=", 100))

				// Verify trend
				Expect(result.Data.Trend.Direction).ToNot(BeEmpty())
				Expect(result.Data.Trend.LastWeekScore).To(BeNumerically(">=", 0))
				Expect(result.Data.Trend.LastWeekScore).To(BeNumerically("<=", 100))
			})

			It("should include detailed health metrics", func() {
				// Act
				resp, err := client.Get("/health/projects/" + projectID)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.ProjectHealthResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())

				// Verify schedule details
				Expect(result.Data.Details.Schedule.ExpectedProgress).To(BeNumerically(">=", 0))
				Expect(result.Data.Details.Schedule.ExpectedProgress).To(BeNumerically("<=", 100))
				Expect(result.Data.Details.Schedule.ActualProgress).To(BeNumerically(">=", 0))
				Expect(result.Data.Details.Schedule.ActualProgress).To(BeNumerically("<=", 100))
				Expect(result.Data.Details.Schedule.DaysUntilDeadline).To(BeNumerically(">=", 0))

				// Verify completion details
				Expect(result.Data.Details.Completion.TotalTasks).To(BeNumerically(">=", 0))
				Expect(result.Data.Details.Completion.Completed).To(BeNumerically(">=", 0))
				Expect(result.Data.Details.Completion.InProgress).To(BeNumerically(">=", 0))
				Expect(result.Data.Details.Completion.CompletionRate).To(BeNumerically(">=", 0))
				Expect(result.Data.Details.Completion.CompletionRate).To(BeNumerically("<=", 100))

				// Verify dependency details
				Expect(result.Data.Details.Dependencies.Total).To(BeNumerically(">=", 0))
				Expect(result.Data.Details.Dependencies.Blocked).To(BeNumerically(">=", 0))
				Expect(result.Data.Details.Dependencies.AtRisk).To(BeNumerically(">=", 0))

				// Verify resource details
				Expect(result.Data.Details.Resources.TeamSize).To(BeNumerically(">=", 0))
				Expect(result.Data.Details.Resources.AvgAllocation).To(BeNumerically(">=", 0))
				Expect(result.Data.Details.Resources.Overallocated).To(BeNumerically(">=", 0))
				Expect(result.Data.Details.Resources.Underutilized).To(BeNumerically(">=", 0))
			})
		})

		Context("When retrieving health trends over time", func() {
			It("should return historical health data with trend analysis", func() {
				// Act
				resp, err := client.Get("/health/trends",
					helpers.WithQueryParam("projectId", projectID),
					helpers.WithQueryParam("days", "30"),
				)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.HealthTrendsResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Success).To(BeTrue())

				// Verify response structure
				Expect(result.Data.ProjectID).To(Equal(projectID))
				Expect(result.Data.TimeRange).To(Equal("30d"))
				Expect(result.Data.Datapoints).ToNot(BeEmpty())

				// Verify datapoints
				for _, point := range result.Data.Datapoints {
					Expect(point.Date).ToNot(BeEmpty())
					Expect(point.HealthScore).To(BeNumerically(">=", 0))
					Expect(point.HealthScore).To(BeNumerically("<=", 100))
					Expect(point.ScheduleHealth).To(BeNumerically(">=", 0))
					Expect(point.ScheduleHealth).To(BeNumerically("<=", 100))
					Expect(point.CompletionHealth).To(BeNumerically(">=", 0))
					Expect(point.CompletionHealth).To(BeNumerically("<=", 100))
				}

				// Verify trend analysis
				Expect(result.Data.Trend.Direction).ToNot(BeEmpty())
				Expect(result.Data.Trend.Slope).ToNot(BeZero())
			})

			It("should provide health predictions", func() {
				// Act
				resp, err := client.Get("/health/trends",
					helpers.WithQueryParam("projectId", projectID),
					helpers.WithQueryParam("days", "30"),
				)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.HealthTrendsResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())

				// Verify prediction exists
				Expect(result.Data.Trend.Prediction.DaysUntilCritical).To(BeNumerically(">=", 0))
				Expect(result.Data.Trend.Prediction.Confidence).To(BeNumerically(">=", 0))
				Expect(result.Data.Trend.Prediction.Confidence).To(BeNumerically("<=", 1))
			})
		})

		Context("When comparing health across multiple projects", func() {
			It("should return bulk project health data", func() {
				// Act
				resp, err := client.Get("/health/projects")

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				// Verify we get a list of project health summaries
				var result helpers.Response[[]dto.ProjectHealthSummary]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Success).To(BeTrue())
				Expect(result.Data).ToNot(BeEmpty())

				// Verify each project has required fields
				for _, project := range result.Data {
					Expect(project.ProjectID).ToNot(BeEmpty())
					Expect(project.Name).ToNot(BeEmpty())
					Expect(project.HealthScore).To(BeNumerically(">=", 0))
					Expect(project.Status).ToNot(BeEmpty())
				}
			})
		})
	})
})
