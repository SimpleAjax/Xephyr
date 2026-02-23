package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/xephyr-ai/xephyr-backend/internal/dto"
)

// AssignmentService defines the interface for assignment-related operations
type AssignmentService interface {
	// GetAssignmentSuggestions returns assignment suggestions for a task
	GetAssignmentSuggestions(ctx context.Context, taskID string, limit int, orgID string) (*dto.AssignmentSuggestionsResponse, error)

	// AssignTask assigns a task to a person
	AssignTask(ctx context.Context, taskID string, req dto.AssignTaskRequest, orgID string, assignedBy uuid.UUID) (*dto.AssignTaskResponse, error)

	// AutoAssignTask auto-assigns a task based on strategy
	AutoAssignTask(ctx context.Context, taskID string, req dto.AutoAssignTaskRequest, orgID string, assignedBy uuid.UUID) (*dto.AssignTaskResponse, error)

	// CheckCompatibility checks person-task compatibility
	CheckCompatibility(ctx context.Context, taskID string, personID string, orgID string) (*dto.AssignmentCompatibilityResponse, error)

	// BulkReassign performs bulk reassignment
	BulkReassign(ctx context.Context, req dto.BulkReassignRequest, orgID string, performedBy uuid.UUID) (*dto.BulkReassignResponse, error)
}

// DummyAssignmentService is a placeholder implementation of AssignmentService
type DummyAssignmentService struct{}

// NewDummyAssignmentService creates a new dummy assignment service
func NewDummyAssignmentService() AssignmentService {
	return &DummyAssignmentService{}
}

// GetAssignmentSuggestions returns dummy suggestions
func (s *DummyAssignmentService) GetAssignmentSuggestions(ctx context.Context, taskID string, limit int, orgID string) (*dto.AssignmentSuggestionsResponse, error) {
	return &dto.AssignmentSuggestionsResponse{
		TaskID:         taskID,
		TaskTitle:      "Checkout Flow Implementation",
		RequiredSkills: []string{"skill-react", "skill-ts", "skill-node"},
		Candidates: []dto.AssignmentCandidate{
			{
				Rank:  1,
				Score: 92,
				Person: dto.PersonInfo{
					ID:        "user-mike",
					Name:      "Mike Rodriguez",
					AvatarURL: "https://example.com/avatar.jpg",
				},
				Breakdown: dto.CandidateBreakdown{
					SkillMatch:      38,
					Availability:    30,
					WorkloadBalance: 18,
					PastPerformance: 6,
				},
				SkillMatchDetails: []dto.SkillMatchDetail{
					{SkillID: "skill-react", Required: true, HasSkill: true, Proficiency: 4, MatchScore: 10},
					{SkillID: "skill-ts", Required: true, HasSkill: true, Proficiency: 4, MatchScore: 10},
					{SkillID: "skill-node", Required: true, HasSkill: true, Proficiency: 3, MatchScore: 8},
				},
				ContextSwitchAnalysis: dto.ContextSwitchAnalysis{
					ActiveProjects:  2,
					CurrentWorkload: 40,
					SwitchPenalty:   5,
					RiskLevel:       "low",
				},
				Warnings:      []string{},
				AIExplanation: "Mike is the best match with 92% compatibility. He has all required skills at high proficiency.",
			},
			{
				Rank:  2,
				Score: 78,
				Person: dto.PersonInfo{
					ID:   "user-alex",
					Name: "Alex Thompson",
				},
				Breakdown: dto.CandidateBreakdown{
					SkillMatch:      30,
					Availability:    10,
					WorkloadBalance: 10,
					PastPerformance: 8,
				},
				Warnings: []string{
					"Currently at 110% allocation",
				},
				AIExplanation: "Alex has strong frontend skills but is currently overallocated.",
			},
		},
		UnassignableReason: nil,
	}, nil
}

// AssignTask performs dummy assignment
func (s *DummyAssignmentService) AssignTask(ctx context.Context, taskID string, req dto.AssignTaskRequest, orgID string, assignedBy uuid.UUID) (*dto.AssignTaskResponse, error) {
	return &dto.AssignTaskResponse{
		TaskID: taskID,
		AssignedTo: dto.AssigneeInfo{
			PersonID: req.PersonID,
			Name:     "Mike Rodriguez",
		},
		PreviousAssignee: nil,
		Assignment: dto.AssignmentInfo{
			AssignedAt: time.Now().UTC(),
			AssignedBy: assignedBy.String(),
		},
		Impact: dto.AssignmentImpact{
			WorkloadUpdated:   true,
			NewAllocation:     65,
			NudgesGenerated:   []string{},
			NotificationsSent: []string{req.PersonID},
		},
	}, nil
}

// AutoAssignTask performs dummy auto-assignment
func (s *DummyAssignmentService) AutoAssignTask(ctx context.Context, taskID string, req dto.AutoAssignTaskRequest, orgID string, assignedBy uuid.UUID) (*dto.AssignTaskResponse, error) {
	return s.AssignTask(ctx, taskID, dto.AssignTaskRequest{
		PersonID: "user-mike",
		Note:     "Auto-assigned using " + req.Strategy + " strategy",
	}, orgID, assignedBy)
}

// CheckCompatibility returns dummy compatibility
func (s *DummyAssignmentService) CheckCompatibility(ctx context.Context, taskID string, personID string, orgID string) (*dto.AssignmentCompatibilityResponse, error) {
	return &dto.AssignmentCompatibilityResponse{
		TaskID:       taskID,
		PersonID:     personID,
		PersonName:   "Mike Rodriguez",
		Score:        92,
		IsCompatible: true,
		Breakdown: dto.CompatibilityBreakdown{
			SkillMatch:      38,
			Availability:    30,
			WorkloadBalance: 18,
			PastPerformance: 6,
		},
		Warnings:      []string{},
		AIExplanation: "Mike is highly compatible with this task.",
	}, nil
}

// BulkReassign performs dummy bulk reassignment
func (s *DummyAssignmentService) BulkReassign(ctx context.Context, req dto.BulkReassignRequest, orgID string, performedBy uuid.UUID) (*dto.BulkReassignResponse, error) {
	results := make([]dto.ReassignmentResult, len(req.Reassignments))
	for i, item := range req.Reassignments {
		results[i] = dto.ReassignmentResult{
			TaskID:   item.TaskID,
			Status:   "success",
			FromUser: item.FromPersonID,
			ToUser:   item.ToPersonID,
		}
	}

	return &dto.BulkReassignResponse{
		Processed:   len(req.Reassignments),
		Succeeded:   len(req.Reassignments),
		Failed:      0,
		Results:     results,
		CompletedAt: time.Now().UTC(),
	}, nil
}
