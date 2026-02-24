'use client';

import { useState } from 'react';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Sparkles, Plus, Check, RefreshCw, Send, User, Bot, ChevronRight, Loader2 } from 'lucide-react';
import { Project, Task } from '@/types';

interface CreateProjectDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onProjectCreate: (project: Project, tasks: Task[]) => void;
  isCreating?: boolean;
}

interface AIMessage {
  role: 'user' | 'assistant';
  content: string;
  suggestions?: ProjectSuggestion;
}

interface ProjectSuggestion {
  project: Partial<Project>;
  tasks: Partial<Task>[];
  reasoning: string[];
}

export function CreateProjectDialog({ open, onOpenChange, onProjectCreate, isCreating }: CreateProjectDialogProps) {
  const [activeTab, setActiveTab] = useState('manual');
  const [aiStep, setAiStep] = useState<'chat' | 'review' | 'editing'>('chat');
  
  // Manual form state
  const [projectName, setProjectName] = useState('');
  const [projectDescription, setProjectDescription] = useState('');
  const [projectPriority, setProjectPriority] = useState(50);
  const [targetDate, setTargetDate] = useState('');
  
  // AI conversation state
  const [messages, setMessages] = useState<AIMessage[]>([
    {
      role: 'assistant',
      content: "Hi! I'm Xephyr AI. I can help you create a new project with tasks and subtasks.\n\nTell me about your project - what are you building? Who's it for? Any specific requirements or deadlines?"
    }
  ]);
  const [inputMessage, setInputMessage] = useState('');
  const [isGenerating, setIsGenerating] = useState(false);
  const [aiSuggestion, setAiSuggestion] = useState<ProjectSuggestion | null>(null);
  
  // Editing state
  const [editingProject, setEditingProject] = useState<Partial<Project> | null>(null);
  const [editingTasks, setEditingTasks] = useState<Partial<Task>[]>([]);

  const resetForm = () => {
    setProjectName('');
    setProjectDescription('');
    setProjectPriority(50);
    setTargetDate('');
    setMessages([{
      role: 'assistant',
      content: "Hi! I'm Xephyr AI. I can help you create a new project with tasks and subtasks.\n\nTell me about your project - what are you building? Who's it for? Any specific requirements or deadlines?"
    }]);
    setAiStep('chat');
    setAiSuggestion(null);
    setInputMessage('');
    setEditingProject(null);
    setEditingTasks([]);
  };

  const handleClose = () => {
    resetForm();
    onOpenChange(false);
  };

  const handleManualSubmit = () => {
    const newProject: Project = {
      id: `proj-${Date.now()}`,
      name: projectName,
      description: projectDescription,
      status: 'active',
      priority: projectPriority,
      startDate: new Date(),
      targetEndDate: targetDate ? new Date(targetDate) : new Date(Date.now() + 90 * 24 * 60 * 60 * 1000),
      healthScore: 100,
      progress: 0,
    };
    onProjectCreate(newProject, []);
    handleClose();
  };

  const sendMessage = () => {
    if (!inputMessage.trim()) return;
    
    const userMessage: AIMessage = { role: 'user', content: inputMessage };
    setMessages(prev => [...prev, userMessage]);
    setInputMessage('');
    setIsGenerating(true);
    
    setTimeout(() => {
      const lowerInput = inputMessage.toLowerCase();
      
      if (lowerInput.includes('generate') || lowerInput.includes('create') || lowerInput.includes('build') || messages.length > 2) {
        const suggestion = generateAIProjectSuggestion(inputMessage);
        setAiSuggestion(suggestion);
        setEditingProject(suggestion.project);
        setEditingTasks(suggestion.tasks);
        
        const assistantMessage: AIMessage = {
          role: 'assistant',
          content: `I've analyzed your requirements and created a project structure. Here's what I recommend:

**${suggestion.project.name}**
${suggestion.project.description}

**Suggested Tasks (${suggestion.tasks.length}):**
${suggestion.tasks.map((t, i) => `${i + 1}. ${t.title} (${t.estimatedHours}h)`).join('\n')}

Would you like me to show you the full details for review?`,
          suggestions: suggestion
        };
        setMessages(prev => [...prev, assistantMessage]);
        setAiStep('review');
      } else {
        const followUp: AIMessage = {
          role: 'assistant',
          content: "Thanks for the details! To help me create the best project structure, could you tell me:\n\n• What's your target timeline?\n• Are there any specific skills needed?\n• Any key milestones or deadlines?\n\nOr just say 'generate' and I'll create a draft based on what you've shared!"
        };
        setMessages(prev => [...prev, followUp]);
      }
      setIsGenerating(false);
    }, 1500);
  };

  const generateAIProjectSuggestion = (context: string): ProjectSuggestion => {
    const lowerContext = context.toLowerCase();
    
    let projectType = 'general';
    if (lowerContext.includes('mobile') || lowerContext.includes('app')) projectType = 'mobile';
    else if (lowerContext.includes('website') || lowerContext.includes('web')) projectType = 'website';
    else if (lowerContext.includes('ecommerce') || lowerContext.includes('shop')) projectType = 'ecommerce';
    else if (lowerContext.includes('dashboard') || lowerContext.includes('analytics')) projectType = 'dashboard';
    
    const suggestions: Record<string, ProjectSuggestion> = {
      mobile: {
        project: {
          name: 'Mobile Application Development',
          description: 'Cross-platform mobile application with user authentication, core features, and app store deployment.',
          priority: 80,
          targetEndDate: new Date(Date.now() + 120 * 24 * 60 * 60 * 1000),
        },
        tasks: [
          { title: 'Requirements gathering & UX research', estimatedHours: 24, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
          { title: 'UI/UX Design - Core screens', estimatedHours: 40, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
          { title: 'UI/UX Design - User flows & edge cases', estimatedHours: 24, priority: 'medium', status: 'backlog', hierarchyLevel: 2 },
          { title: 'Backend API development', estimatedHours: 56, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
          { title: 'Database schema design', estimatedHours: 16, priority: 'high', status: 'backlog', hierarchyLevel: 2 },
          { title: 'API endpoint implementation', estimatedHours: 40, priority: 'high', status: 'backlog', hierarchyLevel: 2 },
          { title: 'Mobile app development', estimatedHours: 80, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
          { title: 'Authentication integration', estimatedHours: 16, priority: 'medium', status: 'backlog', hierarchyLevel: 2 },
          { title: 'Core feature implementation', estimatedHours: 48, priority: 'high', status: 'backlog', hierarchyLevel: 2 },
          { title: 'Testing & QA', estimatedHours: 32, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
          { title: 'App store submission', estimatedHours: 8, priority: 'medium', status: 'backlog', hierarchyLevel: 1 },
        ],
        reasoning: [
          'Based on mobile app requirements, structured in phases: Design → Backend → Frontend → QA',
          'Estimated 280 total hours (approx. 7 weeks with 1 developer)',
          'Critical path: Design → Backend → Mobile dev → QA',
          'Recommended team: 1 Designer, 1 Backend, 1 Mobile dev'
        ]
      },
      website: {
        project: {
          name: 'Marketing Website Redesign',
          description: 'Modern responsive website with improved SEO, performance optimization, and conversion-focused design.',
          priority: 70,
          targetEndDate: new Date(Date.now() + 60 * 24 * 60 * 60 * 1000),
        },
        tasks: [
          { title: 'Content audit & strategy', estimatedHours: 16, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
          { title: 'Wireframing & prototyping', estimatedHours: 24, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
          { title: 'Visual design - Homepage', estimatedHours: 20, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
          { title: 'Visual design - Interior pages', estimatedHours: 24, priority: 'medium', status: 'backlog', hierarchyLevel: 1 },
          { title: 'Frontend development', estimatedHours: 48, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
          { title: 'Component library setup', estimatedHours: 12, priority: 'medium', status: 'backlog', hierarchyLevel: 2 },
          { title: 'Page implementations', estimatedHours: 36, priority: 'high', status: 'backlog', hierarchyLevel: 2 },
          { title: 'SEO optimization', estimatedHours: 12, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
          { title: 'Performance optimization', estimatedHours: 8, priority: 'medium', status: 'backlog', hierarchyLevel: 1 },
          { title: 'Testing & launch', estimatedHours: 12, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
        ],
        reasoning: [
          'Structured for marketing website lifecycle: Strategy → Design → Dev → Launch',
          'SEO and performance prioritized for marketing impact',
          'Estimated 164 hours (approx. 4 weeks)',
          'Can be executed by 1 designer + 1 developer'
        ]
      },
      ecommerce: {
        project: {
          name: 'E-Commerce Platform',
          description: 'Full-featured online store with product catalog, cart, checkout, and payment integration.',
          priority: 90,
          targetEndDate: new Date(Date.now() + 100 * 24 * 60 * 60 * 1000),
        },
        tasks: [
          { title: 'Requirements & architecture planning', estimatedHours: 20, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
          { title: 'Design system & component library', estimatedHours: 32, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
          { title: 'Product catalog design', estimatedHours: 24, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
          { title: 'Backend development', estimatedHours: 72, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
          { title: 'Database design', estimatedHours: 16, priority: 'high', status: 'backlog', hierarchyLevel: 2 },
          { title: 'Product & inventory APIs', estimatedHours: 24, priority: 'high', status: 'backlog', hierarchyLevel: 2 },
          { title: 'Cart & checkout APIs', estimatedHours: 24, priority: 'critical', status: 'backlog', hierarchyLevel: 2 },
          { title: 'Payment integration', estimatedHours: 8, priority: 'critical', status: 'backlog', hierarchyLevel: 2 },
          { title: 'Frontend development', estimatedHours: 64, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
          { title: 'Product catalog UI', estimatedHours: 24, priority: 'high', status: 'backlog', hierarchyLevel: 2 },
          { title: 'Shopping cart & checkout UI', estimatedHours: 24, priority: 'critical', status: 'backlog', hierarchyLevel: 2 },
          { title: 'Admin dashboard', estimatedHours: 32, priority: 'medium', status: 'backlog', hierarchyLevel: 1 },
          { title: 'Testing & security audit', estimatedHours: 32, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
        ],
        reasoning: [
          'E-commerce complexity requires robust backend-first approach',
          'Critical path through checkout flow (payment is highest risk)',
          'Estimated 300+ hours (8-10 weeks with full team)',
          'Recommended: 1 PM, 1 Designer, 2 Backend, 2 Frontend'
        ]
      },
      dashboard: {
        project: {
          name: 'Analytics Dashboard',
          description: 'Real-time data visualization dashboard with custom widgets, reporting, and user management.',
          priority: 75,
          targetEndDate: new Date(Date.now() + 80 * 24 * 60 * 60 * 1000),
        },
        tasks: [
          { title: 'Data source integration planning', estimatedHours: 16, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
          { title: 'Data pipeline setup', estimatedHours: 32, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
          { title: 'ETL process development', estimatedHours: 24, priority: 'high', status: 'backlog', hierarchyLevel: 2 },
          { title: 'Data warehouse configuration', estimatedHours: 8, priority: 'medium', status: 'backlog', hierarchyLevel: 2 },
          { title: 'Backend API development', estimatedHours: 40, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
          { title: 'Dashboard UI design', estimatedHours: 32, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
          { title: 'Widget library development', estimatedHours: 40, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
          { title: 'Chart components', estimatedHours: 16, priority: 'high', status: 'backlog', hierarchyLevel: 2 },
          { title: 'Data table components', estimatedHours: 12, priority: 'medium', status: 'backlog', hierarchyLevel: 2 },
          { title: 'Report builder feature', estimatedHours: 32, priority: 'medium', status: 'backlog', hierarchyLevel: 1 },
          { title: 'User management & permissions', estimatedHours: 16, priority: 'medium', status: 'backlog', hierarchyLevel: 1 },
          { title: 'Testing & optimization', estimatedHours: 24, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
        ],
        reasoning: [
          'Data-heavy project requires pipeline setup first',
          'Widget reusability prioritized for scalability',
          'Estimated 280 hours (7 weeks with data engineer + frontend)',
          'Critical: Data pipeline must be stable before UI development'
        ]
      },
      general: {
        project: {
          name: 'New Project',
          description: 'Custom software project with planning, development, and deployment phases.',
          priority: 60,
          targetEndDate: new Date(Date.now() + 90 * 24 * 60 * 60 * 1000),
        },
        tasks: [
          { title: 'Discovery & requirements', estimatedHours: 20, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
          { title: 'Technical architecture', estimatedHours: 16, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
          { title: 'UI/UX Design', estimatedHours: 40, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
          { title: 'Backend development', estimatedHours: 60, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
          { title: 'Frontend development', estimatedHours: 56, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
          { title: 'Integration & testing', estimatedHours: 32, priority: 'high', status: 'backlog', hierarchyLevel: 1 },
          { title: 'Deployment & launch', estimatedHours: 16, priority: 'medium', status: 'backlog', hierarchyLevel: 1 },
        ],
        reasoning: [
          'Standard software development lifecycle structure',
          'Estimated 240 hours (6 weeks with standard team)',
          'Flexible for iteration based on specific requirements'
        ]
      }
    };
    
    return suggestions[projectType] || suggestions.general;
  };

  const handleConfirmAIProject = () => {
    if (!editingProject) return;
    
    const newProject: Project = {
      id: `proj-${Date.now()}`,
      name: editingProject.name || 'New Project',
      description: editingProject.description || '',
      status: 'active',
      priority: editingProject.priority || 50,
      startDate: new Date(),
      targetEndDate: editingProject.targetEndDate || new Date(Date.now() + 90 * 24 * 60 * 60 * 1000),
      healthScore: 100,
      progress: 0,
    };
    
    const newTasks: Task[] = editingTasks.map((t, i) => ({
      id: `task-${Date.now()}-${i}`,
      projectId: newProject.id,
      parentTaskId: t.parentTaskId,
      hierarchyLevel: t.hierarchyLevel || 1,
      title: t.title || 'Task',
      description: t.description || '',
      status: 'backlog',
      priority: t.priority || 'medium',
      priorityScore: 50,
      businessValue: 50,
      estimatedHours: t.estimatedHours || 8,
      requiredSkills: t.requiredSkills || [],
      isMilestone: t.isMilestone || false,
      isCriticalPath: t.isCriticalPath || false,
      dependencies: [],
      blockedBy: [],
    }));
    
    onProjectCreate(newProject, newTasks);
    handleClose();
  };

  const updateTaskTitle = (index: number, title: string) => {
    setEditingTasks(prev => prev.map((t, i) => i === index ? { ...t, title } : t));
  };

  const updateTaskHours = (index: number, hours: number) => {
    setEditingTasks(prev => prev.map((t, i) => i === index ? { ...t, estimatedHours: hours } : t));
  };

  const removeTask = (index: number) => {
    setEditingTasks(prev => prev.filter((_, i) => i !== index));
  };

  const addTask = () => {
    setEditingTasks(prev => [...prev, {
      title: 'New Task',
      estimatedHours: 8,
      priority: 'medium',
      status: 'backlog',
      hierarchyLevel: 1
    }]);
  };

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Plus className="w-5 h-5" />
            Create New Project
          </DialogTitle>
          <DialogDescription>
            Create a project manually or let AI help you structure it with tasks
          </DialogDescription>
        </DialogHeader>

        <Tabs value={activeTab} onValueChange={setActiveTab} className="mt-4">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="manual">
              <User className="w-4 h-4 mr-2" />
              Manual
            </TabsTrigger>
            <TabsTrigger value="ai">
              <Sparkles className="w-4 h-4 mr-2" />
              AI Assisted
            </TabsTrigger>
          </TabsList>

          <TabsContent value="manual" className="space-y-4 mt-4">
            <div className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="name">Project Name *</Label>
                <Input 
                  id="name" 
                  value={projectName}
                  onChange={(e) => setProjectName(e.target.value)}
                  placeholder="e.g., Website Redesign"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="desc">Description</Label>
                <Textarea 
                  id="desc" 
                  value={projectDescription}
                  onChange={(e) => setProjectDescription(e.target.value)}
                  placeholder="Describe the project goals and scope..."
                  rows={3}
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="priority">Priority (1-100)</Label>
                  <Input 
                    id="priority" 
                    type="number"
                    min={1}
                    max={100}
                    value={projectPriority}
                    onChange={(e) => setProjectPriority(parseInt(e.target.value))}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="date">Target End Date</Label>
                  <Input 
                    id="date" 
                    type="date"
                    value={targetDate}
                    onChange={(e) => setTargetDate(e.target.value)}
                  />
                </div>
              </div>
            </div>

            <DialogFooter>
              <Button variant="outline" onClick={handleClose} disabled={isCreating}>Cancel</Button>
              <Button onClick={handleManualSubmit} disabled={!projectName || isCreating}>
                {isCreating ? (
                  <>
                    <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                    Creating...
                  </>
                ) : (
                  'Create Project'
                )}
              </Button>
            </DialogFooter>
          </TabsContent>

          <TabsContent value="ai" className="mt-4">
            {aiStep === 'chat' && (
              <div className="space-y-4">
                <Card className="bg-muted/50">
                  <CardContent className="p-4 space-y-4 max-h-[400px] overflow-y-auto">
                    {messages.map((msg, i) => (
                      <div key={i} className={`flex gap-3 ${msg.role === 'user' ? 'flex-row-reverse' : ''}`}>
                        <div className={`w-8 h-8 rounded-full flex items-center justify-center ${
                          msg.role === 'user' ? 'bg-primary' : 'bg-gradient-to-br from-purple-500 to-pink-500'
                        }`}>
                          {msg.role === 'user' ? (
                            <User className="w-4 h-4 text-primary-foreground" />
                          ) : (
                            <Bot className="w-4 h-4 text-white" />
                          )}
                        </div>
                        <div className={`max-w-[80%] p-3 rounded-lg text-sm ${
                          msg.role === 'user' ? 'bg-primary text-primary-foreground' : 'bg-background border'
                        }`}>
                          {msg.content.split('\n').map((line, j) => (
                            <p key={j} className={line.startsWith('**') ? 'font-semibold mt-2' : ''}>
                              {line}
                            </p>
                          ))}
                        </div>
                      </div>
                    ))}
                    {isGenerating && (
                      <div className="flex gap-3">
                        <div className="w-8 h-8 rounded-full bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center">
                          <Bot className="w-4 h-4 text-white" />
                        </div>
                        <div className="p-3 bg-background border rounded-lg">
                          <RefreshCw className="w-4 h-4 animate-spin" />
                        </div>
                      </div>
                    )}
                  </CardContent>
                </Card>

                <div className="flex gap-2">
                  <Input 
                    value={inputMessage}
                    onChange={(e) => setInputMessage(e.target.value)}
                    placeholder="Describe your project or type 'generate' to create..."
                    onKeyDown={(e) => e.key === 'Enter' && sendMessage()}
                  />
                  <Button onClick={sendMessage} disabled={isGenerating || !inputMessage.trim()}>
                    <Send className="w-4 h-4" />
                  </Button>
                </div>

                {aiSuggestion && (
                  <Button 
                    className="w-full" 
                    onClick={() => setAiStep('editing')}
                  >
                    Review & Edit Project Structure
                    <ChevronRight className="w-4 h-4 ml-2" />
                  </Button>
                )}
              </div>
            )}

            {aiStep === 'editing' && editingProject && (
              <div className="space-y-6">
                <div className="space-y-4">
                  <h3 className="font-semibold flex items-center gap-2">
                    <Sparkles className="w-4 h-4 text-primary" />
                    Project Details
                  </h3>
                  
                  <div className="space-y-2">
                    <Label>Project Name</Label>
                    <Input 
                      value={editingProject.name}
                      onChange={(e) => setEditingProject({ ...editingProject, name: e.target.value })}
                    />
                  </div>

                  <div className="space-y-2">
                    <Label>Description</Label>
                    <Textarea 
                      value={editingProject.description}
                      onChange={(e) => setEditingProject({ ...editingProject, description: e.target.value })}
                      rows={2}
                    />
                  </div>

                  <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label>Priority</Label>
                      <Input 
                        type="number"
                        min={1}
                        max={100}
                        value={editingProject.priority}
                        onChange={(e) => setEditingProject({ ...editingProject, priority: parseInt(e.target.value) })}
                      />
                    </div>
                    <div className="space-y-2">
                      <Label>Target Date</Label>
                      <Input 
                        type="date"
                        value={editingProject.targetEndDate?.toISOString().split('T')[0]}
                        onChange={(e) => setEditingProject({ ...editingProject, targetEndDate: new Date(e.target.value) })}
                      />
                    </div>
                  </div>
                </div>

                <Separator />

                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <h3 className="font-semibold flex items-center gap-2">
                      <Check className="w-4 h-4 text-primary" />
                      Tasks ({editingTasks.length})
                    </h3>
                    <Button size="sm" variant="outline" onClick={addTask}>
                      <Plus className="w-4 h-4 mr-1" />
                      Add Task
                    </Button>
                  </div>

                  <div className="space-y-2 max-h-[300px] overflow-y-auto">
                    {editingTasks.map((task, index) => (
                      <Card key={index} className={task.hierarchyLevel === 2 ? 'ml-8 border-l-4 border-l-primary' : ''}>
                        <CardContent className="p-3 space-y-2">
                          <div className="flex items-center gap-2">
                            {task.hierarchyLevel === 2 && (
                              <Badge variant="outline" className="text-xs">Subtask</Badge>
                            )}
                            <Input 
                              value={task.title}
                              onChange={(e) => updateTaskTitle(index, e.target.value)}
                              className="flex-1"
                              placeholder="Task title"
                            />
                            <Input 
                              type="number"
                              min={1}
                              max={200}
                              value={task.estimatedHours}
                              onChange={(e) => updateTaskHours(index, parseInt(e.target.value) || 0)}
                              className="w-20"
                            />
                            <span className="text-sm text-muted-foreground">h</span>
                            <Button 
                              size="icon" 
                              variant="ghost" 
                              className="h-8 w-8 text-destructive"
                              onClick={() => removeTask(index)}
                            >
                              ×
                            </Button>
                          </div>
                        </CardContent>
                      </Card>
                    ))}
                  </div>
                </div>

                {aiSuggestion && (
                  <Card className="bg-muted/50">
                    <CardContent className="p-4">
                      <h4 className="font-medium mb-2 flex items-center gap-2">
                        <Bot className="w-4 h-4 text-primary" />
                        AI Reasoning
                      </h4>
                      <ul className="space-y-1 text-sm text-muted-foreground">
                        {aiSuggestion.reasoning.map((reason, i) => (
                          <li key={i} className="flex items-start gap-2">
                            <span className="text-primary">•</span>
                            {reason}
                          </li>
                        ))}
                      </ul>
                    </CardContent>
                  </Card>
                )}

                <DialogFooter>
                  <Button variant="outline" onClick={() => setAiStep('chat')} disabled={isCreating}>
                    Back to Chat
                  </Button>
                  <Button variant="outline" onClick={handleClose} disabled={isCreating}>
                    Cancel
                  </Button>
                  <Button onClick={handleConfirmAIProject} disabled={isCreating}>
                    {isCreating ? (
                      <>
                        <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                        Creating...
                      </>
                    ) : (
                      <>
                        <Check className="w-4 h-4 mr-2" />
                        Create Project & Tasks
                      </>
                    )}
                  </Button>
                </DialogFooter>
              </div>
            )}
          </TabsContent>
        </Tabs>
      </DialogContent>
    </Dialog>
  );
}
