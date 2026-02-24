'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Progress } from '@/components/ui/progress';
import { Skeleton } from '@/components/ui/skeleton';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { useTeamWorkload } from '@/hooks/api';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
// Workload types based on API response
interface WorkloadTask {
  taskId: string;
  title: string;
  projectId: string;
  estimatedHours: number;
  allocationThisWeek: number;
}

interface WorkloadMember {
  personId: string;
  name: string;
  role?: string;
  status: 'overallocated' | 'underutilized' | 'available' | 'optimal';
  allocation: {
    percentage: number;
    assignedHours: number;
    capacityHours: number;
  };
  availability: {
    thisWeek: number;
    nextWeek: number;
  };
  tasks: WorkloadTask[];
  riskLevel: 'low' | 'medium' | 'high';
}

interface WorkloadSummary {
  totalMembers: number;
  overallocated: number;
  optimal: number;
  available: number;
  underutilized: number;
}

interface WorkloadData {
  members: WorkloadMember[];
  summary: WorkloadSummary;
  teamCapacity: number;
  teamAllocation: number;
  utilizationRate: number;
}
import { 
  TrendingUp, 
  TrendingDown, 
  AlertTriangle, 
  CheckCircle2,
  Clock,
  Calendar,
  Briefcase,
  MoreHorizontal,
  RefreshCw
} from 'lucide-react';

