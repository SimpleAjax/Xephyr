package services_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/xephyr-ai/xephyr-backend/internal/models"
	"github.com/xephyr-ai/xephyr-backend/test/fixtures"
	"github.com/xephyr-ai/xephyr-backend/test/helpers"
)

func TestPriorityService(t *testing.T) {
	helpers.RunSuite(t, "Priority Engine Service")
}

var _ = Describe("Priority Engine", func() {
	
	Context("Given a task on the critical path", func() {
		var task models.Task
		var project models.Project
		var allTasks []models.Task
		
		BeforeEach(func() {
			task = fixtures.NewTask().
				WithID("task-critical").
				WithTitle("Critical Task").
				WithPriority(models.TaskPriorityCritical).
				WithBusinessValue(100).
				OnCriticalPath().
				Build()
			
			project = fixtures.NewProject().
				WithID("proj-test").
				WithPriority(95).
				Build()
			
			allTasks = []models.Task{task}
		})
		
		Context("And the task is due within 3 days", func() {
			BeforeEach(func() {
				dueDate := time.Now().Add(2 * 24 * time.Hour)
				task.DueDate = &dueDate
			})
			
			When("priority is calculated", func() {
				var score int
				
				BeforeEach(func() {
					score = CalculatePriorityScore(task, project, allTasks)
				})
				
				It("should have priority score above 90", func() {
					Expect(score).To(BeNumerically(">", 90))
				})
				
				It("should be flagged as critical urgency", func() {
					Expect(IsCriticalUrgency(score)).To(BeTrue())
				})
				
				It("should rank higher than non-critical tasks", func() {
					normalTask := fixtures.NewTask().
						WithID("task-normal").
						WithPriority(models.TaskPriorityMedium).
						WithBusinessValue(50).
						Build()
					normalScore := CalculatePriorityScore(normalTask, project, allTasks)
					Expect(score).To(BeNumerically(">", normalScore))
				})
			})
		})
		
		Context("And the task is overdue", func() {
			BeforeEach(func() {
				dueDate := time.Now().Add(-5 * 24 * time.Hour)
				task.DueDate = &dueDate
			})
			
			When("priority is calculated", func() {
				It("should have maximum deadline urgency", func() {
					urgency := CalculateDeadlineUrgency(*task.DueDate)
					Expect(urgency).To(Equal(100))
				})
				
				It("should have the highest possible priority score", func() {
					score := CalculatePriorityScore(task, project, allTasks)
					Expect(score).To(Equal(100))
				})
			})
		})
	})
	
	Context("Given a task with multiple dependencies", func() {
		var blockingTask models.Task
		var dependentTask1, dependentTask2 models.Task
		var allTasks []models.Task
		
		BeforeEach(func() {
			blockingTask = fixtures.NewTask().
				WithID("task-blocking").
				WithTitle("Blocking Task").
				WithStatus(models.TaskStatusInProgress).
				WithBusinessValue(80).
				Build()
			
			dependentTask1 = fixtures.NewTask().
				WithID("task-dep1").
				WithTitle("Dependent Task 1").
				WithStatus(models.TaskStatusBacklog).
				WithBusinessValue(70).
				Build()
			
			dependentTask2 = fixtures.NewTask().
				WithID("task-dep2").
				WithTitle("Dependent Task 2").
				WithStatus(models.TaskStatusBacklog).
				WithBusinessValue(60).
				OnCriticalPath().
				Build()
			
			allTasks = []models.Task{blockingTask, dependentTask1, dependentTask2}
		})
		
		When("dependency impact is calculated", func() {
			It("should return higher impact for more blocked tasks", func() {
				// Task blocking multiple high-value tasks
				impact := CalculateDependencyImpact(blockingTask, allTasks)
				Expect(impact).To(BeNumerically(">", 0))
			})
			
			It("should return zero impact when no tasks depend on it", func() {
				impact := CalculateDependencyImpact(dependentTask1, allTasks)
				Expect(impact).To(Equal(0))
			})
		})
	})
	
	Context("Given tasks with different due dates", func() {
		DescribeTable("Deadline urgency calculations",
			func(daysUntilDue int, expectedUrgency int) {
				dueDate := time.Now().Add(time.Duration(daysUntilDue) * 24 * time.Hour)
				urgency := CalculateDeadlineUrgency(dueDate)
				Expect(urgency).To(Equal(expectedUrgency))
			},
			Entry("overdue task", -5, 100),
			Entry("due today", 0, 95),
			Entry("due tomorrow", 1, 90),
			Entry("due within 3 days", 2, 80),
			Entry("due within a week", 5, 65),
			Entry("due within 2 weeks", 10, 45),
			Entry("due within a month", 20, 25),
			Entry("due later", 60, 10),
		)
	})
	
	Context("Given projects with different priorities", func() {
		DescribeTable("Project priority contribution",
			func(projectPriority int, expectedContribution float64) {
				project := fixtures.NewProject().WithPriority(projectPriority).Build()
				contribution := float64(project.Priority) * 0.25 // 25% weight
				Expect(contribution).To(BeNumerically("~", expectedContribution, 0.1))
			},
			Entry("high priority project", 95, 23.75),
			Entry("medium priority project", 50, 12.5),
			Entry("low priority project", 25, 6.25),
		)
	})
	
	Context("Given invalid task data", func() {
		When("priority is calculated", func() {
			It("should return error for nil task", func() {
				_, err := CalculatePriorityScoreSafe(models.Task{}, models.Project{}, nil)
				Expect(err).To(HaveOccurred())
			})
			
			It("should handle task with zero estimated hours gracefully", func() {
				task := fixtures.NewTask().WithEstimatedHours(0).Build()
				project := fixtures.NewProject().Build()
				
				score := CalculatePriorityScore(task, project, []models.Task{task})
				Expect(score).To(BeNumerically(">=", 0))
				Expect(score).To(BeNumerically("<=", 100))
			})
		})
	})
	
	Context("Given multiple tasks to rank", func() {
		var tasks []models.Task
		var project models.Project
		
		BeforeEach(func() {
			project = fixtures.NewProject().WithPriority(80).Build()
			now := time.Now()
			
			tasks = []models.Task{
				fixtures.NewTask().
					WithID("task-a").
					WithPriority(models.TaskPriorityLow).
					WithBusinessValue(30).
					WithDueDate(now.Add(30 * 24 * time.Hour)).
					Build(),
				fixtures.NewTask().
					WithID("task-b").
					WithPriority(models.TaskPriorityCritical).
					WithBusinessValue(100).
					OnCriticalPath().
					WithDueDate(now.Add(2 * 24 * time.Hour)).
					Build(),
				fixtures.NewTask().
					WithID("task-c").
					WithPriority(models.TaskPriorityHigh).
					WithBusinessValue(80).
					WithDueDate(now.Add(7 * 24 * time.Hour)).
					Build(),
			}
		})
		
		When("tasks are ranked by priority", func() {
			var rankings []TaskRanking
			
			BeforeEach(func() {
				rankings = RankTasksByPriority(tasks, project)
			})
			
			It("should return tasks in descending priority order", func() {
				Expect(rankings).To(HaveLen(3))
				Expect(rankings[0].TaskID).To(Equal("task-b")) // Critical path, high value, due soon
				Expect(rankings[1].TaskID).To(Equal("task-c")) // High priority
				Expect(rankings[2].TaskID).To(Equal("task-a")) // Low priority
			})
			
			It("should assign consecutive ranks starting from 1", func() {
				for i, ranking := range rankings {
					Expect(ranking.Rank).To(Equal(i + 1))
				}
			})
		})
	})
})

