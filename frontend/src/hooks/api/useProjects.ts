// React Query hooks for Project-related APIs
// Combines Health, Progress, and Priority APIs for projects

import {
  useQuery,
  useMutation,
  useQueryClient,
  UseQueryOptions,
} from '@tanstack/react-query';
import { healthApi, progressApi, priorityApi } from '@/lib/api';
import {
  ProjectHealth,
  ProjectProgress,
  ProjectRanking,
  HealthTrend,
  ProjectRankingParams,
  HealthTrendParams,
} from '@/lib/api/types';
import { ApiResponse } from '@/lib/api/types';
import { ApiClientError } from '@/lib/api/client';

// Query keys
export const projectKeys = {
  all: ['projects'] as const,
  health: (projectId: string) => [...projectKeys.all, 'health', projectId] as const,
  progress: (projectId: string) => [...projectKeys.all, 'progress', projectId] as const,
  ranking: (projectId: string, params?: ProjectRankingParams) => 
    [...projectKeys.all, 'ranking', projectId, params] as const,
  trends: (params?: HealthTrendParams) => [...projectKeys.all, 'trends', params] as const,
  rollup: (projectId: string) => [...projectKeys.all, 'rollup', projectId] as const,
};

/**
 * Hook to fetch project progress
 */
export function useProjectProgress(
  projectId: string,
  options?: UseQueryOptions<ApiResponse<ProjectProgress>, ApiClientError>
) {
  return useQuery({
    queryKey: projectKeys.progress(projectId),
    queryFn: () => progressApi.getProjectProgress(projectId),
    enabled: !!projectId,
    ...options,
  });
}

/**
 * Hook to fetch project task ranking by priority
 */
export function useProjectRanking(
  projectId: string,
  params?: ProjectRankingParams,
  options?: UseQueryOptions<ApiResponse<ProjectRanking>, ApiClientError>
) {
  return useQuery({
    queryKey: projectKeys.ranking(projectId, params),
    queryFn: () => priorityApi.getProjectRanking(projectId, params),
    enabled: !!projectId,
    ...options,
  });
}

/**
 * Hook to fetch health trends for a project
 */
export function useProjectHealthTrends(
  params?: HealthTrendParams,
  options?: UseQueryOptions<ApiResponse<HealthTrend>, ApiClientError>
) {
  return useQuery({
    queryKey: projectKeys.trends(params),
    queryFn: () => healthApi.getHealthTrends(params),
    ...options,
  });
}

/**
 * Hook to fetch progress rollup (hierarchical view)
 */
export function useProgressRollup(
  projectId: string,
  options?: UseQueryOptions<
    ApiResponse<{
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
    }>,
    ApiClientError
  >
) {
  return useQuery({
    queryKey: projectKeys.rollup(projectId),
    queryFn: () => progressApi.getProgressRollup(projectId),
    enabled: !!projectId,
    ...options,
  });
}
