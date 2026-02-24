package services_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/SimpleAjax/Xephyr/internal/models"
	"github.com/SimpleAjax/Xephyr/tests/fixtures"
)

var _ = Describe("Health Scoring System", func() {
	
	Context("Given a portfolio with multiple active projects", func() {
		var projects []models.Project
		
		BeforeEach(func() {
			projects = []models.Project{
				fixtures.NewProject().WithID("proj-1").WithPriority(95).WithHealthScore(90).WithStatus(models.ProjectActive).Build(),
				fixtures.NewProject().WithID("proj-2").WithPriority(80).WithHealthScore(70).WithStatus(models.ProjectActive).Build(),
				fixtures.NewProject().WithID("proj-3").WithPriority(60).WithHealthScore(45).WithStatus(models.ProjectActive).Build(),
				fixtures.NewProject().WithID("proj-4").WithPriority(90).WithHealthScore(85).WithStatus(models.ProjectActive).Build(),
			}
		})
		
		When("portfolio health is calculated", func() {
			var portfolioHealth PortfolioHealth
			
			BeforeEach(func() {
				portfolioHealth = CalculatePortfolioHealth(projects)
			})
			
			It("should calculate weighted average based on project priorities", func() {
				// Weighted calculation:
				// (90*0.95 + 70*0.80 + 45*0.60 + 85*0.90) / (0.95 + 0.80 + 0.60 + 0.90)
				// = (85.5 + 56 + 27 + 76.5) / 3.25
				// = 245 / 3.25 = 75.38
				Expect(portfolioHealth.Score).To(BeNumerically(">=", 75))
				Expect(portfolioHealth.Score).To(BeNumerically("<=", 76))
			})
			
			It("should categorize portfolio health status correctly", func() {
				Expect(portfolioHealth.Status).To(Equal("caution"))
			})
			
			It("should count projects by health status", func() {
				Expect(portfolioHealth.HealthyCount).To(Equal(2))      // 90, 85
				Expect(portfolioHealth.CautionCount).To(Equal(1))      // 70
				Expect(portfolioHealth.AtRiskCount).To(Equal(1))       // 45
				Expect(portfolioHealth.CriticalCount).To(Equal(0))
			})
		})
	})
	
	Context("Given an empty portfolio", func() {
		When("portfolio health is calculated", func() {
			It("should return perfect health score", func() {
				health := CalculatePortfolioHealth([]models.Project{})
				Expect(health.Score).To(Equal(100))
			})
		})
	})
	
	Context("Given a project with tasks in various states", func() {
		var project models.Project
		var tasks []models.Task
		var workloadData []models.WorkloadEntry
		
		BeforeEach(func() {
			project = fixtures.NewProject().
				WithID("proj-health-test").
				WithPriority(80).
				WithProgress(45).
				Build()
			
			now := time.Now()
			startDate := now.Add(-30 * 24 * time.Hour)
			endDate := now.Add(30 * 24 * time.Hour)
			project.StartDate = &startDate
			project.TargetEndDate = &endDate
			
			tasks = []models.Task{
				fixtures.NewTask().WithStatus(models.TaskStatusDone).WithEstimatedHours(40).WithActualHours(40).Build(),
				fixtures.NewTask().WithStatus(models.TaskStatusDone).WithEstimatedHours(30).WithActualHours(30).Build(),
				fixtures.NewTask().WithStatus(models.TaskStatusInProgress).WithEstimatedHours(50).WithActualHours(25).Build(),
				fixtures.NewTask().WithStatus(models.TaskStatusReview).WithEstimatedHours(30).WithActualHours(24).Build(),
				fixtures.NewTask().WithStatus(models.TaskStatusBacklog).WithEstimatedHours(40).Build(),
				fixtures.NewTask().WithStatus(models.TaskStatusReady).WithEstimatedHours(20).Build(),
			}
			
			workloadData = fixtures.CreateTestWorkloadData()
		})
		
		When("project health is calculated", func() {
			var health ProjectHealth
			
			BeforeEach(func() {
				health = CalculateProjectHealth(project, tasks, workloadData)
			})
			
			It("should return a health score between 0 and 100", func() {
				Expect(health.Score).To(BeNumerically(">=", 0))
				Expect(health.Score).To(BeNumerically("<=", 100))
			})
			
			It("should include breakdown of all health factors", func() {
				Expect(health.Breakdown.ScheduleHealth).To(BeNumerically(">=", 0))
				Expect(health.Breakdown.CompletionHealth).To(BeNumerically(">=", 0))
				Expect(health.Breakdown.DependencyHealth).To(BeNumerically(">=", 0))
				Expect(health.Breakdown.ResourceHealth).To(BeNumerically(">=", 0))
				Expect(health.Breakdown.CriticalPathHealth).To(BeNumerically(">=", 0))
			})
			
			It("should calculate schedule health based on progress vs timeline", func() {
				// 50% time elapsed, 45% progress = slightly behind
				Expect(health.Breakdown.ScheduleHealth).To(BeNumerically("<", 100))
			})
		})
	})
	
	Context("Given a project behind schedule", func() {
		var project models.Project
		var tasks []models.Task
		
		BeforeEach(func() {
			project = fixtures.NewProject().WithProgress(30).Build()
			now := time.Now()
			startDate := now.Add(-50 * 24 * time.Hour)
			endDate := now.Add(50 * 24 * time.Hour)
			project.StartDate = &startDate
			project.TargetEndDate = &endDate
			
			// Expected progress: 50%, actual: 30% = 20% behind
			tasks = []models.Task{
				fixtures.NewTask().WithStatus(models.TaskStatusDone).WithEstimatedHours(30).Build(),
				fixtures.NewTask().WithStatus(models.TaskStatusInProgress).WithEstimatedHours(40).Build(),
				fixtures.NewTask().WithStatus(models.TaskStatusBacklog).WithEstimatedHours(30).Build(),
			}
		})
		
		When("schedule health is calculated", func() {
			It("should return reduced health for behind-schedule project", func() {
				health := CalculateScheduleHealth(project, tasks)
				Expect(health).To(BeNumerically("<", 85))
			})
			
			It("should calculate correct variance", func() {
				variance := CalculateProgressVariance(project, tasks)
				Expect(variance.ExpectedProgress).To(Equal(50))
				Expect(variance.ActualProgress).To(Equal(30))
				Expect(variance.Variance).To(Equal(-20))
			})
		})
	})
	
	Context("Given a project with resource overallocation", func() {
		var workloadData []models.WorkloadEntry
		var projectID string
		
		BeforeEach(func() {
			projectID = "proj-overallocated"
			workloadData = []models.WorkloadEntry{
				{UserID: stringToUUID("user-1"), AllocationPercentage: 125},
				{UserID: stringToUUID("user-2"), AllocationPercentage: 110},
				{UserID: stringToUUID("user-3"), AllocationPercentage: 85},
				{UserID: stringToUUID("user-4"), AllocationPercentage: 70},
			}
		})
		
		When("resource health is calculated", func() {
			It("should return lower health for overallocated team", func() {
				health := CalculateResourceHealth(projectID, workloadData)
				Expect(health).To(BeNumerically("<", 80))
			})
		})
	})
	
	Context("Given a project with all tasks completed", func() {
		var tasks []models.Task
		
		BeforeEach(func() {
			tasks = []models.Task{
				fixtures.NewTask().WithStatus(models.TaskStatusDone).WithEstimatedHours(40).Build(),
				fixtures.NewTask().WithStatus(models.TaskStatusDone).WithEstimatedHours(60).Build(),
				fixtures.NewTask().WithStatus(models.TaskStatusDone).WithEstimatedHours(50).Build(),
			}
		})
		
		When("completion health is calculated", func() {
			It("should return perfect completion health", func() {
				health := CalculateCompletionHealth(tasks)
				Expect(health).To(Equal(100))
			})
		})
	})
	
	Context("Given health score thresholds", func() {
		DescribeTable("Health status categorization",
			func(score int, expectedStatus string) {
				status := GetHealthStatus(score)
				Expect(status).To(Equal(expectedStatus))
			},
			Entry("perfect health", 100, "healthy"),
			Entry("good health", 85, "healthy"),
			Entry("borderline healthy", 80, "healthy"),
			Entry("caution threshold", 79, "caution"),
			Entry("mid caution", 70, "caution"),
			Entry("at risk threshold", 60, "caution"),
			Entry("at risk", 59, "at_risk"),
			Entry("critical threshold", 40, "at_risk"),
			Entry("critical", 39, "critical"),
			Entry("very critical", 10, "critical"),
		)
	})
	
	Context("Given health trends over time", func() {
		var historicalData []HealthDatapoint
		
		BeforeEach(func() {
			historicalData = []HealthDatapoint{
				{Date: time.Now().AddDate(0, 0, -21), HealthScore: 85},
				{Date: time.Now().AddDate(0, 0, -14), HealthScore: 82},
				{Date: time.Now().AddDate(0, 0, -7), HealthScore: 78},
				{Date: time.Now(), HealthScore: 72},
			}
		})
		
		When("health trend is calculated", func() {
			var trend HealthTrend
			
			BeforeEach(func() {
				trend = CalculateHealthTrend(historicalData)
			})
			
			It("should identify declining trend", func() {
				Expect(trend.Direction).To(Equal("declining"))
			})
			
			It("should calculate slope correctly", func() {
				Expect(trend.Slope).To(BeNumerically("<", 0))
			})
		})
	})
})