// ===== Helper Types and Functions (to be implemented) =====

type TaskRanking struct {
	TaskID        string
	PriorityScore int
	Rank          int
}

func CalculatePriorityScore(task models.Task, project models.Project, allTasks []models.Task) int {
	// Placeholder for actual implementation
	// Formula: (project.priority * 0.25) + (task.businessValue * 0.25) + 
	//          (deadlineUrgency * 0.25) + (criticalPathWeight * 0.15) + (dependencyImpact * 0.10)
	
	deadlineUrgency := 50
	if task.DueDate != nil {
		deadlineUrgency = CalculateDeadlineUrgency(*task.DueDate)
	}
	
	criticalPathWeight := 0
	if task.IsCriticalPath {
		criticalPathWeight = 30
	}
	
	dependencyImpact := CalculateDependencyImpact(task, allTasks)
	
	score := float64(project.Priority)*0.25 +
		float64(task.BusinessValue)*0.25 +
		float64(deadlineUrgency)*0.25 +
		float64(criticalPathWeight)*0.15 +
		float64(dependencyImpact)*0.10
	
	if score > 100 {
		return 100
	}
	return int(score)
}

func CalculatePriorityScoreSafe(task models.Task, project models.Project, allTasks []models.Task) (int, error) {
	// Placeholder for actual implementation with error handling
	return CalculatePriorityScore(task, project, allTasks), nil
}

func CalculateDeadlineUrgency(dueDate time.Time) int {
	now := time.Now()
	daysUntilDue := int(dueDate.Sub(now).Hours() / 24)
	
	if daysUntilDue < 0 {
		return 100
	}
	if daysUntilDue == 0 {
		return 95
	}
	if daysUntilDue == 1 {
		return 90
	}
	if daysUntilDue <= 3 {
		return 80
	}
	if daysUntilDue <= 7 {
		return 65
	}
	if daysUntilDue <= 14 {
		return 45
	}
	if daysUntilDue <= 30 {
		return 25
	}
	return 10
}

func CalculateDependencyImpact(task models.Task, allTasks []models.Task) int {
	// Find tasks that depend on this task
	blockedValue := 0.0
	for _, t := range allTasks {
		for _, dep := range t.Dependencies {
			if dep.DependsOnTaskID == task.ID {
				multiplier := 1.0
				if t.IsCriticalPath {
					multiplier = 1.5
				}
				blockedValue += float64(t.BusinessValue) * multiplier
			}
		}
	}
	
	// Normalize to 0-20 scale
	impact := int(blockedValue / 50)
	if impact > 20 {
		return 20
	}
	return impact
}

func IsCriticalUrgency(score int) bool {
	return score >= 90
}

func RankTasksByPriority(tasks []models.Task, project models.Project) []TaskRanking {
	// Calculate scores
	scores := make([]struct {
		task  models.Task
		score int
	}, len(tasks))
	
	for i, task := range tasks {
		scores[i] = struct {
			task  models.Task
			score int
		}{
			task:  task,
			score: CalculatePriorityScore(task, project, tasks),
		}
	}
	
	// Simple bubble sort (in production, use sort.Slice)
	for i := 0; i < len(scores); i++ {
		for j := i + 1; j < len(scores); j++ {
			if scores[j].score > scores[i].score {
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	}
	
	// Create rankings
	rankings := make([]TaskRanking, len(scores))
	for i, s := range scores {
		rankings[i] = TaskRanking{
			TaskID:        s.task.ID.String(),
			PriorityScore: s.score,
			Rank:          i + 1,
		}
	}
	
	return rankings
}

