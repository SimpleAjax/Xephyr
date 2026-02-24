'use client';

import { useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Progress } from '@/components/ui/progress';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { useProjectsList, useTasks, useCreateProject, useCreateTask } from '@/hooks/api';
import { CreateProjectDialog } from '@/components/projects/create-project-dialog';
import { 
  Search, 
  Filter, 
  AlertTriangle, 
  CheckCircle2, 
  Clock,
  Calendar,
  Users,
  ArrowRight,
  MoreHorizontal,
  Plus,
  Sparkles,
  Loader2
} from 'lucide-react';
import Link from 'next/link';
import { Project, Task } from '@/types';

export default function ProjectsPage() {
  const [searchQuery, setSearchQuery] = useState('');
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  
  const { data: projectsData, isLoading: projectsLoading } = useProjectsList();
  const { data: tasksData } = useTasks();
  const createProjectMutation = useCreateProject();
  const createTaskMutation = useCreateTask();
  
  const projects = projectsData?.data?.projects || [];
  const tasks = tasksData?.data?.tasks || [];
  
  const activeProjects = projects.filter((p: Project) => p.status === 'active');
  const atRiskProjects = projects.filter((p: Project) => p.healthScore < 50);
  
  const filteredProjects = projects.filter((project: Project) => {
    const matchesSearch = project.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
                         project.description.toLowerCase().includes(searchQuery.toLowerCase());
    return matchesSearch;
  });

  const getStatusIcon = (healthScore: number) => {
    if (healthScore >= 80) return <CheckCircle2 className="w-5 h-5 text-green-500" />;
    if (healthScore >= 60) return <Clock className="w-5 h-5 text-yellow-500" />;
    return <AlertTriangle className="w-5 h-5 text-destructive" />;
  };

  const getAssignedCount = (projectId: string) => {
    return tasks.filter((t: Task) => t.projectId === projectId && t.assigneeId).length;
  };

  const handleProjectCreate = async (newProject: Project, newTasks: Task[]) => {
    try {
      // Create project first
      const projectResponse = await createProjectMutation.mutateAsync({
        name: newProject.name,
        description: newProject.description,
        priority: newProject.priority,
        startDate: newProject.startDate.toISOString(),
        targetEndDate: newProject.targetEndDate.toISOString(),
      });
      
      const createdProjectId = projectResponse.data?.id;
      
      // Then create tasks if project was created successfully
      if (createdProjectId && newTasks.length > 0) {
        for (const task of newTasks) {
          await createTaskMutation.mutateAsync({
            title: task.title,
            description: task.description,
            projectId: createdProjectId,
            priority: task.priority,
            estimatedHours: task.estimatedHours,
            hierarchyLevel: task.hierarchyLevel || 1,
            parentTaskId: task.parentTaskId,
            requiredSkills: task.requiredSkills,
          });
        }
      }
      
      setCreateDialogOpen(false);
    } catch (error) {
      console.error('Failed to create project:', error);
    }
  };

  if (projectsLoading) {
    return (
      <div className="p-8 h-full flex items-center justify-center">
        <div className="flex items-center gap-2 text-muted-foreground">
          <Loader2 className="w-5 h-5 animate-spin" />
          <span>Loading projects...</span>
        </div>
      </div>
    );
  }

  return (
    <div className="p-8 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-start">
        <div>
          <h1 className="text-3xl font-bold">Projects</h1>
          <p className="text-muted-foreground">Manage and monitor all your projects</p>
        </div>
        <Button onClick={() => setCreateDialogOpen(true)}>
          <Sparkles className="w-4 h-4 mr-2" />
          New Project
        </Button>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Total Projects</CardDescription>
            <CardTitle className="text-2xl">{projects.length}</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Active</CardDescription>
            <CardTitle className="text-2xl">{activeProjects.length}</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>At Risk</CardDescription>
            <CardTitle className="text-2xl text-destructive">{atRiskProjects.length}</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Avg Health</CardDescription>
            <CardTitle className="text-2xl">
              {projects.length > 0 
                ? Math.round(projects.reduce((sum: number, p: Project) => sum + p.healthScore, 0) / projects.length)
                : 0}%
            </CardTitle>
          </CardHeader>
        </Card>
      </div>

      {/* Filters */}
      <div className="flex gap-4">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input 
            placeholder="Search projects..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-10"
          />
        </div>
        <Button variant="outline" size="icon">
          <Filter className="h-4 w-4" />
        </Button>
      </div>

      {/* Projects Tabs */}
      <Tabs defaultValue="all" className="space-y-4">
        <TabsList>
          <TabsTrigger value="all">All Projects</TabsTrigger>
          <TabsTrigger value="active">Active</TabsTrigger>
          <TabsTrigger value="atrisk">
            At Risk
            {atRiskProjects.length > 0 && (
              <Badge variant="destructive" className="ml-2 h-5 px-1.5">
                {atRiskProjects.length}
              </Badge>
            )}
          </TabsTrigger>
          <TabsTrigger value="completed">Completed</TabsTrigger>
        </TabsList>

        <TabsContent value="all" className="space-y-4">
          <div className="grid gap-4">
            {filteredProjects.map((project: Project) => (
              <Card key={project.id}>
                <CardHeader>
                  <div className="flex justify-between items-start">
                    <div className="flex items-start gap-4">
                      {getStatusIcon(project.healthScore)}
                      <div>
                        <CardTitle className="text-lg">{project.name}</CardTitle>
                        <CardDescription className="mt-1 max-w-2xl">
                          {project.description}
                        </CardDescription>
                        <div className="flex items-center gap-2 mt-2">
                          <Badge variant={project.status === 'active' ? 'default' : 'secondary'}>
                            {project.status}
                          </Badge>
                          <Badge variant="outline">Priority: {project.priority}</Badge>
                        </div>
                      </div>
                    </div>
                    <div className="text-right">
                      <div className="text-2xl font-bold">{project.healthScore}%</div>
                      <div className="text-xs text-muted-foreground">Health Score</div>
                    </div>
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <div className="space-y-2">
                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">Progress</span>
                        <span className="font-medium">{project.progress}%</span>
                      </div>
                      <Progress value={project.progress} className="h-2" />
                    </div>
                    
                    <div className="flex items-center justify-between pt-2">
                      <div className="flex items-center gap-4 text-sm text-muted-foreground">
                        <span className="flex items-center">
                          <Calendar className="w-4 h-4 mr-1" />
                          Due {new Date(project.targetEndDate).toLocaleDateString()}
                        </span>
                        <span className="flex items-center">
                          <Users className="w-4 h-4 mr-1" />
                          {getAssignedCount(project.id)} assigned
                        </span>
                      </div>
                      <div className="flex gap-2">
                        <Button variant="outline" size="sm" asChild>
                          <Link href={`/projects/${project.id}`}>
                            View Details
                          </Link>
                        </Button>
                        <Button variant="outline" size="sm" asChild>
                          <Link href={`/kanban?project=${project.id}`}>
                            Board
                          </Link>
                        </Button>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="active" className="space-y-4">
          <div className="grid gap-4">
            {activeProjects.map((project: Project) => (
              <Card key={project.id}>
                <CardHeader>
                  <div className="flex justify-between items-start">
                    <div>
                      <CardTitle className="text-lg">{project.name}</CardTitle>
                      <CardDescription>{project.description}</CardDescription>
                    </div>
                    <Badge variant={project.healthScore >= 80 ? 'default' : 'destructive'}>
                      {project.healthScore}%
                    </Badge>
                  </div>
                </CardHeader>
                <CardContent>
                  <Progress value={project.progress} className="h-2" />
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="atrisk" className="space-y-4">
          {atRiskProjects.length === 0 ? (
            <Card className="border-dashed">
              <CardContent className="p-8 text-center">
                <CheckCircle2 className="w-12 h-12 text-green-500 mx-auto mb-4" />
                <h3 className="text-lg font-medium">No projects at risk!</h3>
                <p className="text-muted-foreground">All your projects are healthy.</p>
              </CardContent>
            </Card>
          ) : (
            <div className="grid gap-4">
              {atRiskProjects.map((project: Project) => (
                <Card key={project.id} className="border-destructive">
                  <CardHeader>
                    <div className="flex justify-between items-start">
                      <div className="flex items-start gap-3">
                        <AlertTriangle className="w-5 h-5 text-destructive mt-1" />
                        <div>
                          <CardTitle className="text-lg">{project.name}</CardTitle>
                          <CardDescription>{project.description}</CardDescription>
                          <div className="flex items-center gap-2 mt-2">
                            <Badge variant="destructive">Health: {project.healthScore}%</Badge>
                            <Badge variant="outline">Progress: {project.progress}%</Badge>
                          </div>
                        </div>
                      </div>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className="flex items-center gap-4">
                      <Button variant="destructive" size="sm">
                        Take Action
                      </Button>
                      <Button variant="outline" size="sm" asChild>
                        <Link href={`/projects/${project.id}`}>
                          View Details
                        </Link>
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </TabsContent>

        <TabsContent value="completed" className="space-y-4">
          <Card className="border-dashed">
            <CardContent className="p-8 text-center">
              <Clock className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
              <h3 className="text-lg font-medium">No completed projects yet</h3>
              <p className="text-muted-foreground">Projects will appear here when marked as completed.</p>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      <CreateProjectDialog 
        open={createDialogOpen} 
        onOpenChange={setCreateDialogOpen}
        onProjectCreate={handleProjectCreate}
        isCreating={createProjectMutation.isPending || createTaskMutation.isPending}
      />
    </div>
  );
}
