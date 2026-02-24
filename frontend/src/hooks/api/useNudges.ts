// React Query hooks for Nudge API

import {
  useQuery,
  useMutation,
  useQueryClient,
  UseQueryOptions,
} from '@tanstack/react-query';
import { nudgeApi } from '@/lib/api';
import {
  Nudge,
  NudgeListResponse,
  NudgeQueryParams,
  NudgeActionRequest,
  NudgeActionResponse,
  NudgeStats,
  NudgeStatus,
} from '@/lib/api/types';
import { ApiResponse } from '@/lib/api/types';
import { ApiClientError } from '@/lib/api/client';

// Query keys
export const nudgeKeys = {
  all: ['nudges'] as const,
  lists: () => [...nudgeKeys.all, 'list'] as const,
  list: (params?: NudgeQueryParams) => [...nudgeKeys.lists(), params] as const,
  details: () => [...nudgeKeys.all, 'detail'] as const,
  detail: (id: string) => [...nudgeKeys.details(), id] as const,
  stats: () => [...nudgeKeys.all, 'stats'] as const,
};

/**
 * Hook to fetch nudges with filters
 */
export function useNudges(
  params?: NudgeQueryParams,
  options?: UseQueryOptions<ApiResponse<NudgeListResponse>, ApiClientError>
) {
  return useQuery({
    queryKey: nudgeKeys.list(params),
    queryFn: () => nudgeApi.getNudges(params),
    ...options,
  });
}

/**
 * Hook to fetch a single nudge
 */
export function useNudge(
  nudgeId: string,
  options?: UseQueryOptions<ApiResponse<Nudge>, ApiClientError>
) {
  return useQuery({
    queryKey: nudgeKeys.detail(nudgeId),
    queryFn: () => nudgeApi.getNudge(nudgeId),
    enabled: !!nudgeId,
    ...options,
  });
}

/**
 * Hook to fetch nudge statistics
 */
export function useNudgeStats(
  period: string = '30d',
  options?: UseQueryOptions<ApiResponse<NudgeStats>, ApiClientError>
) {
  return useQuery({
    queryKey: nudgeKeys.stats(),
    queryFn: () => nudgeApi.getNudgeStats(period),
    ...options,
  });
}

/**
 * Hook to take action on a nudge
 */
export function useNudgeAction() {
  const queryClient = useQueryClient();

  return useMutation<ApiResponse<NudgeActionResponse>, ApiClientError, {
    nudgeId: string;
    action: NudgeActionRequest;
  }>({
    mutationFn: ({ nudgeId, action }) => nudgeApi.takeNudgeAction(nudgeId, action),
    onSuccess: (_, variables) => {
      // Invalidate specific nudge and lists
      queryClient.invalidateQueries({ queryKey: nudgeKeys.detail(variables.nudgeId) });
      queryClient.invalidateQueries({ queryKey: nudgeKeys.lists() });
      queryClient.invalidateQueries({ queryKey: nudgeKeys.stats() });
    },
  });
}

/**
 * Hook to update nudge status
 */
export function useUpdateNudgeStatus() {
  const queryClient = useQueryClient();

  return useMutation<ApiResponse<Nudge>, ApiClientError, {
    nudgeId: string;
    status: 'read' | 'dismissed' | 'acted';
  }>({
    mutationFn: ({ nudgeId, status }) => nudgeApi.updateNudgeStatus(nudgeId, status),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: nudgeKeys.detail(variables.nudgeId) });
      queryClient.invalidateQueries({ queryKey: nudgeKeys.lists() });
    },
  });
}

/**
 * Hook to mark a nudge as read
 */
export function useMarkNudgeAsRead() {
  const updateStatus = useUpdateNudgeStatus();
  
  return (nudgeId: string) => {
    return updateStatus.mutate({ nudgeId, status: 'read' });
  };
}

/**
 * Hook to dismiss a nudge
 */
export function useDismissNudge() {
  const updateStatus = useUpdateNudgeStatus();
  
  return (nudgeId: string) => {
    return updateStatus.mutate({ nudgeId, status: 'dismissed' });
  };
}