export default function TeamPage() {
  // Fetch real workload data from API
  const { 
    data: workloadData, 
    isLoading,
    refetch 
  } = useTeamWorkload(undefined, true);

  const workload = workloadData?.data as WorkloadData | undefined;
  const members = workload?.members || [];
  const summary = workload?.summary;
  
  const overallocated = members.filter(m => m.status === 'overallocated');
  const underutilized = members.filter(m => m.status === 'underutilized');
  const available = members.filter(m => m.status === 'available');
  const healthy = members.filter((m: WorkloadMember) => m.status === 'optimal');

  const getAllocationColor = (percentage: number, status: string) => {
    if (status === 'overallocated' || percentage > 100) return 'text-destructive';
    if (percentage >= 90) return 'text-yellow-600';
    if (percentage >= 70) return 'text-green-600';
    return 'text-blue-500';
  };

  const getAllocationBg = (percentage: number, status: string) => {
    if (status === 'overallocated' || percentage > 100) return 'bg-destructive';
    if (percentage >= 90) return 'bg-yellow-500';
    if (percentage >= 70) return 'bg-green-500';
    return 'bg-blue-500';
  };

  return (
    <div className="p-8 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-start">
        <div>
          <h1 className="text-3xl font-bold">Team & Workload</h1>
          <p className="text-muted-foreground">Monitor team capacity and resource allocation</p>
        </div>
        <div className="flex items-center gap-2">
          <Button 
            variant="outline" 
            size="icon"
            onClick={() => refetch()}
            disabled={isLoading}
          >
            <RefreshCw className={`w-4 h-4 ${isLoading ? 'animate-spin' : ''}`} />
          </Button>
          <Button>
            + Add Member
          </Button>
        </div>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Total Members</CardDescription>
            <CardTitle className="text-2xl">
              {isLoading ? <Skeleton className="h-8 w-16" /> : members.length}
            </CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription className="text-destructive">Overallocated</CardDescription>
            <CardTitle className="text-2xl text-destructive">
              {isLoading ? <Skeleton className="h-8 w-16" /> : summary?.overallocated || overallocated.length}
            </CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription className="text-green-600">Healthy Load</CardDescription>
            <CardTitle className="text-2xl text-green-600">
              {isLoading ? <Skeleton className="h-8 w-16" /> : summary?.optimal || healthy.length}
            </CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription className="text-blue-500">Available</CardDescription>
            <CardTitle className="text-2xl text-blue-500">
              {isLoading ? <Skeleton className="h-8 w-16" /> : (summary?.available || available.length)}
            </CardTitle>
          </CardHeader>
        </Card>
      </div>

      {/* Team Utilization Overview */}
      {!isLoading && workload && (
        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-lg">Team Utilization</CardTitle>
            <CardDescription>
              {workload.teamAllocation}h allocated of {workload.teamCapacity}h capacity
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Overall Utilization</span>
                <span className={`font-medium ${
                  workload.utilizationRate > 1 ? 'text-destructive' : 'text-green-600'
                }`}>
                  {Math.round(workload.utilizationRate * 100)}%
                </span>
              </div>
              <Progress 
                value={Math.min(workload.utilizationRate * 100, 100)} 
                className="h-3"
              />
              {workload.utilizationRate > 1 && (
                <p className="text-xs text-destructive flex items-center">
                  <AlertTriangle className="w-3 h-3 mr-1" />
                  Team is overallocated by {Math.round((workload.utilizationRate - 1) * 100)}%
                </p>
              )}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Team Members */}
      <Tabs defaultValue="all" className="space-y-4">
        <TabsList>
          <TabsTrigger value="all">All Members</TabsTrigger>
          <TabsTrigger value="overallocated" className="text-destructive">
            Overallocated
            {overallocated.length > 0 && (
              <Badge variant="destructive" className="ml-2 h-5 px-1.5">
                {overallocated.length}
              </Badge>
            )}
          </TabsTrigger>
          <TabsTrigger value="available">Available</TabsTrigger>
        </TabsList>

        <TabsContent value="all" className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2">
            {isLoading ? (
              // Loading skeletons
              Array.from({ length: 4 }).map((_, i) => (
                <Card key={i}>
                  <CardHeader>
                    <div className="flex items-center gap-4">
                      <Skeleton className="w-12 h-12 rounded-full" />
                      <div className="flex-1">
                        <Skeleton className="h-5 w-32" />
                        <Skeleton className="h-4 w-20 mt-2" />
                      </div>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <Skeleton className="h-2 w-full" />
                  </CardContent>
                </Card>
              ))
            ) : members.length === 0 ? (
              <Card className="border-dashed col-span-2">
                <CardContent className="p-8 text-center">
                  <Briefcase className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                  <h3 className="text-lg font-medium">No team data</h3>
                  <p className="text-muted-foreground">Workload data will appear here.</p>
                </CardContent>
              </Card>
            ) : (
              members.map((member: WorkloadMember) => {
                const isOverallocated = member.status === 'overallocated';
                const isUnderutilized = member.status === 'underutilized';
                const allocationPercentage = member.allocation.percentage;
                
                return (
                  <Card key={member.personId} className={isOverallocated ? 'border-l-4 border-l-destructive' : ''}>
                    <CardHeader>
                      <div className="flex items-start justify-between">
                        <div className="flex items-center gap-4">
                          <Avatar className="w-12 h-12">
                            <AvatarImage src={`https://api.dicebear.com/7.x/avataaars/svg?seed=${member.name}`} />
                            <AvatarFallback>{member.name.charAt(0)}</AvatarFallback>
                          </Avatar>
                          <div>
                            <CardTitle className="text-lg">{member.name}</CardTitle>
                            <CardDescription className="capitalize">{member.role || 'Team Member'}</CardDescription>
                            <div className="flex items-center gap-2 mt-1">
                              <Badge variant="outline" className="text-xs">
                                {member.tasks?.length || 0} tasks
                              </Badge>
                            </div>
                          </div>
                        </div>
                        <div className="text-right">
                          <div className={`text-2xl font-bold ${getAllocationColor(allocationPercentage, member.status)}`}>
                            {allocationPercentage}%
                          </div>
                          <div className="text-xs text-muted-foreground">Allocated</div>
                        </div>
                      </div>
                    </CardHeader>
                    <CardContent className="space-y-4">
                      {/* Capacity Bar */}
                      <div className="space-y-2">
                        <div className="flex justify-between text-sm">
                          <span className="text-muted-foreground">Weekly Capacity</span>
                          <span className="font-medium">
                            {member.allocation.assignedHours}h of {member.allocation.capacityHours}h
                          </span>
                        </div>
                        <Progress 
                          value={Math.min(allocationPercentage, 100)} 
                          className="h-2"
                        />
                        {isOverallocated && (
                          <p className="text-xs text-destructive flex items-center">
                            <AlertTriangle className="w-3 h-3 mr-1" />
                            {allocationPercentage - 100}% over capacity
                          </p>
                        )}
                        {member.riskLevel !== 'low' && (
                          <p className={`text-xs flex items-center ${
                            member.riskLevel === 'high' ? 'text-destructive' : 'text-yellow-600'
                          }`}>
                            <AlertTriangle className="w-3 h-3 mr-1" />
                            {member.riskLevel} risk
                          </p>
                        )}
                      </div>

                      {/* Stats Grid */}
                      <div className="grid grid-cols-3 gap-4 pt-2">
                        <div className="text-center">
                          <p className="text-lg font-semibold">{member.tasks?.length || 0}</p>
                          <p className="text-xs text-muted-foreground">Tasks</p>
                        </div>
                        <div className="text-center">
                          <p className="text-lg font-semibold">{member.allocation.assignedHours}h</p>
                          <p className="text-xs text-muted-foreground">Est. Hours</p>
                        </div>
                        <div className="text-center">
                          <p className="text-lg font-semibold">{member.availability.thisWeek}h</p>
                          <p className="text-xs text-muted-foreground">Available</p>
                        </div>
                      </div>

                      {/* Tasks List (first 2) */}
                      {member.tasks && member.tasks.length > 0 && (
                        <div className="pt-2">
                          <p className="text-xs text-muted-foreground mb-2">Current Tasks</p>
                          <div className="space-y-1">
                            {member.tasks.slice(0, 2).map((task: WorkloadTask) => (
                              <div key={task.taskId} className="text-sm truncate">
                                • {task.title}
                              </div>
                            ))}
                            {member.tasks.length > 2 && (
                              <p className="text-xs text-muted-foreground">
                                +{member.tasks.length - 2} more
                              </p>
                            )}
                          </div>
                        </div>
                      )}

                      {/* Skills */}
                      <div className="pt-2">
                        <p className="text-xs text-muted-foreground mb-2">Skills</p>
                        <div className="flex flex-wrap gap-1">
                          <Badge variant="outline" className="text-xs text-muted-foreground">
                            Skills coming from API...
                          </Badge>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                );
              })
            )}
          </div>
        </TabsContent>

        <TabsContent value="overallocated" className="space-y-4">
          {isLoading ? (
            <div className="grid gap-4 md:grid-cols-2">
              {Array.from({ length: 2 }).map((_, i) => (
                <Card key={i}>
                  <CardHeader>
                    <Skeleton className="h-6 w-32" />
                    <Skeleton className="h-4 w-20 mt-2" />
                  </CardHeader>
                </Card>
              ))}
            </div>
          ) : overallocated.length === 0 ? (
            <Card className="border-dashed">
              <CardContent className="p-8 text-center">
                <CheckCircle2 className="w-12 h-12 text-green-500 mx-auto mb-4" />
                <h3 className="text-lg font-medium">No overallocation!</h3>
                <p className="text-muted-foreground">All team members are within capacity.</p>
              </CardContent>
            </Card>
          ) : (
            <div className="grid gap-4 md:grid-cols-2">
              {overallocated.map((member: WorkloadMember) => {
                return (
                  <Card key={member.personId} className="border-destructive">
                    <CardHeader>
                      <div className="flex items-center gap-4">
                        <Avatar className="w-12 h-12">
                          <AvatarImage src={`https://api.dicebear.com/7.x/avataaars/svg?seed=${member.name}`} />
                          <AvatarFallback>{member.name.charAt(0)}</AvatarFallback>
                        </Avatar>
                        <div className="flex-1">
                          <CardTitle className="text-lg">{member.name}</CardTitle>
                          <div className="flex items-center gap-2 text-destructive">
                            <TrendingUp className="w-4 h-4" />
                            <span className="font-medium">{member.allocation.percentage}% allocated</span>
                          </div>
                        </div>
                      </div>
                    </CardHeader>
                    <CardContent>
                      <p className="text-sm text-muted-foreground mb-4">
                        {member.name} has {member.allocation.assignedHours} hours assigned this week,
                        which is {member.allocation.percentage - 100}% over their {member.allocation.capacityHours}-hour capacity.
                      </p>
                      <div className="space-y-2 mb-4">
                        {member.tasks?.slice(0, 3).map((task: WorkloadTask) => (
                          <div key={task.taskId} className="text-sm flex justify-between">
                            <span className="truncate">{task.title}</span>
                            <span className="text-muted-foreground">{task.allocationThisWeek}h</span>
                          </div>
                        ))}
                      </div>
                      <div className="flex gap-2">
                        <Button variant="destructive" size="sm">
                          Reassign Tasks
                        </Button>
                        <Button variant="outline" size="sm">
                          View Workload
                        </Button>
                      </div>
                    </CardContent>
                  </Card>
                );
              })}
            </div>
          )}
        </TabsContent>

        <TabsContent value="available" className="space-y-4">
          {isLoading ? (
            <div className="grid gap-4 md:grid-cols-2">
              {Array.from({ length: 2 }).map((_, i) => (
                <Card key={i}>
                  <CardHeader>
                    <Skeleton className="h-6 w-32" />
                    <Skeleton className="h-4 w-20 mt-2" />
                  </CardHeader>
                </Card>
              ))}
            </div>
          ) : available.length === 0 && underutilized.length === 0 ? (
            <Card className="border-dashed">
              <CardContent className="p-8 text-center">
                <Briefcase className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                <h3 className="text-lg font-medium">No available capacity</h3>
                <p className="text-muted-foreground">All team members are at full capacity.</p>
              </CardContent>
            </Card>
          ) : (
            <div className="grid gap-4 md:grid-cols-2">
              {[...available, ...underutilized].map((member: WorkloadMember) => {
                return (
                  <Card key={member.personId}>
                    <CardHeader>
                      <div className="flex items-center gap-4">
                        <Avatar className="w-12 h-12">
                          <AvatarImage src={`https://api.dicebear.com/7.x/avataaars/svg?seed=${member.name}`} />
                          <AvatarFallback>{member.name.charAt(0)}</AvatarFallback>
                        </Avatar>
                        <div className="flex-1">
                          <CardTitle className="text-lg">{member.name}</CardTitle>
                          <div className="flex items-center gap-2 text-blue-500">
                            <TrendingDown className="w-4 h-4" />
                            <span className="font-medium">{member.allocation.percentage}% allocated</span>
                          </div>
                        </div>
                      </div>
                    </CardHeader>
                    <CardContent>
                      <p className="text-sm text-muted-foreground mb-4">
                        {member.name} has {member.availability.thisWeek} hours available this week.
                      </p>
                      <div className="flex gap-2">
                        <Button size="sm">
                          Assign Task
                        </Button>
                        <Button variant="outline" size="sm">
                          View Profile
                        </Button>
                      </div>
                    </CardContent>
                  </Card>
                );
              })}
            </div>
          )}
        </TabsContent>
      </Tabs>
    </div>
  );
}
