// Dependency API Module
// Handles task dependencies and critical path management

import { apiClient } from './client';
import {
  TaskDependencies,
  CreateDependencyRequest,
  Dependency,
  CriticalPath,
  ValidateDependencyRequest,
  ValidateDependencyResponse,
} from './types';

/**
 * Get dependencies for a task
 */
export async function getTaskDependencies(
  taskId: string,
  includeIndirect: boolean = false
) {
  return apiClient.get<TaskDependencies>(`dependencies/tasks/${taskId}`, {
    includeIndirect,
  });
}

/**
 * Create a new dependency
 */
export async function createDependency(request: CreateDependencyRequest) {
  return apiClient.post<Dependency>('dependencies', request);
}

/**
 * Remove a dependency
 */
export async function deleteDependency(dependencyId: string) {
  return apiClient.delete<void>(`dependencies/${dependencyId}`);
}

/**
 * Get critical path for a project
 */
export async function getCriticalPath(projectId: string) {
  return apiClient.get<CriticalPath>(`dependencies/critical-path/${projectId}`);
}

/**
 * Validate a potential dependency (check for cycles, etc.)
 */
export async function validateDependency(
  request: ValidateDependencyRequest
) {
  return apiClient.post<ValidateDependencyResponse>(
    'dependencies/validate',
    request
  );
}

/**
 * Get dependency graph for a project
 */
export async function getDependencyGraph(projectId: string) {
  return apiClient.get<{
    projectId: string;
    nodes: Array<{
      id: string;
      title: string;
      status: string;
    }>;
    edges: Array<{
      from: string;
      to: string;
      type: string;
    }>;
  }>(`dependencies/graph/${projectId}`);
}

// Export as namespace
export const dependencyApi = {
  getTaskDependencies,
  createDependency,
  deleteDependency,
  getCriticalPath,
  validateDependency,
  getDependencyGraph,
};

export default dependencyApi;
