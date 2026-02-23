import { create } from 'zustand';
import { Task, TaskStatus, Nudge, NudgeStatus, Project } from '@/types';
import { mockData } from '@/data/mock-data';

interface AppState {
  // Nudges
  nudges: Nudge[];
  unreadNudgesCount: number;
  markNudgeAsRead: (id: string) => void;
  dismissNudge: (id: string) => void;
  takeNudgeAction: (id: string) => void;
  
  // Projects
  projects: Project[];
  createProject: (project: Project) => void;
  updateProject: (projectId: string, updates: Partial<Project>) => void;
  
  // Tasks
  tasks: Task[];
  updateTaskStatus: (taskId: string, status: TaskStatus) => void;
  assignTask: (taskId: string, personId: string) => void;
  createTask: (task: Partial<Task>) => void;
  updateTask: (taskId: string, updates: Partial<Task>) => void;
  
  // UI State
  selectedProjectId: string | null;
  setSelectedProjectId: (id: string | null) => void;
  
  // Refresh data
  refreshNudges: () => void;
  refreshProjects: () => void;
  refreshTasks: () => void;
}

export const useAppStore = create<AppState>((set, get) => ({
  // Initial state
  nudges: mockData.getNudges(),
  unreadNudgesCount: mockData.getNudges().filter(n => n.status === 'unread').length,
  projects: mockData.getProjects(),
  tasks: mockData.getAllTasks(),
  selectedProjectId: null,

  // Nudge actions
  markNudgeAsRead: (id: string) => {
    set((state) => ({
      nudges: state.nudges.map(n => 
        n.id === id ? { ...n, status: 'read' as NudgeStatus } : n
      ),
      unreadNudgesCount: Math.max(0, state.unreadNudgesCount - 1),
    }));
  },

  dismissNudge: (id: string) => {
    const nudge = get().nudges.find(n => n.id === id);
    set((state) => ({
      nudges: state.nudges.map(n => 
        n.id === id ? { ...n, status: 'dismissed' as NudgeStatus } : n
      ),
      unreadNudgesCount: nudge?.status === 'unread' 
        ? Math.max(0, state.unreadNudgesCount - 1) 
        : state.unreadNudgesCount,
    }));
  },

  takeNudgeAction: (id: string) => {
    const nudge = get().nudges.find(n => n.id === id);
    set((state) => ({
      nudges: state.nudges.map(n => 
        n.id === id ? { ...n, status: 'acted' as NudgeStatus } : n
      ),
      unreadNudgesCount: nudge?.status === 'unread' 
        ? Math.max(0, state.unreadNudgesCount - 1) 
        : state.unreadNudgesCount,
    }));
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

  // Task actions
  updateTaskStatus: (taskId: string, status: TaskStatus) => {
    set((state) => ({
      tasks: state.tasks.map(t => 
        t.id === taskId ? { ...t, status } : t
      ),
    }));
  },

  assignTask: (taskId: string, personId: string) => {
    set((state) => ({
      tasks: state.tasks.map(t => 
        t.id === taskId ? { ...t, assigneeId: personId } : t
      ),
    }));
  },

  createTask: (task: Partial<Task>) => {
    const newTask: Task = {
      id: `task-${Date.now()}`,
      projectId: task.projectId || mockData.getProjects()[0].id,
      title: task.title || 'New Task',
      description: task.description || '',
      status: 'backlog',
      priority: task.priority || 'medium',
      priorityScore: 50,
      businessValue: 50,
      estimatedHours: task.estimatedHours || 8,
      requiredSkills: task.requiredSkills || [],
      isMilestone: task.isMilestone || false,
      isCriticalPath: task.isCriticalPath || false,
      hierarchyLevel: task.hierarchyLevel || 1,
      parentTaskId: task.parentTaskId,
      assigneeId: task.assigneeId,
      dependencies: [],
      blockedBy: [],
      ...task,
    };
    
    set((state) => ({
      tasks: [...state.tasks, newTask],
    }));
  },

  updateTask: (taskId: string, updates: Partial<Task>) => {
    set((state) => ({
      tasks: state.tasks.map(t => 
        t.id === taskId ? { ...t, ...updates } : t
      ),
    }));
  },

  // UI actions
  setSelectedProjectId: (id: string | null) => {
    set({ selectedProjectId: id });
  },

  refreshNudges: () => {
    const nudges = mockData.getNudges();
    set({ 
      nudges,
      unreadNudgesCount: nudges.filter(n => n.status === 'unread').length,
    });
  },

  refreshProjects: () => {
    set({ projects: mockData.getProjects() });
  },

  refreshTasks: () => {
    set({ tasks: mockData.getAllTasks() });
  },
}));
