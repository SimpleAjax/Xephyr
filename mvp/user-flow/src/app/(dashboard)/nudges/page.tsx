'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { useAppStore } from '@/store/app-store';
import { mockData } from '@/data/mock-data';
import { 
  Bell, 
  CheckCircle2, 
  XCircle,
  TrendingUp,
  Clock,
  AlertCircle,
  Users,
  AlertTriangle,
  Zap,
  FolderKanban,
  Sparkles,
  Trash2,
  Eye
} from 'lucide-react';

const nudgeTypeConfig = {
  overload: { icon: TrendingUp, color: 'text-orange-500', bg: 'bg-orange-50', label: 'Overload' },
  delay_risk: { icon: Clock, color: 'text-yellow-500', bg: 'bg-yellow-50', label: 'Delay Risk' },
  skill_gap: { icon: AlertCircle, color: 'text-purple-500', bg: 'bg-purple-50', label: 'Skill Gap' },
  unassigned: { icon: Users, color: 'text-blue-500', bg: 'bg-blue-50', label: 'Unassigned' },
  blocked: { icon: AlertTriangle, color: 'text-red-500', bg: 'bg-red-50', label: 'Blocked' },
  conflict: { icon: Zap, color: 'text-pink-500', bg: 'bg-pink-50', label: 'Conflict' },
  dependency_block: { icon: FolderKanban, color: 'text-cyan-500', bg: 'bg-cyan-50', label: 'Dependency' },
};

