// API Response Types - Based on API Design Document

export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: ApiError;
  meta?: ResponseMeta;
}

export interface ApiError {
  code: string;
  message: string;
  details?: Record<string, string[]>;
}

export interface ResponseMeta {
  page?: number;
  perPage?: number;
  total?: number;
  timestamp: string;
  requestId: string;
  hasMore?: boolean;
  nextCursor?: string;
}

// ==================== PRIORITY API TYPES ====================

export interface TaskPriority {
  taskId: string;
  universalPriorityScore: number;
  rank: number;
  breakdown?: {
    projectPriority: number;
    businessValue: number;
    deadlineUrgency: number;
    criticalPathWeight: number;
    dependencyImpact: number;
  };
  factors?: {
    projectPriority: number;
    businessValue: number;
    deadlineUrgency: number;
    isOnCriticalPath: boolean;
    blockedTasksCount: number;
    daysUntilDue: number;
  };
  calculatedAt: string;
}

export interface BulkPriorityRequest {
  taskIds: string[];
  includeBreakdown?: boolean;
}

export interface BulkPriorityResponse {
  priorities: Array<{
    taskId: string;
    universalPriorityScore: number;
    rank: number;
  }>;
  sortedOrder: string[];
}

export interface ProjectRanking {
  projectId: string;
  rankings: Array<{
    rank: number;
    taskId: string;
    title: string;
    priorityScore: number;
    status: string;
    assigneeId: string | null;
  }>;
  total: number;
}

export interface RecalculateRequest {
  scope: 'project' | 'organization' | 'task';
  projectId?: string;
  async?: boolean;
}

export interface RecalculateResponse {
  recalculated?: number;
  duration?: string;
  affectedTasks?: string[];
  jobId?: string;
  status?: string;
  estimatedDuration?: string;
}

// ==================== HEALTH API TYPES ====================

export interface PortfolioHealth {
  portfolioHealthScore: number;
  status: 'healthy' | 'caution' | 'at_risk' | 'critical';
  summary: {
    totalProjects: number;
    healthy: number;
    caution: number;
    atRisk: number;
    critical: number;
  };
  projects: Array<{
    projectId: string;
    name: string;
    healthScore: number;
    status: string;
    priority: number;
    progress: number;
    trend: 'improving' | 'stable' | 'worsening';
  }>;
  calculatedAt: string;
}

export interface ProjectHealth {
  projectId: string;
  projectName: string;
  healthScore: number;
  status: string;
  breakdown: {
    scheduleHealth: number;
    completionHealth: number;
    dependencyHealth: number;
    resourceHealth: number;
    criticalPathHealth: number;
  };
  details: {
    schedule: {
      expectedProgress: number;
      actualProgress: number;
      variance: number;
      daysUntilDeadline: number;
    };
    completion: {
      totalTasks: number;
      completed: number;
      inProgress: number;
      completionRate: number;
    };
    dependencies: {
      total: number;
      blocked: number;
      atRisk: number;
    };
    resources: {
      teamSize: number;
      avgAllocation: number;
      overallocated: number;
      underutilized: number;
    };
  };
  trend: {
    direction: 'improving' | 'stable' | 'worsening';
    change: number;
    lastWeekScore: number;
  };
  calculatedAt: string;
}

export interface HealthTrend {
  projectId: string;
  timeRange: string;
  datapoints: Array<{
    date: string;
    healthScore: number;
    scheduleHealth: number;
    completionHealth: number;
  }>;
  trend: {
    slope: number;
    direction: 'improving' | 'declining' | 'stable';
    prediction?: {
      daysUntilCritical?: number;
      confidence: number;
    };
  };
}

// ==================== NUDGE API TYPES ====================

export type NudgeType = 'overload' | 'delay_risk' | 'skill_gap' | 'unassigned' | 'blocked' | 'conflict' | 'dependency_block';
export type NudgeSeverity = 'low' | 'medium' | 'high';
export type NudgeStatus = 'unread' | 'read' | 'dismissed' | 'acted';

