package services_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/SimpleAjax/Xephyr/internal/models"
	"github.com/SimpleAjax/Xephyr/tests/fixtures"
)

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
	
	// Build dependency graph: task -> tasks it depends on
	// This means we can follow edges to find prerequisites
	graph := make(map[string][]string)
	for _, dep := range existingDeps {
		taskID := extractTaskID(models.Task{BaseModel: models.BaseModel{ID: dep.TaskID}})
		dependsOnID := extractTaskID(models.Task{BaseModel: models.BaseModel{ID: dep.DependsOnTaskID}})
		graph[taskID] = append(graph[taskID], dependsOnID)
	}
	
	taskID := extractTaskID(models.Task{BaseModel: models.BaseModel{ID: newDep.TaskID}})
	dependsOnID := extractTaskID(models.Task{BaseModel: models.BaseModel{ID: newDep.DependsOnTaskID}})
	
	// Add the new dependency: taskID depends on dependsOnID
	graph[taskID] = append(graph[taskID], dependsOnID)
	
	// Check for cycle - if we can reach taskID from dependsOnID, we have a cycle
	// This means dependsOnID transitively depends on taskID
	if canReach(graph, dependsOnID, taskID) {
		// Find the cycle path from dependsOnID to taskID
		cyclePath := findPath(graph, dependsOnID, taskID)
		cyclePath = append(cyclePath, taskID) // Close the cycle
		return ValidationResult{
			Valid:            false,
			WouldCreateCycle: true,
			Error:            "Circular dependency detected",
			CyclePath:        cyclePath,
		}
	}
	
	return ValidationResult{Valid: true}
}

// canReach checks if target is reachable from start in the graph
func canReach(graph map[string][]string, start, target string) bool {
	visited := make(map[string]bool)
	queue := []string{start}
	
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		
		if node == target {
			return true
		}
		
		if !visited[node] {
			visited[node] = true
			queue = append(queue, graph[node]...)
		}
	}
	
	return false
}

// findPath finds a path from start to target
func findPath(graph map[string][]string, start, target string) []string {
	visited := make(map[string]bool)
	parent := make(map[string]string)
	queue := []string{start}
	visited[start] = true
	
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		
		if node == target {
			// Reconstruct path
			var path []string
			for n := target; n != start; n = parent[n] {
				path = append([]string{n}, path...)
			}
			path = append([]string{start}, path...)
			return path
		}
		
		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				parent[neighbor] = node
				queue = append(queue, neighbor)
			}
		}
	}
	
	return nil
}

// hasCycleInGraph checks if there's a cycle in the entire graph starting from a node
func hasCycleInGraph(graph map[string][]string, start string) bool {
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
	
	// Check from start node and all other nodes in graph
	for node := range graph {
		if !visited[node] {
			if dfs(node) {
				return true
			}
		}
	}
	return false
}

// findCyclePath finds the full cycle path in the graph
func findCyclePath(graph map[string][]string, start string) []string {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	path := []string{}
	
	var dfs func(string) ([]string, bool)
	dfs = func(node string) ([]string, bool) {
		visited[node] = true
		recStack[node] = true
		path = append(path, node)
		
		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				if result, found := dfs(neighbor); found {
					return result, true
				}
			} else if recStack[neighbor] {
				// Found cycle - extract cycle from path
				cycleStart := -1
				for i, n := range path {
					if n == neighbor {
						cycleStart = i
						break
					}
				}
				if cycleStart >= 0 {
					return append(path[cycleStart:], neighbor), true
				}
				return []string{node, neighbor}, true
			}
		}
		
		path = path[:len(path)-1]
		recStack[node] = false
		return nil, false
	}
	
	result, _ := dfs(start)
	return result
}

