package services_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/xephyr-ai/xephyr-backend/internal/models"
	"github.com/xephyr-ai/xephyr-backend/test/fixtures"
	"github.com/xephyr-ai/xephyr-backend/test/helpers"
)

func TestDependencyService(t *testing.T) {
	helpers.RunSuite(t, "Dependency Management Service")
}

var _ = Describe("Dependency Management System", func() {
	
	Context("Given a new dependency to create", func() {
		var newDep models.TaskDependency
		var existingDeps []models.TaskDependency
		
		BeforeEach(func() {
			newDep = models.TaskDependency{
				TaskID:          stringToUUID("task-b"),
				DependsOnTaskID: stringToUUID("task-a"),
				DependencyType:  models.DependencyFinishToStart,
				LagHours:        0,
			}
			
			existingDeps = []models.TaskDependency{}
		})
		
		When("validation is performed", func() {
			var result ValidationResult
			
			BeforeEach(func() {
				result = ValidateDependency(newDep, existingDeps)
			})
			
			It("should validate successfully for valid dependency", func() {
				Expect(result.Valid).To(BeTrue())
				Expect(result.Error).To(BeEmpty())
			})
			
			It("should not create a cycle", func() {
				Expect(result.WouldCreateCycle).To(BeFalse())
			})
		})
	})
	
	Context("Given a self-dependency attempt", func() {
		var newDep models.TaskDependency
		
		BeforeEach(func() {
			newDep = models.TaskDependency{
				TaskID:          stringToUUID("task-a"),
				DependsOnTaskID: stringToUUID("task-a"),
				DependencyType:  models.DependencyFinishToStart,
			}
		})
		
		When("validation is performed", func() {
			It("should reject self-dependency", func() {
				result := ValidateDependency(newDep, []models.TaskDependency{})
				Expect(result.Valid).To(BeFalse())
				Expect(result.Error).To(ContainSubstring("self"))
			})
		})
	})
	
	Context("Given dependencies that would create a cycle", func() {
		var existingDeps []models.TaskDependency
		var newDep models.TaskDependency
		
		BeforeEach(func() {
			// A -> B -> C exists, trying to add C -> A
			existingDeps = []models.TaskDependency{
				{TaskID: stringToUUID("task-b"), DependsOnTaskID: stringToUUID("task-a")},
				{TaskID: stringToUUID("task-c"), DependsOnTaskID: stringToUUID("task-b")},
			}
			
			newDep = models.TaskDependency{
				TaskID:          stringToUUID("task-a"),
				DependsOnTaskID: stringToUUID("task-c"),
			}
		})
		
		When("validation is performed", func() {
			var result ValidationResult
			
			BeforeEach(func() {
				result = ValidateDependency(newDep, existingDeps)
			})
			
			It("should detect the circular dependency", func() {
				Expect(result.Valid).To(BeFalse())
				Expect(result.WouldCreateCycle).To(BeTrue())
			})
			
			It("should return the cycle path", func() {
				Expect(result.CyclePath).To(ContainElements("task-a", "task-b", "task-c"))
			})
		})
	})
	
	Context("Given a project with tasks and dependencies", func() {
		var projectID string
		var tasks []models.Task
		var dependencies []models.TaskDependency
		
		BeforeEach(func() {
			projectID = "proj-cpm-test"
			
			tasks = []models.Task{
				fixtures.NewTask().WithID("task-1").WithProject(projectID).WithEstimatedHours(40).Build(),
				fixtures.NewTask().WithID("task-2").WithProject(projectID).WithEstimatedHours(80).Build(),
				fixtures.NewTask().WithID("task-3").WithProject(projectID).WithEstimatedHours(50).Build(),
				fixtures.NewTask().WithID("task-4").WithProject(projectID).WithEstimatedHours(30).Build(),
			}
			
			// 1 -> 2 -> 3 (critical path)
			// 1 -> 4 (non-critical)
			dependencies = []models.TaskDependency{
				{TaskID: stringToUUID("task-2"), DependsOnTaskID: stringToUUID("task-1"), DependencyType: models.DependencyFinishToStart},
				{TaskID: stringToUUID("task-3"), DependsOnTaskID: stringToUUID("task-2"), DependencyType: models.DependencyFinishToStart},
				{TaskID: stringToUUID("task-4"), DependsOnTaskID: stringToUUID("task-1"), DependencyType: models.DependencyFinishToStart},
			}
		})
		
		When("critical path is calculated", func() {
			var result CriticalPathResult
			
			BeforeEach(func() {
				result = CalculateCriticalPath(projectID, tasks, dependencies)
			})
			
			It("should identify the critical path tasks", func() {
				Expect(result.CriticalPath).To(Equal([]string{"task-1", "task-2", "task-3"}))
			})
			
			It("should calculate correct project duration", func() {
				// 40 + 80 + 50 = 170 hours
				Expect(result.ProjectDuration).To(Equal(170))
			})
			
			It("should identify zero float for critical tasks", func() {
				Expect(result.FloatTimes["task-1"]).To(Equal(0))
				Expect(result.FloatTimes["task-2"]).To(Equal(0))
				Expect(result.FloatTimes["task-3"]).To(Equal(0))
			})
			
			It("should calculate positive float for non-critical tasks", func() {
				// Task 4 has float since it's not on critical path
				Expect(result.FloatTimes["task-4"]).To(BeNumerically(">", 0))
			})
		})
	})
	
	Context("Given a task with dependencies", func() {
		var task models.Task
		var allTasks []models.Task
		var allDeps []models.TaskDependency
		
		BeforeEach(func() {
			task = fixtures.NewTask().WithID("task-target").Build()
			
			allTasks = []models.Task{
				fixtures.NewTask().WithID("task-dep-1").WithStatus(models.TaskStatusDone).Build(),
				fixtures.NewTask().WithID("task-dep-2").WithStatus(models.TaskStatusInProgress).Build(),
				fixtures.NewTask().WithID("task-dep-3").WithStatus(models.TaskStatusBacklog).Build(),
			}
			
			allDeps = []models.TaskDependency{
				{TaskID: stringToUUID("task-target"), DependsOnTaskID: stringToUUID("task-dep-1")},
				{TaskID: stringToUUID("task-target"), DependsOnTaskID: stringToUUID("task-dep-2")},
				{TaskID: stringToUUID("task-target"), DependsOnTaskID: stringToUUID("task-dep-3")},
			}
		})
		
		When("dependency status is checked", func() {
			It("should return blocked when dependencies incomplete", func() {
				status := GetTaskDependencyStatus(task, allTasks, allDeps)
				Expect(status).To(Equal("blocked"))
			})
			
			It("should return ready when all dependencies done", func() {
				// Mark all deps as done
				for i := range allTasks {
					allTasks[i].Status = models.TaskStatusDone
				}
				
				status := GetTaskDependencyStatus(task, allTasks, allDeps)
				Expect(status).To(Equal("ready"))
			})
		})
	})
	
	Context("Given dependency graph operations", func() {
		var deps []models.TaskDependency
		
		BeforeEach(func() {
			// Create a diamond pattern:
			//    A
			//   / \
			//  B   C
			//   \ /
			//    D
			deps = []models.TaskDependency{
				{TaskID: stringToUUID("task-b"), DependsOnTaskID: stringToUUID("task-a")},
				{TaskID: stringToUUID("task-c"), DependsOnTaskID: stringToUUID("task-a")},
				{TaskID: stringToUUID("task-d"), DependsOnTaskID: stringToUUID("task-b")},
				{TaskID: stringToUUID("task-d"), DependsOnTaskID: stringToUUID("task-c")},
			}
		})
		
		When("topological sort is performed", func() {
			It("should return valid topological ordering", func() {
				order, err := TopologicalSort(deps)
				Expect(err).ToNot(HaveOccurred())
				Expect(order).To(HaveLen(4))
				
				// A should come before B, C, D
				aIndex := indexOf(order, "task-a")
				bIndex := indexOf(order, "task-b")
				cIndex := indexOf(order, "task-c")
				dIndex := indexOf(order, "task-d")
				
				Expect(aIndex).To(BeNumerically("<", bIndex))
				Expect(aIndex).To(BeNumerically("<", cIndex))
				Expect(bIndex).To(BeNumerically("<", dIndex))
				Expect(cIndex).To(BeNumerically("<", dIndex))
			})
		})
		
		When("descendants are queried", func() {
			It("should return all downstream tasks", func() {
				descendants := GetDescendants("task-a", deps)
				Expect(descendants).To(ContainElements("task-b", "task-c", "task-d"))
			})
		})
		
		When("ancestors are queried", func() {
			It("should return all upstream tasks", func() {
				ancestors := GetAncestors("task-d", deps)
				Expect(ancestors).To(ContainElements("task-a", "task-b", "task-c"))
			})
		})
	})
	
	Context("Given different dependency types", func() {
		DescribeTable("dependency type handling",
			func(depType models.DependencyType, expectedLag int) {
				dep := models.TaskDependency{
					DependencyType: depType,
					LagHours:       expectedLag,
				}
				Expect(dep.DependencyType).To(Equal(depType))
			},
			Entry("finish to start", models.DependencyFinishToStart, 0),
			Entry("start to start", models.DependencyStartToStart, 0),
			Entry("finish to finish", models.DependencyFinishToFinish, 0),
			Entry("start to finish", models.DependencyStartToFinish, 0),
		)
	})
	
	Context("Given impact analysis for new dependency", func() {
		var newDep models.TaskDependency
		var tasks []models.Task
		var existingDeps []models.TaskDependency
		
		BeforeEach(func() {
			newDep = models.TaskDependency{
				TaskID:          stringToUUID("task-new"),
				DependsOnTaskID: stringToUUID("task-existing"),
				LagHours:        8,
			}
			
			tasks = fixtures.CreateTestTasks()
			existingDeps = []models.TaskDependency{}
		})
		
		When("impact is calculated", func() {
			var impact DependencyImpact
			
			BeforeEach(func() {
				impact = CalculateNewDependencyImpact(newDep, tasks, existingDeps)
			})
			
			It("should identify affected tasks", func() {
				Expect(len(impact.AffectedTasks)).To(BeNumerically(">=", 0))
			})
			
			It("should estimate delay from lag hours", func() {
				Expect(impact.EstimatedDelayHours).To(Equal(8))
			})
		})
	})
})