export default function NudgesPage() {
  const nudges = useAppStore((state) => state.nudges);
  const markNudgeAsRead = useAppStore((state) => state.markNudgeAsRead);
  const dismissNudge = useAppStore((state) => state.dismissNudge);
  const takeNudgeAction = useAppStore((state) => state.takeNudgeAction);
  
  const unreadNudges = nudges.filter(n => n.status === 'unread');
  const readNudges = nudges.filter(n => n.status === 'read' || n.status === 'acted');
  const dismissedNudges = nudges.filter(n => n.status === 'dismissed');
  
  const highSeverity = unreadNudges.filter(n => n.severity === 'high');
  const mediumSeverity = unreadNudges.filter(n => n.severity === 'medium');
  const lowSeverity = unreadNudges.filter(n => n.severity === 'low');

  const renderNudgeCard = (nudge: typeof nudges[0], showActions = true) => {
    const config = nudgeTypeConfig[nudge.type];
    const Icon = config.icon;
    
    const relatedProject = nudge.relatedProjectId 
      ? mockData.getProjectById(nudge.relatedProjectId)
      : null;
    const relatedPerson = nudge.relatedPersonId
      ? mockData.getUserById(nudge.relatedPersonId)
      : null;

    return (
      <Card 
        key={nudge.id} 
        className={`border-l-4 ${
          nudge.severity === 'high' ? 'border-l-destructive' : 
          nudge.severity === 'medium' ? 'border-l-yellow-500' : 
          'border-l-blue-500'
        }`}
      >
        <CardHeader>
          <div className="flex justify-between items-start">
            <div className="flex items-start gap-4">
              <div className={`p-2 rounded-lg ${config.bg}`}>
                <Icon className={`w-5 h-5 ${config.color}`} />
              </div>
              <div className="flex-1">
                <CardTitle className="text-lg">{nudge.title}</CardTitle>
                <CardDescription className="mt-1">{nudge.description}</CardDescription>
                
                {/* Nudge Meta */}
                <div className="flex items-center gap-2 mt-3">
                  <Badge variant="outline" className="flex items-center gap-1">
                    {config.label}
                  </Badge>
                  <Badge 
                    variant={nudge.severity === 'high' ? 'destructive' : 'secondary'}
                    className="text-xs"
                  >
                    {nudge.severity}
                  </Badge>
                  {nudge.status !== 'unread' && (
                    <Badge variant="outline" className="text-xs">
                      {nudge.status}
                    </Badge>
                  )}
                </div>

                {/* Related Entities */}
                <div className="flex items-center gap-4 mt-2 text-sm text-muted-foreground">
                  {relatedProject && (
                    <span className="flex items-center gap-1">
                      <FolderKanban className="w-3 h-3" />
                      {relatedProject.name}
                    </span>
                  )}
                  {relatedPerson && (
                    <span className="flex items-center gap-1">
                      <Users className="w-3 h-3" />
                      {relatedPerson.name}
                    </span>
                  )}
                  <span>
                    {nudge.createdAt.toLocaleDateString()}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* AI Explanation */}
          <div className="bg-muted p-4 rounded-lg">
            <div className="flex items-center gap-2 mb-2">
              <Sparkles className="w-4 h-4 text-primary" />
              <span className="text-sm font-medium">AI Analysis</span>
            </div>
            <p className="text-sm text-muted-foreground">{nudge.aiExplanation}</p>
          </div>

          {/* Suggested Action */}
          {nudge.suggestedAction && (
            <div className="flex items-start gap-2 p-3 bg-primary/5 rounded-lg border border-primary/20">
              <CheckCircle2 className="w-4 h-4 text-primary mt-0.5" />
              <div>
                <p className="text-sm font-medium">Suggested Action</p>
                <p className="text-sm text-muted-foreground">{nudge.suggestedAction}</p>
              </div>
            </div>
          )}

          {/* Actions */}
          {showActions && nudge.status === 'unread' && (
            <div className="flex gap-2 pt-2">
              <Button 
                onClick={() => takeNudgeAction(nudge.id)}
                className="flex-1"
              >
                <CheckCircle2 className="w-4 h-4 mr-2" />
                Take Action
              </Button>
              <Button 
                variant="outline" 
                onClick={() => markNudgeAsRead(nudge.id)}
              >
                <Eye className="w-4 h-4 mr-2" />
                Mark Read
              </Button>
              <Button 
                variant="outline" 
                onClick={() => dismissNudge(nudge.id)}
                className="text-destructive hover:text-destructive"
              >
                <XCircle className="w-4 h-4" />
              </Button>
            </div>
          )}
        </CardContent>
      </Card>
    );
  };

  return (
    <div className="p-8 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-start">
        <div>
          <h1 className="text-3xl font-bold">AI Nudges</h1>
          <p className="text-muted-foreground">Proactive alerts for potential issues</p>
        </div>
        <div className="flex items-center gap-4">
          <Button variant="outline">
            Configure Alerts
          </Button>
        </div>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Total Nudges</CardDescription>
            <CardTitle className="text-2xl">{nudges.length}</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription className="text-destructive">High Priority</CardDescription>
            <CardTitle className="text-2xl text-destructive">{highSeverity.length}</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription className="text-yellow-600">Medium Priority</CardDescription>
            <CardTitle className="text-2xl text-yellow-600">{mediumSeverity.length}</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Unread</CardDescription>
            <CardTitle className="text-2xl">{unreadNudges.length}</CardTitle>
          </CardHeader>
        </Card>
      </div>

      {/* Nudges by Type Summary */}
      <div className="flex flex-wrap gap-2">
        {Object.entries(nudgeTypeConfig).map(([type, config]) => {
          const count = unreadNudges.filter(n => n.type === type).length;
          if (count === 0) return null;
          const Icon = config.icon;
          return (
            <Badge key={type} variant="outline" className={`flex items-center gap-1 px-3 py-1 ${config.bg}`}>
              <Icon className={`w-3 h-3 ${config.color}`} />
              {config.label}: {count}
            </Badge>
          );
        })}
      </div>

      {/* Nudges Tabs */}
      <Tabs defaultValue="unread" className="space-y-4">
        <TabsList>
          <TabsTrigger value="unread">
            Unread
            {unreadNudges.length > 0 && (
              <Badge variant="destructive" className="ml-2 h-5 px-1.5">
                {unreadNudges.length}
              </Badge>
            )}
          </TabsTrigger>
          <TabsTrigger value="all">All Nudges</TabsTrigger>
          <TabsTrigger value="history">History</TabsTrigger>
        </TabsList>

        <TabsContent value="unread" className="space-y-4">
          {unreadNudges.length === 0 ? (
            <Card className="border-dashed">
              <CardContent className="p-8 text-center">
                <CheckCircle2 className="w-12 h-12 text-green-500 mx-auto mb-4" />
                <h3 className="text-lg font-medium">All caught up!</h3>
                <p className="text-muted-foreground">No unread nudges. Great job!</p>
              </CardContent>
            </Card>
          ) : (
            <div className="space-y-4">
              {highSeverity.length > 0 && (
                <div>
                  <h3 className="text-sm font-semibold text-destructive mb-3 flex items-center gap-2">
                    <AlertTriangle className="w-4 h-4" />
                    High Priority ({highSeverity.length})
                  </h3>
                  <div className="space-y-3">
                    {highSeverity.map(nudge => renderNudgeCard(nudge))}
                  </div>
                </div>
              )}
              
              {mediumSeverity.length > 0 && (
                <div>
                  <h3 className="text-sm font-semibold text-yellow-600 mb-3 flex items-center gap-2">
                    <Clock className="w-4 h-4" />
                    Medium Priority ({mediumSeverity.length})
                  </h3>
                  <div className="space-y-3">
                    {mediumSeverity.map(nudge => renderNudgeCard(nudge))}
                  </div>
                </div>
              )}
              
              {lowSeverity.length > 0 && (
                <div>
                  <h3 className="text-sm font-semibold text-blue-600 mb-3 flex items-center gap-2">
                    <Bell className="w-4 h-4" />
                    Low Priority ({lowSeverity.length})
                  </h3>
                  <div className="space-y-3">
                    {lowSeverity.map(nudge => renderNudgeCard(nudge))}
                  </div>
                </div>
              )}
            </div>
          )}
        </TabsContent>

        <TabsContent value="all" className="space-y-4">
          <div className="space-y-4">
            {nudges.map(nudge => renderNudgeCard(nudge, nudge.status === 'unread'))}
          </div>
        </TabsContent>

        <TabsContent value="history" className="space-y-4">
          <div className="space-y-4">
            {readNudges.map(nudge => renderNudgeCard(nudge, false))}
            {dismissedNudges.map(nudge => renderNudgeCard(nudge, false))}
            {readNudges.length === 0 && dismissedNudges.length === 0 && (
              <Card className="border-dashed">
                <CardContent className="p-8 text-center">
                  <Clock className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                  <h3 className="text-lg font-medium">No history yet</h3>
                  <p className="text-muted-foreground">Resolved nudges will appear here.</p>
                </CardContent>
              </Card>
            )}
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
}
