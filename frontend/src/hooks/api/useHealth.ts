// React Query hooks for Health API

import {
  useQuery,
  useMutation,
  useQueryClient,
  UseQueryOptions,
} from '@tanstack/react-query';
import { healthApi } from '@/lib/api';
import { PortfolioHealth, ProjectHealth, HealthTrend, HealthTrendParams } from '@/lib/api/types';
import { ApiResponse } from '@/lib/api/types';
import { ApiClientError } from '@/lib/api/client';

// Query keys
export const healthKeys = {
  all: ['health'] as const,
  portfolio: () => [...healthKeys.all, 'portfolio'] as const,
  project: (projectId: string) => [...healthKeys.all, 'project', projectId] as const,
  trends: (params?: HealthTrendParams) => [...healthKeys.all, 'trends', params] as const,
};

/**
 * Hook to fetch portfolio health
 */
export function usePortfolioHealth(
  options?: UseQueryOptions<ApiResponse<PortfolioHealth>, ApiClientError>
) {
  return useQuery({
    queryKey: healthKeys.portfolio(),
    queryFn: () => healthApi.getPortfolioHealth(),
    ...options,
  });
}

/**
 * Hook to fetch project health
 */
export function useProjectHealth(
  projectId: string,
  options?: UseQueryOptions<ApiResponse<ProjectHealth>, ApiClientError>
) {
  return useQuery({
    queryKey: healthKeys.project(projectId),
    queryFn: () => healthApi.getProjectHealth(projectId),
    enabled: !!projectId,
    ...options,
  });
}

/**
 * Hook to fetch health trends
 */
export function useHealthTrends(
  params?: HealthTrendParams,
  options?: UseQueryOptions<ApiResponse<HealthTrend>, ApiClientError>
) {
  return useQuery({
    queryKey: healthKeys.trends(params),
    queryFn: () => healthApi.getHealthTrends(params),
    ...options,
  });
}

/**
 * Hook to prefetch project health (useful for navigation)
 */
export function usePrefetchProjectHealth() {
  const queryClient = useQueryClient();
  
  return (projectId: string) => {
    queryClient.prefetchQuery({
      queryKey: healthKeys.project(projectId),
      queryFn: () => healthApi.getProjectHealth(projectId),
    });
  };
}
