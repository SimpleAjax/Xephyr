// Health API Module
// Handles portfolio and project health monitoring

import { apiClient } from './client';
import {
  PortfolioHealth,
  ProjectHealth,
  HealthTrend,
  HealthTrendParams,
} from './types';

/**
 * Get portfolio-wide health overview
 */
export async function getPortfolioHealth() {
  return apiClient.get<PortfolioHealth>('health/portfolio');
}

/**
 * Get detailed health metrics for a specific project
 */
export async function getProjectHealth(
  projectId: string,
  includeBreakdown: boolean = true
) {
  return apiClient.get<ProjectHealth>(`health/projects/${projectId}`, {
    includeBreakdown,
  });
}

/**
 * Get health data for multiple projects
 */
export async function getBulkProjectHealth(projectIds?: string[]) {
  const params: Record<string, string> = {};
  if (projectIds && projectIds.length > 0) {
    params.projectIds = projectIds.join(',');
  }
  return apiClient.get<{ projects: ProjectHealth[] }>('health/projects', params);
}

/**
 * Get health trend data over time
 */
export async function getHealthTrends(params?: HealthTrendParams) {
  const queryParams: Record<string, string> = {};
  if (params?.projectId) queryParams.projectId = params.projectId;
  if (params?.days) queryParams.days = String(params.days);
  return apiClient.get<HealthTrend>('health/trends', queryParams);
}

// Export as namespace
export const healthApi = {
  getPortfolioHealth,
  getProjectHealth,
  getBulkProjectHealth,
  getHealthTrends,
};

export default healthApi;