export interface Nudge {
  id: string;
  type: NudgeType;
  severity: NudgeSeverity;
  status: NudgeStatus;
  title: string;
  description: string;
  aiExplanation: string;
  suggestedAction?: {
    type: string;
    description: string;
    targetTaskId?: string;
    suggestedAssigneeId?: string;
  };
  relatedEntities?: {
    projectId?: string;
    taskId?: string;
    personId?: string;
  };
  metrics?: {
    allocationPercentage?: number;
    assignedTasks?: number;
    totalHours?: number;
  };
  criticalityScore: number;
  createdAt: string;
  expiresAt?: string;
}

export interface NudgeListResponse {
  nudges: Nudge[];
  summary: {
    total: number;
    unread: number;
    bySeverity: Record<string, number>;
    byType: Record<string, number>;
  };
}

export interface NudgeActionRequest {
  actionType: 'accept_suggestion' | 'dismiss' | 'custom_action' | 'ask_alternatives' | 'snooze';
  parameters?: Record<string, any>;
}

export interface NudgeActionResponse {
  nudgeId: string;
  actionTaken: string;
  result: {
    taskReassigned?: boolean;
    fromUserId?: string;
    toUserId?: string;
    taskId?: string;
  };
  nudgeStatus: NudgeStatus;
  followUpNudges: string[];
  completedAt: string;
}

export interface NudgeStats {
  period: string;
  generated: number;
  acted: number;
  dismissed: number;
  expired: number;
  actionRate: number;
  avgTimeToAction: string;
  byType: Record<string, { generated: number; acted: number }>;
}

// ==================== PROGRESS API TYPES ====================

export interface ProjectProgress {
  projectId: string;
  progressPercentage: number;
  calculationMethod: string;
  breakdown: {
    byStatus: Record<string, { count: number; hours: number; percentage: number }>;
    byHierarchy: {
      tasks: { total: number; completed: number };
      subtasks: { total: number; completed: number };
    };
  };
  variance: {
    expectedProgress: number;
    actualProgress: number;
    variance: number;
    status: 'ahead_of_schedule' | 'on_schedule' | 'behind_schedule';
  };
  milestones: Array<{
    taskId: string;
    title: string;
    status: string;
    completedAt?: string;
  }>;
  calculatedAt: string;
}

export interface TaskProgress {
  taskId: string;
  status: string;
  progressPercentage: number;
  actualHours: number;
  estimatedHours: number;
  remainingHours: number;
}

export interface UpdateProgressRequest {
  status?: string;
  progressPercentage?: number;
  actualHours?: number;
  note?: string;
}

export interface UpdateProgressResponse {
  taskId: string;
  previousStatus: string;
  newStatus: string;
  progressPercentage: number;
  actualHours: number;
  estimatedHours: number;
  remainingHours: number;
  affected: {
    parentProgressUpdated: boolean;
    projectProgressUpdated: boolean;
    dependentsNotified: string[];
  };
}

// ==================== DEPENDENCY API TYPES ====================

export interface TaskDependencies {
  taskId: string;
  dependencies: {
    direct: Array<{
      dependencyId: string;
      dependsOnTaskId: string;
      dependencyType: string;
      lagHours: number;
      status: string;
      isBlocking: boolean;
    }>;
    indirect: Array<{
      path: string[];
      depth: number;
    }>;
  };
  dependents: {
    direct: Array<{
      taskId: string;
      dependencyType: string;
      isBlocked: boolean;
    }>;
    indirect: Array<{
      path: string[];
      depth: number;
    }>;
  };
  chainAnalysis: {
    longestChain: number;
    criticalPathPosition: 'on_path' | 'near_path' | 'off_path';
    floatHours: number;
  };
}

export interface CreateDependencyRequest {
  taskId: string;
  dependsOnTaskId: string;
  dependencyType: 'finish_to_start' | 'start_to_start' | 'finish_to_finish' | 'start_to_finish';
  lagHours?: number;
}

export interface Dependency {
  dependencyId: string;
  taskId: string;
  dependsOnTaskId: string;
  dependencyType: string;
  lagHours: number;
  createdAt: string;
  validation: {
    valid: boolean;
    wouldCreateCycle: boolean;
  };
  impact?: {
    criticalPathChanged: boolean;
    affectedTasks: string[];
    newProjectEndDate?: string;
  };
}

