// React Query hooks for Projects CRUD API
// Note: useProjects.ts exists for health/progress analytics, this is for basic CRUD

import {
  useQuery,
  useMutation,
  useQueryClient,
  UseQueryOptions,
} from '@tanstack/react-query';
import { projectsApi } from '@/lib/api';
import { Project } from '@/types';
import { ApiResponse } from '@/lib/api/types';
import { ApiClientError } from '@/lib/api/client';

// Query keys
export const projectCrudKeys = {
  all: ['projects', 'crud'] as const,
  lists: () => [...projectCrudKeys.all, 'list'] as const,
  list: (status?: string) => [...projectCrudKeys.lists(), status] as const,
  details: () => [...projectCrudKeys.all, 'detail'] as const,
  detail: (id: string) => [...projectCrudKeys.details(), id] as const,
  team: (id: string) => [...projectCrudKeys.detail(id), 'team'] as const,
};

/**
 * Hook to fetch all projects
 */
export function useProjectsList(
  status?: string,
  options?: UseQueryOptions<ApiResponse<{ projects: Project[]; total: number }>, ApiClientError>
) {
  return useQuery({
    queryKey: projectCrudKeys.list(status),
    queryFn: () => projectsApi.getProjects(status),
    ...options,
  });
}

/**
 * Hook to fetch a single project
 */
export function useProject(
  projectId: string,
  options?: UseQueryOptions<ApiResponse<Project>, ApiClientError>
) {
  return useQuery({
    queryKey: projectCrudKeys.detail(projectId),
    queryFn: () => projectsApi.getProject(projectId),
    enabled: !!projectId,
    ...options,
  });
}

/**
 * Hook to create a new project
 */
export function useCreateProject() {
  const queryClient = useQueryClient();

  return useMutation<
    ApiResponse<Project>,
    ApiClientError,
    Parameters<typeof projectsApi.createProject>[0]
  >({
    mutationFn: (request) => projectsApi.createProject(request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: projectCrudKeys.lists() });
      queryClient.invalidateQueries({ queryKey: ['projects'] });
    },
  });
}

/**
 * Hook to update a project
 */
export function useUpdateProject() {
  const queryClient = useQueryClient();

  return useMutation<
    ApiResponse<Project>,
    ApiClientError,
    { projectId: string; updates: Parameters<typeof projectsApi.updateProject>[1] }
  >({
    mutationFn: ({ projectId, updates }) => projectsApi.updateProject(projectId, updates),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: projectCrudKeys.detail(variables.projectId) });
      queryClient.invalidateQueries({ queryKey: projectCrudKeys.lists() });
      queryClient.invalidateQueries({ queryKey: ['projects'] });
    },
  });
}

/**
 * Hook to delete a project
 */
export function useDeleteProject() {
  const queryClient = useQueryClient();

  return useMutation<ApiResponse<void>, ApiClientError, string>({
    mutationFn: (projectId) => projectsApi.deleteProject(projectId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: projectCrudKeys.lists() });
      queryClient.invalidateQueries({ queryKey: ['projects'] });
    },
  });
}

/**
 * Hook to fetch project team
 */
export function useProjectTeam(
  projectId: string,
  options?: UseQueryOptions<
    ApiResponse<{ users: Array<{ id: string; name: string; role: string }> }>,
    ApiClientError
  >
) {
  return useQuery({
    queryKey: projectCrudKeys.team(projectId),
    queryFn: () => projectsApi.getProjectTeam(projectId),
    enabled: !!projectId,
    ...options,
  });
}