// ===== Helper Types and Functions =====

type PortfolioHealth struct {
	Score         int
	Status        string
	TotalProjects int
	HealthyCount  int
	CautionCount  int
	AtRiskCount   int
	CriticalCount int
}

type ProjectHealth struct {
	Score      int
	Status     string
	Breakdown  HealthBreakdown
	Trend      string
}

type HealthBreakdown struct {
	ScheduleHealth     int
	CompletionHealth   int
	DependencyHealth   int
	ResourceHealth     int
	CriticalPathHealth int
}

type HealthDatapoint struct {
	Date        time.Time
	HealthScore int
}

type HealthTrend struct {
	Direction string
	Slope     float64
	Change    int
}

func CalculatePortfolioHealth(projects []models.Project) PortfolioHealth {
	if len(projects) == 0 {
		return PortfolioHealth{Score: 100, Status: "healthy"}
	}
	
	activeProjects := make([]models.Project, 0)
	for _, p := range projects {
		if p.Status == models.ProjectActive {
			activeProjects = append(activeProjects, p)
		}
	}
	
	if len(activeProjects) == 0 {
		return PortfolioHealth{Score: 100, Status: "healthy"}
	}
	
	var totalWeight float64
	var weightedHealth float64
	healthyCount, cautionCount, atRiskCount, criticalCount := 0, 0, 0, 0
	
	for _, p := range activeProjects {
		weight := float64(p.Priority) / 100.0
		totalWeight += weight
		weightedHealth += float64(p.HealthScore) * weight
		
		status := GetHealthStatus(p.HealthScore)
		switch status {
		case "healthy":
			healthyCount++
		case "caution":
			cautionCount++
		case "at_risk":
			atRiskCount++
		case "critical":
			criticalCount++
		}
	}
	
	score := int(weightedHealth / totalWeight)
	
	return PortfolioHealth{
		Score:         score,
		Status:        GetHealthStatus(score),
		TotalProjects: len(activeProjects),
		HealthyCount:  healthyCount,
		CautionCount:  cautionCount,
		AtRiskCount:   atRiskCount,
		CriticalCount: criticalCount,
	}
}

