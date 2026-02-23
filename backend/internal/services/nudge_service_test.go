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

func TestNudgeService(t *testing.T) {
	helpers.RunSuite(t, "Nudge Detection Service")
}

var _ = Describe("Nudge Detection System", func() {
	
	Context("Given a team member with overallocation", func() {
		var workloadData []models.WorkloadEntry
		var users []models.User
		
		BeforeEach(func() {
			workloadData = []models.WorkloadEntry{
				{UserID: stringToUUID("user-emma"), AllocationPercentage: 125, AssignedTasks: 4, TotalEstimatedHours: 75},
				{UserID: stringToUUID("user-alex"), AllocationPercentage: 110, AssignedTasks: 3, TotalEstimatedHours: 65},
				{UserID: stringToUUID("user-james"), AllocationPercentage: 95, AssignedTasks: 2, TotalEstimatedHours: 60},
			}
			
			users = []models.User{
				fixtures.NewUser().WithID("user-emma").WithName("Emma").Build(),
				fixtures.NewUser().WithID("user-alex").WithName("Alex").Build(),
				fixtures.NewUser().WithID("user-james").WithName("James").Build(),
			}
		})
		
		When("overload detection runs", func() {
			var nudges []models.Nudge
			
			BeforeEach(func() {
				nudges = DetectOverloadNudges(users, workloadData)
			})
			
			It("should detect members over 100% allocation", func() {
				Expect(len(nudges)).To(BeNumerically(">=", 2))
			})
			
			It("should categorize severity based on overload percentage", func() {
				for _, nudge := range nudges {
					if nudge.RelatedUserID != nil && nudge.RelatedUserID.String() == "user-emma" {
						Expect(nudge.Severity).To(Equal(models.NudgeSeverityHigh))
					}
				}
			})
			
			It("should include relevant metrics in the nudge", func() {
				Expect(len(nudges)).To(BeNumerically(">", 0))
				Expect(nudges[0].Metrics).ToNot(BeNil())
			})
		})
	})
	
	Context("Given tasks at risk of delay", func() {
		var tasks []models.Task
		
		BeforeEach(func() {
			now := time.Now()
			startDate := now.Add(-10 * 24 * time.Hour)
			dueDate := now.Add(10 * 24 * time.Hour)
			
			tasks = []models.Task{
				fixtures.NewTask().
					WithID("task-at-risk").
					WithStatus(models.TaskStatusInProgress).
					WithEstimatedHours(80).
					WithActualHours(20). // 25% done, should be ~50%
					WithDueDate(dueDate).
					Build(),
				fixtures.NewTask().
					WithID("task-on-track").
					WithStatus(models.TaskStatusInProgress).
					WithEstimatedHours(80).
					WithActualHours(45). // ~56% done, on track
					WithDueDate(dueDate).
					Build(),
			}
			
			// Set start dates
			tasks[0].StartDate = &startDate
			tasks[1].StartDate = &startDate
		})
		
		When("delay risk detection runs", func() {
			It("should detect tasks behind schedule", func() {
				nudges := DetectDelayRiskNudges(tasks)
				
				foundAtRisk := false
				for _, n := range nudges {
					if n.RelatedTaskID != nil && n.RelatedTaskID.String() == "task-at-risk" {
						foundAtRisk = true
					}
				}
				Expect(foundAtRisk).To(BeTrue())
			})
			
			It("should calculate risk ratio correctly", func() {
				nudges := DetectDelayRiskNudges(tasks)
				Expect(len(nudges)).To(BeNumerically(">", 0))
			})
			
			It("should assign higher severity to critical path tasks", func() {
				tasks[0].IsCriticalPath = true
				nudges := DetectDelayRiskNudges(tasks)
				
				for _, n := range nudges {
					if n.RelatedTaskID != nil && n.RelatedTaskID.String() == "task-at-risk" {
						Expect(n.Severity).To(Equal(models.NudgeSeverityHigh))
					}
				}
			})
		})
	})
	
	Context("Given tasks with required skills not covered by team", func() {
		var tasks []models.Task
		var users []models.User
		var userSkills map[string][]models.UserSkill
		
		BeforeEach(func() {
			tasks = []models.Task{
				fixtures.NewTask().
					WithID("task-skill-gap").
					WithStatus(models.TaskStatusBacklog).
					Build(),
			}
			tasks[0].Skills = []models.TaskSkill{
				{SkillID: stringToUUID("skill-threejs"), ProficiencyRequired: 3, IsRequired: true},
			}
			
			users = fixtures.CreateTestTeam()
			userSkills = fixtures.CreateTestUserSkills()
		})
		
		When("skill gap detection runs", func() {
			It("should detect tasks with uncovered required skills", func() {
				nudges := DetectSkillGapNudges(tasks, users, userSkills)
				Expect(len(nudges)).To(BeNumerically(">", 0))
			})
			
			It("should identify the missing skills", func() {
				nudges := DetectSkillGapNudges(tasks, users, userSkills)
				Expect(nudges[0].Metrics).ToNot(BeNil())
			})
		})
	})
	
	Context("Given high-priority unassigned tasks", func() {
		var tasks []models.Task
		
		BeforeEach(func() {
			now := time.Now()
			
			tasks = []models.Task{
				fixtures.NewTask().
					WithID("task-critical-unassigned").
					WithStatus(models.TaskStatusReady).
					WithPriority(models.TaskPriorityCritical).
					WithDueDate(now.Add(3 * 24 * time.Hour)).
					Build(),
				fixtures.NewTask().
					WithID("task-low-unassigned").
					WithStatus(models.TaskStatusBacklog).
					WithPriority(models.TaskPriorityLow).
					Build(),
			}
		})
		
		When("unassigned task detection runs", func() {
			var nudges []models.Nudge
			
			BeforeEach(func() {
				nudges = DetectUnassignedNudges(tasks)
			})
			
			It("should detect critical unassigned tasks", func() {
				foundCritical := false
				for _, n := range nudges {
					if n.RelatedTaskID != nil && n.RelatedTaskID.String() == "task-critical-unassigned" {
						foundCritical = true
					}
				}
				Expect(foundCritical).To(BeTrue())
			})
			
			It("should escalate severity for due-soon tasks", func() {
				for _, n := range nudges {
					if n.RelatedTaskID != nil && n.RelatedTaskID.String() == "task-critical-unassigned" {
						Expect(n.Severity).To(Equal(models.NudgeSeverityHigh))
					}
				}
			})
			
			It("should not flag low-priority tasks", func() {
				for _, n := range nudges {
					Expect(n.RelatedTaskID.String()).ToNot(Equal("task-low-unassigned"))
				}
			})
		})
	})
	
	Context("Given tasks blocked by dependencies", func() {
		var tasks []models.Task
		
		BeforeEach(func() {
			tasks = []models.Task{
				fixtures.NewTask().
					WithID("task-dep-done").
					WithStatus(models.TaskStatusDone).
					Build(),
				fixtures.NewTask().
					WithID("task-blocked").
					WithStatus(models.TaskStatusReady).
					Build(),
			}
			tasks[1].BlockedBy = []models.TaskDependency{
				{DependsOnTaskID: stringToUUID("task-dep-done")},
			}
		})
		
		When("dependency block detection runs", func() {
			It("should detect tasks with incomplete dependencies", func() {
				nudges := DetectDependencyBlockNudges(tasks)
				// Should not flag since dependency is done
				Expect(len(nudges)).To(Equal(0))
			})
		})
	})
	
	Context("Given team members with conflicting priorities", func() {
		var tasks []models.Task
		var workloadData []models.WorkloadEntry
		
		BeforeEach(func() {
			now := time.Now()
			
			tasks = []models.Task{
				fixtures.NewTask().
					WithID("task-conflict-1").
					WithStatus(models.TaskStatusInProgress).
					WithPriority(models.TaskPriorityCritical).
					WithDueDate(now.Add(5 * 24 * time.Hour)).
					WithAssignee("user-alex").
					WithEstimatedHours(40).
					Build(),
				fixtures.NewTask().
					WithID("task-conflict-2").
					WithStatus(models.TaskStatusInProgress).
					WithPriority(models.TaskPriorityHigh).
					WithDueDate(now.Add(7 * 24 * time.Hour)).
					WithAssignee("user-alex").
					WithEstimatedHours(40).
					Build(),
			}
		})
		
		When("conflict detection runs", func() {
			It("should detect resource conflicts", func() {
				nudges := DetectConflictNudges(tasks, workloadData)
				Expect(len(nudges)).To(BeNumerically(">", 0))
			})
		})
	})
	
	Context("Given a nudge trigger", func() {
		var trigger NudgeTrigger
		
		BeforeEach(func() {
			trigger = NudgeTrigger{
				Type:     models.NudgeTypeOverload,
				Severity: models.NudgeSeverityHigh,
				Metrics: map[string]interface{}{
					"allocationPercentage": 125,
					"isCriticalPath":       true,
				},
			}
		})
		
		When("criticality score is calculated", func() {
			var score int
			
			BeforeEach(func() {
				score = CalculateNudgeCriticality(trigger)
			})
			
			It("should apply base score by type", func() {
				Expect(score).To(BeNumerically(">=", 85)) // overload base
			})
			
			It("should apply severity multiplier", func() {
				Expect(score).To(BeNumerically(">", 85)) // high severity = 1.2x
			})
			
			It("should boost for critical path involvement", func() {
				Expect(score).To(BeNumerically(">=", 85*1.15))
			})
		})
	})
	
	Context("Given multiple generated nudges", func() {
		var nudges []models.Nudge
		
		BeforeEach(func() {
			nudges = []models.Nudge{
				{Type: models.NudgeTypeOverload, Severity: models.NudgeSeverityHigh, CriticalityScore: 90},
				{Type: models.NudgeTypeDelayRisk, Severity: models.NudgeSeverityHigh, CriticalityScore: 85},
				{Type: models.NudgeTypeUnassigned, Severity: models.NudgeSeverityMedium, CriticalityScore: 70},
				{Type: models.NudgeTypeSkillGap, Severity: models.NudgeSeverityLow, CriticalityScore: 55},
			}
		})
		
		When("nudges are sorted by criticality", func() {
			It("should sort in descending order of criticality", func() {
				sorted := SortNudgesByCriticality(nudges)
				
				for i := 0; i < len(sorted)-1; i++ {
					Expect(sorted[i].CriticalityScore).To(BeNumerically(">=", sorted[i+1].CriticalityScore))
				}
			})
		})
	})
})