export interface CriticalPath {
  projectId: string;
  criticalPath: {
    taskIds: string[];
    totalDuration: number;
    tasks: Array<{
      taskId: string;
      title: string;
      estimatedHours: number;
      earliestStart: string;
      earliestFinish: string;
      latestStart: string;
      latestFinish: string;
      floatHours: number;
    }>;
  };
  nonCriticalTasks: Array<{
    taskId: string;
    title: string;
    floatHours: number;
  }>;
  projectDuration: number;
  calculatedAt: string;
}

export interface ValidateDependencyRequest {
  taskId: string;
  dependsOnTaskId: string;
  dependencyType: string;
}

export interface ValidateDependencyResponse {
  valid: boolean;
  wouldCreateCycle: boolean;
  warnings: Array<{
    type: string;
    message: string;
  }>;
  impact: {
    estimatedDelay: number;
    affectedTasks: number;
  };
}

// ==================== ASSIGNMENT API TYPES ====================

export interface AssignmentSuggestion {
  taskId: string;
  taskTitle: string;
  requiredSkills: string[];
  candidates: Array<{
    rank: number;
    person: {
      id: string;
      name: string;
      avatarUrl?: string;
    };
    score: number;
    breakdown: {
      skillMatch: number;
      availability: number;
      workloadBalance: number;
      pastPerformance: number;
    };
    skillMatchDetails: Array<{
      skillId: string;
      required: boolean;
      hasSkill: boolean;
      proficiency: number;
      matchScore: number;
    }>;
    contextSwitchAnalysis: {
      activeProjects: number;
      currentWorkload: number;
      switchPenalty: number;
      riskLevel: 'low' | 'medium' | 'high';
    };
    warnings: string[];
    aiExplanation: string;
  }>;
  unassignableReason: string | null;
}

export interface AssignTaskRequest {
  personId: string;
  note?: string;
  skipSuggestion?: boolean;
}

export interface AssignTaskResponse {
  taskId: string;
  assignedTo: {
    personId: string;
    name: string;
  };
  previousAssignee: {
    personId: string;
    name: string;
  } | null;
  assignment: {
    assignedAt: string;
    assignedBy: string;
  };
  impact: {
    workloadUpdated: boolean;
    newAllocation: number;
    nudgesGenerated: string[];
    notificationsSent: string[];
  };
}

export interface AutoAssignRequest {
  strategy: 'best_match' | 'balanced_workload' | 'fastest_completion';
  constraints?: {
    maxAllocation?: number;
    requiredProficiency?: number;
  };
}

export interface CompatibilityCheck {
  taskId: string;
  personId: string;
  compatible: boolean;
  score: number;
  reasons: string[];
}

export interface BulkReassignRequest {
  reassignments: Array<{
    taskId: string;
    fromPersonId: string;
    toPersonId: string;
  }>;
  reason?: string;
}

// ==================== SCENARIO API TYPES ====================

export type ScenarioChangeType = 'employee_leave' | 'scope_change' | 'reallocation' | 'priority_shift' | 'deadline_change';
export type ScenarioStatus = 'draft' | 'pending' | 'simulated' | 'applied' | 'rejected';

export interface Scenario {
  scenarioId: string;
  title: string;
  description: string;
  changeType: ScenarioChangeType;
  status: ScenarioStatus;
  proposedChanges: Record<string, any>;
  impactAnalysis?: ScenarioImpactAnalysis;
  simulationStatus: 'pending' | 'running' | 'completed' | 'failed';
  createdAt: string;
  history?: Array<{
    action: string;
    timestamp: string;
    userId: string;
  }>;
}

export interface CreateScenarioRequest {
  title: string;
  description: string;
  changeType: ScenarioChangeType;
  proposedChanges: Record<string, any>;
}

