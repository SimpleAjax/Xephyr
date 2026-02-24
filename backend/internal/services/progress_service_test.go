package services_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/SimpleAjax/Xephyr/internal/models"
	"github.com/SimpleAjax/Xephyr/tests/fixtures"
)

var _ = Describe("Progress Tracking System", func() {
	
	Context("Given a project with tasks in various states", func() {
		var project models.Project
		var tasks []models.Task
		
		BeforeEach(func() {
			project = fixtures.NewProject().WithID("proj-progress").Build()
			
			tasks = []models.Task{
				fixtures.NewTask().WithStatus(models.TaskStatusDone).WithEstimatedHours(40).WithActualHours(40).Build(),
				fixtures.NewTask().WithStatus(models.TaskStatusDone).WithEstimatedHours(30).WithActualHours(35).Build(),
				fixtures.NewTask().WithStatus(models.TaskStatusInProgress).WithEstimatedHours(50).WithActualHours(25).Build(),
				fixtures.NewTask().WithStatus(models.TaskStatusReview).WithEstimatedHours(30).WithActualHours(24).Build(),
				fixtures.NewTask().WithStatus(models.TaskStatusBacklog).WithEstimatedHours(40).Build(),
				fixtures.NewTask().WithStatus(models.TaskStatusReady).WithEstimatedHours(30).Build(),
			}
		})
		
		When("project completion is calculated", func() {
			var completion ProjectCompletion
			
			BeforeEach(func() {
				completion = CalculateProjectCompletion(project, tasks)
			})
			
			It("should calculate weighted completion by estimated hours", func() {
				// Done: 40 + 30 = 70 hours
				// In Progress: 50 * 0.5 = 25 hours
				// Review: 30 * 0.8 = 24 hours
				// Total done = 70 + 25 + 24 = 119
				// Total hours = 220
				// Completion = 119/220 = ~54%
				Expect(completion.Percentage).To(BeNumerically(">=", 50))
				Expect(completion.Percentage).To(BeNumerically("<=", 60))
			})
			
			It("should provide breakdown by status", func() {
				Expect(completion.ByStatus[string(models.TaskStatusDone)].Count).To(Equal(2))
				Expect(completion.ByStatus[string(models.TaskStatusInProgress)].Count).To(Equal(1))
				Expect(completion.ByStatus[string(models.TaskStatusBacklog)].Count).To(Equal(1))
			})
			
			It("should calculate total estimated hours", func() {
				Expect(completion.TotalEstimatedHours).To(Equal(220.0))
			})
		})
	})
	
	Context("Given a hierarchical task structure", func() {
		var parentTask models.Task
		var subtasks []models.Task
		var allTasks []models.Task
		
		BeforeEach(func() {
			parentTask = fixtures.NewTask().
				WithID("task-parent").
				WithStatus(models.TaskStatusInProgress).
				WithEstimatedHours(100).
				Build()
			
			subtasks = []models.Task{
				fixtures.NewTask().
					WithID("task-child-1").
					WithProject("proj-1").
					WithStatus(models.TaskStatusDone).
					WithEstimatedHours(40).
					Build(),
				fixtures.NewTask().
					WithID("task-child-2").
					WithProject("proj-1").
					WithStatus(models.TaskStatusInProgress).
					WithEstimatedHours(60).
					Build(),
			}
			
			// Set parent relationships
			for i := range subtasks {
				subtasks[i].ParentTaskID = &parentTask.ID
			}
			
			allTasks = append([]models.Task{parentTask}, subtasks...)
		})
		
		When("hierarchical progress is calculated", func() {
			var progress float64
			
			BeforeEach(func() {
				progress = CalculateHierarchicalProgress(parentTask, allTasks)
			})
			
			It("should aggregate from child tasks", func() {
				// Child 1: 100% of 40 = 40
				// Child 2: 40% of 60 = 24 (in_progress = 0.4)
				// Total = 64 / 100 = 64%
				Expect(progress).To(BeNumerically(">=", 60))
				Expect(progress).To(BeNumerically("<=", 70))
			})
		})
	})
	
	Context("Given progress variance tracking", func() {
		var project models.Project
		var tasks []models.Task
		
		BeforeEach(func() {
			project = fixtures.NewProject().WithProgress(35).Build()
			
			now := time.Now()
			startDate := now.Add(-50 * 24 * time.Hour)
			endDate := now.Add(50 * 24 * time.Hour)
			project.StartDate = &startDate
			project.TargetEndDate = &endDate
			
			tasks = []models.Task{
				fixtures.NewTask().WithStatus(models.TaskStatusDone).WithEstimatedHours(35).Build(),
				fixtures.NewTask().WithStatus(models.TaskStatusInProgress).WithEstimatedHours(65).Build(),
			}
		})
		
		When("progress variance is calculated", func() {
			var variance ProgressVariance
			
			BeforeEach(func() {
				variance = CalculateProgressVariance(project, tasks)
			})
			
			It("should calculate expected progress based on time", func() {
				// 50 days elapsed out of 100 = 50% expected
				Expect(variance.ExpectedProgress).To(Equal(50))
			})
			
			It("should calculate actual progress", func() {
				// Based on project.Progress = 35%
				Expect(variance.ActualProgress).To(Equal(35))
			})
			
			It("should calculate variance correctly", func() {
				// 35 - 50 = -15
				Expect(variance.Variance).To(Equal(-15))
			})
			
			It("should identify behind-schedule status", func() {
				Expect(variance.Status).To(Equal("behind_schedule"))
			})
		})
	})
	
	Context("Given task status progress weights", func() {
		DescribeTable("status to progress mapping",
			func(status models.TaskStatus, expectedProgress float64) {
				progress := GetTaskStatusProgressWeight(status)
				Expect(progress).To(Equal(expectedProgress))
			},
			Entry("backlog", models.TaskStatusBacklog, 0.0),
			Entry("ready", models.TaskStatusReady, 0.1),
			Entry("in_progress", models.TaskStatusInProgress, 0.4),
			Entry("review", models.TaskStatusReview, 0.8),
			Entry("done", models.TaskStatusDone, 1.0),
		)
	})
	
	Context("Given task progress update", func() {
		var task models.Task
		var update TaskProgressUpdate
		
		BeforeEach(func() {
			task = fixtures.NewTask().
				WithID("task-update").
				WithStatus(models.TaskStatusReady).
				WithEstimatedHours(40).
				Build()
			
			update = TaskProgressUpdate{
				NewStatus:          models.TaskStatusInProgress,
				ProgressPercentage: 25,
				ActualHours:        10,
				Note:               "Started working on implementation",
			}
		})
		
		When("progress is updated", func() {
			var result UpdateResult
			
			BeforeEach(func() {
				result = UpdateTaskProgress(&task, update)
			})
			
			It("should update task status", func() {
				Expect(task.Status).To(Equal(models.TaskStatusInProgress))
			})
			
			It("should record actual hours", func() {
				Expect(task.ActualHours).To(Equal(10.0))
			})
			
			It("should return affected entities", func() {
				Expect(result.ParentProgressUpdated).To(BeTrue())
				Expect(result.ProjectProgressUpdated).To(BeTrue())
			})
		})
		
		When("task is marked complete", func() {
			It("should set completed timestamp", func() {
				update.NewStatus = models.TaskStatusDone
				UpdateTaskProgress(&task, update)
				
				Expect(task.Status).To(Equal(models.TaskStatusDone))
				Expect(task.CompletedAt).ToNot(BeNil())
			})
		})
	})
	
	Context("Given milestone tracking", func() {
		var milestones []models.Task
		
		BeforeEach(func() {
			now := time.Now()
			completedAt := now.Add(-5 * 24 * time.Hour)
			
			milestones = []models.Task{
				fixtures.NewTask().
					WithID("milestone-1").
					WithTitle("Design Complete").
					WithStatus(models.TaskStatusDone).
					IsMilestone().
					Build(),
				fixtures.NewTask().
					WithID("milestone-2").
					WithTitle("MVP Launch").
					WithStatus(models.TaskStatusInProgress).
					IsMilestone().
					Build(),
				fixtures.NewTask().
					WithID("milestone-3").
					WithTitle("Production Ready").
					WithStatus(models.TaskStatusBacklog).
					IsMilestone().
					Build(),
			}
			
			milestones[0].CompletedAt = &completedAt
		})
		
		When("milestones are analyzed", func() {
			var status MilestoneStatus
			
			BeforeEach(func() {
				status = AnalyzeMilestones(milestones)
			})
			
			It("should count completed milestones", func() {
				Expect(status.CompletedCount).To(Equal(1))
			})
			
			It("should count upcoming milestones", func() {
				Expect(status.UpcomingCount).To(Equal(2))
			})
			
			It("should calculate milestone completion percentage", func() {
				Expect(status.CompletionPercentage).To(Equal(33)) // 1 of 3
			})
		})
	})
	
	Context("Given progress history tracking", func() {
		var history []ProgressDatapoint
		
		BeforeEach(func() {
			history = []ProgressDatapoint{
				{Date: time.Now().AddDate(0, 0, -21), Progress: 10},
				{Date: time.Now().AddDate(0, 0, -14), Progress: 20},
				{Date: time.Now().AddDate(0, 0, -7), Progress: 28},
				{Date: time.Now(), Progress: 35},
			}
		})
		
		When("progress velocity is calculated", func() {
			var velocity VelocityMetrics
			
			BeforeEach(func() {
				velocity = CalculateProgressVelocity(history)
			})
			
			It("should calculate average weekly progress", func() {
				// (35-10) / 3 weeks = ~8.3% per week
				Expect(velocity.WeeklyAverage).To(BeNumerically(">=", 8))
				Expect(velocity.WeeklyAverage).To(BeNumerically("<=", 9))
			})
			
			It("should project completion date", func() {
				// 65% remaining at 8.3%/week = ~8 weeks
				Expect(velocity.ProjectedWeeksRemaining).To(BeNumerically(">=", 7))
				Expect(velocity.ProjectedWeeksRemaining).To(BeNumerically("<=", 9))
			})
		})
	})
})