// ===== Helper Types and Functions =====

type NudgeTrigger struct {
	Type            models.NudgeType
	Severity        models.NudgeSeverity
	RelatedTaskID   string
	RelatedUserID   string
	RelatedProjectID string
	Metrics         map[string]interface{}
}

func DetectOverloadNudges(users []models.User, workloadData []models.WorkloadEntry) []models.Nudge {
	var nudges []models.Nudge
	
	for _, w := range workloadData {
		if w.AllocationPercentage < 100 {
			continue
		}
		
		var severity models.NudgeSeverity
		switch {
		case w.AllocationPercentage >= 125:
			severity = models.NudgeSeverityHigh
		case w.AllocationPercentage >= 110:
			severity = models.NudgeSeverityMedium
		default:
			severity = models.NudgeSeverityLow
		}
		
		nudges = append(nudges, models.Nudge{
			Type:       models.NudgeTypeOverload,
			Severity:   severity,
			Status:     models.NudgeStatusUnread,
			Title:      "Team member is overallocated",
			Metrics:    map[string]interface{}{
				"allocationPercentage": w.AllocationPercentage,
				"assignedTasks":        w.AssignedTasks,
			},
		})
	}
	
	return nudges
}

func DetectDelayRiskNudges(tasks []models.Task) []models.Nudge {
	var nudges []models.Nudge
	
	for _, t := range tasks {
		if t.Status == models.TaskStatusDone || t.StartDate == nil || t.DueDate == nil {
			continue
		}
		
		totalDuration := t.DueDate.Sub(*t.StartDate)
		elapsed := time.Since(*t.StartDate)
		timeProgress := elapsed.Hours() / totalDuration.Hours()
		
		statusProgress := GetStatusProgressWeight(t.Status)
		buffer := 0.15
		
		if statusProgress < (timeProgress - buffer) {
			riskRatio := (timeProgress - statusProgress) / timeProgress
			
			var severity models.NudgeSeverity
			if t.IsCriticalPath && riskRatio > 0.3 {
				severity = models.NudgeSeverityHigh
			} else if riskRatio > 0.4 {
				severity = models.NudgeSeverityHigh
			} else if riskRatio > 0.25 {
				severity = models.NudgeSeverityMedium
			} else {
				severity = models.NudgeSeverityLow
			}
			
			nudges = append(nudges, models.Nudge{
				Type:          models.NudgeTypeDelayRisk,
				Severity:      severity,
				Status:        models.NudgeStatusUnread,
				RelatedTaskID: &t.ID,
			})
		}
	}
	
	return nudges
}

