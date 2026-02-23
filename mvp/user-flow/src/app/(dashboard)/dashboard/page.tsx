'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Progress } from '@/components/ui/progress';
import { mockData } from '@/data/mock-data';
import { useAppStore } from '@/store/app-store';
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
  AlertCircle
} from 'lucide-react';
import Link from 'next/link';

export default function DashboardPage() {
  const projects = mockData.getActiveProjects();
  const nudges = useAppStore((state) => state.nudges);
  const users = mockData.getUsers();
  const workloadData = mockData.getWorkloadData();
  
  const unreadNudges = nudges.filter(n => n.status === 'unread');
  const highSeverityNudges = unreadNudges.filter(n => n.severity === 'high');
  const atRiskProjects = mockData.getAtRiskProjects();
  const overallocationCount = workloadData.filter(w => w.allocationPercentage > 100).length;
  const portfolioHealth = mockData.getPortfolioHealthScore();

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

  return (
    <div className="p-8 space-y-8">
      {/* Header */}
      <div className="flex justify-between items-start">
        <div>
          <h1 className="text-3xl font-bold">Dashboard</h1>
          <p className="text-muted-foreground">Welcome back, Sarah! Here's what's happening with your portfolio.</p>
        </div>
        <div className="flex items-center space-x-4">
          <div className="text-right">
            <p className="text-sm text-muted-foreground">Portfolio Health</p>
            <p className={`text-2xl font-bold ${getHealthColor(portfolioHealth)}`}>
              {portfolioHealth}%
            </p>
          </div>
          <div className={`w-12 h-12 rounded-full ${getHealthBg(portfolioHealth)} bg-opacity-20 flex items-center justify-center`}>
            {portfolioHealth >= 80 ? (
              <CheckCircle2 className={`w-6 h-6 ${getHealthColor(portfolioHealth)}`} />
            ) : portfolioHealth >= 60 ? (
              <Clock className={`w-6 h-6 ${getHealthColor(portfolioHealth)}`} />
            ) : (
              <AlertTriangle className={`w-6 h-6 ${getHealthColor(portfolioHealth)}`} />
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
            <CardTitle className="text-3xl">{projects.length}</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center text-sm">
              {atRiskProjects.length > 0 ? (
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
            <CardTitle className="text-3xl">{users.length}</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center text-sm">
              {overallocationCount > 0 ? (
                <span className="text-destructive flex items-center">
                  <TrendingUp className="w-4 h-4 mr-1" />
                  {overallocationCount} overallocated
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
            <CardTitle className="text-3xl">{unreadNudges.length}</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center text-sm">
              {highSeverityNudges.length > 0 ? (
                <span className="text-destructive flex items-center">
                  <AlertCircle className="w-4 h-4 mr-1" />
                  {highSeverityNudges.length} need attention
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
            <CardTitle className="text-3xl">{mockData.getTasksByStatus('in_progress').length}</CardTitle>
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
            {projects.map((project) => (
              <Card key={project.id} className={project.healthScore < 50 ? 'border-l-4 border-l-destructive' : ''}>
                <CardHeader className="pb-3">
                  <div className="flex justify-between items-start">
                    <div>
                      <CardTitle className="text-lg">{project.name}</CardTitle>
                      <CardDescription className="mt-1 line-clamp-1">
                        {project.description}
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
                      <span>Due {project.targetEndDate.toLocaleDateString()}</span>
                      <span>Priority: {project.priority}/100</span>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
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
            {unreadNudges.slice(0, 5).map((nudge) => (
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
                          {nudge.createdAt.toLocaleDateString()}
                        </span>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
            
            {unreadNudges.length === 0 && (
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
            {atRiskProjects.map((project) => (
              <Card key={project.id} className="border-destructive/50">
                <CardHeader>
                  <div className="flex justify-between items-start">
                    <CardTitle className="text-lg">{project.name}</CardTitle>
                    <Badge variant="destructive">{project.healthScore}% Health</Badge>
                  </div>
                  <CardDescription>{project.description}</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="flex items-center gap-4 text-sm">
                    <span className="flex items-center text-destructive">
                      <TrendingDown className="w-4 h-4 mr-1" />
                      Needs attention
                    </span>
                    <Button size="sm" variant="outline" asChild>
                      <Link href={`/projects?project=${project.id}`}>
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