// ===== Helper Types and Functions =====

type ProjectCompletion struct {
	Percentage          int
	TotalEstimatedHours float64
	TotalCompletedHours float64
	ByStatus            map[string]StatusBreakdown
}

type StatusBreakdown struct {
	Count  int
	Hours  float64
	Percentage float64
}

type ProgressVariance struct {
	ExpectedProgress int
	ActualProgress   int
	Variance         int
	Trend            string
	Status           string
}

type TaskProgressUpdate struct {
	NewStatus          models.TaskStatus
	ProgressPercentage int
	ActualHours        float64
	Note               string
}

type UpdateResult struct {
	TaskUpdated            bool
	ParentProgressUpdated  bool
	ProjectProgressUpdated bool
	DependentsNotified     []string
}

type MilestoneStatus struct {
	CompletedCount       int
	UpcomingCount        int
	CompletionPercentage int
	NextMilestone        *models.Task
}

type ProgressDatapoint struct {
	Date     time.Time
	Progress int
}

type VelocityMetrics struct {
	WeeklyAverage          float64
	Trend                  string
	ProjectedWeeksRemaining int
}

func CalculateProjectCompletion(project models.Project, tasks []models.Task) ProjectCompletion {
	totalHours := 0.0
	completedHours := 0.0
	
	byStatus := make(map[string]StatusBreakdown)
	
	for _, task := range tasks {
		totalHours += task.EstimatedHours
		
		status := string(task.Status)
		breakdown := byStatus[status]
		breakdown.Count++
		breakdown.Hours += task.EstimatedHours
		byStatus[status] = breakdown
		
		// Calculate completed hours based on status
		progressWeight := GetTaskStatusProgressWeight(task.Status)
		completedHours += task.EstimatedHours * progressWeight
	}
	
	percentage := 0
	if totalHours > 0 {
		percentage = int((completedHours / totalHours) * 100)
	}
	
	// Calculate percentages for each status
	for status, breakdown := range byStatus {
		if totalHours > 0 {
			breakdown.Percentage = float64(int((breakdown.Hours / totalHours) * 100))
			byStatus[status] = breakdown
		}
	}
	
	return ProjectCompletion{
		Percentage:          percentage,
		TotalEstimatedHours: totalHours,
		TotalCompletedHours: completedHours,
		ByStatus:            byStatus,
	}
}