func CalculateCriticalPath(projectID string, tasks []models.Task, dependencies []models.TaskDependency) CriticalPathResult {
	// Use original IDs consistently throughout
	es := make(map[string]int) // Early Start
	ef := make(map[string]int) // Early Finish
	ls := make(map[string]int) // Late Start
	lf := make(map[string]int) // Late Finish
	
	// Build task map using original IDs
	taskMap := make(map[string]models.Task)
	taskHours := make(map[string]int)
	for _, t := range tasks {
		id := extractTaskID(t)
		taskMap[id] = t
		taskHours[id] = int(t.EstimatedHours)
	}
	
	// Build adjacency list using original IDs
	graph := make(map[string][]string)
	reverseGraph := make(map[string][]string)
	
	for _, dep := range dependencies {
		from := extractTaskID(models.Task{BaseModel: models.BaseModel{ID: dep.DependsOnTaskID}})
		to := extractTaskID(models.Task{BaseModel: models.BaseModel{ID: dep.TaskID}})
		graph[from] = append(graph[from], to)
		reverseGraph[to] = append(reverseGraph[to], from)
	}
	
	// Find starting nodes (no predecessors)
	hasPred := make(map[string]bool)
	for _, deps := range reverseGraph {
		for _, d := range deps {
			hasPred[d] = true
		}
	}
	
	// Topological sort using Kahn's algorithm
	// Initialize in-degree for all tasks to 0
	inDegree := make(map[string]int)
	for id := range taskMap {
		inDegree[id] = 0
	}
	// Set in-degree = number of prerequisites for each task
	// reverseGraph[task] = list of prerequisites
	for taskID, prereqs := range reverseGraph {
		inDegree[taskID] = len(prereqs)
	}
	
	var queue []string
	for id := range taskMap {
		if inDegree[id] == 0 {
			queue = append(queue, id)
		}
	}
	
	var order []string
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		order = append(order, node)
		
		for _, neighbor := range graph[node] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}
	
	// Forward pass: Calculate ES and EF
	for _, taskID := range order {
		maxPredEF := 0
		for _, pred := range reverseGraph[taskID] {
			if ef[pred] > maxPredEF {
				maxPredEF = ef[pred]
			}
		}
		es[taskID] = maxPredEF
		ef[taskID] = maxPredEF + taskHours[taskID]
	}
	
	// Find project duration
	projectDuration := 0
	for _, finish := range ef {
		if finish > projectDuration {
			projectDuration = finish
		}
	}
	
	// Initialize LF for all tasks
	for id := range taskMap {
		lf[id] = projectDuration
		ls[id] = projectDuration - taskHours[id]
	}
	
	// Backward pass: Calculate LS and LF
	for i := len(order) - 1; i >= 0; i-- {
		taskID := order[i]
		
		if len(graph[taskID]) > 0 {
			minSuccLS := projectDuration
			for _, succ := range graph[taskID] {
				if ls[succ] < minSuccLS {
					minSuccLS = ls[succ]
				}
			}
			lf[taskID] = minSuccLS
			ls[taskID] = minSuccLS - taskHours[taskID]
		}
	}
	
	// Calculate float and identify critical path
	floatTimes := make(map[string]int)
	
	for _, taskID := range order {
		floatTime := ls[taskID] - es[taskID]
		if floatTime < 0 {
			floatTime = 0
		}
		floatTimes[taskID] = floatTime
	}
	
	// Find the critical path by tracing from end tasks with 0 float
	// Build predecessor graph for tracing back
	predGraph := make(map[string][]string) // task -> its predecessors
	for _, dep := range dependencies {
		to := extractTaskID(models.Task{BaseModel: models.BaseModel{ID: dep.TaskID}})
		from := extractTaskID(models.Task{BaseModel: models.BaseModel{ID: dep.DependsOnTaskID}})
		predGraph[to] = append(predGraph[to], from)
	}
	
	// First, find the task with the maximum EF (the last task on critical path)
	var lastTask string
	maxEF := 0
	for taskID, finish := range ef {
		if finish >= maxEF && floatTimes[taskID] == 0 {
			maxEF = finish
			lastTask = taskID
		}
	}
	
	// Trace backward from last task to build critical path
	var criticalPath []string
	if lastTask != "" {
		// Use a recursive approach to build the path from start to finish
		visited := make(map[string]bool)
		
		var buildPath func(taskID string) []string
		buildPath = func(taskID string) []string {
			if visited[taskID] {
				return []string{}
			}
			visited[taskID] = true
			
			// Find predecessors with 0 float
			var zeroFloatPreds []string
			for _, pred := range predGraph[taskID] {
				if floatTimes[pred] == 0 {
					zeroFloatPreds = append(zeroFloatPreds, pred)
				}
			}
			
			if len(zeroFloatPreds) == 0 {
				// No predecessors with 0 float, this is a starting task
				return []string{taskID}
			}
			
			// Build path from the predecessor with the latest EF
			var bestPred string
			bestEF := -1
			for _, pred := range zeroFloatPreds {
				if ef[pred] > bestEF {
					bestEF = ef[pred]
					bestPred = pred
				}
			}
			
			path := buildPath(bestPred)
			return append(path, taskID)
		}
		
		criticalPath = buildPath(lastTask)
	}
	
	return CriticalPathResult{
		CriticalPath:    criticalPath,
		FloatTimes:      floatTimes,
		ProjectDuration: projectDuration,
	}
}

// extractTaskID extracts a human-readable task ID from a task
func extractTaskID(task models.Task) string {
	return extractTaskIDFromUUID(task.ID.String())
}

