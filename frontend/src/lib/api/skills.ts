// Skills API Module
// Handles skill catalog operations

import { apiClient } from './client';
import { Skill } from '@/types';

export interface SkillListResponse {
  skills: Skill[];
  total: number;
}

export interface SkillCoverage {
  skillId: string;
  skillName: string;
  usersWithSkill: number;
  coveragePercentage: number;
  avgProficiency: number;
}

export interface SkillCoverageResponse {
  coverage: SkillCoverage[];
}

export interface SkillGap {
  skillId: string;
  skillName: string;
  requiredByTasks: number;
}

export interface SkillGapsResponse {
  gaps: SkillGap[];
}

/**
 * Get all skills
 */
export async function getSkills(category?: string) {
  return apiClient.get<SkillListResponse>('skills', category ? { category } : undefined);
}

/**
 * Get a single skill by ID
 */
export async function getSkill(skillId: string) {
  return apiClient.get<Skill>(`skills/${skillId}`);
}

/**
 * Get skill coverage across team
 */
export async function getSkillCoverage() {
  return apiClient.get<SkillCoverageResponse>('skills/coverage');
}

/**
 * Get skill gaps (skills required by tasks but not available)
 */
export async function getSkillGaps() {
  return apiClient.get<SkillGapsResponse>('skills/gaps');
}

/**
 * Get skills by category
 */
export async function getSkillsByCategory(category: string) {
  return apiClient.get<SkillListResponse>('skills', { category });
}

// Export as namespace
export const skillsApi = {
  getSkills,
  getSkill,
  getSkillCoverage,
  getSkillGaps,
  getSkillsByCategory,
};

export default skillsApi;
