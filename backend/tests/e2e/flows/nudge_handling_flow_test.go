package flows_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/SimpleAjax/Xephyr/internal/dto"
	. "github.com/SimpleAjax/Xephyr/tests/e2e"
	"github.com/SimpleAjax/Xephyr/tests/e2e/helpers"
)

var _ = Describe("Nudge Handling Flow", func() {
	var client *helpers.APIClient
	var testNudgeID string

	BeforeEach(func() {
		client = helpers.NewAPIClient(Config.BaseURL, Config.APIToken, Config.OrganizationID)
		client.WithContext(TestCtx)
	})

	Describe("Given nudges are generated based on project conditions", func() {
		Context("When listing nudges with filters", func() {
			It("should return nudges with summary statistics", func() {
				// Act
				resp, err := client.Get("/nudges",
					helpers.WithQueryParam("limit", "20"),
				)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.NudgeListResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Success).To(BeTrue())

				// Verify nudges list
				Expect(result.Data.Nudges).ToNot(BeNil())

				// Verify summary
				Expect(result.Data.Summary.Total).To(BeNumerically(">=", 0))
				Expect(result.Data.Summary.Unread).To(BeNumerically(">=", 0))
				Expect(result.Data.Summary.BySeverity).ToNot(BeNil())
				Expect(result.Data.Summary.ByType).ToNot(BeNil())
			})

			It("should filter nudges by status", func() {
				// Act
				resp, err := client.Get("/nudges",
					helpers.WithQueryParam("status", "unread"),
					helpers.WithQueryParam("limit", "20"),
				)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.NudgeListResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())

				// Verify all returned nudges have unread status
				for _, nudge := range result.Data.Nudges {
					Expect(nudge.Status).To(Equal("unread"))
				}
			})

			It("should filter nudges by severity", func() {
				// Act
				resp, err := client.Get("/nudges",
					helpers.WithQueryParam("severity", "high"),
					helpers.WithQueryParam("limit", "20"),
				)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.NudgeListResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())

				// Verify all returned nudges have high severity
				for _, nudge := range result.Data.Nudges {
					Expect(nudge.Severity).To(Equal("high"))
				}
			})

			It("should filter nudges by project", func() {
				// Act
				resp, err := client.Get("/nudges",
					helpers.WithQueryParam("projectId", "proj-mobile"),
					helpers.WithQueryParam("limit", "20"),
				)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.NudgeListResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())

				// Verify all returned nudges are for the specified project
				for _, nudge := range result.Data.Nudges {
					Expect(nudge.RelatedEntities.ProjectID).ToNot(BeNil())
					Expect(*nudge.RelatedEntities.ProjectID).To(Equal("proj-mobile"))
				}
			})
		})

		Context("When retrieving a single nudge", func() {
			BeforeEach(func() {
				// First, get a list of nudges to find a valid ID
				resp, err := client.Get("/nudges",
					helpers.WithQueryParam("limit", "1"),
				)
				Expect(err).ToNot(HaveOccurred())

				var result helpers.Response[dto.NudgeListResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())

				if len(result.Data.Nudges) > 0 {
					testNudgeID = result.Data.Nudges[0].ID
				}
			})

			It("should return complete nudge details with AI explanation", func() {
				if testNudgeID == "" {
					Skip("No nudge available for testing")
				}

				// Act
				resp, err := client.Get("/nudges/" + testNudgeID)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.NudgeDetailResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Success).To(BeTrue())

				// Verify nudge details
				Expect(result.Data.ID).To(Equal(testNudgeID))
				Expect(result.Data.Type).ToNot(BeEmpty())
				Expect(result.Data.Severity).ToNot(BeEmpty())
				Expect(result.Data.Status).ToNot(BeEmpty())
				Expect(result.Data.Title).ToNot(BeEmpty())
				Expect(result.Data.Description).ToNot(BeEmpty())
				Expect(result.Data.AIExplanation).ToNot(BeEmpty())
				Expect(result.Data.CriticalityScore).To(BeNumerically(">=", 0))
				Expect(result.Data.CriticalityScore).To(BeNumerically("<=", 100))

				// Verify suggested action
				Expect(result.Data.SuggestedAction.Type).ToNot(BeEmpty())
				Expect(result.Data.SuggestedAction.Description).ToNot(BeEmpty())

				// Verify related entities
				Expect(result.Data.RelatedEntities).ToNot(BeNil())

				// Verify metrics
				Expect(result.Data.Metrics.AllocationPercentage).To(BeNumerically(">=", 0))

				// Verify history
				Expect(result.Data.History).ToNot(BeNil())
			})
		})

		Context("When taking action on a nudge", func() {
			BeforeEach(func() {
				// Get a list of nudges to find a valid ID
				resp, err := client.Get("/nudges",
					helpers.WithQueryParam("limit", "1"),
				)
				Expect(err).ToNot(HaveOccurred())

				var result helpers.Response[dto.NudgeListResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())

				if len(result.Data.Nudges) > 0 {
					testNudgeID = result.Data.Nudges[0].ID
				}
			})

			It("should accept the suggested action", func() {
				if testNudgeID == "" {
					Skip("No nudge available for testing")
				}

				// Arrange
				actionReq := dto.NudgeActionRequest{
					ActionType: "accept_suggestion",
					Parameters: map[string]interface{}{
						"reassignTo": "user-rachel",
					},
				}

				// Act
				resp, err := client.Post("/nudges/"+testNudgeID+"/actions", actionReq)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.NudgeActionResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Success).To(BeTrue())

				// Verify action result
				Expect(result.Data.NudgeID).To(Equal(testNudgeID))
				Expect(result.Data.ActionTaken).To(Equal("accept_suggestion"))
				Expect(result.Data.NudgeStatus).To(Equal("acted"))
				Expect(result.Data.CompletedAt).ToNot(BeZero())
			})

			It("should dismiss a nudge without action", func() {
				if testNudgeID == "" {
					Skip("No nudge available for testing")
				}

				// Arrange
				actionReq := dto.NudgeActionRequest{
					ActionType: "dismiss",
				}

				// Act
				resp, err := client.Post("/nudges/"+testNudgeID+"/actions", actionReq)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.NudgeActionResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Success).To(BeTrue())

				Expect(result.Data.NudgeID).To(Equal(testNudgeID))
				Expect(result.Data.ActionTaken).To(Equal("dismiss"))
				Expect(result.Data.NudgeStatus).To(Equal("dismissed"))
			})
		})

		Context("When updating nudge status", func() {
			BeforeEach(func() {
				// Get a list of nudges to find a valid ID
				resp, err := client.Get("/nudges",
					helpers.WithQueryParam("limit", "1"),
				)
				Expect(err).ToNot(HaveOccurred())

				var result helpers.Response[dto.NudgeListResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())

				if len(result.Data.Nudges) > 0 {
					testNudgeID = result.Data.Nudges[0].ID
				}
			})

			It("should mark nudge as read", func() {
				if testNudgeID == "" {
					Skip("No nudge available for testing")
				}

				// Arrange
				statusReq := dto.UpdateNudgeStatusRequest{
					Status: "read",
				}

				// Act
				resp, err := client.Patch("/nudges/"+testNudgeID+"/status", statusReq)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})
		})

		Context("When retrieving nudge statistics", func() {
			It("should return comprehensive nudge statistics", func() {
				// Act
				resp, err := client.Get("/nudges/stats",
					helpers.WithQueryParam("period", "30d"),
				)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.NudgeStatsResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Success).To(BeTrue())

				// Verify statistics
				Expect(result.Data.Period).To(Equal("30d"))
				Expect(result.Data.Generated).To(BeNumerically(">=", 0))
				Expect(result.Data.Acted).To(BeNumerically(">=", 0))
				Expect(result.Data.Dismissed).To(BeNumerically(">=", 0))
				Expect(result.Data.Expired).To(BeNumerically(">=", 0))
				Expect(result.Data.ActionRate).To(BeNumerically(">=", 0))
				Expect(result.Data.ActionRate).To(BeNumerically("<=", 1))
				Expect(result.Data.AvgTimeToAction).ToNot(BeEmpty())

				// Verify by type statistics
				Expect(result.Data.ByType).ToNot(BeNil())
				for nudgeType, stats := range result.Data.ByType {
					Expect(nudgeType).ToNot(BeEmpty())
					Expect(stats.Generated).To(BeNumerically(">=", 0))
					Expect(stats.Acted).To(BeNumerically(">=", 0))
				}
			})
		})

		Context("When triggering manual nudge generation", func() {
			It("should queue nudge generation for organization scope", func() {
				// Arrange
				generateReq := dto.GenerateNudgesRequest{
					Scope: "organization",
					Types: []string{"overload", "delay_risk"},
					Async: true,
				}

				// Act
				resp, err := client.Post("/nudges/generate", generateReq)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusAccepted))
			})

			It("should generate nudges for specific project", func() {
				// Arrange
				projectID := "proj-mobile"
				generateReq := dto.GenerateNudgesRequest{
					Scope:     "project",
					ProjectID: &projectID,
					Async:     false,
				}

				// Act
				resp, err := client.Post("/nudges/generate", generateReq)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Or(Equal(http.StatusOK), Equal(http.StatusAccepted)))
			})
		})
	})
})
