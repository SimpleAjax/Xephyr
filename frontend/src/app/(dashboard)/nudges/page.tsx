'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Skeleton } from '@/components/ui/skeleton';
import { useNudges, useNudgeAction, useUpdateNudgeStatus } from '@/hooks/api/useNudges';
import { useNudgeStats } from '@/hooks/api/useNudges';
import { Nudge, NudgeType } from '@/lib/api/types';
import { toast } from 'sonner';
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
  Eye,
  RefreshCw
} from 'lucide-react';

const nudgeTypeConfig: Record<NudgeType, { icon: React.ElementType; color: string; bg: string; label: string }> = {
  overload: { icon: TrendingUp, color: 'text-orange-500', bg: 'bg-orange-50', label: 'Overload' },
  delay_risk: { icon: Clock, color: 'text-yellow-500', bg: 'bg-yellow-50', label: 'Delay Risk' },
  skill_gap: { icon: AlertCircle, color: 'text-purple-500', bg: 'bg-purple-50', label: 'Skill Gap' },
  unassigned: { icon: Users, color: 'text-blue-500', bg: 'bg-blue-50', label: 'Unassigned' },
  blocked: { icon: AlertTriangle, color: 'text-red-500', bg: 'bg-red-50', label: 'Blocked' },
  conflict: { icon: Zap, color: 'text-pink-500', bg: 'bg-pink-50', label: 'Conflict' },
  dependency_block: { icon: FolderKanban, color: 'text-cyan-500', bg: 'bg-cyan-50', label: 'Dependency' },
};