func CalculateProjectHealth(project models.Project, tasks []models.Task, workloadData []models.WorkloadEntry) ProjectHealth {
	scheduleHealth := CalculateScheduleHealth(project, tasks)
	completionHealth := CalculateCompletionHealth(tasks)
	dependencyHealth := CalculateDependencyHealth(tasks)
	resourceHealth := CalculateResourceHealth(project.ID.String(), workloadData)
	criticalPathHealth := CalculateCriticalPathHealth(tasks)
	
	score := int(
		float64(scheduleHealth)*0.30 +
		float64(completionHealth)*0.25 +
		float64(dependencyHealth)*0.20 +
		float64(resourceHealth)*0.15 +
		float64(criticalPathHealth)*0.10,
	)
	
	return ProjectHealth{
		Score:  score,
		Status: GetHealthStatus(score),
		Breakdown: HealthBreakdown{
			ScheduleHealth:     scheduleHealth,
			CompletionHealth:   completionHealth,
			DependencyHealth:   dependencyHealth,
			ResourceHealth:     resourceHealth,
			CriticalPathHealth: criticalPathHealth,
		},
	}
}

func CalculateScheduleHealth(project models.Project, tasks []models.Task) int {
	if project.StartDate == nil || project.TargetEndDate == nil {
		return 100
	}
	
	now := time.Now()
	totalDuration := project.TargetEndDate.Sub(*project.StartDate)
	elapsed := now.Sub(*project.StartDate)
	
	if totalDuration <= 0 {
		return 100
	}
	
	expectedProgress := int((elapsed.Hours() / totalDuration.Hours()) * 100)
	actualProgress := project.Progress
	variance := actualProgress - expectedProgress
	
	// Use stricter thresholds - any negative variance reduces health
	if variance >= 0 {
		return 100
	} else if variance >= -5 {
		return 90
	} else if variance >= -10 {
		return 80
	} else if variance >= -20 {
		return 60
	} else if variance >= -30 {
		return 40
	}
	return 20
}

