'use client';

import { useState } from 'react';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Card, CardContent } from '@/components/ui/card';
import { useProjectsList, useSkills, useUsers } from '@/hooks/api';
import { useCreateTask } from '@/hooks/api/useTasks';
import { TaskPriority, Project, Skill, User as UserType } from '@/types';
import { 
  Sparkles, 
  Plus, 
  X, 
  Send, 
  User, 
  Bot, 
  RefreshCw,
  Check,
  ChevronRight,
  GripVertical,
  Clock,
  Loader2
} from 'lucide-react';

interface CreateTaskDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  defaultProjectId?: string;
  parentTaskId?: string;
}

interface AIMessage {
  role: 'user' | 'assistant';
  content: string;
}

interface AISubtask {
  title: string;
  estimatedHours: number;
  description?: string;
}

export function CreateTaskDialog({ open, onOpenChange, defaultProjectId, parentTaskId }: CreateTaskDialogProps) {
  const { data: projectsData } = useProjectsList();
  const { data: skillsData } = useSkills();
  const { data: usersData } = useUsers();
  const createTaskMutation = useCreateTask();
  
  const projects: Project[] = projectsData?.data?.projects || [];
  const skills: Skill[] = skillsData?.data?.skills || [];
  const users: UserType[] = usersData?.data?.users || [];
  
  const [activeTab, setActiveTab] = useState('manual');
  const [aiStep, setAiStep] = useState<'chat' | 'review'>('chat');
  
  // Manual form state
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [projectId, setProjectId] = useState(defaultProjectId || projects[0]?.id || '');
  const [priority, setPriority] = useState<TaskPriority>('medium');
  const [estimatedHours, setEstimatedHours] = useState(8);
  const [selectedSkills, setSelectedSkills] = useState<string[]>([]);
  const [assigneeId, setAssigneeId] = useState<string>('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  
  // AI state
  const [messages, setMessages] = useState<AIMessage[]>([
    {
      role: 'assistant',
      content: "Hi! I can help you break down this task into subtasks. Describe the task and I'll suggest a structured approach with estimated hours for each step."
    }
  ]);
  const [inputMessage, setInputMessage] = useState('');
  const [isGenerating, setIsGenerating] = useState(false);
  const [aiSubtasks, setAiSubtasks] = useState<AISubtask[]>([]);
  const [taskTitle, setTaskTitle] = useState('');

  const resetForm = () => {
    setTitle('');
    setDescription('');
    setProjectId(defaultProjectId || projects[0]?.id || '');
    setPriority('medium');
    setEstimatedHours(8);
    setSelectedSkills([]);
    setAssigneeId('');
    setMessages([{
      role: 'assistant',
      content: "Hi! I can help you break down this task into subtasks. Describe the task and I'll suggest a structured approach with estimated hours for each step."
    }]);
    setAiStep('chat');
    setAiSubtasks([]);
    setInputMessage('');
    setTaskTitle('');
    setIsSubmitting(false);
  };

  const handleClose = () => {
    resetForm();
    onOpenChange(false);
  };

  const handleManualSubmit = async () => {
    setIsSubmitting(true);
    try {
      await createTaskMutation.mutateAsync({
        title,
        description,
        projectId,
        priority,
        estimatedHours,
        requiredSkills: selectedSkills,
        assigneeId: assigneeId || undefined,
        parentTaskId,
        hierarchyLevel: parentTaskId ? 2 : 1,
      });
      handleClose();
    } catch (error) {
      console.error('Failed to create task:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const sendMessage = () => {
    if (!inputMessage.trim()) return;
    
    const userMessage: AIMessage = { role: 'user', content: inputMessage };
    setMessages(prev => [...prev, userMessage]);
    setInputMessage('');
    setIsGenerating(true);
    
    setTimeout(() => {
      const lowerInput = inputMessage.toLowerCase();
      
      // Extract task title from first message or use default
      if (!taskTitle && inputMessage.length > 10) {
        setTaskTitle(inputMessage.split('.')[0].substring(0, 50));
      }
      
      if (lowerInput.includes('generate') || lowerInput.includes('break') || lowerInput.includes('create') || messages.length > 2) {
        const subtasks = generateAISubtasks(inputMessage);
        setAiSubtasks(subtasks);
        
        const totalHours = subtasks.reduce((sum, s) => sum + s.estimatedHours, 0);
        
        const assistantMessage: AIMessage = {
          role: 'assistant',
          content: `I've broken this down into ${subtasks.length} subtasks with an estimated total of ${totalHours} hours:\n\n${subtasks.map((s, i) => `${i + 1}. ${s.title} (${s.estimatedHours}h)`).join('\n')}\n\nWould you like to review and customize these subtasks?`
        };
        setMessages(prev => [...prev, assistantMessage]);
        setAiStep('review');
      } else {
        const followUp: AIMessage = {
          role: 'assistant',
          content: "I understand. To create the best breakdown, could you tell me:\n\n• What's the main objective?\n• Are there any specific steps you already have in mind?\n• Any dependencies or blockers to consider?\n\nOr just say 'break it down' and I'll generate subtasks based on what you've shared!"
        };
        setMessages(prev => [...prev, followUp]);
      }
      setIsGenerating(false);
    }, 1500);
  };

  const generateAISubtasks = (context: string): AISubtask[] => {
    const lowerContext = context.toLowerCase();
    
    if (lowerContext.includes('design') || lowerContext.includes('ui') || lowerContext.includes('figma')) {
      return [
        { title: 'Research & reference gathering', estimatedHours: 4, description: 'Collect inspiration and analyze competitors' },
        { title: 'User flow mapping', estimatedHours: 6, description: 'Map out user journeys and interactions' },
        { title: 'Wireframe creation', estimatedHours: 8, description: 'Low-fidelity wireframes for all screens' },
        { title: 'Visual design - Core screens', estimatedHours: 16, description: 'High-fidelity designs for main flows' },
        { title: 'Visual design - Edge cases', estimatedHours: 8, description: 'Empty states, error states, loading screens' },
        { title: 'Design review & handoff', estimatedHours: 4, description: 'Review with stakeholders and prepare specs' },
      ];
    }
    
    if (lowerContext.includes('api') || lowerContext.includes('backend') || lowerContext.includes('server')) {
      return [
        { title: 'API design & documentation', estimatedHours: 6, description: 'Design endpoints and document with OpenAPI' },
        { title: 'Database schema design', estimatedHours: 8, description: 'Design tables, relationships, and indexes' },
        { title: 'Authentication & authorization', estimatedHours: 8, description: 'Implement auth middleware and permissions' },
        { title: 'Core endpoint implementation', estimatedHours: 16, description: 'Build main CRUD endpoints' },
        { title: 'Business logic implementation', estimatedHours: 16, description: 'Implement domain-specific logic' },
        { title: 'Testing & API validation', estimatedHours: 8, description: 'Unit tests and integration tests' },
      ];
    }
    
    if (lowerContext.includes('frontend') || lowerContext.includes('react') || lowerContext.includes('component')) {
      return [
        { title: 'Component structure planning', estimatedHours: 4, description: 'Plan component hierarchy and props' },
        { title: 'Core component development', estimatedHours: 16, description: 'Build reusable UI components' },
        { title: 'State management setup', estimatedHours: 6, description: 'Configure stores and data flow' },
        { title: 'API integration', estimatedHours: 8, description: 'Connect frontend to backend APIs' },
        { title: 'Form validation & error handling', estimatedHours: 6, description: 'Implement validation logic' },
        { title: 'Responsive styling', estimatedHours: 8, description: 'Mobile and tablet adaptations' },
      ];
    }
    
    // Default generic breakdown
    return [
      { title: 'Requirements clarification', estimatedHours: 2, description: 'Clarify scope and acceptance criteria' },
      { title: 'Research & planning', estimatedHours: 4, description: 'Research solutions and create plan' },
      { title: 'Implementation - Phase 1', estimatedHours: 8, description: 'Core functionality development' },
      { title: 'Implementation - Phase 2', estimatedHours: 8, description: 'Additional features and edge cases' },
      { title: 'Testing & QA', estimatedHours: 4, description: 'Manual testing and bug fixes' },
      { title: 'Review & documentation', estimatedHours: 2, description: 'Code review and documentation' },
    ];
  };

  const handleCreateWithSubtasks = async () => {
    setIsSubmitting(true);
    try {
      // Create parent task
      const parentResponse = await createTaskMutation.mutateAsync({
        title: taskTitle || 'New Task',
        description: description || '',
        projectId,
        priority,
        estimatedHours: aiSubtasks.reduce((sum, s) => sum + s.estimatedHours, 0),
        requiredSkills: selectedSkills,
        assigneeId: assigneeId || undefined,
        hierarchyLevel: 1,
      });
      
      const parentTaskId = parentResponse.data?.id;
      
      // Create subtasks
      if (parentTaskId) {
        for (const subtask of aiSubtasks) {
          await createTaskMutation.mutateAsync({
            title: subtask.title,
            description: subtask.description || '',
            projectId,
            priority,
            estimatedHours: subtask.estimatedHours,
            parentTaskId,
            hierarchyLevel: 2,
            requiredSkills: selectedSkills,
          });
        }
      }
      
      handleClose();
    } catch (error) {
      console.error('Failed to create task with subtasks:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const updateSubtaskTitle = (index: number, title: string) => {
    setAiSubtasks(prev => prev.map((s, i) => i === index ? { ...s, title } : s));
  };

  const updateSubtaskHours = (index: number, hours: number) => {
    setAiSubtasks(prev => prev.map((s, i) => i === index ? { ...s, estimatedHours: hours } : s));
  };

  const removeSubtask = (index: number) => {
    setAiSubtasks(prev => prev.filter((_, i) => i !== index));
  };

  const addSubtask = () => {
    setAiSubtasks(prev => [...prev, { title: 'New Subtask', estimatedHours: 4 }]);
  };

  const toggleSkill = (skillId: string) => {
    setSelectedSkills(prev => 
      prev.includes(skillId) 
        ? prev.filter(id => id !== skillId)
        : [...prev, skillId]
    );
  };

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Plus className="w-5 h-5" />
            {parentTaskId ? 'Add Subtask' : 'Create New Task'}
          </DialogTitle>
          <DialogDescription>
            {parentTaskId 
              ? 'Add a subtask to break down the work further' 
              : 'Create a task manually or let AI help break it down into subtasks'}
          </DialogDescription>
        </DialogHeader>

        <Tabs value={activeTab} onValueChange={setActiveTab} className="mt-4">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="manual">
              <User className="w-4 h-4 mr-2" />
              Manual
            </TabsTrigger>
            <TabsTrigger value="ai" disabled={!!parentTaskId}>
              <Sparkles className="w-4 h-4 mr-2" />
              AI Breakdown
            </TabsTrigger>
          </TabsList>

          <TabsContent value="manual" className="space-y-4 mt-4">
            <div className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="title">Task Title *</Label>
                <Input 
                  id="title" 
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  placeholder="e.g., Implement user authentication"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="description">Description</Label>
                <Textarea 
                  id="description" 
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  placeholder="Describe the task requirements..."
                  rows={3}
                />
              </div>

              <div className="grid grid-cols-3 gap-4">
                <div className="space-y-2">
                  <Label>Project</Label>
                  <Select value={projectId} onValueChange={setProjectId}>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      {projects.map((project: Project) => (
                        <SelectItem key={project.id} value={project.id}>
                          {project.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                <div className="space-y-2">
                  <Label>Priority</Label>
                  <Select value={priority} onValueChange={(v: TaskPriority) => setPriority(v)}>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="low">Low</SelectItem>
                      <SelectItem value="medium">Medium</SelectItem>
                      <SelectItem value="high">High</SelectItem>
                      <SelectItem value="critical">Critical</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div className="space-y-2">
                  <Label>Estimated Hours</Label>
                  <Input 
                    type="number"
                    min={1}
                    max={200}
                    value={estimatedHours}
                    onChange={(e) => setEstimatedHours(parseInt(e.target.value) || 0)}
                  />
                </div>
              </div>

              <div className="space-y-2">
                <Label>Assignee</Label>
                <Select value={assigneeId} onValueChange={setAssigneeId}>
                  <SelectTrigger>
                    <SelectValue placeholder="Unassigned" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="">Unassigned</SelectItem>
                    {users.filter((u: UserType) => u.role === 'member').map((user: UserType) => (
                      <SelectItem key={user.id} value={user.id}>
                        {user.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label>Required Skills</Label>
                <div className="flex flex-wrap gap-2 p-3 border rounded-lg">
                  {skills.map((skill: Skill) => (
                    <Badge 
                      key={skill.id}
                      variant={selectedSkills.includes(skill.id) ? 'default' : 'outline'}
                      className="cursor-pointer"
                      onClick={() => toggleSkill(skill.id)}
                    >
                      {skill.name}
                      {selectedSkills.includes(skill.id) && (
                        <X className="w-3 h-3 ml-1" />
                      )}
                    </Badge>
                  ))}
                </div>
              </div>
            </div>

            <DialogFooter>
              <Button variant="outline" onClick={handleClose} disabled={isSubmitting}>Cancel</Button>
              <Button onClick={handleManualSubmit} disabled={!title || isSubmitting}>
                {isSubmitting ? (
                  <>
                    <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                    Creating...
                  </>
                ) : (
                  <>{parentTaskId ? 'Create Subtask' : 'Create Task'}</>
                )}
              </Button>
            </DialogFooter>
          </TabsContent>

          <TabsContent value="ai" className="mt-4">
            {aiStep === 'chat' && (
              <div className="space-y-4">
                <Card className="bg-muted/50">
                  <CardContent className="p-4 space-y-4 max-h-[350px] overflow-y-auto">
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
                            <p key={j} className={line.startsWith('**') || line.match(/^\d+\./) ? 'font-semibold mt-1' : ''}>
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
                    placeholder="Describe the task you want to break down..."
                    onKeyDown={(e) => e.key === 'Enter' && sendMessage()}
                    disabled={isGenerating}
                  />
                  <Button onClick={sendMessage} disabled={isGenerating || !inputMessage.trim()}>
                    <Send className="w-4 h-4" />
                  </Button>
                </div>

                {aiSubtasks.length > 0 && (
                  <Button 
                    className="w-full" 
                    onClick={() => setAiStep('review')}
                  >
                    Review Subtasks
                    <ChevronRight className="w-4 h-4 ml-2" />
                  </Button>
                )}
              </div>
            )}

            {aiStep === 'review' && aiSubtasks.length > 0 && (
              <div className="space-y-6">
                <div className="space-y-4">
                  <div className="space-y-2">
                    <Label>Task Title</Label>
                    <Input 
                      value={taskTitle}
                      onChange={(e) => setTaskTitle(e.target.value)}
                      placeholder="Main task title"
                    />
                  </div>

                  <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label>Project</Label>
                      <Select value={projectId} onValueChange={setProjectId}>
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          {projects.map((project: Project) => (
                            <SelectItem key={project.id} value={project.id}>
                              {project.name}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                    <div className="space-y-2">
                      <Label>Priority</Label>
                      <Select value={priority} onValueChange={(v: TaskPriority) => setPriority(v)}>
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="low">Low</SelectItem>
                          <SelectItem value="medium">Medium</SelectItem>
                          <SelectItem value="high">High</SelectItem>
                          <SelectItem value="critical">Critical</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                  </div>
                </div>

                <Separator />

                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <h3 className="font-semibold flex items-center gap-2">
                      <Sparkles className="w-4 h-4 text-primary" />
                      AI-Generated Subtasks ({aiSubtasks.length})
                    </h3>
                    <Button size="sm" variant="outline" onClick={addSubtask}>
                      <Plus className="w-4 h-4 mr-1" />
                      Add Subtask
                    </Button>
                  </div>

                  <div className="space-y-2 max-h-[250px] overflow-y-auto">
                    {aiSubtasks.map((subtask: AISubtask, index: number) => (
                      <Card key={index} className="border-l-4 border-l-primary">
                        <CardContent className="p-3">
                          <div className="flex items-center gap-2">
                            <GripVertical className="w-4 h-4 text-muted-foreground" />
                            <Input 
                              value={subtask.title}
                              onChange={(e) => updateSubtaskTitle(index, e.target.value)}
                              className="flex-1"
                              placeholder="Subtask title"
                            />
                            <Input 
                              type="number"
                              min={1}
                              max={100}
                              value={subtask.estimatedHours}
                              onChange={(e) => updateSubtaskHours(index, parseInt(e.target.value) || 0)}
                              className="w-20"
                            />
                            <span className="text-sm text-muted-foreground">h</span>
                            <Button 
                              size="icon" 
                              variant="ghost" 
                              className="h-8 w-8 text-destructive"
                              onClick={() => removeSubtask(index)}
                            >
                              <X className="w-4 h-4" />
                            </Button>
                          </div>
                          {subtask.description && (
                            <p className="text-xs text-muted-foreground mt-1 ml-6">
                              {subtask.description}
                            </p>
                          )}
                        </CardContent>
                      </Card>
                    ))}
                  </div>

                  <Card className="bg-muted/50">
                    <CardContent className="p-4 flex items-center justify-between">
                      <span className="text-sm text-muted-foreground">Total Estimated Hours</span>
                      <span className="text-lg font-semibold">
                        {aiSubtasks.reduce((sum, s) => sum + s.estimatedHours, 0)}h
                      </span>
                    </CardContent>
                  </Card>
                </div>

                <DialogFooter>
                  <Button variant="outline" onClick={() => setAiStep('chat')} disabled={isSubmitting}>
                    Back to Chat
                  </Button>
                  <Button variant="outline" onClick={handleClose} disabled={isSubmitting}>
                    Cancel
                  </Button>
                  <Button onClick={handleCreateWithSubtasks} disabled={!taskTitle || isSubmitting}>
                    {isSubmitting ? (
                      <>
                        <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                        Creating...
                      </>
                    ) : (
                      <>
                        <Check className="w-4 h-4 mr-2" />
                        Create Task & Subtasks
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
