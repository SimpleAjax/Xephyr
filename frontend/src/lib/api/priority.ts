// Priority API Module
// Handles task priority calculations and rankings

import { apiClient } from './client';
import {
  TaskPriority,
  BulkPriorityRequest,
  BulkPriorityResponse,
  ProjectRanking,
  RecalculateRequest,
  RecalculateResponse,
  ProjectRankingParams,
} from './types';

/**
 * Get priority details for a single task
 */
export async function getTaskPriority(taskId: string) {
  return apiClient.get<TaskPriority>(`priorities/tasks/${taskId}`);
}

/**
 * Get priorities for multiple tasks in bulk
 */
export async function getBulkPriorities(request: BulkPriorityRequest) {
  return apiClient.post<BulkPriorityResponse>('priorities/tasks/bulk', request);
}

/**
 * Get task priority ranking for a project
 */
export async function getProjectRanking(
  projectId: string,
  params?: ProjectRankingParams
) {
  const queryParams: Record<string, string> = {};
  if (params?.status) queryParams.status = params.status;
  if (params?.assigneeId) queryParams.assigneeId = params.assigneeId;
  if (params?.minScore) queryParams.minScore = String(params.minScore);
  if (params?.limit) queryParams.limit = String(params.limit);
  if (params?.offset) queryParams.offset = String(params.offset);
  return apiClient.get<ProjectRanking>(
    `priorities/projects/${projectId}/ranking`,
    queryParams
  );
}

/**
 * Trigger priority recalculation
 */
export async function recalculatePriorities(request: RecalculateRequest) {
  return apiClient.post<RecalculateResponse>('priorities/recalculate', request);
}

// Export as namespace
export const priorityApi = {
  getTaskPriority,
  getBulkPriorities,
  getProjectRanking,
  recalculatePriorities,
};

export default priorityApi;
