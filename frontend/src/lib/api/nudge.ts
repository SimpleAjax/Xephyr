// Nudge API Module
// Handles AI-generated nudges and recommendations

import { apiClient } from './client';
import {
  Nudge,
  NudgeListResponse,
  NudgeQueryParams,
  NudgeActionRequest,
  NudgeActionResponse,
  NudgeStats,
} from './types';

/**
 * List nudges with filters
 */
export async function getNudges(params?: NudgeQueryParams) {
  // Handle array params
  const queryParams: Record<string, string> = {};
  
  if (params) {
    if (params.status) {
      queryParams.status = Array.isArray(params.status) 
        ? params.status.join(',') 
        : params.status;
    }
    if (params.severity) {
      queryParams.severity = Array.isArray(params.severity) 
        ? params.severity.join(',') 
        : params.severity;
    }
    if (params.type) {
      queryParams.type = Array.isArray(params.type) 
        ? params.type.join(',') 
        : params.type;
    }
    if (params.projectId) queryParams.projectId = params.projectId;
    if (params.personId) queryParams.personId = params.personId;
    if (params.limit) queryParams.limit = String(params.limit);
    if (params.offset) queryParams.offset = String(params.offset);
    if (params.cursor) queryParams.cursor = params.cursor;
  }

  return apiClient.get<NudgeListResponse>('nudges', queryParams);
}

/**
 * Get a single nudge by ID
 */
export async function getNudge(nudgeId: string) {
  return apiClient.get<Nudge>(`nudges/${nudgeId}`);
}

/**
 * Take action on a nudge
 */
export async function takeNudgeAction(
  nudgeId: string,
  request: NudgeActionRequest
) {
  return apiClient.post<NudgeActionResponse>(
    `nudges/${nudgeId}/actions`,
    request
  );
}

/**
 * Update nudge status (mark as read, etc.)
 */
export async function updateNudgeStatus(
  nudgeId: string,
  status: 'read' | 'dismissed' | 'acted'
) {
  return apiClient.patch<Nudge>(`nudges/${nudgeId}/status`, { status });
}

/**
 * Trigger manual nudge generation
 */
export async function generateNudges(params: {
  scope: 'project' | 'organization';
  projectId?: string;
  types?: string[];
  async?: boolean;
}) {
  return apiClient.post<{ jobId?: string; status: string }>(
    'nudges/generate',
    params
  );
}

/**
 * Get nudge statistics
 */
export async function getNudgeStats(period: string = '30d') {
  return apiClient.get<NudgeStats>('nudges/stats', { period });
}

// Export as namespace
export const nudgeApi = {
  getNudges,
  getNudge,
  takeNudgeAction,
  updateNudgeStatus,
  generateNudges,
  getNudgeStats,
};

export default nudgeApi;