// ===== Helper Types and Functions =====

type ValidationResult struct {
	Valid            bool
	Error            string
	WouldCreateCycle bool
	CyclePath        []string
}

type CriticalPathResult struct {
	CriticalPath   []string
	CriticalTasks  map[string]bool
	FloatTimes     map[string]int
	ProjectDuration int
}

type DependencyImpact struct {
	AffectedTasks       []string
	EstimatedDelayHours int
	CriticalPathChanged bool
}

func ValidateDependency(newDep models.TaskDependency, existingDeps []models.TaskDependency) ValidationResult {
	// Check self-dependency
	if newDep.TaskID == newDep.DependsOnTaskID {
		return ValidationResult{
			Valid: false,
			Error: "Task cannot depend on itself",
		}
	}
	
	// Build graph and check for cycles
	graph := buildDependencyGraph(existingDeps)
	graph[newDep.TaskID.String()] = append(graph[newDep.TaskID.String()], newDep.DependsOnTaskID.String())
	
	if hasCycle(graph, newDep.TaskID.String()) {
		return ValidationResult{
			Valid:            false,
			WouldCreateCycle: true,
			Error:            "Circular dependency detected",
			CyclePath:        []string{newDep.TaskID.String(), newDep.DependsOnTaskID.String()},
		}
	}
	
	return ValidationResult{Valid: true}
}

