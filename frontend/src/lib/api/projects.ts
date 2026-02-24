// Projects API Module
// Handles project CRUD operations

import { apiClient } from './client';
import { Project } from '@/types';

export interface CreateProjectRequest {
  name: string;
  description: string;
  priority?: number;
  startDate?: string;
  targetEndDate?: string;
}

export interface UpdateProjectRequest {
  name?: string;
  description?: string;
  status?: 'active' | 'paused' | 'completed' | 'archived';
  priority?: number;
  targetEndDate?: string;
}

export interface ProjectListResponse {
  projects: Project[];
  total: number;
}

/**
 * Get all projects
 */
export async function getProjects(status?: string) {
  return apiClient.get<ProjectListResponse>('projects', status ? { status } : undefined);
}

/**
 * Get a single project by ID
 */
export async function getProject(projectId: string) {
  return apiClient.get<Project>(`projects/${projectId}`);
}

/**
 * Create a new project
 */
export async function createProject(request: CreateProjectRequest) {
  return apiClient.post<Project>('projects', request);
}

/**
 * Update a project
 */
export async function updateProject(projectId: string, request: UpdateProjectRequest) {
  return apiClient.patch<Project>(`projects/${projectId}`, request);
}

/**
 * Delete a project
 */
export async function deleteProject(projectId: string) {
  return apiClient.delete<void>(`projects/${projectId}`);
}

/**
 * Get project team members
 */
export async function getProjectTeam(projectId: string) {
  return apiClient.get<{ users: Array<{ id: string; name: string; role: string }> }>(
    `projects/${projectId}/team`
  );
}

// Export as namespace
export const projectsApi = {
  getProjects,
  getProject,
  createProject,
  updateProject,
  deleteProject,
  getProjectTeam,
};

export default projectsApi;
