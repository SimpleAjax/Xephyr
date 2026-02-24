// Progress API Module
// Handles project and task progress tracking

import { apiClient } from './client';
import {
  ProjectProgress,
  TaskProgress,
  UpdateProgressRequest,
  UpdateProgressResponse,
} from './types';

/**
 * Get project progress overview
 */
export async function getProjectProgress(projectId: string) {
  return apiClient.get<ProjectProgress>(`progress/projects/${projectId}`);
}

/**
 * Get detailed progress for a task
 */
export async function getTaskProgress(taskId: string) {
  return apiClient.get<TaskProgress>(`progress/tasks/${taskId}`);
}

/**
 * Update task progress
 */
export async function updateTaskProgress(
  taskId: string,
  request: UpdateProgressRequest
) {
  return apiClient.post<UpdateProgressResponse>(
    `progress/tasks/${taskId}/update`,
    request
  );
}

/**
 * Get hierarchical progress roll-up for a project
 */
export async function getProgressRollup(projectId: string) {
  return apiClient.get<{
    projectId: string;
    rollup: Array<{
      taskId: string;
      title: string;
      progress: number;
      weight: number;
      children?: Array<{
        taskId: string;
        title: string;
        progress: number;
      }>;
    }>;
  }>(`progress/rollups/${projectId}`);
}

// Export as namespace
export const progressApi = {
  getProjectProgress,
  getTaskProgress,
  updateTaskProgress,
  getProgressRollup,
};

export default progressApi;