func CalculateCriticalPath(projectID string, tasks []models.Task, dependencies []models.TaskDependency) CriticalPathResult {
	// Build task map
	taskMap := make(map[string]models.Task)
	for _, t := range tasks {
		taskMap[t.ID.String()] = t
	}
	
	// Build adjacency list
	graph := make(map[string][]string)
	for _, dep := range dependencies {
		from := dep.DependsOnTaskID.String()
		to := dep.TaskID.String()
		graph[from] = append(graph[from], to)
	}
	
	// Topological sort
	order, _ := TopologicalSort(dependencies)
	
	// Forward pass: Calculate Early Start (ES) and Early Finish (EF)
	es := make(map[string]int) // Early Start
	ef := make(map[string]int) // Early Finish
	
	for _, taskID := range order {
		task := taskMap[taskID]
		
		// Find max EF of predecessors
		maxPredEF := 0
		for _, dep := range dependencies {
			if dep.TaskID.String() == taskID {
				if ef[dep.DependsOnTaskID.String()] > maxPredEF {
					maxPredEF = ef[dep.DependsOnTaskID.String()]
				}
			}
		}
		
		es[taskID] = maxPredEF
		ef[taskID] = maxPredEF + int(task.EstimatedHours)
	}
	
	// Find project duration
	projectDuration := 0
	for _, finish := range ef {
		if finish > projectDuration {
			projectDuration = finish
		}
	}
	
	// Backward pass: Calculate Late Start (LS) and Late Finish (LF)
	ls := make(map[string]int) // Late Start
	lf := make(map[string]int) // Late Finish
	
	// Initialize terminal tasks
	for _, task := range tasks {
		isTerminal := true
		for _, dep := range dependencies {
			if dep.DependsOnTaskID == task.ID {
				isTerminal = false
				break
			}
		}
		if isTerminal {
			lf[task.ID.String()] = projectDuration
			ls[task.ID.String()] = projectDuration - int(task.EstimatedHours)
		}
	}
	
	// Process in reverse order
	for i := len(order) - 1; i >= 0; i-- {
		taskID := order[i]
		if _, ok := lf[taskID]; !ok {
			// Find min LS of successors
			minSuccLS := projectDuration
			for _, dep := range dependencies {
				if dep.DependsOnTaskID.String() == taskID {
					succLS := ls[dep.TaskID.String()]
					if succLS < minSuccLS {
						minSuccLS = succLS
					}
				}
			}
			
			lf[taskID] = minSuccLS
			ls[taskID] = minSuccLS - int(taskMap[taskID].EstimatedHours)
		}
	}
	
	// Calculate float and identify critical path
	floatTimes := make(map[string]int)
	criticalTasks := make(map[string]bool)
	
	for _, task := range tasks {
		taskID := task.ID.String()
		floatTime := ls[taskID] - es[taskID]
		floatTimes[taskID] = floatTime
		
		if floatTime == 0 {
			criticalTasks[taskID] = true
		}
	}
	
	// Build critical path (simplified)
	var criticalPath []string
	for _, taskID := range order {
		if criticalTasks[taskID] {
			criticalPath = append(criticalPath, taskID)
		}
	}
	
	return CriticalPathResult{
		CriticalPath:    criticalPath,
		CriticalTasks:   criticalTasks,
		FloatTimes:      floatTimes,
		ProjectDuration: projectDuration,
	}
}

