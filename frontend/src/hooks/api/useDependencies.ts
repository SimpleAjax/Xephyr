// React Query hooks for Dependency API

import {
  useQuery,
  useMutation,
  useQueryClient,
  UseQueryOptions,
} from '@tanstack/react-query';
import { dependencyApi } from '@/lib/api';
import {
  TaskDependencies,
  CreateDependencyRequest,
  Dependency,
  CriticalPath,
  ValidateDependencyResponse,
} from '@/lib/api/types';
import { ApiResponse } from '@/lib/api/types';
import { ApiClientError } from '@/lib/api/client';

// Query keys
export const dependencyKeys = {
  all: ['dependencies'] as const,
  task: (taskId: string) => [...dependencyKeys.all, 'task', taskId] as const,
  criticalPath: (projectId: string) => 
    [...dependencyKeys.all, 'criticalPath', projectId] as const,
  graph: (projectId: string) => 
    [...dependencyKeys.all, 'graph', projectId] as const,
  validation: () => [...dependencyKeys.all, 'validation'] as const,
};

/**
 * Hook to fetch dependencies for a task
 */
export function useTaskDependencies(
  taskId: string,
  includeIndirect: boolean = false,
  options?: UseQueryOptions<ApiResponse<TaskDependencies>, ApiClientError>
) {
  return useQuery({
    queryKey: [...dependencyKeys.task(taskId), includeIndirect],
    queryFn: () => dependencyApi.getTaskDependencies(taskId, includeIndirect),
    enabled: !!taskId,
    ...options,
  });
}

/**
 * Hook to fetch critical path for a project
 */
export function useCriticalPath(
  projectId: string,
  options?: UseQueryOptions<ApiResponse<CriticalPath>, ApiClientError>
) {
  return useQuery({
    queryKey: dependencyKeys.criticalPath(projectId),
    queryFn: () => dependencyApi.getCriticalPath(projectId),
    enabled: !!projectId,
    ...options,
  });
}

/**
 * Hook to fetch dependency graph
 */
export function useDependencyGraph(
  projectId: string,
  options?: UseQueryOptions<
    ApiResponse<{
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
    }>,
    ApiClientError
  >
) {
  return useQuery({
    queryKey: dependencyKeys.graph(projectId),
    queryFn: () => dependencyApi.getDependencyGraph(projectId),
    enabled: !!projectId,
    ...options,
  });
}

/**
 * Hook to create a dependency
 */
export function useCreateDependency() {
  const queryClient = useQueryClient();

  return useMutation<ApiResponse<Dependency>, ApiClientError, CreateDependencyRequest>({
    mutationFn: (request) => dependencyApi.createDependency(request),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: dependencyKeys.task(variables.taskId) });
      queryClient.invalidateQueries({ queryKey: dependencyKeys.task(variables.dependsOnTaskId) });
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
      queryClient.invalidateQueries({ queryKey: ['projects'] });
    },
  });
}

/**
 * Hook to delete a dependency
 */
export function useDeleteDependency() {
  const queryClient = useQueryClient();

  return useMutation<ApiResponse<void>, ApiClientError, {
    dependencyId: string;
    taskId?: string;
  }>({
    mutationFn: ({ dependencyId }) => dependencyApi.deleteDependency(dependencyId),
    onSuccess: (_, variables) => {
      if (variables.taskId) {
        queryClient.invalidateQueries({ queryKey: dependencyKeys.task(variables.taskId) });
      }
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
      queryClient.invalidateQueries({ queryKey: ['projects'] });
    },
  });
}

/**
 * Hook to validate a potential dependency
 */
export function useValidateDependency(
  taskId?: string,
  dependsOnTaskId?: string,
  dependencyType?: string,
  options?: UseQueryOptions<ApiResponse<ValidateDependencyResponse>, ApiClientError>
) {
  return useQuery({
    queryKey: [...dependencyKeys.validation(), taskId, dependsOnTaskId, dependencyType],
    queryFn: () => 
      dependencyApi.validateDependency({
        taskId: taskId!,
        dependsOnTaskId: dependsOnTaskId!,
        dependencyType: dependencyType!,
      }),
    enabled: !!taskId && !!dependsOnTaskId && !!dependencyType,
    ...options,
  });
}
