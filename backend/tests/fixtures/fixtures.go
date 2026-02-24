package fixtures

import (
	"time"

	"github.com/google/uuid"
	"github.com/SimpleAjax/Xephyr/internal/models"
)

// parseUUID safely parses a UUID string or generates a deterministic UUID from any string
func parseUUID(s string) uuid.UUID {
	if id, err := uuid.Parse(s); err == nil {
		return id
	}
	// Generate deterministic UUID from string
	return uuid.NewMD5(uuid.NameSpaceOID, []byte(s))
}

// ===== Test Data Builders =====

type UserBuilder struct {
	user models.User
}

func NewUser() *UserBuilder {
	return &UserBuilder{
		user: models.User{
			BaseModel: models.BaseModel{
				ID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			},
			Email:      "test@example.com",
			Name:       "Test User",
			Timezone:   "America/New_York",
			HourlyRate: 75.0,
			IsActive:   true,
		},
	}
}

func (b *UserBuilder) WithID(id string) *UserBuilder {
	b.user.ID = parseUUID(id)
	return b
}

func (b *UserBuilder) WithName(name string) *UserBuilder {
	b.user.Name = name
	return b
}

func (b *UserBuilder) WithEmail(email string) *UserBuilder {
	b.user.Email = email
	return b
}

func (b *UserBuilder) WithHourlyRate(rate float64) *UserBuilder {
	b.user.HourlyRate = rate
	return b
}

func (b *UserBuilder) Build() models.User {
	return b.user
}

type ProjectBuilder struct {
	project models.Project
}

func NewProject() *ProjectBuilder {
	now := time.Now()
	endDate := now.Add(90 * 24 * time.Hour)
	
	return &ProjectBuilder{
		project: models.Project{
			BaseModel: models.BaseModel{
				ID: uuid.MustParse("660e8400-e29b-41d4-a716-446655440001"),
			},
			OrganizationID: uuid.MustParse("770e8400-e29b-41d4-a716-446655440002"),
			Name:           "Test Project",
			Description:    "A test project for testing",
			Status:         models.ProjectActive,
			Priority:       75,
			HealthScore:    80,
			Progress:       40,
			StartDate:      &now,
			TargetEndDate:  &endDate,
		},
	}
}

func (b *ProjectBuilder) WithID(id string) *ProjectBuilder {
	b.project.ID = parseUUID(id)
	return b
}

func (b *ProjectBuilder) WithName(name string) *ProjectBuilder {
	b.project.Name = name
	return b
}

func (b *ProjectBuilder) WithPriority(priority int) *ProjectBuilder {
	b.project.Priority = priority
	return b
}

func (b *ProjectBuilder) WithHealthScore(score int) *ProjectBuilder {
	b.project.HealthScore = score
	return b
}

func (b *ProjectBuilder) WithProgress(progress int) *ProjectBuilder {
	b.project.Progress = progress
	return b
}

func (b *ProjectBuilder) WithDates(start, end time.Time) *ProjectBuilder {
	b.project.StartDate = &start
	b.project.TargetEndDate = &end
	return b
}

func (b *ProjectBuilder) WithStatus(status models.ProjectStatus) *ProjectBuilder {
	b.project.Status = status
	return b
}

func (b *ProjectBuilder) Build() models.Project {
	return b.project
}

type TaskBuilder struct {
	task models.Task
}

func NewTask() *TaskBuilder {
	now := time.Now()
	dueDate := now.Add(14 * 24 * time.Hour)
	
	return &TaskBuilder{
		task: models.Task{
			BaseModel: models.BaseModel{
				ID: uuid.MustParse("880e8400-e29b-41d4-a716-446655440003"),
			},
			ProjectID:      uuid.MustParse("660e8400-e29b-41d4-a716-446655440001"),
			Title:          "Test Task",
			Description:    "A test task for testing",
			Status:         models.TaskStatusBacklog,
			Priority:       models.TaskPriorityMedium,
			PriorityScore:  50,
			BusinessValue:  60,
			EstimatedHours: 40,
			DueDate:        &dueDate,
			HierarchyLevel: 1,
			IsCriticalPath: false,
			RiskScore:      30,
		},
	}
}

func (b *TaskBuilder) WithID(id string) *TaskBuilder {
	b.task.ID = parseUUID(id)
	return b
}

func (b *TaskBuilder) WithProject(projectID string) *TaskBuilder {
	b.task.ProjectID = parseUUID(projectID)
	return b
}

func (b *TaskBuilder) WithTitle(title string) *TaskBuilder {
	b.task.Title = title
	return b
}

func (b *TaskBuilder) WithStatus(status models.TaskStatus) *TaskBuilder {
	b.task.Status = status
	return b
}