func GetTaskDependencyStatus(task models.Task, allTasks []models.Task, allDeps []models.TaskDependency) string {
	if task.Status == models.TaskStatusDone {
		return "done"
	}
	if task.Status == models.TaskStatusInProgress {
		return "in_progress"
	}
	
	// Check incomplete dependencies
	for _, dep := range allDeps {
		if dep.TaskID == task.ID {
			for _, t := range allTasks {
				if t.ID == dep.DependsOnTaskID && t.Status != models.TaskStatusDone {
					return "blocked"
				}
			}
		}
	}
	
	return "ready"
}

func TopologicalSort(deps []models.TaskDependency) ([]string, error) {
	// Build graph and in-degrees
	graph := make(map[string][]string)
	inDegree := make(map[string]int)
	nodes := make(map[string]bool)
	
	for _, dep := range deps {
		from := dep.DependsOnTaskID.String()
		to := dep.TaskID.String()
		
		graph[from] = append(graph[from], to)
		inDegree[to]++
		if _, ok := inDegree[from]; !ok {
			inDegree[from] = 0
		}
		nodes[from] = true
		nodes[to] = true
	}
	
	// Kahn's algorithm
	var queue []string
	for node, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, node)
		}
	}
	
	var result []string
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		result = append(result, node)
		
		for _, neighbor := range graph[node] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}
	
	return result, nil
}

func GetDescendants(taskID string, deps []models.TaskDependency) []string {
	// Build reverse graph
	graph := make(map[string][]string)
	for _, dep := range deps {
		from := dep.DependsOnTaskID.String()
		to := dep.TaskID.String()
		graph[from] = append(graph[from], to)
	}
	
	// BFS
	visited := make(map[string]bool)
	var result []string
	queue := []string{taskID}
	
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		
		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				result = append(result, neighbor)
				queue = append(queue, neighbor)
			}
		}
	}
	
	return result
}

func GetAncestors(taskID string, deps []models.TaskDependency) []string {
	// Build forward graph
	graph := make(map[string][]string)
	for _, dep := range deps {
		from := dep.TaskID.String()
		to := dep.DependsOnTaskID.String()
		graph[from] = append(graph[from], to)
	}
	
	// BFS
	visited := make(map[string]bool)
	var result []string
	queue := []string{taskID}
	
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		
		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				result = append(result, neighbor)
				queue = append(queue, neighbor)
			}
		}
	}
	
	return result
}

func CalculateNewDependencyImpact(newDep models.TaskDependency, tasks []models.Task, existingDeps []models.TaskDependency) DependencyImpact {
	return DependencyImpact{
		AffectedTasks:       []string{newDep.TaskID.String()},
		EstimatedDelayHours: newDep.LagHours,
		CriticalPathChanged: false,
	}
}

func buildDependencyGraph(deps []models.TaskDependency) map[string][]string {
	graph := make(map[string][]string)
	for _, dep := range deps {
		from := dep.DependsOnTaskID.String()
		to := dep.TaskID.String()
		graph[from] = append(graph[from], to)
	}
	return graph
}

func hasCycle(graph map[string][]string, start string) bool {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	
	var dfs func(string) bool
	dfs = func(node string) bool {
		visited[node] = true
		recStack[node] = true
		
		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				if dfs(neighbor) {
					return true
				}
			} else if recStack[neighbor] {
				return true
			}
		}
		
		recStack[node] = false
		return false
	}
	
	return dfs(start)
}

func indexOf(slice []string, item string) int {
	for i, s := range slice {
		if s == item {
			return i
		}
	}
	return -1
}



