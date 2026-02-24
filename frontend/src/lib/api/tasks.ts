// Tasks API Module
// Handles task CRUD operations

import { apiClient } from './client';
import { Task, TaskStatus } from '@/types';

export interface CreateTaskRequest {
  title: string;
  description?: string;
  projectId: string;
  parentTaskId?: string;
  hierarchyLevel?: number;
  priority?: 'low' | 'medium' | 'high' | 'critical';
  estimatedHours?: number;
  dueDate?: string;
  assigneeId?: string;
  requiredSkills?: string[];
  isMilestone?: boolean;
}

export interface UpdateTaskRequest {
  title?: string;
  description?: string;
  status?: TaskStatus;
  priority?: 'low' | 'medium' | 'high' | 'critical';
  estimatedHours?: number;
  actualHours?: number;
  dueDate?: string;
  assigneeId?: string | null;
  requiredSkills?: string[];
  isMilestone?: boolean;
}

export interface TaskListResponse {
  tasks: Task[];
  total: number;
}

export interface UpdateTaskStatusRequest {
  status: TaskStatus;
  progressPercentage?: number;
  note?: string;
}

/**
 * Get all tasks
 */
export async function getTasks(projectId?: string) {
  return apiClient.get<TaskListResponse>('tasks', projectId ? { projectId } : undefined);
}

/**
 * Get a single task by ID
 */
export async function getTask(taskId: string) {
  return apiClient.get<Task>(`tasks/${taskId}`);
}

/**
 * Get tasks by project
 */
export async function getTasksByProject(projectId: string) {
  return apiClient.get<TaskListResponse>('tasks', { projectId });
}

/**
 * Create a new task
 */
export async function createTask(request: CreateTaskRequest) {
  return apiClient.post<Task>('tasks', request);
}

/**
 * Update a task
 */
export async function updateTask(taskId: string, request: UpdateTaskRequest) {
  return apiClient.patch<Task>(`tasks/${taskId}`, request);
}

/**
 * Update task status
 */
export async function updateTaskStatus(taskId: string, request: UpdateTaskStatusRequest) {
  return apiClient.post<Task>(`tasks/${taskId}/status`, request);
}

/**
 * Assign task to a user
 */
export async function assignTaskToUser(taskId: string, personId: string, note?: string) {
  return apiClient.post<Task>(`tasks/${taskId}/assign`, { personId, note });
}

/**
 * Unassign task
 */
export async function unassignTask(taskId: string) {
  return apiClient.post<Task>(`tasks/${taskId}/unassign`, {});
}

/**
 * Delete a task
 */
export async function deleteTask(taskId: string) {
  return apiClient.delete<void>(`tasks/${taskId}`);
}

/**
 * Get unassigned tasks
 */
export async function getUnassignedTasks(projectId?: string) {
  return apiClient.get<TaskListResponse>('tasks/unassigned', projectId ? { projectId } : undefined);
}

// Export as namespace
export const tasksApi = {
  getTasks,
  getTask,
  getTasksByProject,
  createTask,
  updateTask,
  updateTaskStatus,
  assignTask: assignTaskToUser,
  unassignTask,
  deleteTask,
  getUnassignedTasks,
};

export default tasksApi;
