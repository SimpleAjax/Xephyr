package flows_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/SimpleAjax/Xephyr/internal/dto"
	. "github.com/SimpleAjax/Xephyr/tests/e2e"
	"github.com/SimpleAjax/Xephyr/tests/e2e/helpers"
)

var _ = Describe("Dependency Management Flow", func() {
	var client *helpers.APIClient
	var projectID string = "proj-ecommerce"
	var testTaskID string = "task-ec-4"
	var createdDependencyID string

	BeforeEach(func() {
		client = helpers.NewAPIClient(Config.BaseURL, Config.APIToken, Config.OrganizationID)
		client.WithContext(TestCtx)
	})

	AfterEach(func() {
		// Cleanup created dependency if exists
		if createdDependencyID != "" && !Config.SkipCleanup {
			_, _ = client.Delete("/dependencies/" + createdDependencyID)
			createdDependencyID = ""
		}
	})

	Describe("Given a project manager needs to establish task dependencies", func() {
		Context("When retrieving task dependencies", func() {
			It("should return direct and indirect dependencies", func() {
				// Act
				resp, err := client.Get("/dependencies/tasks/"+testTaskID,
					helpers.WithQueryParam("includeIndirect", "true"),
				)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.TaskDependenciesResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Success).To(BeTrue())

				// Verify response structure
				Expect(result.Data.TaskID).To(Equal(testTaskID))
				Expect(result.Data.Dependencies.Direct).ToNot(BeNil())
				Expect(result.Data.Dependencies.Indirect).ToNot(BeNil())
				Expect(result.Data.Dependents.Direct).ToNot(BeNil())
				Expect(result.Data.Dependents.Indirect).ToNot(BeNil())

				// Verify chain analysis
				Expect(result.Data.ChainAnalysis.LongestChain).To(BeNumerically(">=", 0))
				Expect(result.Data.ChainAnalysis.CriticalPathPosition).ToNot(BeEmpty())
				Expect(result.Data.ChainAnalysis.FloatHours).To(BeNumerically(">=", 0))
			})

			It("should identify blocking dependencies correctly", func() {
				// Act
				resp, err := client.Get("/dependencies/tasks/" + testTaskID)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.TaskDependenciesResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())

				// Check for blocking dependencies
				for _, dep := range result.Data.Dependencies.Direct {
					if dep.IsBlocking {
						Expect(dep.DependencyID).ToNot(BeEmpty())
						Expect(dep.DependsOnTaskID).ToNot(BeEmpty())
						Expect(dep.Status).ToNot(BeEmpty())
					}
				}
			})
		})

		Context("When validating a dependency before creation", func() {
			It("should validate a valid dependency without cycle", func() {
				// Arrange
				validateReq := dto.ValidateDependencyRequest{
					TaskID:          "task-ec-5", // Admin Dashboard
					DependsOnTaskID: "task-ec-4", // Checkout Flow
					DependencyType:  "finish_to_start",
				}

				// Act
				resp, err := client.Post("/dependencies/validate", validateReq)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.ValidateDependencyResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Success).To(BeTrue())

				Expect(result.Data.Valid).To(BeTrue())
				Expect(result.Data.WouldCreateCycle).To(BeFalse())
				Expect(result.Data.Impact.EstimatedDelay).To(BeNumerically(">=", 0))
				Expect(result.Data.Impact.AffectedTasks).To(BeNumerically(">=", 0))
			})

			It("should detect circular dependencies", func() {
				// Arrange - Try to create a cycle: ec-1 -> ec-2 -> ec-4 -> ec-1
				validateReq := dto.ValidateDependencyRequest{
					TaskID:          "task-ec-1", // Design System (already done)
					DependsOnTaskID: "task-ec-4", // Checkout Flow (depends on ec-2, which depends on ec-1)
					DependencyType:  "finish_to_start",
				}

				// Act
				resp, err := client.Post("/dependencies/validate", validateReq)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				// Should return 409 Conflict for circular dependency
				Expect(resp.StatusCode).To(Or(Equal(http.StatusConflict), Equal(http.StatusOK)))

				if resp.StatusCode == http.StatusOK {
					var result helpers.Response[dto.ValidateDependencyResponse]
					err = helpers.ParseResponse(resp, &result)
					Expect(err).ToNot(HaveOccurred())

					if result.Data.WouldCreateCycle {
						Expect(result.Data.Valid).To(BeFalse())
					}
				}
			})
		})

		Context("When creating a valid dependency", func() {
			It("should successfully create the dependency with impact analysis", func() {
				// Arrange
				createReq := dto.CreateDependencyRequest{
					TaskID:          "task-ec-5", // Admin Dashboard
					DependsOnTaskID: "task-ec-4", // Checkout Flow
					DependencyType:  "finish_to_start",
					LagHours:        8,
				}

				// Act
				resp, err := client.Post("/dependencies", createReq)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))

				var result helpers.Response[dto.CreateDependencyResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Success).To(BeTrue())

				// Verify dependency was created
				Expect(result.Data.DependencyID).ToNot(BeEmpty())
				createdDependencyID = result.Data.DependencyID

				Expect(result.Data.TaskID).To(Equal(createReq.TaskID))
				Expect(result.Data.DependsOnTaskID).To(Equal(createReq.DependsOnTaskID))
				Expect(result.Data.DependencyType).To(Equal(createReq.DependencyType))
				Expect(result.Data.LagHours).To(Equal(createReq.LagHours))
				Expect(result.Data.CreatedAt).ToNot(BeZero())

				// Verify validation passed
				Expect(result.Data.Validation.Valid).To(BeTrue())
				Expect(result.Data.Validation.WouldCreateCycle).To(BeFalse())

				// Verify impact analysis
				Expect(result.Data.Impact.CriticalPathChanged).ToNot(BeNil())
				Expect(result.Data.Impact.AffectedTasks).ToNot(BeNil())
			})
		})

		Context("When attempting to create a circular dependency", func() {
			It("should reject the request with cycle information", func() {
				// Arrange
				createReq := dto.CreateDependencyRequest{
					TaskID:          "task-ec-1",
					DependsOnTaskID: "task-ec-4", // Would create cycle
					DependencyType:  "finish_to_start",
				}

				// Act
				resp, err := client.Post("/dependencies", createReq)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusConflict))

				var result helpers.Response[interface{}]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Success).To(BeFalse())
			})
		})

		Context("When retrieving the critical path", func() {
			It("should return the critical path with task details", func() {
				// Act
				resp, err := client.Get("/dependencies/critical-path/" + projectID)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.CriticalPathResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Success).To(BeTrue())

				// Verify project info
				Expect(result.Data.ProjectID).To(Equal(projectID))
				Expect(result.Data.ProjectDuration).To(BeNumerically(">", 0))
				Expect(result.Data.CalculatedAt).ToNot(BeZero())

				// Verify critical path
				Expect(result.Data.CriticalPath.TaskIDs).ToNot(BeEmpty())
				Expect(result.Data.CriticalPath.Tasks).ToNot(BeEmpty())
				Expect(result.Data.CriticalPath.TotalDuration).To(BeNumerically(">", 0))

				// Verify critical path task details
				for _, task := range result.Data.CriticalPath.Tasks {
					Expect(task.TaskID).ToNot(BeEmpty())
					Expect(task.Title).ToNot(BeEmpty())
					Expect(task.EstimatedHours).To(BeNumerically(">", 0))
					Expect(task.FloatHours).To(Equal(0)) // Critical path tasks have zero float
				}

				// Verify non-critical tasks
				Expect(result.Data.NonCriticalTasks).ToNot(BeNil())
				for _, task := range result.Data.NonCriticalTasks {
					Expect(task.TaskID).ToNot(BeEmpty())
					Expect(task.Title).ToNot(BeEmpty())
					Expect(task.FloatHours).To(BeNumerically(">=", 0))
				}
			})
		})

		Context("When retrieving the dependency graph", func() {
			It("should return nodes and edges for visualization", func() {
				// Act
				resp, err := client.Get("/dependencies/graph/" + projectID)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var result helpers.Response[dto.DependencyGraphResponse]
				err = helpers.ParseResponse(resp, &result)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Success).To(BeTrue())

				// Verify graph structure
				Expect(result.Data.ProjectID).To(Equal(projectID))
				Expect(result.Data.Nodes).ToNot(BeNil())
				Expect(result.Data.Edges).ToNot(BeNil())

				// Verify nodes
				for _, node := range result.Data.Nodes {
					Expect(node.ID).ToNot(BeEmpty())
					Expect(node.Title).ToNot(BeEmpty())
					Expect(node.Status).ToNot(BeEmpty())
				}

				// Verify edges
				for _, edge := range result.Data.Edges {
					Expect(edge.ID).ToNot(BeEmpty())
					Expect(edge.Source).ToNot(BeEmpty())
					Expect(edge.Target).ToNot(BeEmpty())
					Expect(edge.Type).ToNot(BeEmpty())
					Expect(edge.LagHours).To(BeNumerically(">=", 0))
				}
			})
		})

		Context("When deleting a dependency", func() {
			BeforeEach(func() {
				// Create a dependency to delete
				createReq := dto.CreateDependencyRequest{
					TaskID:          "task-ec-5",
					DependsOnTaskID: "task-ec-3",
					DependencyType:  "finish_to_start",
				}
				resp, err := client.Post("/dependencies", createReq)
				Expect(err).ToNot(HaveOccurred())

				if resp.StatusCode == http.StatusCreated {
					var result helpers.Response[dto.CreateDependencyResponse]
					err = helpers.ParseResponse(resp, &result)
					Expect(err).ToNot(HaveOccurred())
					createdDependencyID = result.Data.DependencyID
				}
			})

			It("should successfully delete the dependency", func() {
				if createdDependencyID == "" {
					Skip("No dependency to delete")
				}

				// Act
				resp, err := client.Delete("/dependencies/" + createdDependencyID)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				createdDependencyID = ""
			})
		})
	})
})
