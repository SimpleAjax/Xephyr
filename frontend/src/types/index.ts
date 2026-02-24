// User Types
export interface User {
  id: string;
  email: string;
  name: string;
  avatarUrl?: string;
  role: 'admin' | 'pm' | 'member';
  hourlyRate?: number;
  timezone: string;
}

// Organization Types
export interface Organization {
  id: string;
  name: string;
  slug: string;
  plan: 'free' | 'pro' | 'enterprise';
}

// Skill Types
export interface Skill {
  id: string;
  name: string;
  category: string;
}

export interface PersonSkill {
  skillId: string;
  proficiency: 1 | 2 | 3 | 4;
}

// Project Types
export interface Project {
  id: string;
  name: string;
  description: string;
  status: 'active' | 'paused' | 'completed' | 'archived';
  priority: number;
  startDate: Date;
  targetEndDate: Date;
  healthScore: number;
  progress: number;
}

// Task Types
export type TaskStatus = 'backlog' | 'ready' | 'in_progress' | 'review' | 'done';
export type TaskPriority = 'low' | 'medium' | 'high' | 'critical';

export interface Task {
  id: string;
  projectId: string;
  parentTaskId?: string;
  hierarchyLevel: number;
  title: string;
  description: string;
  status: TaskStatus;
  priority: TaskPriority;
  priorityScore: number;
  businessValue: number;
  estimatedHours: number;
  actualHours?: number;
  startDate?: Date;
  dueDate?: Date;
  assigneeId?: string;
  requiredSkills: string[];
  isMilestone: boolean;
  isCriticalPath: boolean;
  riskScore?: number;
  dependencies: string[];
  blockedBy: string[];
}

// Dependency Types
export interface TaskDependency {
  id: string;
  taskId: string;
  dependsOnTaskId: string;
  dependencyType: 'finish_to_start' | 'start_to_start' | 'finish_to_finish';
  lagHours: number;
}

// Nudge Types
export type NudgeType = 'overload' | 'delay_risk' | 'skill_gap' | 'unassigned' | 'blocked' | 'conflict' | 'dependency_block';
export type NudgeSeverity = 'low' | 'medium' | 'high';
export type NudgeStatus = 'unread' | 'read' | 'dismissed' | 'acted';

export interface Nudge {
  id: string;
  type: NudgeType;
  severity: NudgeSeverity;
  title: string;
  description: string;
  relatedProjectId?: string;
  relatedTaskId?: string;
  relatedPersonId?: string;
  aiExplanation: string;
  suggestedAction?: string;
  status: 'unread' | 'read' | 'dismissed' | 'acted';
  createdAt: Date;
}

// Assignment Types
export interface AssignmentSuggestion {
  personId: string;
  score: number;
  reasons: string[];
  warnings?: string[];
}

// Workload Types
export interface WorkloadEntry {
  personId: string;
  personName: string;
  allocationPercentage: number;
  assignedTasks: number;
  totalEstimatedHours: number;
  availabilityThisWeek: number;
  availabilityNextWeek: number;
}

// Scenario Types
export interface Scenario {
  id: string;
  title: string;
  description: string;
  changeType: 'employee_leave' | 'scope_change' | 'reallocation' | 'priority_shift';
  status: 'pending' | 'approved' | 'rejected' | 'modified';
  proposedChanges: Record<string, any>;
  impactAnalysis?: ImpactAnalysis;
  createdAt: Date;
}

export interface ImpactAnalysis {
  affectedProjects: string[];
  affectedTasks: string[];
  delayHoursTotal: number;
  costImpact: number;
  recommendations: string[];
}
