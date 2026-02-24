// React Query hooks for Workload API

import {
  useQuery,
  useQueryClient,
  UseQueryOptions,
} from '@tanstack/react-query';
import { workloadApi } from '@/lib/api';
import {
  TeamWorkload,
  PersonWorkload,
  WorkloadForecast,
} from '@/lib/api/types';
import { ApiResponse } from '@/lib/api/types';
import { ApiClientError } from '@/lib/api/client';

// Query keys
export const workloadKeys = {
  all: ['workload'] as const,
  team: (week?: string) => [...workloadKeys.all, 'team', week] as const,
  person: (personId: string, week?: string) => 
    [...workloadKeys.all, 'person', personId, week] as const,
  forecast: (personId: string, weeks?: number) => 
    [...workloadKeys.all, 'forecast', personId, weeks] as const,
  analytics: (period?: string) => [...workloadKeys.all, 'analytics', period] as const,
};

/**
 * Hook to fetch team workload
 */
export function useTeamWorkload(
  week?: string,
  includeForecast: boolean = false,
  options?: UseQueryOptions<ApiResponse<TeamWorkload>, ApiClientError>
) {
  return useQuery({
    queryKey: workloadKeys.team(week),
    queryFn: () => workloadApi.getTeamWorkload(week, includeForecast),
    ...options,
  });
}

/**
 * Hook to fetch individual person workload
 */
export function usePersonWorkload(
  personId: string,
  week?: string,
  options?: UseQueryOptions<ApiResponse<PersonWorkload>, ApiClientError>
) {
  return useQuery({
    queryKey: workloadKeys.person(personId, week),
    queryFn: () => workloadApi.getPersonWorkload(personId, week),
    enabled: !!personId,
    ...options,
  });
}

/**
 * Hook to fetch workload forecast
 */
export function useWorkloadForecast(
  personId: string,
  weeks: number = 8,
  options?: UseQueryOptions<ApiResponse<WorkloadForecast>, ApiClientError>
) {
  return useQuery({
    queryKey: workloadKeys.forecast(personId, weeks),
    queryFn: () => workloadApi.getWorkloadForecast(personId, weeks),
    enabled: !!personId,
    ...options,
  });
}

/**
 * Hook to fetch workload analytics
 */
export function useWorkloadAnalytics(
  period: string = '30d',
  options?: UseQueryOptions<
    ApiResponse<{
      period: string;
      avgUtilization: number;
      overallocatedDays: number;
      riskTrend: 'improving' | 'stable' | 'worsening';
      byPerson: Array<{
        personId: string;
        name: string;
        avgAllocation: number;
        peakAllocation: number;
        riskDays: number;
      }>;
    }>,
    ApiClientError
  >
) {
  return useQuery({
    queryKey: workloadKeys.analytics(period),
    queryFn: () => workloadApi.getWorkloadAnalytics(period),
    ...options,
  });
}

/**
 * Hook to prefetch workload data
 */
export function usePrefetchWorkload() {
  const queryClient = useQueryClient();
  
  return (week?: string) => {
    queryClient.prefetchQuery({
      queryKey: workloadKeys.team(week),
      queryFn: () => workloadApi.getTeamWorkload(week),
    });
  };
}
