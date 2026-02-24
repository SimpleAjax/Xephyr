// Users API Module
// Handles user/team member operations

import { apiClient } from './client';
import { User, Skill } from '@/types';

export interface UserSkill {
  skill: Skill;
  proficiency: number;
}

export interface UserSkillsResponse {
  userId: string;
  skills: UserSkill[];
}

export interface UserListResponse {
  users: User[];
  total: number;
}

export interface UserWorkload {
  personId: string;
  personName: string;
  allocationPercentage: number;
  assignedTasks: number;
  totalEstimatedHours: number;
  availabilityThisWeek: number;
  availabilityNextWeek: number;
}

/**
 * Get all users/team members
 */
export async function getUsers(role?: string) {
  return apiClient.get<UserListResponse>('users', role ? { role } : undefined);
}

/**
 * Get a single user by ID
 */
export async function getUser(userId: string) {
  return apiClient.get<User>(`users/${userId}`);
}

/**
 * Get user's skills
 */
export async function getUserSkills(userId: string) {
  return apiClient.get<UserSkillsResponse>(`users/${userId}/skills`);
}

/**
 * Get current user's workload
 */
export async function getUserWorkload(userId: string) {
  return apiClient.get<UserWorkload>(`workload/people/${userId}`);
}

/**
 * Get users by role
 */
export async function getUsersByRole(role: User['role']) {
  return apiClient.get<UserListResponse>('users', { role });
}

// Export as namespace
export const usersApi = {
  getUsers,
  getUser,
  getUserSkills,
  getUserWorkload,
  getUsersByRole,
};

export default usersApi;