func (b *TaskBuilder) WithPriority(priority models.TaskPriority) *TaskBuilder {
	b.task.Priority = priority
	return b
}

func (b *TaskBuilder) WithPriorityScore(score int) *TaskBuilder {
	b.task.PriorityScore = score
	return b
}

func (b *TaskBuilder) WithBusinessValue(value int) *TaskBuilder {
	b.task.BusinessValue = value
	return b
}

func (b *TaskBuilder) WithEstimatedHours(hours float64) *TaskBuilder {
	b.task.EstimatedHours = hours
	return b
}

func (b *TaskBuilder) WithActualHours(hours float64) *TaskBuilder {
	b.task.ActualHours = hours
	return b
}

func (b *TaskBuilder) WithDueDate(date time.Time) *TaskBuilder {
	b.task.DueDate = &date
	return b
}

func (b *TaskBuilder) WithAssignee(userID string) *TaskBuilder {
	uid := parseUUID(userID)
	b.task.AssigneeID = &uid
	return b
}

func (b *TaskBuilder) OnCriticalPath() *TaskBuilder {
	b.task.IsCriticalPath = true
	return b
}

func (b *TaskBuilder) IsMilestone() *TaskBuilder {
	b.task.IsMilestone = true
	return b
}

func (b *TaskBuilder) WithDependencies(deps []string) *TaskBuilder {
	b.task.Dependencies = make([]models.TaskDependency, len(deps))
	for i, dep := range deps {
		b.task.Dependencies[i] = models.TaskDependency{
			DependsOnTaskID: uuid.MustParse(dep),
		}
	}
	return b
}

func (b *TaskBuilder) Build() models.Task {
	return b.task
}

type SkillBuilder struct {
	skill models.Skill
}

func NewSkill() *SkillBuilder {
	return &SkillBuilder{
		skill: models.Skill{
			BaseModel: models.BaseModel{
				ID: uuid.MustParse("990e8400-e29b-41d4-a716-446655440004"),
			},
			Name:     "Test Skill",
			Category: "Development",
		},
	}
}

func (b *SkillBuilder) WithID(id string) *SkillBuilder {
	b.skill.ID = parseUUID(id)
	return b
}

func (b *SkillBuilder) WithName(name string) *SkillBuilder {
	b.skill.Name = name
	return b
}

func (b *SkillBuilder) WithCategory(category string) *SkillBuilder {
	b.skill.Category = category
	return b
}

func (b *SkillBuilder) Build() models.Skill {
	return b.skill
}

// ===== Predefined Test Teams =====

func CreateTestTeam() []models.User {
	return []models.User{
		NewUser().WithID("user-sarah").WithName("Sarah Chen").WithEmail("sarah@example.com").WithHourlyRate(85).Build(),
		NewUser().WithID("user-mike").WithName("Mike Rodriguez").WithEmail("mike@example.com").WithHourlyRate(95).Build(),
		NewUser().WithID("user-alex").WithName("Alex Thompson").WithEmail("alex@example.com").WithHourlyRate(75).Build(),
		NewUser().WithID("user-emma").WithName("Emma Wilson").WithEmail("emma@example.com").WithHourlyRate(80).Build(),
		NewUser().WithID("user-james").WithName("James Kim").WithEmail("james@example.com").WithHourlyRate(90).Build(),
		NewUser().WithID("user-lisa").WithName("Lisa Park").WithEmail("lisa@example.com").WithHourlyRate(70).Build(),
		NewUser().WithID("user-david").WithName("David Martinez").WithEmail("david@example.com").WithHourlyRate(78).Build(),
		NewUser().WithID("user-rachel").WithName("Rachel Green").WithEmail("rachel@example.com").WithHourlyRate(82).Build(),
	}
}

// ===== Predefined Test Skills =====

func CreateTestSkills() []models.Skill {
	return []models.Skill{
		NewSkill().WithID("skill-react").WithName("React").WithCategory("Frontend").Build(),
		NewSkill().WithID("skill-ts").WithName("TypeScript").WithCategory("Frontend").Build(),
		NewSkill().WithID("skill-next").WithName("Next.js").WithCategory("Frontend").Build(),
		NewSkill().WithID("skill-node").WithName("Node.js").WithCategory("Backend").Build(),
		NewSkill().WithID("skill-python").WithName("Python").WithCategory("Backend").Build(),
		NewSkill().WithID("skill-postgres").WithName("PostgreSQL").WithCategory("Backend").Build(),
		NewSkill().WithID("skill-graphql").WithName("GraphQL").WithCategory("Backend").Build(),
		NewSkill().WithID("skill-figma").WithName("Figma").WithCategory("Design").Build(),
		NewSkill().WithID("skill-ui").WithName("UI Design").WithCategory("Design").Build(),
		NewSkill().WithID("skill-aws").WithName("AWS").WithCategory("DevOps").Build(),
		NewSkill().WithID("skill-docker").WithName("Docker").WithCategory("DevOps").Build(),
		NewSkill().WithID("skill-threejs").WithName("Three.js").WithCategory("Frontend").Build(),
	}
}

