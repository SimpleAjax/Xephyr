// React Query hooks for Skills API

import {
  useQuery,
  useQueryClient,
  UseQueryOptions,
} from '@tanstack/react-query';
import { skillsApi } from '@/lib/api';
import { Skill } from '@/types';
import { ApiResponse } from '@/lib/api/types';
import { ApiClientError } from '@/lib/api/client';

// Query keys
export const skillKeys = {
  all: ['skills'] as const,
  lists: () => [...skillKeys.all, 'list'] as const,
  list: (category?: string) => [...skillKeys.lists(), category] as const,
  details: () => [...skillKeys.all, 'detail'] as const,
  detail: (id: string) => [...skillKeys.details(), id] as const,
  coverage: () => [...skillKeys.all, 'coverage'] as const,
  gaps: () => [...skillKeys.all, 'gaps'] as const,
};

/**
 * Hook to fetch all skills
 */
export function useSkills(
  category?: string,
  options?: UseQueryOptions<ApiResponse<{ skills: Skill[]; total: number }>, ApiClientError>
) {
  return useQuery({
    queryKey: skillKeys.list(category),
    queryFn: () => skillsApi.getSkills(category),
    ...options,
  });
}

/**
 * Hook to fetch a single skill
 */
export function useSkill(
  skillId: string,
  options?: UseQueryOptions<ApiResponse<Skill>, ApiClientError>
) {
  return useQuery({
    queryKey: skillKeys.detail(skillId),
    queryFn: () => skillsApi.getSkill(skillId),
    enabled: !!skillId,
    ...options,
  });
}

/**
 * Hook to fetch skill coverage across team
 */
export function useSkillCoverage(
  options?: UseQueryOptions<
    ApiResponse<{
      coverage: Array<{
        skillId: string;
        skillName: string;
        usersWithSkill: number;
        coveragePercentage: number;
        avgProficiency: number;
      }>;
    }>,
    ApiClientError
  >
) {
  return useQuery({
    queryKey: skillKeys.coverage(),
    queryFn: () => skillsApi.getSkillCoverage(),
    ...options,
  });
}

/**
 * Hook to fetch skill gaps
 */
export function useSkillGaps(
  options?: UseQueryOptions<
    ApiResponse<{
      gaps: Array<{
        skillId: string;
        skillName: string;
        requiredByTasks: number;
      }>;
    }>,
    ApiClientError
  >
) {
  return useQuery({
    queryKey: skillKeys.gaps(),
    queryFn: () => skillsApi.getSkillGaps(),
    ...options,
  });
}
