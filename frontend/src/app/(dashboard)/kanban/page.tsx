'use client';

import { useState, useEffect } from 'react';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { useAppStore } from '@/store/app-store';
import { useTasks, useUsers, useProjectsList } from '@/hooks/api';
import { Task, TaskStatus, Project, User } from '@/types';
import { 
  MoreHorizontal, 
  Calendar, 
  Clock, 
  AlertCircle,
  CheckCircle2,
  AlertTriangle,
  GripVertical,
  Plus,
  Loader2
} from 'lucide-react';

const columns: { id: TaskStatus; title: string; color: string }[] = [
  { id: 'backlog', title: 'Backlog', color: 'bg-slate-100' },
  { id: 'ready', title: 'Ready', color: 'bg-blue-50' },
  { id: 'in_progress', title: 'In Progress', color: 'bg-yellow-50' },
  { id: 'review', title: 'Review', color: 'bg-purple-50' },
  { id: 'done', title: 'Done', color: 'bg-green-50' },
];

const priorityColors = {
  low: 'bg-slate-100 text-slate-700',
  medium: 'bg-blue-100 text-blue-700',
  high: 'bg-orange-100 text-orange-700',
  critical: 'bg-red-100 text-red-700',
};

export default function KanbanPage() {
  const [selectedProject, setSelectedProject] = useState<string>('all');
  const { data: tasksData, isLoading: tasksLoading } = useTasks(selectedProject === 'all' ? undefined : selectedProject);
  const { data: projectsData, isLoading: projectsLoading } = useProjectsList();
  const { data: usersData } = useUsers();
  const updateTaskStatus = useAppStore((state) => state.updateTaskStatus);
  const setTasks = useAppStore((state) => state.setTasks);
  
  const tasks: Task[] = tasksData?.data?.tasks || [];
  const projects: Project[] = projectsData?.data?.projects || [];
  const users: User[] = usersData?.data?.users || [];
  
  // Sync tasks with store
  useEffect(() => {
    if (tasks.length > 0) {
      setTasks(tasks);
    }
  }, [tasks, setTasks]);
  
  const filteredTasks = tasks.filter(task => 
    selectedProject === 'all' || task.projectId === selectedProject
  );

  const getTasksByStatus = (status: TaskStatus) => {
    return filteredTasks.filter(task => task.status === status);
  };

  const getProjectName = (projectId: string) => {
    return projects.find((p: Project) => p.id === projectId)?.name || 'Unknown';
  };

  const getAssigneeName = (assigneeId?: string) => {
    if (!assigneeId) return null;
    return users.find((u: User) => u.id === assigneeId)?.name;
  };

  const handleDragStart = (e: React.DragEvent, taskId: string) => {
    e.dataTransfer.setData('taskId', taskId);
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
  };

  const handleDrop = (e: React.DragEvent, status: TaskStatus) => {
    e.preventDefault();
    const taskId = e.dataTransfer.getData('taskId');
    if (taskId) {
      updateTaskStatus(taskId, status);
    }
  };

  if (tasksLoading || projectsLoading) {
    return (
      <div className="p-8 h-full flex items-center justify-center">
        <div className="flex items-center gap-2 text-muted-foreground">
          <Loader2 className="w-5 h-5 animate-spin" />
          <span>Loading kanban board...</span>
        </div>
      </div>
    );
  }

  return (
    <div className="p-8 h-full flex flex-col">
      {/* Header */}
      <div className="flex justify-between items-start mb-6">
        <div>
          <h1 className="text-3xl font-bold">Kanban Board</h1>
          <p className="text-muted-foreground">Drag and drop tasks to update status</p>
        </div>
        <div className="flex items-center gap-4">
          <Select value={selectedProject} onValueChange={setSelectedProject}>
            <SelectTrigger className="w-[250px]">
              <SelectValue placeholder="Filter by project" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Projects</SelectItem>
              {projects.map((project: Project) => (
                <SelectItem key={project.id} value={project.id}>
                  {project.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <Button>
            <Plus className="w-4 h-4 mr-2" />
            Add Task
          </Button>
        </div>
      </div>

      {/* Kanban Board */}
      <div className="flex-1 overflow-x-auto">
        <div className="flex gap-4 h-full min-w-max">
          {columns.map((column) => (
            <div 
              key={column.id}
              className="w-80 flex flex-col"
              onDragOver={handleDragOver}
              onDrop={(e) => handleDrop(e, column.id)}
            >
              {/* Column Header */}
              <div className={`${column.color} rounded-t-lg p-3 flex items-center justify-between`}>
                <div className="flex items-center gap-2">
                  <h3 className="font-semibold">{column.title}</h3>
                  <Badge variant="secondary" className="text-xs">
                    {getTasksByStatus(column.id).length}
                  </Badge>
                </div>
                <Button variant="ghost" size="icon" className="h-6 w-6">
                  <Plus className="w-4 h-4" />
                </Button>
              </div>
              
              {/* Column Content */}
              <div className={`${column.color} rounded-b-lg flex-1 p-2 space-y-2 min-h-[200px]`}>
                {getTasksByStatus(column.id).map((task) => (
                  <Card 
                    key={task.id} 
                    className="cursor-move hover:shadow-md transition-shadow"
                    draggable
                    onDragStart={(e) => handleDragStart(e, task.id)}
                  >
                    <CardContent className="p-3 space-y-3">
                      {/* Task Header */}
                      <div className="flex items-start justify-between">
                        <div className="flex items-center gap-2">
                          <GripVertical className="w-4 h-4 text-muted-foreground" />
                          {task.isCriticalPath && (
                            <AlertTriangle className="w-4 h-4 text-destructive" />
                          )}
                        </div>
                        <Button variant="ghost" size="icon" className="h-6 w-6">
                          <MoreHorizontal className="w-4 h-4" />
                        </Button>
                      </div>
                      
                      {/* Task Title */}
                      <h4 className="font-medium text-sm line-clamp-2">{task.title}</h4>
                      
                      {/* Task Meta */}
                      <div className="flex flex-wrap gap-1">
                        <Badge variant="secondary" className={`text-xs ${priorityColors[task.priority]}`}>
                          {task.priority}
                        </Badge>
                        {task.isMilestone && (
                          <Badge variant="outline" className="text-xs">
                            Milestone
                          </Badge>
                        )}
                      </div>
                      
                      {/* Task Footer */}
                      <div className="flex items-center justify-between text-xs text-muted-foreground">
                        <div className="flex items-center gap-2">
                          {task.assigneeId ? (
                            <div className="flex items-center gap-1">
                              <div className="w-5 h-5 rounded-full bg-primary/20 flex items-center justify-center text-xs font-medium">
                                {getAssigneeName(task.assigneeId)?.charAt(0)}
                              </div>
                              <span className="truncate max-w-[80px]">
                                {getAssigneeName(task.assigneeId)}
                              </span>
                            </div>
                          ) : (
                            <span className="text-destructive flex items-center gap-1">
                              <AlertCircle className="w-3 h-3" />
                              Unassigned
                            </span>
                          )}
                        </div>
                        
                        {task.dueDate && (
                          <div className="flex items-center gap-1">
                            <Calendar className="w-3 h-3" />
                            {new Date(task.dueDate).toLocaleDateString(undefined, { month: 'short', day: 'numeric' })}
                          </div>
                        )}
                      </div>
                      
                      {/* Project Badge */}
                      {selectedProject === 'all' && (
                        <div className="pt-1">
                          <Badge variant="outline" className="text-xs w-full justify-center">
                            {getProjectName(task.projectId)}
                          </Badge>
                        </div>
                      )}
                    </CardContent>
                  </Card>
                ))}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