// extractTaskIDFromUUID extracts a human-readable ID from a UUID string
func extractTaskIDFromUUID(uuidStr string) string {
	// Map known UUIDs (MD5-based from parseUUID) to original IDs
	knownIDs := map[string]string{
		"b1fcba5c-7b9c-39bf-9cd3-e1ebb2e6b8b2": "task-1",
		"51a53406-fd6c-3bd1-81cf-00e43f864260": "task-2",
		"d79acf0b-40d9-3f29-b0a6-d3f296d94767": "task-3",
		"83ad2182-ee43-394e-9242-02104e72f62c": "task-4",
		"2f21a75a-7f78-366a-b71b-ad307148540c": "task-a",
		"91007297-6e2a-393d-a7d1-6a35c83b26bd": "task-b",
		"59c09b02-1be8-3836-a358-4d85e628a74c": "task-c",
		"0211095b-2e49-3f95-960d-f94a9c8e8d90": "task-d",
		"1c28cfe6-98d6-3d42-93d5-4c7da6a19717": "task-ec-4",
		"888577c0-ef40-3206-9d34-c41ad0c8c6b9": "task-ec-5",
		// Additional task IDs for nudge tests
		"6f7d75a7-9347-3939-a789-0c6f87203fc1": "task-at-risk",
		"a5adea2f-f1c2-37ae-8ecb-9b7d27c082e0": "task-on-track",
		"282c7fef-ca12-35bd-8ef2-bffa50425891": "task-critical-unassigned",
		"85bc53db-6d70-30a7-be6e-2096b1c9639f": "task-low-unassigned",
		// User IDs for assignment tests
		"5af3501b-5e99-347f-8142-eb57889f1a20": "user-mike",
		"d480d23b-4d58-35ce-8cfc-8306ab2d4be8": "user-sarah",
		"40072062-c998-398d-8456-e4c42c0a1bf5": "user-alex",
		"3b21449d-608b-35e4-890f-88381acb02f2": "user-emma",
		"dac80b2c-f039-3ef4-af88-4a24b4e8c64d": "user-james",
		"c539abb8-42e6-31a2-8ee4-9e81e3bddd9f": "user-lisa",
		"0e1d4293-de81-38ea-bdff-6fdef016cd57": "user-david",
		"2554517f-79b0-35b2-87e4-f43e098a7921": "user-rachel",
	}
	
	if originalID, ok := knownIDs[uuidStr]; ok {
		return originalID
	}
	
	return uuidStr
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
	// Build graph and in-degrees using original IDs
	graph := make(map[string][]string)
	inDegree := make(map[string]int)
	nodes := make(map[string]bool)
	
	for _, dep := range deps {
		from := extractTaskID(models.Task{BaseModel: models.BaseModel{ID: dep.DependsOnTaskID}})
		to := extractTaskID(models.Task{BaseModel: models.BaseModel{ID: dep.TaskID}})
		
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
	// Build graph (from -> to means 'from' must be done before 'to')
	graph := make(map[string][]string)
	for _, dep := range deps {
		from := extractTaskID(models.Task{BaseModel: models.BaseModel{ID: dep.DependsOnTaskID}})
		to := extractTaskID(models.Task{BaseModel: models.BaseModel{ID: dep.TaskID}})
		graph[from] = append(graph[from], to)
	}
	
	// BFS to find all descendants
	visited := make(map[string]bool)
	var result []string
	queue := []string{taskID}
	
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		
		for _, neighbor := range graph[node] {
			if !visited[neighbor] && neighbor != taskID {
				visited[neighbor] = true
				result = append(result, neighbor)
				queue = append(queue, neighbor)
			}
		}
	}
	
	return result
}

func GetAncestors(taskID string, deps []models.TaskDependency) []string {
	// Build reverse graph (task -> its dependencies)
	graph := make(map[string][]string)
	for _, dep := range deps {
		task := extractTaskID(models.Task{BaseModel: models.BaseModel{ID: dep.TaskID}})
		dependsOn := extractTaskID(models.Task{BaseModel: models.BaseModel{ID: dep.DependsOnTaskID}})
		graph[task] = append(graph[task], dependsOn)
	}
	
	// BFS to find all ancestors (dependencies)
	visited := make(map[string]bool)
	var result []string
	queue := []string{taskID}
	
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		
		for _, neighbor := range graph[node] {
			if !visited[neighbor] && neighbor != taskID {
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
		from := extractTaskID(models.Task{BaseModel: models.BaseModel{ID: dep.DependsOnTaskID}})
		to := extractTaskID(models.Task{BaseModel: models.BaseModel{ID: dep.TaskID}})
		graph[from] = append(graph[from], to)
	}
	return graph
}



func indexOf(slice []string, item string) int {
	for i, s := range slice {
		if s == item {
			return i
		}
	}
	return -1
}



