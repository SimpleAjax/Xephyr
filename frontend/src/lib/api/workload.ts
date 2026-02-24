// Workload API Module
// Handles team and individual workload management

import { apiClient } from './client';
import {
  TeamWorkload,
  PersonWorkload,
  WorkloadForecast,
} from './types';

/**
 * Get team workload overview
 */
export async function getTeamWorkload(
  week?: string,
  includeForecast: boolean = false
) {
  const params: Record<string, string> = { includeForecast: String(includeForecast) };
  if (week) {
    params.week = week;
  }
  return apiClient.get<TeamWorkload>('workload/team', params);
}

/**
 * Get individual person workload
 */
export async function getPersonWorkload(
  personId: string,
  week?: string
) {
  const params: Record<string, string> = {};
  if (week) {
    params.week = week;
  }
  return apiClient.get<PersonWorkload>(`workload/people/${personId}`, params);
}

/**
 * Get workload forecast for a person
 */
export async function getWorkloadForecast(
  personId: string,
  weeks: number = 8
) {
  return apiClient.get<WorkloadForecast>('workload/forecast', {
    personId,
    weeks,
  });
}

/**
 * Get workload analytics
 */
export async function getWorkloadAnalytics(period: string = '30d') {
  return apiClient.get<{
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
  }>('workload/analytics', { period });
}

// Export as namespace
export const workloadApi = {
  getTeamWorkload,
  getPersonWorkload,
  getWorkloadForecast,
  getWorkloadAnalytics,
};

export default workloadApi;