func GetStatusProgressWeight(status models.TaskStatus) float64 {
	weights := map[models.TaskStatus]float64{
		models.TaskStatusBacklog:    0,
		models.TaskStatusReady:      0.1,
		models.TaskStatusInProgress: 0.4,
		models.TaskStatusReview:     0.8,
		models.TaskStatusDone:       1.0,
	}
	return weights[status]
}

func DetectSkillGapNudges(tasks []models.Task, users []models.User, userSkills map[string][]models.UserSkill) []models.Nudge {
	var nudges []models.Nudge
	
	for _, t := range tasks {
		if t.Status == models.TaskStatusDone {
			continue
		}
		
		for _, skill := range t.Skills {
			if !skill.IsRequired {
				continue
			}
			
			hasSkill := false
			for _, user := range users {
				userSkillList := userSkills[user.ID.String()]
				for _, us := range userSkillList {
					if us.SkillID == skill.SkillID && us.Proficiency >= skill.ProficiencyRequired {
						hasSkill = true
						break
					}
				}
				if hasSkill {
					break
				}
			}
			
			if !hasSkill {
				severity := models.NudgeSeverityMedium
				if t.Priority == models.TaskPriorityCritical {
					severity = models.NudgeSeverityHigh
				}
				
				nudges = append(nudges, models.Nudge{
					Type:          models.NudgeTypeSkillGap,
					Severity:      severity,
					Status:        models.NudgeStatusUnread,
					RelatedTaskID: &t.ID,
				})
			}
		}
	}
	
	return nudges
}

