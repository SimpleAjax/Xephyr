// React Query hooks for Scenario API

import {
  useQuery,
  useMutation,
  useQueryClient,
  UseQueryOptions,
} from '@tanstack/react-query';
import { scenarioApi } from '@/lib/api';
import {
  Scenario,
  CreateScenarioRequest,
  SimulateResponse,
  ApplyScenarioResponse,
} from '@/lib/api/types';
import { ApiResponse } from '@/lib/api/types';
import { ApiClientError } from '@/lib/api/client';

// Query keys
export const scenarioKeys = {
  all: ['scenarios'] as const,
  lists: () => [...scenarioKeys.all, 'list'] as const,
  list: (status?: string) => [...scenarioKeys.lists(), status] as const,
  details: () => [...scenarioKeys.all, 'detail'] as const,
  detail: (id: string) => [...scenarioKeys.details(), id] as const,
  simulation: (id: string) => [...scenarioKeys.all, 'simulation', id] as const,
};

/**
 * Hook to fetch all scenarios
 */
export function useScenarios(
  status?: string,
  options?: UseQueryOptions<ApiResponse<{ scenarios: Scenario[] }>, ApiClientError>
) {
  return useQuery({
    queryKey: scenarioKeys.list(status),
    queryFn: () => scenarioApi.getScenarios(status),
    ...options,
  });
}

/**
 * Hook to fetch a single scenario
 */
export function useScenario(
  scenarioId: string,
  options?: UseQueryOptions<ApiResponse<Scenario>, ApiClientError>
) {
  return useQuery({
    queryKey: scenarioKeys.detail(scenarioId),
    queryFn: () => scenarioApi.getScenario(scenarioId),
    enabled: !!scenarioId,
    ...options,
  });
}

/**
 * Hook to create a new scenario
 */
export function useCreateScenario() {
  const queryClient = useQueryClient();

  return useMutation<ApiResponse<Scenario>, ApiClientError, CreateScenarioRequest>({
    mutationFn: (request) => scenarioApi.createScenario(request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: scenarioKeys.lists() });
    },
  });
}

/**
 * Hook to simulate a scenario
 */
export function useSimulateScenario() {
  const queryClient = useQueryClient();

  return useMutation<ApiResponse<SimulateResponse>, ApiClientError, {
    scenarioId: string;
    depth?: 'quick' | 'full';
    includeRecommendations?: boolean;
  }>({
    mutationFn: ({ scenarioId, depth, includeRecommendations }) =>
      scenarioApi.simulateScenario(scenarioId, { depth, includeRecommendations }),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: scenarioKeys.detail(variables.scenarioId) });
      queryClient.invalidateQueries({ queryKey: scenarioKeys.simulation(variables.scenarioId) });
    },
  });
}

/**
 * Hook to apply a scenario
 */
export function useApplyScenario() {
  const queryClient = useQueryClient();

  return useMutation<ApiResponse<ApplyScenarioResponse>, ApiClientError, {
    scenarioId: string;
    applyRecommendations?: boolean;
    notifyStakeholders?: boolean;
  }>({
    mutationFn: ({ scenarioId, applyRecommendations, notifyStakeholders }) =>
      scenarioApi.applyScenario(scenarioId, { applyRecommendations, notifyStakeholders }),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: scenarioKeys.detail(variables.scenarioId) });
      queryClient.invalidateQueries({ queryKey: scenarioKeys.lists() });
      // Also invalidate related data
      queryClient.invalidateQueries({ queryKey: ['projects'] });
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
      queryClient.invalidateQueries({ queryKey: ['workload'] });
      queryClient.invalidateQueries({ queryKey: ['health'] });
    },
  });
}

/**
 * Hook to reject a scenario
 */
export function useRejectScenario() {
  const queryClient = useQueryClient();

  return useMutation<ApiResponse<Scenario>, ApiClientError, {
    scenarioId: string;
    reason?: string;
  }>({
    mutationFn: ({ scenarioId, reason }) => scenarioApi.rejectScenario(scenarioId, reason),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: scenarioKeys.detail(variables.scenarioId) });
      queryClient.invalidateQueries({ queryKey: scenarioKeys.lists() });
    },
  });
}
