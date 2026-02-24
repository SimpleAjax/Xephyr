'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Progress } from '@/components/ui/progress';
import { useSkills, useUsers, useSkillGaps } from '@/hooks/api';
import { Skill, User } from '@/types';
import { 
  Search, 
  Plus, 
  GraduationCap,
  Users,
  AlertTriangle,
  CheckCircle2,
  Code,
  Server,
  Palette,
  Cloud,
  Megaphone,
  Settings,
  Loader2
} from 'lucide-react';
import { useState, useMemo } from 'react';

const categoryIcons: Record<string, React.ReactNode> = {
  Frontend: <Code className="w-4 h-4" />,
  Backend: <Server className="w-4 h-4" />,
  Design: <Palette className="w-4 h-4" />,
  DevOps: <Cloud className="w-4 h-4" />,
  Marketing: <Megaphone className="w-4 h-4" />,
  Management: <Settings className="w-4 h-4" />,
};

export default function SkillsPage() {
  const [searchQuery, setSearchQuery] = useState('');
  const { data: skillsData, isLoading: skillsLoading } = useSkills();
  const { data: usersData } = useUsers();
  const { data: skillGapsData } = useSkillGaps();
  
  const skills = skillsData?.data?.skills || [];
  const users = usersData?.data?.users || [];
  const skillGaps = skillGapsData?.data?.gaps || [];
  
  const filteredSkills = skills.filter((skill: Skill) => 
    skill.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    skill.category.toLowerCase().includes(searchQuery.toLowerCase())
  );

  // Group skills by category
  const skillsByCategory = useMemo(() => {
    return filteredSkills.reduce((acc: Record<string, Skill[]>, skill: Skill) => {
      if (!acc[skill.category]) acc[skill.category] = [];
      acc[skill.category].push(skill);
      return acc;
    }, {} as Record<string, Skill[]>);
  }, [filteredSkills]);

  // Calculate skill coverage (mock calculation based on available data)
  const getSkillCoverage = (skillId: string) => {
    // This is a simplified calculation - in real implementation, 
    // you'd fetch this from the API
    const usersWithSkill = users.filter((user: User) => {
      // Mock: assume some users have each skill
      return user.id.charCodeAt(0) % 3 === 0 || user.id.charCodeAt(0) % 5 === 0;
    });
    const count = Math.max(1, usersWithSkill.length % users.length);
    return {
      count,
      percentage: Math.round((count / Math.max(1, users.length)) * 100),
      avgProficiency: Math.floor(Math.random() * 2) + 2, // Mock: 2-4 proficiency
    };
  };

  // Check if skill is a gap
  const isSkillGap = (skillId: string) => {
    return skillGaps.some((gap: { skillId: string }) => gap.skillId === skillId);
  };

  if (skillsLoading) {
    return (
      <div className="p-8 h-full flex items-center justify-center">
        <div className="flex items-center gap-2 text-muted-foreground">
          <Loader2 className="w-5 h-5 animate-spin" />
          <span>Loading skills...</span>
        </div>
      </div>
    );
  }

  return (
    <div className="p-8 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-start">
        <div>
          <h1 className="text-3xl font-bold">Skills Catalog</h1>
          <p className="text-muted-foreground">Manage team skills and identify gaps</p>
        </div>
        <Button>
          <Plus className="w-4 h-4 mr-2" />
          Add Skill
        </Button>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Total Skills</CardDescription>
            <CardTitle className="text-2xl">{skills.length}</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Categories</CardDescription>
            <CardTitle className="text-2xl">{Object.keys(skillsByCategory).length}</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription className="text-destructive">Skill Gaps</CardDescription>
            <CardTitle className="text-2xl text-destructive">{skillGaps.length}</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Avg Coverage</CardDescription>
            <CardTitle className="text-2xl">
              {skills.length > 0 
                ? Math.round(skills.reduce((sum: number, s: Skill) => sum + getSkillCoverage(s.id).percentage, 0) / skills.length)
                : 0}%
            </CardTitle>
          </CardHeader>
        </Card>
      </div>

      {/* Skill Gaps Alert */}
      {skillGaps.length > 0 && (
        <Card className="border-destructive">
          <CardHeader>
            <div className="flex items-center gap-2 text-destructive">
              <AlertTriangle className="w-5 h-5" />
              <CardTitle className="text-lg">Skill Gaps Detected</CardTitle>
            </div>
            <CardDescription>
              The following skills are required by tasks but no team member has them:
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-wrap gap-2">
              {skillGaps.map((gap: { skillId: string; skillName: string }) => (
                <Badge key={gap.skillId} variant="destructive">
                  {gap.skillName}
                </Badge>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Search */}
      <div className="relative max-w-sm">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
        <Input 
          placeholder="Search skills..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="pl-10"
        />
      </div>

      {/* Skills by Category */}
      <div className="space-y-6">
        {Object.entries(skillsByCategory).map(([category, categorySkills]: [string, Skill[]]) => (
          <div key={category}>
            <h2 className="text-lg font-semibold mb-4 flex items-center gap-2">
              {categoryIcons[category] || <GraduationCap className="w-4 h-4" />}
              {category}
            </h2>
            <div className="grid gap-4 md:grid-cols-3">
              {categorySkills.map((skill: Skill) => {
                const coverage = getSkillCoverage(skill.id);
                const isGap = isSkillGap(skill.id);
                
                return (
                  <Card key={skill.id} className={isGap ? 'border-destructive/50' : ''}>
                    <CardHeader className="pb-2">
                      <div className="flex justify-between items-start">
                        <CardTitle className="text-base">{skill.name}</CardTitle>
                        {isGap ? (
                          <Badge variant="destructive" className="text-xs">Gap</Badge>
                        ) : (
                          <Badge variant="outline" className="text-xs">
                            <Users className="w-3 h-3 mr-1" />
                            {coverage.count}
                          </Badge>
                        )}
                      </div>
                    </CardHeader>
                    <CardContent className="space-y-3">
                      <div className="space-y-1">
                        <div className="flex justify-between text-sm">
                          <span className="text-muted-foreground">Team Coverage</span>
                          <span className="font-medium">{coverage.percentage}%</span>
                        </div>
                        <Progress 
                          value={coverage.percentage} 
                          className="h-1.5"
                        />
                      </div>
                      
                      {coverage.count > 0 && (
                        <div className="flex justify-between text-sm">
                          <span className="text-muted-foreground">Avg Proficiency</span>
                          <span className="font-medium">{coverage.avgProficiency}/4</span>
                        </div>
                      )}
                      
                      {/* Users with this skill - mock display */}
                      <div className="pt-2">
                        <div className="flex -space-x-2">
                          {users
                            .filter((_: User, i: number) => i % 3 === 0 || i % 5 === 0)
                            .slice(0, Math.min(4, coverage.count))
                            .map((user: User) => (
                              <div 
                                key={user.id}
                                className="w-7 h-7 rounded-full bg-primary/20 flex items-center justify-center text-xs font-medium border-2 border-background"
                                title={user.name}
                              >
                                {user.name.charAt(0)}
                              </div>
                            ))}
                          {coverage.count > 4 && (
                            <div className="w-7 h-7 rounded-full bg-muted flex items-center justify-center text-xs font-medium border-2 border-background">
                              +{coverage.count - 4}
                            </div>
                          )}
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                );
              })}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
