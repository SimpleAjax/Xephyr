'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Progress } from '@/components/ui/progress';
import { Skeleton } from '@/components/ui/skeleton';
import { usePortfolioHealth, useTeamWorkload } from '@/hooks/api';
import { useNudges } from '@/hooks/api/useNudges';
import { 
  AlertTriangle, 
  CheckCircle2, 
  Clock, 
  Users, 
  FolderKanban,
  TrendingUp,
  TrendingDown,
  ArrowRight,
  Zap,
  AlertCircle,
  RefreshCw
} from 'lucide-react';
import Link from 'next/link';

export default function DashboardPage() {
  // Fetch real data from APIs
  const { 
    data: portfolioData, 
    isLoading: portfolioLoading,
    error: portfolioError,
    refetch: refetchPortfolio 
  } = usePortfolioHealth();
  
  const { 
    data: nudgesData, 
    isLoading: nudgesLoading,
    refetch: refetchNudges 
  } = useNudges({ limit: 20 });
  
  const { 
    data: workloadData, 
    isLoading: workloadLoading,
    refetch: refetchWorkload 
  } = useTeamWorkload();

  // Get data from responses
  const portfolioHealth = portfolioData?.data;
  const nudgesResponse = nudgesData?.data;
  const workloadResponse = workloadData?.data;

  // Derived data
  const unreadNudges = nudgesResponse?.nudges?.filter((n: { status: string }) => n.status === 'unread') || [];
  const highSeverityNudges = unreadNudges.filter((n: { severity: string }) => n.severity === 'high');
  const atRiskProjects = portfolioHealth?.projects?.filter((p: { status: string }) => p.status === 'at_risk' || p.status === 'critical') || [];
  const overallocatedCount = workloadResponse?.summary?.overallocated || 0;
  
  // Users and tasks - will come from backend APIs when ready
  const users: { id: string; name: string }[] = [];
  const tasks: { id: string; status: string }[] = [];

  const getHealthColor = (score: number) => {
    if (score >= 80) return 'text-green-500';
    if (score >= 60) return 'text-yellow-500';
    return 'text-red-500';
  };

  const getHealthBg = (score: number) => {
    if (score >= 80) return 'bg-green-500';
    if (score >= 60) return 'bg-yellow-500';
    return 'bg-red-500';
  };

  const getHealthStatus = (score: number) => {
    if (score >= 80) return 'healthy';
    if (score >= 60) return 'caution';
    if (score >= 40) return 'at_risk';
    return 'critical';
  };

  const isLoading = portfolioLoading || nudgesLoading || workloadLoading;

  return (
    <div className="p-8 space-y-8">
      {/* Header */}
      <div className="flex justify-between items-start">
        <div>
          <h1 className="text-3xl font-bold">Dashboard</h1>
          <p className="text-muted-foreground">Welcome back! Here&apos;s what&apos;s happening with your portfolio.</p>
        </div>
        <div className="flex items-center space-x-4">
          <Button 
            variant="ghost" 
            size="icon" 
            onClick={() => {
              refetchPortfolio();
              refetchNudges();
              refetchWorkload();
            }}
            disabled={isLoading}
          >
            <RefreshCw className={`w-4 h-4 ${isLoading ? 'animate-spin' : ''}`} />
          </Button>
          <div className="text-right">
            <p className="text-sm text-muted-foreground">Portfolio Health</p>
            {portfolioLoading ? (
              <Skeleton className="h-8 w-16 ml-auto" />
            ) : portfolioError ? (
              <p className="text-sm text-destructive">Error</p>
            ) : (
              <p className={`text-2xl font-bold ${getHealthColor(portfolioHealth?.portfolioHealthScore || 0)}`}>
                {portfolioHealth?.portfolioHealthScore || 0}%
              </p>
            )}
          </div>
          <div className={`w-12 h-12 rounded-full ${getHealthBg(portfolioHealth?.portfolioHealthScore || 0)} bg-opacity-20 flex items-center justify-center`}>
            {(portfolioHealth?.portfolioHealthScore || 0) >= 80 ? (
              <CheckCircle2 className={`w-6 h-6 ${getHealthColor(portfolioHealth?.portfolioHealthScore || 0)}`} />
            ) : (portfolioHealth?.portfolioHealthScore || 0) >= 60 ? (
              <Clock className={`w-6 h-6 ${getHealthColor(portfolioHealth?.portfolioHealthScore || 0)}`} />
            ) : (
              <AlertTriangle className={`w-6 h-6 ${getHealthColor(portfolioHealth?.portfolioHealthScore || 0)}`} />
            )}
          </div>
        </div>
      </div>

      {/* Stats Row */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardDescription className="flex items-center">
              <FolderKanban className="w-4 h-4 mr-2" />
              Active Projects
            </CardDescription>
            <CardTitle className="text-3xl">
              {portfolioLoading ? <Skeleton className="h-9 w-12" /> : portfolioHealth?.summary?.totalProjects || 0}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center text-sm">
              {portfolioLoading ? (
                <Skeleton className="h-4 w-24" />
              ) : atRiskProjects.length > 0 ? (
                <span className="text-destructive flex items-center">
                  <AlertTriangle className="w-4 h-4 mr-1" />
                  {atRiskProjects.length} at risk
                </span>
              ) : (
                <span className="text-green-500 flex items-center">
                  <CheckCircle2 className="w-4 h-4 mr-1" />
                  All healthy
                </span>
              )}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardDescription className="flex items-center">
              <Users className="w-4 h-4 mr-2" />
              Team Members
            </CardDescription>
            <CardTitle className="text-3xl">
              {workloadLoading ? <Skeleton className="h-9 w-12" /> : workloadResponse?.members?.length || users.length}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center text-sm">
              {workloadLoading ? (
                <Skeleton className="h-4 w-24" />
              ) : overallocatedCount > 0 ? (
                <span className="text-destructive flex items-center">
                  <TrendingUp className="w-4 h-4 mr-1" />
                  {overallocatedCount} overallocated
                </span>
              ) : (
                <span className="text-green-500 flex items-center">
                  <CheckCircle2 className="w-4 h-4 mr-1" />
                  Well balanced
                </span>
              )}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardDescription className="flex items-center">
              <Zap className="w-4 h-4 mr-2" />
              AI Nudges
            </CardDescription>
            <CardTitle className="text-3xl">
              {nudgesLoading ? <Skeleton className="h-9 w-12" /> : nudgesResponse?.summary?.unread || 0}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center text-sm">
              {nudgesLoading ? (
                <Skeleton className="h-4 w-24" />
              ) : highSeverityNudges.length > 0 ? (
                <span className="text-destructive flex items-center">
                  <AlertCircle className="w-4 h-4 mr-1" />
                  {nudgesResponse?.summary?.bySeverity?.high || 0} need attention
                </span>
              ) : (
                <span className="text-green-500 flex items-center">
                  <CheckCircle2 className="w-4 h-4 mr-1" />
                  All caught up
                </span>
              )}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardDescription className="flex items-center">
              <Clock className="w-4 h-4 mr-2" />
              This Week
            </CardDescription>
            <CardTitle className="text-3xl">
              {tasks.filter(t => t.status === 'in_progress').length}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center text-sm text-muted-foreground">
              tasks in progress
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        {/* Projects Section */}
        <div className="lg:col-span-2 space-y-4">
          <div className="flex items-center justify-between">
            <h2 className="text-xl font-semibold">Active Projects</h2>
            <Button variant="ghost" size="sm" asChild>
              <Link href="/projects">
                View all
                <ArrowRight className="w-4 h-4 ml-1" />
              </Link>
            </Button>
          </div>
          
          <div className="grid gap-4">
            {portfolioLoading ? (
              // Loading skeletons
              Array.from({ length: 3 }).map((_, i) => (
                <Card key={i}>
                  <CardHeader className="pb-3">
                    <Skeleton className="h-6 w-48" />
                    <Skeleton className="h-4 w-full mt-2" />
                  </CardHeader>
                  <CardContent>
                    <Skeleton className="h-2 w-full" />
                  </CardContent>
                </Card>
              ))
            ) : portfolioHealth?.projects && portfolioHealth.projects.length > 0 ? (
              portfolioHealth.projects.map((project: { projectId: string; name: string; healthScore: number; progress: number; trend: string; priority: number }) => (
                <Card 
                  key={project.projectId} 
                  className={project.healthScore < 50 ? 'border-l-4 border-l-destructive' : ''}
                >
                  <CardHeader className="pb-3">
                    <div className="flex justify-between items-start">
                      <div>
                        <CardTitle className="text-lg">{project.name}</CardTitle>
                        <CardDescription className="mt-1 line-clamp-1">
                          {project.projectId}
                        </CardDescription>
                      </div>
                      <Badge 
                        variant={project.healthScore >= 80 ? 'default' : project.healthScore >= 60 ? 'secondary' : 'destructive'}
                        className="ml-2"
                      >
                        {project.healthScore}%
                      </Badge>
                    </div>
                  </CardHeader>
                  <CardContent className="pt-0">
                    <div className="space-y-2">
                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">Progress</span>
                        <span className="font-medium">{project.progress}%</span>
                      </div>
                      <Progress value={project.progress} className="h-2" />
                      <div className="flex justify-between text-xs text-muted-foreground pt-1">
                        <span className="capitalize">{project.trend} trend</span>
                        <span>Priority: {project.priority}/100</span>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              ))
            ) : (
              <Card className="border-dashed">
                <CardContent className="p-6 text-center">
                  <FolderKanban className="w-8 h-8 text-muted-foreground mx-auto mb-2" />
                  <p className="text-sm text-muted-foreground">No active projects</p>
                </CardContent>
              </Card>
            )}
          </div>
        </div>

        {/* Nudges Section */}
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <h2 className="text-xl font-semibold">AI Nudges</h2>
            <Button variant="ghost" size="sm" asChild>
              <Link href="/nudges">
                View all
                <ArrowRight className="w-4 h-4 ml-1" />
              </Link>
            </Button>
          </div>
          
          <div className="space-y-3">
            {nudgesLoading ? (
              // Loading skeletons
              Array.from({ length: 3 }).map((_, i) => (
                <Card key={i}>
                  <CardContent className="p-4">
                    <Skeleton className="h-4 w-full" />
                    <Skeleton className="h-3 w-3/4 mt-2" />
                  </CardContent>
                </Card>
              ))
            ) : unreadNudges.length > 0 ? (
              unreadNudges.slice(0, 5).map((nudge: { id: string; severity: 'high' | 'medium' | 'low'; type: string; title: string; description: string; createdAt: string }) => (
                <Card 
                  key={nudge.id} 
                  className={`border-l-4 ${
                    nudge.severity === 'high' ? 'border-l-destructive' : 
                    nudge.severity === 'medium' ? 'border-l-yellow-500' : 
                    'border-l-blue-500'
                  }`}
                >
                  <CardContent className="p-4">
                    <div className="flex items-start gap-3">
                      <div className={`mt-0.5 ${
                        nudge.severity === 'high' ? 'text-destructive' : 
                        nudge.severity === 'medium' ? 'text-yellow-500' : 
                        'text-blue-500'
                      }`}>
                        {nudge.type === 'overload' && <TrendingUp className="w-4 h-4" />}
                        {nudge.type === 'delay_risk' && <Clock className="w-4 h-4" />}
                        {nudge.type === 'skill_gap' && <AlertCircle className="w-4 h-4" />}
                        {nudge.type === 'unassigned' && <Users className="w-4 h-4" />}
                        {nudge.type === 'blocked' && <AlertTriangle className="w-4 h-4" />}
                        {nudge.type === 'conflict' && <Zap className="w-4 h-4" />}
                        {nudge.type === 'dependency_block' && <FolderKanban className="w-4 h-4" />}
                      </div>
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium line-clamp-2">{nudge.title}</p>
                        <p className="text-xs text-muted-foreground mt-1 line-clamp-2">
                          {nudge.description}
                        </p>
                        <div className="flex items-center gap-2 mt-2">
                          <Badge variant={nudge.severity === 'high' ? 'destructive' : 'secondary'} className="text-xs">
                            {nudge.severity}
                          </Badge>
                          <span className="text-xs text-muted-foreground">
                            {new Date(nudge.createdAt).toLocaleDateString()}
                          </span>
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              ))
            ) : (
              <Card className="border-dashed">
                <CardContent className="p-6 text-center">
                  <CheckCircle2 className="w-8 h-8 text-green-500 mx-auto mb-2" />
                  <p className="text-sm text-muted-foreground">All caught up!</p>
                  <p className="text-xs text-muted-foreground">No pending nudges</p>
                </CardContent>
              </Card>
            )}
          </div>
        </div>
      </div>

      {/* At Risk Section */}
      {atRiskProjects.length > 0 && (
        <div className="space-y-4">
          <h2 className="text-xl font-semibold text-destructive flex items-center">
            <AlertTriangle className="w-5 h-5 mr-2" />
            At-Risk Projects
          </h2>
          <div className="grid gap-4 md:grid-cols-2">
            {atRiskProjects.map((project: { projectId: string; name: string; healthScore: number; trend: string }) => (
              <Card key={project.projectId} className="border-destructive/50">
                <CardHeader>
                  <div className="flex justify-between items-start">
                    <CardTitle className="text-lg">{project.name}</CardTitle>
                    <Badge variant="destructive">{project.healthScore}% Health</Badge>
                  </div>
                  <CardDescription>{project.projectId}</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="flex items-center gap-4 text-sm">
                    <span className="flex items-center text-destructive">
                      <TrendingDown className="w-4 h-4 mr-1" />
                      {project.trend} trend
                    </span>
                    <Button size="sm" variant="outline" asChild>
                      <Link href={`/projects?project=${project.projectId}`}>
                        View details
                      </Link>
                    </Button>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