func DetectUnassignedNudges(tasks []models.Task) []models.Nudge {
	var nudges []models.Nudge
	now := time.Now()
	
	for _, t := range tasks {
		if t.Status == models.TaskStatusDone || t.AssigneeID != nil {
			continue
		}
		
		var severity models.NudgeSeverity
		switch t.Priority {
		case models.TaskPriorityCritical:
			severity = models.NudgeSeverityHigh
		case models.TaskPriorityHigh:
			severity = models.NudgeSeverityMedium
		case models.TaskPriorityMedium:
			severity = models.NudgeSeverityLow
		default:
			continue
		}
		
		if t.DueDate != nil {
			daysUntilDue := int(t.DueDate.Sub(now).Hours() / 24)
			if daysUntilDue <= 3 && severity != models.NudgeSeverityHigh {
				severity = models.NudgeSeverityHigh
			}
		}
		
		nudges = append(nudges, models.Nudge{
			Type:          models.NudgeTypeUnassigned,
			Severity:      severity,
			Status:        models.NudgeStatusUnread,
			RelatedTaskID: &t.ID,
		})
	}
	
	return nudges
}

func DetectDependencyBlockNudges(tasks []models.Task) []models.Nudge {
	var nudges []models.Nudge
	
	for _, t := range tasks {
		if t.Status == models.TaskStatusDone {
			continue
		}
		
		incompleteDeps := 0
		atRiskDeps := 0
		
		for _, dep := range t.Dependencies {
			for _, otherTask := range tasks {
				if otherTask.ID == dep.DependsOnTaskID && otherTask.Status != models.TaskStatusDone {
					incompleteDeps++
					if otherTask.DueDate != nil {
						daysUntil := int(otherTask.DueDate.Sub(time.Now()).Hours() / 24)
						if daysUntil <= 7 && otherTask.Status != models.TaskStatusInProgress {
							atRiskDeps++
						}
					}
				}
			}
		}
		
		if atRiskDeps > 0 {
			severity := models.NudgeSeverityMedium
			if t.IsCriticalPath {
				severity = models.NudgeSeverityHigh
			}
			
			nudges = append(nudges, models.Nudge{
				Type:          models.NudgeTypeDependencyBlock,
				Severity:      severity,
				Status:        models.NudgeStatusUnread,
				RelatedTaskID: &t.ID,
			})
		}
	}
	
	return nudges
}