// ===== Test User Skills Mapping =====

func CreateTestUserSkills() map[string][]models.UserSkill {
	return map[string][]models.UserSkill{
		"user-sarah": {
			{SkillID: parseUUID("skill-figma"), Proficiency: 3},
		},
		"user-mike": {
			{SkillID: parseUUID("skill-react"), Proficiency: 4},
			{SkillID: parseUUID("skill-ts"), Proficiency: 4},
			{SkillID: parseUUID("skill-node"), Proficiency: 3},
			{SkillID: parseUUID("skill-aws"), Proficiency: 4},
		},
		"user-alex": {
			{SkillID: parseUUID("skill-react"), Proficiency: 3},
			{SkillID: parseUUID("skill-ts"), Proficiency: 3},
			{SkillID: parseUUID("skill-next"), Proficiency: 4},
		},
		"user-emma": {
			{SkillID: parseUUID("skill-figma"), Proficiency: 4},
			{SkillID: parseUUID("skill-ui"), Proficiency: 4},
		},
		"user-james": {
			{SkillID: parseUUID("skill-node"), Proficiency: 4},
			{SkillID: parseUUID("skill-python"), Proficiency: 4},
			{SkillID: parseUUID("skill-postgres"), Proficiency: 4},
			{SkillID: parseUUID("skill-graphql"), Proficiency: 3},
		},
		"user-lisa": {
			{SkillID: parseUUID("skill-react"), Proficiency: 3},
			{SkillID: parseUUID("skill-ts"), Proficiency: 2},
		},
		"user-david": {
			{SkillID: parseUUID("skill-docker"), Proficiency: 4},
			{SkillID: parseUUID("skill-aws"), Proficiency: 3},
		},
		"user-rachel": {
			{SkillID: parseUUID("skill-figma"), Proficiency: 3},
			{SkillID: parseUUID("skill-ui"), Proficiency: 3},
		},
	}
}

// ===== Predefined Test Projects =====

func CreateTestProjects() []models.Project {
	now := time.Now()
	
	return []models.Project{
		NewProject().
			WithID("proj-ecommerce").
			WithName("E-Commerce Platform").
			WithPriority(95).
			WithHealthScore(72).
			WithProgress(45).
			WithDates(now.Add(-30*24*time.Hour), now.Add(60*24*time.Hour)).
			Build(),
		NewProject().
			WithID("proj-mobile").
			WithName("Fitness App").
			WithPriority(88).
			WithHealthScore(45).
			WithProgress(25).
			WithDates(now.Add(-15*24*time.Hour), now.Add(75*24*time.Hour)).
			Build(),
		NewProject().
			WithID("proj-saas").
			WithName("SaaS Dashboard").
			WithPriority(75).
			WithHealthScore(85).
			WithProgress(60).
			WithDates(now.Add(-10*24*time.Hour), now.Add(50*24*time.Hour)).
			Build(),
	}
}

// ===== Predefined Test Tasks =====

