// Assignment API Module
// Handles task assignments and AI suggestions

import { apiClient } from './client';
import {
  AssignmentSuggestion,
  AssignTaskRequest,
  AssignTaskResponse,
  AutoAssignRequest,
  CompatibilityCheck,
  BulkReassignRequest,
} from './types';

/**
 * Get assignment suggestions for a task
 */
export async function getAssignmentSuggestions(
  taskId: string,
  limit: number = 3
) {
  return apiClient.get<AssignmentSuggestion>('assignments/suggestions', {
    taskId,
    limit,
  });
}

/**
 * Assign a task to a person
 */
export async function assignTask(taskId: string, request: AssignTaskRequest) {
  return apiClient.post<AssignTaskResponse>(
    `assignments/tasks/${taskId}/assign`,
    request
  );
}

/**
 * Auto-assign task using AI
 */
export async function autoAssignTask(
  taskId: string,
  request?: AutoAssignRequest
) {
  return apiClient.post<AssignTaskResponse>(
    `assignments/tasks/${taskId}/auto-assign`,
    request || { strategy: 'best_match' }
  );
}

/**
 * Check compatibility between person and task
 */
export async function checkCompatibility(taskId: string, personId: string) {
  return apiClient.get<CompatibilityCheck>('assignments/compatibility', {
    taskId,
    personId,
  });
}

/**
 * Bulk reassign tasks
 */
export async function bulkReassign(request: BulkReassignRequest) {
  return apiClient.post<{ processed: number; succeeded: number; failed: number }>(
    'assignments/bulk-reassign',
    request
  );
}

// Export as namespace
export const assignmentApi = {
  getAssignmentSuggestions,
  assignTask,
  autoAssignTask,
  checkCompatibility,
  bulkReassign,
};

export default assignmentApi;
