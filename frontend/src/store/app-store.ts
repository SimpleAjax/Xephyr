import { create } from 'zustand';
import { Task, TaskStatus, Nudge, NudgeStatus, Project } from '@/types';
import { 
  nudgeApi, 
  tasksApi, 
  projectsApi,
  type UpdateTaskRequest 
} from '@/lib/api';
import { useQueryClient } from '@tanstack/react-query';

interface AppState {
  // Nudges (local state for UI, synced with API)
  nudges: Nudge[];
  unreadNudgesCount: number;
  setNudges: (nudges: Nudge[]) => void;
  markNudgeAsRead: (id: string) => Promise<void>;
  dismissNudge: (id: string) => Promise<void>;
  takeNudgeAction: (id: string) => Promise<void>;
  
  // Projects (local cache)
  projects: Project[];
  setProjects: (projects: Project[]) => void;
  createProject: (project: Project) => void;
  updateProject: (projectId: string, updates: Partial<Project>) => void;
  
  // Tasks (local cache)
  tasks: Task[];
  setTasks: (tasks: Task[]) => void;
  updateTaskStatus: (taskId: string, status: TaskStatus) => Promise<void>;
  assignTask: (taskId: string, personId: string) => Promise<void>;
  createTask: (task: Partial<Task>) => Promise<void>;
  updateTask: (taskId: string, updates: Partial<Task>) => Promise<void>;
  
  // UI State
  selectedProjectId: string | null;
  setSelectedProjectId: (id: string | null) => void;
  
  // Loading states
  isLoading: boolean;
  error: string | null;
}

export const useAppStore = create<AppState>((set, get) => ({
  // Initial state
  nudges: [],
  unreadNudgesCount: 0,
  projects: [],
  tasks: [],
  selectedProjectId: null,
  isLoading: false,
  error: null,

  // Setters for data from React Query
  setNudges: (nudges) => {
    set({ 
      nudges,
      unreadNudgesCount: nudges.filter(n => n.status === 'unread').length,
    });
  },

  setProjects: (projects) => set({ projects }),
  
  setTasks: (tasks) => set({ tasks }),

  // Nudge actions (sync with API)
  markNudgeAsRead: async (id: string) => {
    try {
      await nudgeApi.updateNudgeStatus(id, 'read');
      set((state) => ({
        nudges: state.nudges.map(n => 
          n.id === id ? { ...n, status: 'read' as NudgeStatus } : n
        ),
        unreadNudgesCount: Math.max(0, state.unreadNudgesCount - 1),
      }));
    } catch (error) {
      console.error('Failed to mark nudge as read:', error);
    }
  },

  dismissNudge: async (id: string) => {
    try {
      const nudge = get().nudges.find(n => n.id === id);
      await nudgeApi.updateNudgeStatus(id, 'dismissed');
      set((state) => ({
        nudges: state.nudges.map(n => 
          n.id === id ? { ...n, status: 'dismissed' as NudgeStatus } : n
        ),
        unreadNudgesCount: nudge?.status === 'unread' 
          ? Math.max(0, state.unreadNudgesCount - 1) 
          : state.unreadNudgesCount,
      }));
    } catch (error) {
      console.error('Failed to dismiss nudge:', error);
    }
  },

  takeNudgeAction: async (id: string) => {
    try {
      const nudge = get().nudges.find(n => n.id === id);
      await nudgeApi.takeNudgeAction(id, { actionType: 'accept_suggestion' });
      set((state) => ({
        nudges: state.nudges.map(n => 
          n.id === id ? { ...n, status: 'acted' as NudgeStatus } : n
        ),
        unreadNudgesCount: nudge?.status === 'unread' 
          ? Math.max(0, state.unreadNudgesCount - 1) 
          : state.unreadNudgesCount,
      }));
    } catch (error) {
      console.error('Failed to take nudge action:', error);
    }
  },

  // Project actions
  createProject: (project: Project) => {
    set((state) => ({
      projects: [...state.projects, project],
    }));
  },

  updateProject: (projectId: string, updates: Partial<Project>) => {
    set((state) => ({
      projects: state.projects.map(p => 
        p.id === projectId ? { ...p, ...updates } : p
      ),
    }));
  },

  // Task actions (sync with API)
  updateTaskStatus: async (taskId: string, status: TaskStatus) => {
    try {
      await tasksApi.updateTaskStatus(taskId, { status });
      set((state) => ({
        tasks: state.tasks.map(t => 
          t.id === taskId ? { ...t, status } : t
        ),
      }));
    } catch (error) {
      console.error('Failed to update task status:', error);
    }
  },

  assignTask: async (taskId: string, personId: string) => {
    try {
      await tasksApi.assignTask(taskId, personId);
      set((state) => ({
        tasks: state.tasks.map(t => 
          t.id === taskId ? { ...t, assigneeId: personId } : t
        ),
      }));
    } catch (error) {
      console.error('Failed to assign task:', error);
    }
  },

  createTask: async (task: Partial<Task>) => {
    try {
      const response = await tasksApi.createTask({
        title: task.title || 'New Task',
        description: task.description || '',
        projectId: task.projectId || get().projects[0]?.id || '',
        priority: task.priority || 'medium',
        estimatedHours: task.estimatedHours || 8,
        requiredSkills: task.requiredSkills || [],
        assigneeId: task.assigneeId,
        parentTaskId: task.parentTaskId,
        hierarchyLevel: task.hierarchyLevel || 1,
        isMilestone: task.isMilestone || false,
      });
      
      if (response.data) {
        set((state) => ({
          tasks: [...state.tasks, response.data!],
        }));
      }
    } catch (error) {
      console.error('Failed to create task:', error);
    }
  },

  updateTask: async (taskId: string, updates: Partial<Task>) => {
    try {
      const apiUpdates: UpdateTaskRequest = {};
      if (updates.title !== undefined) apiUpdates.title = updates.title;
      if (updates.description !== undefined) apiUpdates.description = updates.description;
      if (updates.status !== undefined) apiUpdates.status = updates.status;
      if (updates.priority !== undefined) apiUpdates.priority = updates.priority;
      if (updates.estimatedHours !== undefined) apiUpdates.estimatedHours = updates.estimatedHours;
      if (updates.assigneeId !== undefined) apiUpdates.assigneeId = updates.assigneeId;
      if (updates.requiredSkills !== undefined) apiUpdates.requiredSkills = updates.requiredSkills;
      if (updates.isMilestone !== undefined) apiUpdates.isMilestone = updates.isMilestone;
      
      await tasksApi.updateTask(taskId, apiUpdates);
      set((state) => ({
        tasks: state.tasks.map(t => 
          t.id === taskId ? { ...t, ...updates } : t
        ),
      }));
    } catch (error) {
      console.error('Failed to update task:', error);
    }
  },

  // UI actions
  setSelectedProjectId: (id: string | null) => {
    set({ selectedProjectId: id });
  },
}));
