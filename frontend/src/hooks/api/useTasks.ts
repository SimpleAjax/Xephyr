// React Query hooks for Tasks API

import {
  useQuery,
  useMutation,
  useQueryClient,
  UseQueryOptions,
} from '@tanstack/react-query';
import { tasksApi } from '@/lib/api';
import { Task, TaskStatus } from '@/types';
import { ApiResponse } from '@/lib/api/types';
import { ApiClientError } from '@/lib/api/client';

// Query keys
export const taskKeys = {
  all: ['tasks'] as const,
  lists: () => [...taskKeys.all, 'list'] as const,
  list: (projectId?: string) => [...taskKeys.lists(), projectId] as const,
  details: () => [...taskKeys.all, 'detail'] as const,
  detail: (id: string) => [...taskKeys.details(), id] as const,
  unassigned: () => [...taskKeys.all, 'unassigned'] as const,
};

/**
 * Hook to fetch all tasks
 */
export function useTasks(
  projectId?: string,
  options?: UseQueryOptions<ApiResponse<{ tasks: Task[]; total: number }>, ApiClientError>
) {
  return useQuery({
    queryKey: taskKeys.list(projectId),
    queryFn: () => tasksApi.getTasks(projectId),
    ...options,
  });
}

/**
 * Hook to fetch a single task
 */
export function useTask(
  taskId: string,
  options?: UseQueryOptions<ApiResponse<Task>, ApiClientError>
) {
  return useQuery({
    queryKey: taskKeys.detail(taskId),
    queryFn: () => tasksApi.getTask(taskId),
    enabled: !!taskId,
    ...options,
  });
}

/**
 * Hook to fetch tasks by project
 */
export function useTasksByProject(
  projectId: string,
  options?: UseQueryOptions<ApiResponse<{ tasks: Task[]; total: number }>, ApiClientError>
) {
  return useQuery({
    queryKey: [...taskKeys.list(projectId), 'byProject'],
    queryFn: () => tasksApi.getTasksByProject(projectId),
    enabled: !!projectId,
    ...options,
  });
}

/**
 * Hook to fetch unassigned tasks
 */
export function useUnassignedTasks(
  projectId?: string,
  options?: UseQueryOptions<ApiResponse<{ tasks: Task[]; total: number }>, ApiClientError>
) {
  return useQuery({
    queryKey: taskKeys.unassigned(),
    queryFn: () => tasksApi.getUnassignedTasks(projectId),
    ...options,
  });
}

/**
 * Hook to create a new task
 */
export function useCreateTask() {
  const queryClient = useQueryClient();

  return useMutation<
    ApiResponse<Task>,
    ApiClientError,
    Parameters<typeof tasksApi.createTask>[0]
  >({
    mutationFn: (request) => tasksApi.createTask(request),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: taskKeys.lists() });
      if (variables.projectId) {
        queryClient.invalidateQueries({ queryKey: taskKeys.list(variables.projectId) });
      }
    },
  });
}

/**
 * Hook to update a task
 */
export function useUpdateTask() {
  const queryClient = useQueryClient();

  return useMutation<
    ApiResponse<Task>,
    ApiClientError,
    { taskId: string; updates: Parameters<typeof tasksApi.updateTask>[1] }
  >({
    mutationFn: ({ taskId, updates }) => tasksApi.updateTask(taskId, updates),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: taskKeys.detail(variables.taskId) });
      queryClient.invalidateQueries({ queryKey: taskKeys.lists() });
    },
  });
}

/**
 * Hook to update task status
 */
export function useUpdateTaskStatus() {
  const queryClient = useQueryClient();

  return useMutation<
    ApiResponse<Task>,
    ApiClientError,
    { taskId: string; status: TaskStatus; progressPercentage?: number; note?: string }
  >({
    mutationFn: ({ taskId, status, progressPercentage, note }) =>
      tasksApi.updateTaskStatus(taskId, { status, progressPercentage, note }),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: taskKeys.detail(variables.taskId) });
      queryClient.invalidateQueries({ queryKey: taskKeys.lists() });
    },
  });
}

/**
 * Hook to assign a task to a user
 */
export function useAssignTaskToUser() {
  const queryClient = useQueryClient();

  return useMutation<
    ApiResponse<Task>,
    ApiClientError,
    { taskId: string; personId: string; note?: string }
  >({
    mutationFn: ({ taskId, personId, note }) => tasksApi.assignTask(taskId, personId, note),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: taskKeys.detail(variables.taskId) });
      queryClient.invalidateQueries({ queryKey: taskKeys.lists() });
      queryClient.invalidateQueries({ queryKey: ['workload'] });
    },
  });
}

/**
 * Hook to unassign a task
 */
export function useUnassignTaskFromUser() {
  const queryClient = useQueryClient();

  return useMutation<ApiResponse<Task>, ApiClientError, string>({
    mutationFn: (taskId) => tasksApi.unassignTask(taskId),
    onSuccess: (_, taskId) => {
      queryClient.invalidateQueries({ queryKey: taskKeys.detail(taskId) });
      queryClient.invalidateQueries({ queryKey: taskKeys.lists() });
      queryClient.invalidateQueries({ queryKey: ['workload'] });
    },
  });
}

/**
 * Hook to delete a task
 */
export function useDeleteTask() {
  const queryClient = useQueryClient();

  return useMutation<ApiResponse<void>, ApiClientError, string>({
    mutationFn: (taskId) => tasksApi.deleteTask(taskId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: taskKeys.lists() });
    },
  });
}
