// React Query hooks for Users API

import {
  useQuery,
  useMutation,
  useQueryClient,
  UseQueryOptions,
} from '@tanstack/react-query';
import { usersApi } from '@/lib/api';
import { User } from '@/types';
import { UserWorkload } from '@/lib/api/users';
import { ApiResponse } from '@/lib/api/types';
import { ApiClientError } from '@/lib/api/client';

// Query keys
export const userKeys = {
  all: ['users'] as const,
  lists: () => [...userKeys.all, 'list'] as const,
  list: (role?: string) => [...userKeys.lists(), role] as const,
  details: () => [...userKeys.all, 'detail'] as const,
  detail: (id: string) => [...userKeys.details(), id] as const,
  skills: (id: string) => [...userKeys.detail(id), 'skills'] as const,
  workload: (id: string) => [...userKeys.detail(id), 'workload'] as const,
};

/**
 * Hook to fetch all users
 */
export function useUsers(
  role?: string,
  options?: UseQueryOptions<ApiResponse<{ users: User[]; total: number }>, ApiClientError>
) {
  return useQuery({
    queryKey: userKeys.list(role),
    queryFn: () => usersApi.getUsers(role),
    ...options,
  });
}

/**
 * Hook to fetch a single user
 */
export function useUser(
  userId: string,
  options?: UseQueryOptions<ApiResponse<User>, ApiClientError>
) {
  return useQuery({
    queryKey: userKeys.detail(userId),
    queryFn: () => usersApi.getUser(userId),
    enabled: !!userId,
    ...options,
  });
}

/**
 * Hook to fetch users by role
 */
export function useUsersByRole(
  role: User['role'],
  options?: UseQueryOptions<ApiResponse<{ users: User[]; total: number }>, ApiClientError>
) {
  return useQuery({
    queryKey: [...userKeys.list(role), 'byRole'],
    queryFn: () => usersApi.getUsersByRole(role),
    ...options,
  });
}

/**
 * Hook to fetch user's skills
 */
export function useUserSkills(
  userId: string,
  options?: UseQueryOptions<
    ApiResponse<{ userId: string; skills: Array<{ skill: { id: string; name: string; category: string }; proficiency: number }> }>,
    ApiClientError
  >
) {
  return useQuery({
    queryKey: userKeys.skills(userId),
    queryFn: () => usersApi.getUserSkills(userId),
    enabled: !!userId,
    ...options,
  });
}

/**
 * Hook to fetch user's workload
 */
export function useUserWorkload(
  userId: string,
  options?: UseQueryOptions<
    ApiResponse<UserWorkload>,
    ApiClientError
  >
) {
  return useQuery({
    queryKey: userKeys.workload(userId),
    queryFn: () => usersApi.getUserWorkload(userId),
    enabled: !!userId,
    ...options,
  });
}
