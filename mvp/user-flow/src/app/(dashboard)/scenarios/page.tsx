'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { mockData } from '@/data/mock-data';
import { 
  Sparkles, 
  Plus, 
  Clock, 
  CheckCircle2, 
  XCircle,
  AlertCircle,
  Users,
  Briefcase,
  TrendingUp,
  DollarSign,
  ArrowRight
} from 'lucide-react';

export default function ScenariosPage() {
  const scenarios = mockData.getScenarios();
  
  const pendingScenarios = scenarios.filter(s => s.status === 'pending');
  const approvedScenarios = scenarios.filter(s => s.status === 'approved');
  const rejectedScenarios = scenarios.filter(s => s.status === 'rejected');

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'pending': return <Clock className="w-5 h-5 text-yellow-500" />;
      case 'approved': return <CheckCircle2 className="w-5 h-5 text-green-500" />;
      case 'rejected': return <XCircle className="w-5 h-5 text-destructive" />;
      default: return <AlertCircle className="w-5 h-5" />;
    }
  };

  const getChangeTypeIcon = (type: string) => {
    switch (type) {
      case 'employee_leave': return <Users className="w-4 h-4" />;
      case 'scope_change': return <Briefcase className="w-4 h-4" />;
      case 'reallocation': return <TrendingUp className="w-4 h-4" />;
      case 'priority_shift': return <AlertCircle className="w-4 h-4" />;
      default: return <Sparkles className="w-4 h-4" />;
    }
  };

  const getChangeTypeLabel = (type: string) => {
    return type.split('_').map(word => 
      word.charAt(0).toUpperCase() + word.slice(1)
    ).join(' ');
  };

  return (
    <div className="p-8 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-start">
        <div>
          <h1 className="text-3xl font-bold">What-If Scenarios</h1>
          <p className="text-muted-foreground">Simulate changes and see impact before committing</p>
        </div>
        <Button>
          <Plus className="w-4 h-4 mr-2" />
          New Scenario
        </Button>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Total Scenarios</CardDescription>
            <CardTitle className="text-2xl">{scenarios.length}</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription className="text-yellow-600">Pending</CardDescription>
            <CardTitle className="text-2xl text-yellow-600">{pendingScenarios.length}</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription className="text-green-600">Approved</CardDescription>
            <CardTitle className="text-2xl text-green-600">{approvedScenarios.length}</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Avg Delay Impact</CardDescription>
            <CardTitle className="text-2xl">
              {Math.round(scenarios.reduce((sum, s) => sum + (s.impactAnalysis?.delayHoursTotal || 0), 0) / scenarios.length)}h
            </CardTitle>
          </CardHeader>
        </Card>
      </div>

      {/* Scenarios Tabs */}
      <Tabs defaultValue="pending" className="space-y-4">
        <TabsList>
          <TabsTrigger value="pending">
            Pending
            {pendingScenarios.length > 0 && (
              <Badge variant="secondary" className="ml-2 h-5 px-1.5">
                {pendingScenarios.length}
              </Badge>
            )}
          </TabsTrigger>
          <TabsTrigger value="approved">Approved</TabsTrigger>
          <TabsTrigger value="rejected">Rejected</TabsTrigger>
          <TabsTrigger value="all">All</TabsTrigger>
        </TabsList>

        <TabsContent value="pending" className="space-y-4">
          {pendingScenarios.length === 0 ? (
            <Card className="border-dashed">
              <CardContent className="p-8 text-center">
                <Sparkles className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                <h3 className="text-lg font-medium">No pending scenarios</h3>
                <p className="text-muted-foreground">Create a scenario to see potential impacts.</p>
              </CardContent>
            </Card>
          ) : (
            <div className="grid gap-4">
              {pendingScenarios.map((scenario) => (
                <Card key={scenario.id} className="border-l-4 border-l-yellow-500">
                  <CardHeader>
                    <div className="flex justify-between items-start">
                      <div className="flex items-start gap-4">
                        <div className="mt-1">{getStatusIcon(scenario.status)}</div>
                        <div>
                          <CardTitle className="text-lg">{scenario.title}</CardTitle>
                          <CardDescription className="mt-1 max-w-2xl">
                            {scenario.description}
                          </CardDescription>
                          <div className="flex items-center gap-2 mt-2">
                            <Badge variant="outline" className="flex items-center gap-1">
                              {getChangeTypeIcon(scenario.changeType)}
                              {getChangeTypeLabel(scenario.changeType)}
                            </Badge>
                            <Badge variant="secondary">Pending Review</Badge>
                          </div>
                        </div>
                      </div>
                      <div className="text-right">
                        <div className="text-2xl font-bold text-destructive">
                          +{scenario.impactAnalysis?.delayHoursTotal}h
                        </div>
                        <div className="text-xs text-muted-foreground">Delay Impact</div>
                      </div>
                    </div>
                  </CardHeader>
                  <CardContent>
                    {scenario.impactAnalysis && (
                      <div className="space-y-4">
                        <div className="grid grid-cols-3 gap-4 p-4 bg-muted rounded-lg">
                          <div>
                            <p className="text-sm text-muted-foreground">Affected Projects</p>
                            <p className="text-lg font-semibold">{scenario.impactAnalysis.affectedProjects.length}</p>
                          </div>
                          <div>
                            <p className="text-sm text-muted-foreground">Affected Tasks</p>
                            <p className="text-lg font-semibold">{scenario.impactAnalysis.affectedTasks.length}</p>
                          </div>
                          <div>
                            <p className="text-sm text-muted-foreground">Cost Impact</p>
                            <p className="text-lg font-semibold">${scenario.impactAnalysis.costImpact.toLocaleString()}</p>
                          </div>
                        </div>
                        
                        <div>
                          <p className="text-sm font-medium mb-2">AI Recommendations:</p>
                          <ul className="space-y-1">
                            {scenario.impactAnalysis.recommendations.map((rec, i) => (
                              <li key={i} className="text-sm text-muted-foreground flex items-start gap-2">
                                <Sparkles className="w-4 h-4 text-primary mt-0.5" />
                                {rec}
                              </li>
                            ))}
                          </ul>
                        </div>

                        <div className="flex gap-2 pt-2">
                          <Button className="flex-1">
                            <CheckCircle2 className="w-4 h-4 mr-2" />
                            Approve
                          </Button>
                          <Button variant="outline" className="flex-1">
                            Modify
                          </Button>
                          <Button variant="outline" className="text-destructive hover:text-destructive">
                            <XCircle className="w-4 h-4" />
                          </Button>
                        </div>
                      </div>
                    )}
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </TabsContent>

        <TabsContent value="approved" className="space-y-4">
          {approvedScenarios.length === 0 ? (
            <Card className="border-dashed">
              <CardContent className="p-8 text-center">
                <CheckCircle2 className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                <h3 className="text-lg font-medium">No approved scenarios</h3>
                <p className="text-muted-foreground">Approved scenarios will appear here.</p>
              </CardContent>
            </Card>
          ) : (
            <div className="grid gap-4">
              {approvedScenarios.map((scenario) => (
                <Card key={scenario.id} className="border-l-4 border-l-green-500">
                  <CardHeader>
                    <div className="flex justify-between items-start">
                      <div className="flex items-center gap-4">
                        <CheckCircle2 className="w-5 h-5 text-green-500" />
                        <div>
                          <CardTitle className="text-lg">{scenario.title}</CardTitle>
                          <CardDescription>{scenario.description}</CardDescription>
                        </div>
                      </div>
                      <Badge>Approved</Badge>
                    </div>
                  </CardHeader>
                </Card>
              ))}
            </div>
          )}
        </TabsContent>

        <TabsContent value="rejected" className="space-y-4">
          {rejectedScenarios.length === 0 ? (
            <Card className="border-dashed">
              <CardContent className="p-8 text-center">
                <XCircle className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                <h3 className="text-lg font-medium">No rejected scenarios</h3>
                <p className="text-muted-foreground">Rejected scenarios will appear here.</p>
              </CardContent>
            </Card>
          ) : (
            <div className="grid gap-4">
              {rejectedScenarios.map((scenario) => (
                <Card key={scenario.id} className="border-l-4 border-l-destructive opacity-60">
                  <CardHeader>
                    <div className="flex justify-between items-start">
                      <div className="flex items-center gap-4">
                        <XCircle className="w-5 h-5 text-destructive" />
                        <div>
                          <CardTitle className="text-lg">{scenario.title}</CardTitle>
                          <CardDescription>{scenario.description}</CardDescription>
                        </div>
                      </div>
                      <Badge variant="destructive">Rejected</Badge>
                    </div>
                  </CardHeader>
                </Card>
              ))}
            </div>
          )}
        </TabsContent>

        <TabsContent value="all" className="space-y-4">
          <div className="grid gap-4">
            {scenarios.map((scenario) => (
              <Card key={scenario.id}>
                <CardHeader>
                  <div className="flex justify-between items-start">
                    <div className="flex items-center gap-4">
                      {getStatusIcon(scenario.status)}
                      <div>
                        <CardTitle className="text-lg">{scenario.title}</CardTitle>
                        <CardDescription>{scenario.description}</CardDescription>
                        <div className="flex items-center gap-2 mt-2">
                          <Badge variant="outline">{getChangeTypeLabel(scenario.changeType)}</Badge>
                          <Badge variant={scenario.status === 'approved' ? 'default' : scenario.status === 'rejected' ? 'destructive' : 'secondary'}>
                            {scenario.status}
                          </Badge>
                        </div>
                      </div>
                    </div>
                    {scenario.impactAnalysis && (
                      <div className="text-right">
                        <div className="text-xl font-bold">
                          +{scenario.impactAnalysis.delayHoursTotal}h
                        </div>
                        <div className="text-xs text-muted-foreground">
                          ${scenario.impactAnalysis.costImpact.toLocaleString()}
                        </div>
                      </div>
                    )}
                  </div>
                </CardHeader>
              </Card>
            ))}
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
}