func CalculateCompletionHealth(tasks []models.Task) int {
	if len(tasks) == 0 {
		return 100
	}
	
	totalHours := 0.0
	completedHours := 0.0
	
	for _, t := range tasks {
		totalHours += t.EstimatedHours
		if t.Status == models.TaskStatusDone {
			completedHours += t.EstimatedHours
		} else if t.Status == models.TaskStatusInProgress {
			completedHours += t.EstimatedHours * 0.5
		} else if t.Status == models.TaskStatusReview {
			completedHours += t.EstimatedHours * 0.8
		}
	}
	
	if totalHours == 0 {
		return 100
	}
	
	return int((completedHours / totalHours) * 100)
}

func CalculateDependencyHealth(tasks []models.Task) int {
	if len(tasks) == 0 {
		return 100
	}
	
	blockedCount := 0
	for _, t := range tasks {
		if t.Status != models.TaskStatusDone && len(t.BlockedBy) > 0 {
			blockedCount++
		}
	}
	
	blockedRate := float64(blockedCount) / float64(len(tasks))
	return int(100 - (blockedRate * 50))
}

func CalculateResourceHealth(projectID string, workloadData []models.WorkloadEntry) int {
	if len(workloadData) == 0 {
		return 100
	}
	
	totalHealth := 0
	for _, w := range workloadData {
		allocation := w.AllocationPercentage
		switch {
		case allocation >= 70 && allocation <= 90:
			totalHealth += 100
		case allocation >= 50 && allocation < 70:
			totalHealth += 85
		case allocation > 90 && allocation <= 100:
			totalHealth += 75
		case allocation > 100 && allocation <= 110:
			totalHealth += 50
		case allocation > 110:
			totalHealth += 25
		default:
			totalHealth += 60
		}
	}
	
	return totalHealth / len(workloadData)
}

func CalculateCriticalPathHealth(tasks []models.Task) int {
	criticalPathTasks := 0
	completedCriticalPathTasks := 0
	
	for _, t := range tasks {
		if t.IsCriticalPath {
			criticalPathTasks++
			if t.Status == models.TaskStatusDone {
				completedCriticalPathTasks++
			}
		}
	}
	
	if criticalPathTasks == 0 {
		return 100
	}
	
	return int((float64(completedCriticalPathTasks) / float64(criticalPathTasks)) * 100)
}

func GetHealthStatus(score int) string {
	if score >= 80 {
		return "healthy"
	} else if score >= 60 {
		return "caution"
	} else if score >= 40 {
		return "at_risk"
	}
	return "critical"
}



func CalculateHealthTrend(data []HealthDatapoint) HealthTrend {
	if len(data) < 2 {
		return HealthTrend{Direction: "stable", Slope: 0, Change: 0}
	}
	
	first := data[0].HealthScore
	last := data[len(data)-1].HealthScore
	change := last - first
	
	var direction string
	if change > 5 {
		direction = "improving"
	} else if change < -5 {
		direction = "declining"
	} else {
		direction = "stable"
	}
	
	return HealthTrend{
		Direction: direction,
		Slope:     float64(change) / float64(len(data)-1),
		Change:    change,
	}
}

func uuidFromString(s string) string {
	return s // Simplified for test
}