export default function NudgesPage() {
  // Fetch real nudges from API
  const { 
    data: nudgesData, 
    isLoading: nudgesLoading,
    refetch: refetchNudges 
  } = useNudges({ limit: 100 });
  
  const { 
    data: statsData, 
    isLoading: statsLoading 
  } = useNudgeStats('30d');

  const nudgeAction = useNudgeAction();
  const updateStatus = useUpdateNudgeStatus();

  const nudges = nudgesData?.data?.nudges || [];
  const summary = nudgesData?.data?.summary;
  const stats = statsData?.data;
  
  const unreadNudges = nudges.filter((n: Nudge) => n.status === 'unread');
  const readNudges = nudges.filter((n: Nudge) => n.status === 'read' || n.status === 'acted');
  const dismissedNudges = nudges.filter((n: Nudge) => n.status === 'dismissed');
  
  const highSeverity = unreadNudges.filter((n: Nudge) => n.severity === 'high');
  const mediumSeverity = unreadNudges.filter((n: Nudge) => n.severity === 'medium');
  const lowSeverity = unreadNudges.filter((n: Nudge) => n.severity === 'low');

  const handleTakeAction = async (nudgeId: string) => {
    try {
      await nudgeAction.mutateAsync({
        nudgeId,
        action: { actionType: 'accept_suggestion' }
      });
      toast.success('Action taken successfully');
    } catch (error) {
      toast.error('Failed to take action');
    }
  };

  const handleMarkAsRead = async (nudgeId: string) => {
    try {
      await updateStatus.mutateAsync({ nudgeId, status: 'read' });
      toast.success('Marked as read');
    } catch (error) {
      toast.error('Failed to mark as read');
    }
  };

  const handleDismiss = async (nudgeId: string) => {
    try {
      await updateStatus.mutateAsync({ nudgeId, status: 'dismissed' });
      toast.success('Nudge dismissed');
    } catch (error) {
      toast.error('Failed to dismiss nudge');
    }
  };

  const renderNudgeCard = (nudge: Nudge, showActions = true) => {
    const config = nudgeTypeConfig[nudge.type];
    const Icon = config.icon;
    
    // Related entities info comes from nudge data or can be fetched via API
    const relatedProjectName = nudge.relatedEntities?.projectId;
    const relatedPersonName = nudge.relatedEntities?.personId;

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
                  {relatedProjectName && (
                    <span className="flex items-center gap-1">
                      <FolderKanban className="w-3 h-3" />
                      Project: {relatedProjectName.slice(0, 8)}...
                    </span>
                  )}
                  {relatedPersonName && (
                    <span className="flex items-center gap-1">
                      <Users className="w-3 h-3" />
                      User: {relatedPersonName.slice(0, 8)}...
                    </span>
                  )}
                  <span>
                    {new Date(nudge.createdAt).toLocaleDateString()}
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
                <p className="text-sm text-muted-foreground">{nudge.suggestedAction.description}</p>
              </div>
            </div>
          )}

          {/* Actions */}
          {showActions && nudge.status === 'unread' && (
            <div className="flex gap-2 pt-2">
              <Button 
                onClick={() => handleTakeAction(nudge.id)}
                className="flex-1"
                disabled={nudgeAction.isPending}
              >
                <CheckCircle2 className="w-4 h-4 mr-2" />
                Take Action
              </Button>
              <Button 
                variant="outline" 
                onClick={() => handleMarkAsRead(nudge.id)}
                disabled={updateStatus.isPending}
              >
                <Eye className="w-4 h-4 mr-2" />
                Mark Read
              </Button>
              <Button 
                variant="outline" 
                onClick={() => handleDismiss(nudge.id)}
                disabled={updateStatus.isPending}
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

  const isLoading = nudgesLoading || statsLoading;

  return (
    <div className="p-8 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-start">
        <div>
          <h1 className="text-3xl font-bold">AI Nudges</h1>
          <p className="text-muted-foreground">Proactive alerts for potential issues</p>
        </div>
        <div className="flex items-center gap-4">
          <Button 
            variant="outline" 
            onClick={() => refetchNudges()}
            disabled={isLoading}
          >
            <RefreshCw className={`w-4 h-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
        </div>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Total Nudges</CardDescription>
            <CardTitle className="text-2xl">
              {statsLoading ? <Skeleton className="h-8 w-16" /> : stats?.generated || nudges.length}
            </CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription className="text-destructive">High Priority</CardDescription>
            <CardTitle className="text-2xl text-destructive">
              {nudgesLoading ? <Skeleton className="h-8 w-16" /> : summary?.bySeverity?.high || highSeverity.length}
            </CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription className="text-yellow-600">Medium Priority</CardDescription>
            <CardTitle className="text-2xl text-yellow-600">
              {nudgesLoading ? <Skeleton className="h-8 w-16" /> : summary?.bySeverity?.medium || mediumSeverity.length}
            </CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Unread</CardDescription>
            <CardTitle className="text-2xl">
              {nudgesLoading ? <Skeleton className="h-8 w-16" /> : summary?.unread || unreadNudges.length}
            </CardTitle>
          </CardHeader>
        </Card>
      </div>

      {/* Nudges by Type Summary */}
      {isLoading ? (
        <div className="flex flex-wrap gap-2">
          <Skeleton className="h-6 w-24" />
          <Skeleton className="h-6 w-24" />
          <Skeleton className="h-6 w-24" />
        </div>
      ) : (
        <div className="flex flex-wrap gap-2">
          {Object.entries(nudgeTypeConfig).map(([type, config]) => {
            const count = unreadNudges.filter((n: Nudge) => n.type === type).length;
            const summaryCount = summary?.byType?.[type];
            const displayCount = summaryCount !== undefined ? summaryCount : count;
            if (displayCount === 0) return null;
            const Icon = config.icon;
            return (
              <Badge key={type} variant="outline" className={`flex items-center gap-1 px-3 py-1 ${config.bg}`}>
                <Icon className={`w-3 h-3 ${config.color}`} />
                {config.label}: {displayCount}
              </Badge>
            );
          })}
        </div>
      )}

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
          {nudgesLoading ? (
            <div className="space-y-4">
              {Array.from({ length: 3 }).map((_, i) => (
                <Card key={i}>
                  <CardHeader>
                    <Skeleton className="h-6 w-3/4" />
                    <Skeleton className="h-4 w-full mt-2" />
                  </CardHeader>
                  <CardContent>
                    <Skeleton className="h-16 w-full" />
                  </CardContent>
                </Card>
              ))}
            </div>
          ) : unreadNudges.length === 0 ? (
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
                    {highSeverity.map((nudge: Nudge) => renderNudgeCard(nudge))}
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
                    {mediumSeverity.map((nudge: Nudge) => renderNudgeCard(nudge))}
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
                    {lowSeverity.map((nudge: Nudge) => renderNudgeCard(nudge))}
                  </div>
                </div>
              )}
            </div>
          )}
        </TabsContent>

        <TabsContent value="all" className="space-y-4">
          {nudgesLoading ? (
            <div className="space-y-4">
              {Array.from({ length: 5 }).map((_, i) => (
                <Card key={i}>
                  <CardHeader>
                    <Skeleton className="h-6 w-3/4" />
                    <Skeleton className="h-4 w-full mt-2" />
                  </CardHeader>
                  <CardContent>
                    <Skeleton className="h-16 w-full" />
                  </CardContent>
                </Card>
              ))}
            </div>
          ) : (
            <div className="space-y-4">
              {nudges.map((nudge: Nudge) => renderNudgeCard(nudge, nudge.status === 'unread'))}
            </div>
          )}
        </TabsContent>

        <TabsContent value="history" className="space-y-4">
          <div className="space-y-4">
            {readNudges.map((nudge: Nudge) => renderNudgeCard(nudge, false))}
            {dismissedNudges.map((nudge: Nudge) => renderNudgeCard(nudge, false))}
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
