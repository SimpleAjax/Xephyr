// React Query hooks for Assignment API

import {
  useQuery,
  useMutation,
  useQueryClient,
  UseQueryOptions,
} from '@tanstack/react-query';
import { assignmentApi } from '@/lib/api';
import {
  AssignmentSuggestion,
  AssignTaskRequest,
  AssignTaskResponse,
  CompatibilityCheck,
} from '@/lib/api/types';
import { ApiResponse } from '@/lib/api/types';
import { ApiClientError } from '@/lib/api/client';

// Query keys
export const assignmentKeys = {
  all: ['assignments'] as const,
  suggestions: (taskId: string) => [...assignmentKeys.all, 'suggestions', taskId] as const,
  compatibility: (taskId: string, personId: string) => 
    [...assignmentKeys.all, 'compatibility', taskId, personId] as const,
};

/**
 * Hook to fetch assignment suggestions for a task
 */
export function useAssignmentSuggestions(
  taskId: string,
  limit: number = 3,
  options?: UseQueryOptions<ApiResponse<AssignmentSuggestion>, ApiClientError>
) {
  return useQuery({
    queryKey: assignmentKeys.suggestions(taskId),
    queryFn: () => assignmentApi.getAssignmentSuggestions(taskId, limit),
    enabled: !!taskId,
    ...options,
  });
}

/**
 * Hook to check compatibility between person and task
 */
export function useCompatibilityCheck(
  taskId: string,
  personId: string,
  options?: UseQueryOptions<ApiResponse<CompatibilityCheck>, ApiClientError>
) {
  return useQuery({
    queryKey: assignmentKeys.compatibility(taskId, personId),
    queryFn: () => assignmentApi.checkCompatibility(taskId, personId),
    enabled: !!taskId && !!personId,
    ...options,
  });
}

/**
 * Hook to assign a task
 */
export function useAssignTask() {
  const queryClient = useQueryClient();

  return useMutation<ApiResponse<AssignTaskResponse>, ApiClientError, {
    taskId: string;
    request: AssignTaskRequest;
  }>({
    mutationFn: ({ taskId, request }) => assignmentApi.assignTask(taskId, request),
    onSuccess: (_, variables) => {
      // Invalidate related queries
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
      queryClient.invalidateQueries({ queryKey: ['projects'] });
      queryClient.invalidateQueries({ queryKey: assignmentKeys.suggestions(variables.taskId) });
      queryClient.invalidateQueries({ queryKey: ['workload'] });
      queryClient.invalidateQueries({ queryKey: ['nudges'] });
    },
  });
}

/**
 * Hook to auto-assign a task
 */
export function useAutoAssignTask() {
  const queryClient = useQueryClient();

  return useMutation<ApiResponse<AssignTaskResponse>, ApiClientError, {
    taskId: string;
    strategy?: 'best_match' | 'balanced_workload' | 'fastest_completion';
  }>({
    mutationFn: ({ taskId, strategy = 'best_match' }) => 
      assignmentApi.autoAssignTask(taskId, { strategy }),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
      queryClient.invalidateQueries({ queryKey: ['projects'] });
      queryClient.invalidateQueries({ queryKey: ['workload'] });
      queryClient.invalidateQueries({ queryKey: ['nudges'] });
    },
  });
}