export interface ScenarioImpactAnalysis {
  affectedProjects: Array<{
    projectId: string;
    name: string;
    impact: 'low' | 'medium' | 'high';
    delayDays: number;
    affectedTasks: string[];
  }>;
  affectedTasks: Array<{
    taskId: string;
    title: string;
    originalDueDate: string;
    newDueDate: string;
    delayDays: number;
    reason: string;
    suggestedReassignment?: {
      toPersonId: string;
      compatibility: number;
    };
  }>;
  timelineComparison: {
    originalEndDate: string;
    newEndDate: string;
    totalDelayDays: number;
  };
  costAnalysis: {
    totalCost: number;
    breakdown: Array<{
      category: string;
      amount: number;
    }>;
    confidence: number;
  };
  resourceImpacts: Array<{
    personId: string;
    currentAllocation: number;
    newAllocation: number;
    risk: string;
  }>;
}

export interface SimulateRequest {
  depth?: 'quick' | 'full';
  includeRecommendations?: boolean;
}

export interface SimulateResponse {
  scenarioId: string;
  simulationStatus: string;
  impactAnalysis: ScenarioImpactAnalysis;
  aiRecommendations: Array<{
    priority: number;
    action: string;
    reasoning: string;
    estimatedImpact: string;
  }>;
  calculatedAt: string;
  simulationDuration: string;
}

export interface ApplyScenarioRequest {
  applyRecommendations?: boolean;
  selectedRecommendations?: number[];
  notifyStakeholders?: boolean;
}

export interface ApplyScenarioResponse {
  scenarioId: string;
  status: string;
  appliedAt: string;
  appliedBy: string;
  changes: {
    tasksReassigned: Array<{
      taskId: string;
      from: string;
      to: string;
    }>;
    datesAdjusted: Array<{
      taskId: string;
      originalDueDate: string;
      newDueDate: string;
    }>;
    notificationsSent: number;
  };
  followUp: {
    nudgesCreated: string[];
    calendarEventsCreated: boolean;
  };
}

// ==================== WORKLOAD API TYPES ====================

export interface TeamWorkload {
  weekStarting: string;
  teamCapacity: number;
  teamAllocation: number;
  utilizationRate: number;
  members: Array<{
    personId: string;
    name: string;
    role?: string;
    allocation: {
      percentage: number;
      assignedHours: number;
      capacityHours: number;
    };
    tasks: Array<{
      taskId: string;
      title: string;
      projectId: string;
      estimatedHours: number;
      allocationThisWeek: number;
    }>;
    status: 'overallocated' | 'optimal' | 'available' | 'underutilized';
    riskLevel: 'low' | 'medium' | 'high';
    availability: {
      thisWeek: number;
      nextWeek: number;
    };
  }>;
  summary: {
    overallocated: number;
    optimal: number;
    available: number;
    underutilized: number;
  };
}

export interface PersonWorkload {
  personId: string;
  weekStarting: string;
  allocation: {
    percentage: number;
    assignedHours: number;
    capacityHours: number;
  };
  tasks: Array<{
    taskId: string;
    title: string;
    projectId: string;
    estimatedHours: number;
    allocationThisWeek: number;
  }>;
}

export interface WorkloadForecast {
  personId: string;
  forecast: Array<{
    weekStarting: string;
    allocation: number;
    assignedHours: number;
    tasks: number;
    risk: 'overallocated' | 'optimal' | 'available' | 'underutilized';
  }>;
  riskPeriods: Array<{
    startWeek: string;
    endWeek: string;
    severity: 'low' | 'medium' | 'high';
    reason: string;
  }>;
  recommendations: string[];
}

// ==================== QUERY PARAMETERS ====================

export interface NudgeQueryParams {
  status?: NudgeStatus | NudgeStatus[];
  severity?: NudgeSeverity | NudgeSeverity[];
  type?: NudgeType | NudgeType[];
  projectId?: string;
  personId?: string;
  limit?: number;
  offset?: number;
  cursor?: string;
}

export interface HealthTrendParams {
  projectId?: string;
  days?: number;
}

export interface ProjectRankingParams {
  status?: string;
  assigneeId?: string;
  minScore?: number;
  limit?: number;
  offset?: number;
}