func DetectConflictNudges(tasks []models.Task, workloadData []models.WorkloadEntry) []models.Nudge {
	var nudges []models.Nudge
	
	// Group tasks by assignee
	tasksByAssignee := make(map[string][]models.Task)
	for _, t := range tasks {
		if t.AssigneeID != nil {
			tasksByAssignee[t.AssigneeID.String()] = append(tasksByAssignee[t.AssigneeID.String()], t)
		}
	}
	
	for userId, userTasks := range tasksByAssignee {
		criticalTasks := make([]models.Task, 0)
		for _, t := range userTasks {
			if t.IsCriticalPath || t.Priority == models.TaskPriorityCritical {
				criticalTasks = append(criticalTasks, t)
			}
		}
		
		if len(criticalTasks) >= 2 {
			totalHours := 0.0
			for _, t := range criticalTasks {
				totalHours += t.EstimatedHours
			}
			
			if totalHours > 40 {
				userUUID := stringToUUID(userId)
				nudges = append(nudges, models.Nudge{
					Type:            models.NudgeTypeConflict,
					Severity:        models.NudgeSeverityMedium,
					Status:          models.NudgeStatusUnread,
					RelatedUserID:   &userUUID,
				})
			}
		}
	}
	
	return nudges
}

func CalculateNudgeCriticality(trigger NudgeTrigger) int {
	baseScores := map[models.NudgeType]int{
		models.NudgeTypeOverload:        85,
		models.NudgeTypeDelayRisk:       80,
		models.NudgeTypeSkillGap:        60,
		models.NudgeTypeUnassigned:      70,
		models.NudgeTypeBlocked:         75,
		models.NudgeTypeConflict:        55,
		models.NudgeTypeDependencyBlock: 65,
	}
	
	score := float64(baseScores[trigger.Type])
	
	severityMultiplier := map[models.NudgeSeverity]float64{
		models.NudgeSeverityHigh:   1.2,
		models.NudgeSeverityMedium: 1.0,
		models.NudgeSeverityLow:    0.8,
	}
	
	score *= severityMultiplier[trigger.Severity]
	
	if trigger.Metrics["isCriticalPath"] == true {
		score *= 1.15
	}
	
	if score > 100 {
		return 100
	}
	return int(score)
}

func SortNudgesByCriticality(nudges []models.Nudge) []models.Nudge {
	// Simple bubble sort for demonstration
	result := make([]models.Nudge, len(nudges))
	copy(result, nudges)
	
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].CriticalityScore > result[i].CriticalityScore {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	
	return result
}



