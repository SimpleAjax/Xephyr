// Scenario API Module
// Handles "what-if" scenario simulation

import { apiClient } from './client';
import {
  Scenario,
  CreateScenarioRequest,
  SimulateRequest,
  SimulateResponse,
  ApplyScenarioRequest,
  ApplyScenarioResponse,
} from './types';

/**
 * List all scenarios
 */
export async function getScenarios(status?: string) {
  return apiClient.get<{ scenarios: Scenario[] }>('scenarios', status ? { status } : undefined);
}

/**
 * Get a single scenario by ID
 */
export async function getScenario(scenarioId: string) {
  return apiClient.get<Scenario>(`scenarios/${scenarioId}`);
}

/**
 * Create a new scenario
 */
export async function createScenario(request: CreateScenarioRequest) {
  return apiClient.post<Scenario>('scenarios', request);
}

/**
 * Run simulation on a scenario
 */
export async function simulateScenario(
  scenarioId: string,
  request?: SimulateRequest
) {
  return apiClient.post<SimulateResponse>(
    `scenarios/${scenarioId}/simulate`,
    request || { depth: 'full', includeRecommendations: true }
  );
}

/**
 * Apply a scenario to the project
 */
export async function applyScenario(
  scenarioId: string,
  request?: ApplyScenarioRequest
) {
  return apiClient.post<ApplyScenarioResponse>(
    `scenarios/${scenarioId}/apply`,
    request || { applyRecommendations: false, notifyStakeholders: true }
  );
}

/**
 * Reject a scenario
 */
export async function rejectScenario(scenarioId: string, reason?: string) {
  return apiClient.post<Scenario>(`scenarios/${scenarioId}/reject`, { reason });
}

/**
 * Modify an existing scenario
 */
export async function modifyScenario(
  scenarioId: string,
  updates: Partial<CreateScenarioRequest>
) {
  return apiClient.patch<Scenario>(`scenarios/${scenarioId}/modify`, updates);
}

// Export as namespace
export const scenarioApi = {
  getScenarios,
  getScenario,
  createScenario,
  simulateScenario,
  applyScenario,
  rejectScenario,
  modifyScenario,
};

export default scenarioApi;