func CreateTestTasks() []models.Task {
	now := time.Now()
	
	return []models.Task{
		// E-Commerce Tasks
		NewTask().
			WithID("task-ec-1").
			WithProject("proj-ecommerce").
			WithTitle("Design System").
			WithStatus(models.TaskStatusDone).
			WithPriority(models.TaskPriorityHigh).
			WithPriorityScore(90).
			WithBusinessValue(95).
			WithEstimatedHours(40).
			WithActualHours(42).
			OnCriticalPath().
			Build(),
		NewTask().
			WithID("task-ec-2").
			WithProject("proj-ecommerce").
			WithTitle("Backend API").
			WithStatus(models.TaskStatusInProgress).
			WithPriority(models.TaskPriorityCritical).
			WithPriorityScore(95).
			WithBusinessValue(100).
			WithEstimatedHours(80).
			WithActualHours(24).
			WithAssignee("user-james").
			OnCriticalPath().
			Build(),
		NewTask().
			WithID("task-ec-3").
			WithProject("proj-ecommerce").
			WithTitle("Frontend Catalog").
			WithStatus(models.TaskStatusInProgress).
			WithPriority(models.TaskPriorityHigh).
			WithPriorityScore(85).
			WithBusinessValue(90).
			WithEstimatedHours(60).
			WithAssignee("user-alex").
			OnCriticalPath().
			Build(),
		NewTask().
			WithID("task-ec-4").
			WithProject("proj-ecommerce").
			WithTitle("Checkout Flow").
			WithStatus(models.TaskStatusReady).
			WithPriority(models.TaskPriorityCritical).
			WithPriorityScore(95).
			WithBusinessValue(100).
			WithEstimatedHours(50).
			WithDueDate(now.Add(42 * 24 * time.Hour)).
			OnCriticalPath().
			Build(),
		NewTask().
			WithID("task-ec-5").
			WithProject("proj-ecommerce").
			WithTitle("Admin Dashboard").
			WithStatus(models.TaskStatusBacklog).
			WithPriority(models.TaskPriorityMedium).
			WithPriorityScore(70).
			WithBusinessValue(80).
			WithEstimatedHours(55).
			WithAssignee("user-lisa").
			Build(),
		
		// Fitness App Tasks
		NewTask().
			WithID("task-fit-1").
			WithProject("proj-mobile").
			WithTitle("Mobile App UI").
			WithStatus(models.TaskStatusInProgress).
			WithPriority(models.TaskPriorityHigh).
			WithPriorityScore(88).
			WithBusinessValue(90).
			WithEstimatedHours(70).
			WithAssignee("user-emma").
			OnCriticalPath().
			Build(),
		NewTask().
			WithID("task-fit-2").
			WithProject("proj-mobile").
			WithTitle("Workout Tracking API").
			WithStatus(models.TaskStatusInProgress).
			WithPriority(models.TaskPriorityHigh).
			WithPriorityScore(85).
			WithBusinessValue(85).
			WithEstimatedHours(65).
			WithAssignee("user-james").
			OnCriticalPath().
			Build(),
		NewTask().
			WithID("task-fit-3").
			WithProject("proj-mobile").
			WithTitle("Social Features").
			WithStatus(models.TaskStatusBacklog).
			WithPriority(models.TaskPriorityMedium).
			WithPriorityScore(60).
			WithBusinessValue(70).
			WithEstimatedHours(45).
			Build(),
		NewTask().
			WithID("task-fit-4").
			WithProject("proj-mobile").
			WithTitle("3D Exercise Animations").
			WithStatus(models.TaskStatusBacklog).
			WithPriority(models.TaskPriorityMedium).
			WithPriorityScore(65).
			WithBusinessValue(75).
			WithEstimatedHours(80).
			Build(),
	}
}

// ===== Test Workload Data =====

func CreateTestWorkloadData() []models.WorkloadEntry {
	now := time.Now()
	weekStart := now.Truncate(7 * 24 * time.Hour)
	
	return []models.WorkloadEntry{
		{
			BaseModel:            models.BaseModel{ID: uuid.New()},
			UserID:               parseUUID("user-emma"),
			WeekStart:            weekStart,
			AllocationPercentage: 125,
			AssignedTasks:        4,
			TotalEstimatedHours:  75,
			AvailableHours:       0,
		},
		{
			BaseModel:            models.BaseModel{ID: uuid.New()},
			UserID:               parseUUID("user-james"),
			WeekStart:            weekStart,
			AllocationPercentage: 95,
			AssignedTasks:        2,
			TotalEstimatedHours:  60,
			AvailableHours:       8,
		},
		{
			BaseModel:            models.BaseModel{ID: uuid.New()},
			UserID:               parseUUID("user-alex"),
			WeekStart:            weekStart,
			AllocationPercentage: 110,
			AssignedTasks:        3,
			TotalEstimatedHours:  65,
			AvailableHours:       0,
		},
		{
			BaseModel:            models.BaseModel{ID: uuid.New()},
			UserID:               parseUUID("user-lisa"),
			WeekStart:            weekStart,
			AllocationPercentage: 70,
			AssignedTasks:        2,
			TotalEstimatedHours:  35,
			AvailableHours:       20,
		},
		{
			BaseModel:            models.BaseModel{ID: uuid.New()},
			UserID:               parseUUID("user-rachel"),
			WeekStart:            weekStart,
			AllocationPercentage: 60,
			AssignedTasks:        2,
			TotalEstimatedHours:  30,
			AvailableHours:       25,
		},
		{
			BaseModel:            models.BaseModel{ID: uuid.New()},
			UserID:               parseUUID("user-david"),
			WeekStart:            weekStart,
			AllocationPercentage: 50,
			AssignedTasks:        1,
			TotalEstimatedHours:  25,
			AvailableHours:       30,
		},
		{
			BaseModel:            models.BaseModel{ID: uuid.New()},
			UserID:               parseUUID("user-mike"),
			WeekStart:            weekStart,
			AllocationPercentage: 40,
			AssignedTasks:        1,
			TotalEstimatedHours:  20,
			AvailableHours:       35,
		},
	}
}

// Ptr returns a pointer to the given value
func Ptr[T any](v T) *T {
	return &v
}
