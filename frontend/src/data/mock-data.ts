import { 
  createRealisticTeam, 
  createRealisticProjects, 
  createRealisticTasks, 
  createSkillsCatalog, 
  createRealisticNudges, 
  createWorkloadData,
  createRealisticScenarios,
  teamSkills
} from '@/lib/faker/generators';
import { User, Project, Task, Skill, Nudge, WorkloadEntry, Scenario } from '@/types';

// Generate mock data store
class MockDataStore {
  users: User[] = [];
  projects: Project[] = [];
  tasks: Task[] = [];
  skills: Skill[] = [];
  nudges: Nudge[] = [];
  workloadData: WorkloadEntry[] = [];
  scenarios: Scenario[] = [];
  teamSkills: Record<string, { skillId: string; proficiency: 1 | 2 | 3 | 4 }[]> = {};

  constructor() {
    this.initialize();
  }

  private initialize() {
    this.users = createRealisticTeam();
    this.projects = createRealisticProjects();
    this.skills = createSkillsCatalog();
    this.tasks = createRealisticTasks();
    this.nudges = createRealisticNudges();
    this.workloadData = createWorkloadData(this.users);
    this.scenarios = createRealisticScenarios();
    this.teamSkills = teamSkills;
  }

  // Mutation methods
  addProject(project: Project) {
    this.projects.push(project);
  }

  addTask(task: Task) {
    this.tasks.push(task);
  }

  updateTask(taskId: string, updates: Partial<Task>) {
    const index = this.tasks.findIndex(t => t.id === taskId);
    if (index !== -1) {
      this.tasks[index] = { ...this.tasks[index], ...updates };
    }
  }

  updateProject(projectId: string, updates: Partial<Project>) {
    const index = this.projects.findIndex(p => p.id === projectId);
    if (index !== -1) {
      this.projects[index] = { ...this.projects[index], ...updates };
    }
  }

  // User operations
  getUsers(): User[] {
    return this.users;
  }

  getUserById(id: string): User | undefined {
    return this.users.find(u => u.id === id);
  }

  getUsersByRole(role: User['role']): User[] {
    return this.users.filter(u => u.role === role);
  }

  // Project operations
  getProjects(): Project[] {
    return this.projects;
  }

  getProjectById(id: string): Project | undefined {
    return this.projects.find(p => p.id === id);
  }

  getActiveProjects(): Project[] {
    return this.projects.filter(p => p.status === 'active');
  }

  getAtRiskProjects(): Project[] {
    return this.projects.filter(p => p.healthScore < 50);
  }

  // Task operations
  getTasksByProject(projectId: string): Task[] {
    return this.tasks.filter(t => t.projectId === projectId);
  }

  getAllTasks(): Task[] {
    return this.tasks;
  }

  getTaskById(taskId: string): Task | undefined {
    return this.tasks.find(t => t.id === taskId);
  }

  getUnassignedTasks(): Task[] {
    return this.tasks.filter(t => !t.assigneeId && t.status !== 'done');
  }

  getCriticalPathTasks(projectId?: string): Task[] {
    let tasks = this.tasks.filter(t => t.isCriticalPath);
    if (projectId) {
      tasks = tasks.filter(t => t.projectId === projectId);
    }
    return tasks;
  }

  getTasksByStatus(status: Task['status']): Task[] {
    return this.tasks.filter(t => t.status === status);
  }

  getTasksByAssignee(personId: string): Task[] {
    return this.tasks.filter(t => t.assigneeId === personId);
  }

  // Skills operations
  getSkills(): Skill[] {
    return this.skills;
  }

  getSkillById(id: string): Skill | undefined {
    return this.skills.find(s => s.id === id);
  }

  getUserSkills(userId: string): { skill: Skill; proficiency: number }[] {
    const userSkillIds = this.teamSkills[userId] || [];
    return userSkillIds.map(us => ({
      skill: this.getSkillById(us.skillId)!,
      proficiency: us.proficiency,
    })).filter(us => us.skill);
  }

  getSkillsGap(requiredSkills: string[], userId: string): Skill[] {
    const userSkills = this.getUserSkills(userId).map(us => us.skill.id);
    return this.skills.filter(s => requiredSkills.includes(s.id) && !userSkills.includes(s.id));
  }

  // Nudges operations
  getNudges(): Nudge[] {
    return this.nudges;
  }

  getNudgeById(id: string): Nudge | undefined {
    return this.nudges.find(n => n.id === id);
  }

  getNudgesByType(type: Nudge['type']): Nudge[] {
    return this.nudges.filter(n => n.type === type);
  }

  getUnreadNudges(): Nudge[] {
    return this.nudges.filter(n => n.status === 'unread');
  }

  getHighSeverityNudges(): Nudge[] {
    return this.nudges.filter(n => n.severity === 'high');
  }

  // Workload operations
  getWorkloadData(): WorkloadEntry[] {
    return this.workloadData;
  }

  getWorkloadByPerson(personId: string): WorkloadEntry | undefined {
    return this.workloadData.find(w => w.personId === personId);
  }

  getOverallocatedMembers(): WorkloadEntry[] {
    return this.workloadData.filter(w => w.allocationPercentage > 100);
  }

  getUnderutilizedMembers(): WorkloadEntry[] {
    return this.workloadData.filter(w => w.allocationPercentage < 70);
  }

  // Scenario operations
  getScenarios(): Scenario[] {
    return this.scenarios;
  }

  getScenarioById(id: string): Scenario | undefined {
    return this.scenarios.find(s => s.id === id);
  }

  getPendingScenarios(): Scenario[] {
    return this.scenarios.filter(s => s.status === 'pending');
  }

  // Portfolio metrics
  getPortfolioHealthScore(): number {
    const activeProjects = this.getActiveProjects();
    if (activeProjects.length === 0) return 100;
    const totalHealth = activeProjects.reduce((sum, p) => sum + p.healthScore, 0);
    return Math.round(totalHealth / activeProjects.length);
  }
}

// Export singleton instance
// export const mockData = new MockDataStore();
