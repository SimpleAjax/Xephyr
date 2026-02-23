'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Progress } from '@/components/ui/progress';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { mockData } from '@/data/mock-data';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { 
  TrendingUp, 
  TrendingDown, 
  AlertTriangle, 
  CheckCircle2,
  Clock,
  Calendar,
  Briefcase,
  MoreHorizontal
} from 'lucide-react';

export default function TeamPage() {
  const users = mockData.getUsers();
  const workloadData = mockData.getWorkloadData();
  
  const overallocated = workloadData.filter(w => w.allocationPercentage > 100);
  const underutilized = workloadData.filter(w => w.allocationPercentage < 70);
  const healthy = workloadData.filter(w => w.allocationPercentage >= 70 && w.allocationPercentage <= 100);

  const getAllocationColor = (percentage: number) => {
    if (percentage > 100) return 'text-destructive';
    if (percentage >= 90) return 'text-yellow-600';
    if (percentage >= 70) return 'text-green-600';
    return 'text-blue-500';
  };

  const getAllocationBg = (percentage: number) => {
    if (percentage > 100) return 'bg-destructive';
    if (percentage >= 90) return 'bg-yellow-500';
    if (percentage >= 70) return 'bg-green-500';
    return 'bg-blue-500';
  };

  const getUserTasks = (userId: string) => {
    return mockData.getTasksByAssignee(userId);
  };

  return (
    <div className="p-8 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-start">
        <div>
          <h1 className="text-3xl font-bold">Team & Workload</h1>
          <p className="text-muted-foreground">Monitor team capacity and resource allocation</p>
        </div>
        <Button>
          + Add Member
        </Button>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Total Members</CardDescription>
            <CardTitle className="text-2xl">{users.length}</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription className="text-destructive">Overallocated</CardDescription>
            <CardTitle className="text-2xl text-destructive">{overallocated.length}</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription className="text-green-600">Healthy Load</CardDescription>
            <CardTitle className="text-2xl text-green-600">{healthy.length}</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription className="text-blue-500">Available</CardDescription>
            <CardTitle className="text-2xl text-blue-500">{underutilized.length}</CardTitle>
          </CardHeader>
        </Card>
      </div>

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
            {workloadData.map((workload) => {
              const user = mockData.getUserById(workload.personId);
              if (!user) return null;
              
              const tasks = getUserTasks(user.id);
              const isOverallocated = workload.allocationPercentage > 100;
              const isUnderutilized = workload.allocationPercentage < 70;
              
              return (
                <Card key={user.id} className={isOverallocated ? 'border-l-4 border-l-destructive' : ''}>
                  <CardHeader>
                    <div className="flex items-start justify-between">
                      <div className="flex items-center gap-4">
                        <Avatar className="w-12 h-12">
                          <AvatarImage src={user.avatarUrl} />
                          <AvatarFallback>{user.name.charAt(0)}</AvatarFallback>
                        </Avatar>
                        <div>
                          <CardTitle className="text-lg">{user.name}</CardTitle>
                          <CardDescription className="capitalize">{user.role}</CardDescription>
                          <div className="flex items-center gap-2 mt-1">
                            <Badge variant="outline" className="text-xs">
                              ${user.hourlyRate}/hr
                            </Badge>
                            <Badge variant="outline" className="text-xs">
                              {tasks.length} tasks
                            </Badge>
                          </div>
                        </div>
                      </div>
                      <div className="text-right">
                        <div className={`text-2xl font-bold ${getAllocationColor(workload.allocationPercentage)}`}>
                          {workload.allocationPercentage}%
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
                          {Math.min(workload.allocationPercentage, 100)}% of 40h
                        </span>
                      </div>
                      <Progress 
                        value={Math.min(workload.allocationPercentage, 100)} 
                        className="h-2"
                      />
                      {isOverallocated && (
                        <p className="text-xs text-destructive flex items-center">
                          <AlertTriangle className="w-3 h-3 mr-1" />
                          {workload.allocationPercentage - 100}% over capacity
                        </p>
                      )}
                    </div>

                    {/* Stats Grid */}
                    <div className="grid grid-cols-3 gap-4 pt-2">
                      <div className="text-center">
                        <p className="text-lg font-semibold">{workload.assignedTasks}</p>
                        <p className="text-xs text-muted-foreground">Tasks</p>
                      </div>
                      <div className="text-center">
                        <p className="text-lg font-semibold">{workload.totalEstimatedHours}h</p>
                        <p className="text-xs text-muted-foreground">Est. Hours</p>
                      </div>
                      <div className="text-center">
                        <p className="text-lg font-semibold">{workload.availabilityThisWeek}h</p>
                        <p className="text-xs text-muted-foreground">Available</p>
                      </div>
                    </div>

                    {/* Skills */}
                    <div className="pt-2">
                      <p className="text-xs text-muted-foreground mb-2">Skills</p>
                      <div className="flex flex-wrap gap-1">
                        {mockData.getUserSkills(user.id).slice(0, 4).map(({ skill, proficiency }) => (
                          <Badge key={skill.id} variant="secondary" className="text-xs">
                            {skill.name} ({proficiency})
                          </Badge>
                        ))}
                        {mockData.getUserSkills(user.id).length > 4 && (
                          <Badge variant="outline" className="text-xs">
                            +{mockData.getUserSkills(user.id).length - 4} more
                          </Badge>
                        )}
                      </div>
                    </div>
                  </CardContent>
                </Card>
              );
            })}
          </div>
        </TabsContent>

        <TabsContent value="overallocated" className="space-y-4">
          {overallocated.length === 0 ? (
            <Card className="border-dashed">
              <CardContent className="p-8 text-center">
                <CheckCircle2 className="w-12 h-12 text-green-500 mx-auto mb-4" />
                <h3 className="text-lg font-medium">No overallocation!</h3>
                <p className="text-muted-foreground">All team members are within capacity.</p>
              </CardContent>
            </Card>
          ) : (
            <div className="grid gap-4 md:grid-cols-2">
              {overallocated.map((workload) => {
                const user = mockData.getUserById(workload.personId);
                if (!user) return null;
                
                return (
                  <Card key={user.id} className="border-destructive">
                    <CardHeader>
                      <div className="flex items-center gap-4">
                        <Avatar className="w-12 h-12">
                          <AvatarImage src={user.avatarUrl} />
                          <AvatarFallback>{user.name.charAt(0)}</AvatarFallback>
                        </Avatar>
                        <div className="flex-1">
                          <CardTitle className="text-lg">{user.name}</CardTitle>
                          <div className="flex items-center gap-2 text-destructive">
                            <TrendingUp className="w-4 h-4" />
                            <span className="font-medium">{workload.allocationPercentage}% allocated</span>
                          </div>
                        </div>
                      </div>
                    </CardHeader>
                    <CardContent>
                      <p className="text-sm text-muted-foreground mb-4">
                        {user.name} has {workload.totalEstimatedHours} hours assigned this week,
                        which is {workload.allocationPercentage - 100}% over their 40-hour capacity.
                      </p>
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
          <div className="grid gap-4 md:grid-cols-2">
            {underutilized.map((workload) => {
              const user = mockData.getUserById(workload.personId);
              if (!user) return null;
              
              return (
                <Card key={user.id}>
                  <CardHeader>
                    <div className="flex items-center gap-4">
                      <Avatar className="w-12 h-12">
                        <AvatarImage src={user.avatarUrl} />
                        <AvatarFallback>{user.name.charAt(0)}</AvatarFallback>
                      </Avatar>
                      <div className="flex-1">
                        <CardTitle className="text-lg">{user.name}</CardTitle>
                        <div className="flex items-center gap-2 text-blue-500">
                          <TrendingDown className="w-4 h-4" />
                          <span className="font-medium">{workload.allocationPercentage}% allocated</span>
                        </div>
                      </div>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <p className="text-sm text-muted-foreground mb-4">
                      {user.name} has {workload.availabilityThisWeek} hours available this week.
                    </p>
                    <Button size="sm">
                      Assign Task
                    </Button>
                  </CardContent>
                </Card>
              );
            })}
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
}
