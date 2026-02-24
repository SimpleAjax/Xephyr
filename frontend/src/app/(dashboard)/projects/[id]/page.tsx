'use client';

import { useState, useMemo } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Progress } from '@/components/ui/progress';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { useProject, useTasksByProject, useUsers } from '@/hooks/api';
import { useAppStore } from '@/store/app-store';
import { Task, User } from '@/types';
import { User as UserIcon } from 'lucide-react';
import { 
  ArrowLeft, 
  MoreHorizontal, 
  Plus, 
  Calendar, 
  Clock, 
  AlertTriangle,
  CheckCircle2,
  ChevronRight,
  ChevronDown,
  GripVertical,
  LayoutList,
  Kanban,
  Sparkles,
  FolderKanban,
  Loader2
} from 'lucide-react';
import Link from 'next/link';
import { CreateTaskDialog } from '@/components/tasks/create-task-dialog';

interface ProjectPageProps {
  params: { id: string };
}

export default function ProjectPage({ params }: ProjectPageProps) {
  const { id } = params;
  const [viewMode, setViewMode] = useState<'list' | 'hierarchy'>('hierarchy');
  const [createTaskOpen, setCreateTaskOpen] = useState(false);
  const [expandedTasks, setExpandedTasks] = useState<Set<string>>(new Set());
  
  const { data: projectData, isLoading: projectLoading } = useProject(id);
  const { data: tasksData, isLoading: tasksLoading } = useTasksByProject(id);
  const { data: usersData } = useUsers();
  
  const project = projectData?.data;
  const tasks = tasksData?.data?.tasks || [];
  const users = usersData?.data?.users || [];
  
  if (projectLoading) {
    return (
      <div className="p-8 h-full flex items-center justify-center">
        <div className="flex items-center gap-2 text-muted-foreground">
          <Loader2 className="w-5 h-5 animate-spin" />
          <span>Loading project...</span>
        </div>
      </div>
    );
  }
  
  if (!project) {
    return (
      <div className="p-8">
        <Card>
          <CardContent className="p-8 text-center">
            <AlertTriangle className="w-12 h-12 text-destructive mx-auto mb-4" />
            <h1 className="text-xl font-bold">Project Not Found</h1>
            <p className="text-muted-foreground">The project you're looking for doesn't exist.</p>
            <Button className="mt-4" asChild>
              <Link href="/projects">Back to Projects</Link>
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  const projectTasks = tasks.filter((t: Task) => t.hierarchyLevel === 1);
  const subtasks = tasks.filter((t: Task) => t.hierarchyLevel === 2);
  
  const stats = {
    total: tasks.length,
    completed: tasks.filter((t: Task) => t.status === 'done').length,
    inProgress: tasks.filter((t: Task) => t.status === 'in_progress').length,
    backlog: tasks.filter((t: Task) => t.status === 'backlog').length,
    totalHours: tasks.reduce((sum: number, t: Task) => sum + (t.estimatedHours || 0), 0),
    unassigned: tasks.filter((t: Task) => !t.assigneeId && t.status !== 'done').length,
  };

  const getStatusColor = (status: Task['status']) => {
    const colors: Record<string, string> = {
      backlog: 'bg-slate-100 text-slate-700',
      ready: 'bg-blue-100 text-blue-700',
      in_progress: 'bg-yellow-100 text-yellow-700',
      review: 'bg-purple-100 text-purple-700',
      done: 'bg-green-100 text-green-700',
    };
    return colors[status] || 'bg-slate-100';
  };

  const getPriorityColor = (priority: Task['priority']) => {
    const colors: Record<string, string> = {
      low: 'bg-slate-100 text-slate-700',
      medium: 'bg-blue-100 text-blue-700',
      high: 'bg-orange-100 text-orange-700',
      critical: 'bg-red-100 text-red-700',
    };
    return colors[priority] || 'bg-slate-100';
  };

  const toggleExpand = (taskId: string) => {
    setExpandedTasks(prev => {
      const newSet = new Set(prev);
      if (newSet.has(taskId)) {
        newSet.delete(taskId);
      } else {
        newSet.add(taskId);
      }
      return newSet;
    });
  };

  const getSubtasks = (parentId: string) => {
    return subtasks.filter((t: Task) => t.parentTaskId === parentId);
  };

  const getAssigneeName = (assigneeId?: string) => {
    if (!assigneeId) return null;
    return users.find((u: User) => u.id === assigneeId)?.name;
  };

  const renderTaskRow = (task: Task, isSubtask = false) => {
    const taskSubtasks = getSubtasks(task.id);
    const hasSubtasks = taskSubtasks.length > 0;
    const isExpanded = expandedTasks.has(task.id);

    return (
      <div key={task.id}>
        <div 
          className={`flex items-center gap-3 p-3 hover:bg-muted/50 rounded-lg group ${
            isSubtask ? 'ml-8 border-l-2 border-l-primary/30 pl-4' : ''
          }`}
        >
          {/* Expand/Collapse */}
          {hasSubtasks ? (
            <button 
              onClick={() => toggleExpand(task.id)}
              className="w-5 h-5 flex items-center justify-center text-muted-foreground hover:text-foreground"
            >
              {isExpanded ? <ChevronDown className="w-4 h-4" /> : <ChevronRight className="w-4 h-4" />}
            </button>
          ) : (
            <div className="w-5" />
          )}

          {/* Drag Handle */}
          {!isSubtask && <GripVertical className="w-4 h-4 text-muted-foreground cursor-move" />}

          {/* Task Info */}
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2">
              <p className={`font-medium truncate ${isSubtask ? 'text-sm' : ''}`}>
                {task.title}
              </p>
              {task.isCriticalPath && (
                <Badge variant="outline" className="text-xs text-destructive border-destructive">
                  Critical
                </Badge>
              )}
              {task.isMilestone && (
                <Badge variant="outline" className="text-xs">
                  Milestone
                </Badge>
              )}
            </div>
            {task.description && (
              <p className="text-xs text-muted-foreground truncate">{task.description}</p>
            )}
          </div>

          {/* Meta */}
          <div className="flex items-center gap-4 text-sm">
            <Badge className={`text-xs ${getPriorityColor(task.priority)}`}>
              {task.priority}
            </Badge>
            <Badge variant="outline" className={`text-xs ${getStatusColor(task.status)}`}>
              {task.status.replace('_', ' ')}
            </Badge>
            
            {/* Assignee */}
            {task.assigneeId ? (
              <div className="flex items-center gap-1 text-muted-foreground">
                <div className="w-5 h-5 rounded-full bg-primary/20 flex items-center justify-center text-xs">
                  {getAssigneeName(task.assigneeId)?.charAt(0)}
                </div>
                <span className="hidden sm:inline text-xs">{getAssigneeName(task.assigneeId)}</span>
              </div>
            ) : (
              <Badge variant="secondary" className="text-xs">Unassigned</Badge>
            )}

            {/* Hours */}
            <div className="flex items-center gap-1 text-muted-foreground text-xs">
              <Clock className="w-3 h-3" />
              {task.estimatedHours}h
            </div>

            {/* Due Date */}
            {task.dueDate && (
              <div className="flex items-center gap-1 text-muted-foreground text-xs">
                <Calendar className="w-3 h-3" />
                {new Date(task.dueDate).toLocaleDateString(undefined, { month: 'short', day: 'numeric' })}
              </div>
            )}

            <Button variant="ghost" size="icon" className="h-7 w-7 opacity-0 group-hover:opacity-100">
              <MoreHorizontal className="w-4 h-4" />
            </Button>
          </div>
        </div>

        {/* Subtasks */}
        {hasSubtasks && isExpanded && (
          <div className="mt-1">
            {taskSubtasks.map((subtask: Task) => renderTaskRow(subtask, true))}
          </div>
        )}
      </div>
    );
  };

  return (
    <div className="p-8 space-y-6">
      {/* Header */}
      <div className="flex items-start justify-between">
        <div className="space-y-1">
          <div className="flex items-center gap-2 text-muted-foreground">
            <Button variant="ghost" size="sm" className="h-auto p-0" asChild>
              <Link href="/projects">
                <ArrowLeft className="w-4 h-4 mr-1" />
                Projects
              </Link>
            </Button>
            <ChevronRight className="w-4 h-4" />
            <span>{project.name}</span>
          </div>
          <h1 className="text-3xl font-bold">{project.name}</h1>
          <p className="text-muted-foreground max-w-2xl">{project.description}</p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" asChild>
            <Link href={`/kanban?project=${project.id}`}>
              <Kanban className="w-4 h-4 mr-2" />
              Kanban View
            </Link>
          </Button>
          <Button onClick={() => setCreateTaskOpen(true)}>
            <Plus className="w-4 h-4 mr-2" />
            Add Task
          </Button>
        </div>
      </div>

      {/* Project Stats */}
      <div className="grid gap-4 md:grid-cols-6">
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Progress</CardDescription>
            <CardTitle className="text-2xl">{project.progress}%</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Health Score</CardDescription>
            <CardTitle className={`text-2xl ${project.healthScore >= 80 ? 'text-green-500' : project.healthScore >= 60 ? 'text-yellow-500' : 'text-destructive'}`}>
              {project.healthScore}%
            </CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Tasks</CardDescription>
            <CardTitle className="text-2xl">{stats.total}</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Estimated Hours</CardDescription>
            <CardTitle className="text-2xl">{stats.totalHours}h</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Unassigned</CardDescription>
            <CardTitle className={`text-2xl ${stats.unassigned > 0 ? 'text-destructive' : 'text-green-500'}`}>
              {stats.unassigned}
            </CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Due Date</CardDescription>
            <CardTitle className="text-lg">
              {new Date(project.targetEndDate).toLocaleDateString()}
            </CardTitle>
          </CardHeader>
        </Card>
      </div>

      {/* Tasks Section */}
      <Tabs defaultValue="tasks" className="space-y-4">
        <div className="flex items-center justify-between">
          <TabsList>
            <TabsTrigger value="tasks">
              <LayoutList className="w-4 h-4 mr-2" />
              Tasks
            </TabsTrigger>
            <TabsTrigger value="team">
              <UserIcon className="w-4 h-4 mr-2" />
              Team
            </TabsTrigger>
          </TabsList>

          <div className="flex items-center gap-2">
            <Button 
              variant={viewMode === 'hierarchy' ? 'secondary' : 'ghost'} 
              size="sm"
              onClick={() => setViewMode('hierarchy')}
            >
              <FolderKanban className="w-4 h-4 mr-2" />
              Hierarchy
            </Button>
            <Button 
              variant={viewMode === 'list' ? 'secondary' : 'ghost'} 
              size="sm"
              onClick={() => setViewMode('list')}
            >
              <LayoutList className="w-4 h-4 mr-2" />
              List
            </Button>
          </div>
        </div>

        <TabsContent value="tasks" className="space-y-4">
          <Card>
            <CardHeader className="pb-3">
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle>Project Tasks</CardTitle>
                  <CardDescription>
                    {tasksLoading ? 'Loading...' : `${projectTasks.length} main tasks, ${subtasks.length} subtasks`}
                  </CardDescription>
                </div>
                <Button size="sm" variant="outline" onClick={() => setCreateTaskOpen(true)}>
                  <Sparkles className="w-4 h-4 mr-2" />
                  AI Generate Tasks
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              {tasksLoading ? (
                <div className="flex items-center justify-center py-12">
                  <Loader2 className="w-5 h-5 animate-spin text-muted-foreground" />
                </div>
              ) : projectTasks.length === 0 ? (
                <div className="text-center py-12">
                  <FolderKanban className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                  <h3 className="font-medium">No tasks yet</h3>
                  <p className="text-sm text-muted-foreground mb-4">
                    Get started by adding tasks to this project
                  </p>
                  <Button onClick={() => setCreateTaskOpen(true)}>
                    <Plus className="w-4 h-4 mr-2" />
                    Add First Task
                  </Button>
                </div>
              ) : (
                <div className="space-y-1">
                  {projectTasks.map(task => renderTaskRow(task))}
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="team" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Project Team</CardTitle>
              <CardDescription>Team members assigned to this project</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid gap-4 md:grid-cols-3">
                {users
                  .filter((u: User) => tasks.some((t: Task) => t.assigneeId === u.id))
                  .map((user: User) => {
                    const userTasks = tasks.filter((t: Task) => t.assigneeId === user.id);
                    const completedTasks = userTasks.filter((t: Task) => t.status === 'done').length;
                    
                    return (
                      <Card key={user.id}>
                        <CardContent className="p-4">
                          <div className="flex items-center gap-3">
                            <div className="w-10 h-10 rounded-full bg-primary/20 flex items-center justify-center font-medium">
                              {user.name.charAt(0)}
                            </div>
                            <div>
                              <p className="font-medium">{user.name}</p>
                              <p className="text-sm text-muted-foreground capitalize">{user.role}</p>
                            </div>
                          </div>
                          <div className="flex items-center gap-4 mt-4 text-sm">
                            <span>{userTasks.length} tasks</span>
                            <span className="text-green-500">{completedTasks} done</span>
                          </div>
                        </CardContent>
                      </Card>
                    );
                  })}
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      <CreateTaskDialog 
        open={createTaskOpen} 
        onOpenChange={setCreateTaskOpen}
        defaultProjectId={project.id}
      />
    </div>
  );
}