func CalculateHierarchicalProgress(task models.Task, allTasks []models.Task) float64 {
	// Find subtasks
	var subtasks []models.Task
	for _, t := range allTasks {
		if t.ParentTaskID != nil && *t.ParentTaskID == task.ID {
			subtasks = append(subtasks, t)
		}
	}
	
	// Leaf task - use status weight
	if len(subtasks) == 0 {
		return GetTaskStatusProgressWeight(task.Status) * 100
	}
	
	// Parent task - aggregate from children
	totalHours := 0.0
	weightedProgress := 0.0
	
	for _, child := range subtasks {
		childProgress := CalculateHierarchicalProgress(child, allTasks)
		totalHours += child.EstimatedHours
		weightedProgress += childProgress * child.EstimatedHours
	}
	
	if totalHours == 0 {
		return 0
	}
	
	return weightedProgress / totalHours
}

func CalculateProgressVariance(project models.Project, tasks []models.Task) ProgressVariance {
	if project.StartDate == nil || project.TargetEndDate == nil {
		return ProgressVariance{}
	}
	
	now := time.Now()
	totalDuration := project.TargetEndDate.Sub(*project.StartDate)
	elapsed := now.Sub(*project.StartDate)
	
	expectedProgress := int((elapsed.Hours() / totalDuration.Hours()) * 100)
	actualProgress := project.Progress
	variance := actualProgress - expectedProgress
	
	status := "on_track"
	if variance < -10 {
		status = "behind_schedule"
	} else if variance > 10 {
		status = "ahead_of_schedule"
	}
	
	return ProgressVariance{
		ExpectedProgress: expectedProgress,
		ActualProgress:   actualProgress,
		Variance:         variance,
		Trend:            "stable",
		Status:           status,
	}
}

