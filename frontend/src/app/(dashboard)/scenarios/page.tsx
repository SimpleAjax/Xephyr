'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { useScenarios } from '@/hooks/api';
import { Scenario } from '@/lib/api/types';
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
  Loader2
} from 'lucide-react';

// Map API scenario status to UI status
const mapScenarioStatus = (status: Scenario['status']): 'pending' | 'approved' | 'rejected' => {
  switch (status) {
    case 'applied':
      return 'approved';
    case 'rejected':
      return 'rejected';
    case 'draft':
    case 'pending':
    case 'simulated':
    default:
      return 'pending';
  }
};

export default function ScenariosPage() {
  const { data: scenariosData, isLoading } = useScenarios();
  
  const scenarios = scenariosData?.data?.scenarios || [];
  
  // Mapped scenario type with UI-specific status
  type MappedScenario = {
    status: 'pending' | 'approved' | 'rejected';
    scenarioId: string;
    title: string;
    description: string;
    changeType: Scenario['changeType'];
    proposedChanges: Record<string, any>;
    impactAnalysis?: Scenario['impactAnalysis'] & {
      aiRecommendations?: Array<{ action: string }>;
    };
    simulationStatus: Scenario['simulationStatus'];
    createdAt: string;
    history?: Scenario['history'];
  };

  const mappedScenarios: MappedScenario[] = scenarios.map(s => ({
    ...s,
    status: mapScenarioStatus(s.status),
  }));
  
  const pendingScenarios = mappedScenarios.filter(s => s.status === 'pending');
  const approvedScenarios = mappedScenarios.filter(s => s.status === 'approved');
  const rejectedScenarios = mappedScenarios.filter(s => s.status === 'rejected');

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

  // Calculate total delay from impact analysis
  const calculateDelay = (scenario: { impactAnalysis?: MappedScenario['impactAnalysis'] }): number => {
    if (!scenario.impactAnalysis) return 0;
    // Use timelineComparison if available
    if ('timelineComparison' in scenario.impactAnalysis && scenario.impactAnalysis.timelineComparison) {
      return scenario.impactAnalysis.timelineComparison.totalDelayDays * 24;
    }
    // Fallback: sum affected tasks delay
    return scenario.impactAnalysis.affectedTasks?.reduce((sum, t) => sum + ((t as { delayDays?: number }).delayDays || 0), 0) * 24 || 0;
  };

  // Calculate cost impact
  const calculateCost = (scenario: { impactAnalysis?: MappedScenario['impactAnalysis'] }): number => {
    if (!scenario.impactAnalysis?.costAnalysis) return 0;
    return scenario.impactAnalysis.costAnalysis.totalCost || 0;
  };

  if (isLoading) {
    return (
      <div className="p-8 h-full flex items-center justify-center">
        <div className="flex items-center gap-2 text-muted-foreground">
          <Loader2 className="w-5 h-5 animate-spin" />
          <span>Loading scenarios...</span>
        </div>
      </div>
    );
  }

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
            <CardTitle className="text-2xl">{mappedScenarios.length}</CardTitle>
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
              {mappedScenarios.length > 0 
                ? Math.round(mappedScenarios.reduce((sum, s) => sum + calculateDelay(s as unknown as Scenario), 0) / mappedScenarios.length)
                : 0}h
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
              {pendingScenarios.map((scenario: MappedScenario) => (
                <Card key={scenario.scenarioId} className="border-l-4 border-l-yellow-500">
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
                          +{calculateDelay(scenario)}h
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
                            <p className="text-lg font-semibold">{scenario.impactAnalysis.affectedProjects?.length || 0}</p>
                          </div>
                          <div>
                            <p className="text-sm text-muted-foreground">Affected Tasks</p>
                            <p className="text-lg font-semibold">{scenario.impactAnalysis.affectedTasks?.length || 0}</p>
                          </div>
                          <div>
                            <p className="text-sm text-muted-foreground">Cost Impact</p>
                            <p className="text-lg font-semibold">${calculateCost(scenario).toLocaleString()}</p>
                          </div>
                        </div>
                        
                        {scenario.impactAnalysis.aiRecommendations && (
                          <div>
                            <p className="text-sm font-medium mb-2">AI Recommendations:</p>
                            <ul className="space-y-1">
                              {scenario.impactAnalysis.aiRecommendations.map((rec, i) => (
                                <li key={i} className="text-sm text-muted-foreground flex items-start gap-2">
                                  <Sparkles className="w-4 h-4 text-primary mt-0.5" />
                                  {rec.action}
                                </li>
                              ))}
                            </ul>
                          </div>
                        )}

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
              {approvedScenarios.map((scenario: MappedScenario) => (
                <Card key={scenario.scenarioId} className="border-l-4 border-l-green-500">
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
              {rejectedScenarios.map((scenario: MappedScenario) => (
                <Card key={scenario.scenarioId} className="border-l-4 border-l-destructive opacity-60">
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
            {mappedScenarios.map((scenario: MappedScenario) => (
              <Card key={scenario.scenarioId}>
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
                          +{calculateDelay(scenario)}h
                        </div>
                        <div className="text-xs text-muted-foreground">
                          ${calculateCost(scenario).toLocaleString()}
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