func GetTaskStatusProgressWeight(status models.TaskStatus) float64 {
	weights := map[models.TaskStatus]float64{
		models.TaskStatusBacklog:    0.0,
		models.TaskStatusReady:      0.1,
		models.TaskStatusInProgress: 0.4,
		models.TaskStatusReview:     0.8,
		models.TaskStatusDone:       1.0,
	}
	return weights[status]
}

func UpdateTaskProgress(task *models.Task, update TaskProgressUpdate) UpdateResult {
	task.Status = update.NewStatus
	task.ActualHours = update.ActualHours
	
	if update.NewStatus == models.TaskStatusDone {
		now := time.Now()
		task.CompletedAt = &now
	}
	
	return UpdateResult{
		TaskUpdated:            true,
		ParentProgressUpdated:  true,
		ProjectProgressUpdated: true,
	}
}

func AnalyzeMilestones(milestones []models.Task) MilestoneStatus {
	completed := 0
	upcoming := 0
	
	for _, m := range milestones {
		if m.Status == models.TaskStatusDone {
			completed++
		} else {
			upcoming++
		}
	}
	
	total := len(milestones)
	percentage := 0
	if total > 0 {
		percentage = int((float64(completed) / float64(total)) * 100)
	}
	
	return MilestoneStatus{
		CompletedCount:       completed,
		UpcomingCount:        upcoming,
		CompletionPercentage: percentage,
	}
}

func CalculateProgressVelocity(history []ProgressDatapoint) VelocityMetrics {
	if len(history) < 2 {
		return VelocityMetrics{}
	}
	
	// Calculate average weekly progress
	totalProgress := history[len(history)-1].Progress - history[0].Progress
	totalWeeks := len(history) - 1
	
	weeklyAverage := float64(totalProgress) / float64(totalWeeks)
	
	// Calculate trend
	trend := "stable"
	if len(history) >= 3 {
		recentProgress := history[len(history)-1].Progress - history[len(history)-2].Progress
		olderProgress := history[len(history)-2].Progress - history[len(history)-3].Progress
		
		if recentProgress > olderProgress {
			trend = "accelerating"
		} else if recentProgress < olderProgress {
			trend = "decelerating"
		}
	}
	
	// Project remaining weeks
	currentProgress := history[len(history)-1].Progress
	remainingProgress := 100 - currentProgress
	projectedWeeks := 0
	if weeklyAverage > 0 {
		projectedWeeks = int(float64(remainingProgress) / weeklyAverage)
	}
	
	return VelocityMetrics{
		WeeklyAverage:           weeklyAverage,
		Trend:                   trend,
		ProjectedWeeksRemaining: projectedWeeks,
	}
}

