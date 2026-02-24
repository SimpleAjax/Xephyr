'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { useAppStore } from '@/store/app-store';
import { 
  LayoutDashboard, 
  FolderKanban, 
  ListTodo, 
  Users, 
  Bell,
  Settings,
  Sparkles,
  GraduationCap,
  GitBranch
} from 'lucide-react';

const navigation = [
  { name: 'Dashboard', href: '/dashboard', icon: LayoutDashboard },
  { name: 'Projects', href: '/projects', icon: FolderKanban },
  { name: 'Kanban Board', href: '/kanban', icon: ListTodo },
  { name: 'Team & Workload', href: '/team', icon: Users },
  { name: 'Skills', href: '/skills', icon: GraduationCap },
  { name: 'Scenarios', href: '/scenarios', icon: Sparkles },
  { name: 'Nudges', href: '/nudges', icon: Bell, badge: 'unreadNudges' },
  { name: 'Settings', href: '/settings', icon: Settings },
];

export function Sidebar() {
  const pathname = usePathname();
  const unreadNudges = useAppStore((state) => state.unreadNudgesCount);

  return (
    <div className="w-64 border-r bg-card flex flex-col">
      <div className="p-6">
        <Link href="/dashboard" className="flex items-center space-x-2">
          <div className="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
            <GitBranch className="w-5 h-5 text-primary-foreground" />
          </div>
          <span className="text-xl font-bold">Xephyr</span>
        </Link>
        <p className="text-xs text-muted-foreground mt-1 ml-10">AI Project Manager</p>
      </div>

      <nav className="flex-1 px-3 space-y-1">
        {navigation.map((item) => {
          const isActive = pathname === item.href || pathname.startsWith(`${item.href}/`);
          const badgeCount = item.badge === 'unreadNudges' ? unreadNudges : 0;
          
          return (
            <Button
              key={item.name}
              variant={isActive ? 'secondary' : 'ghost'}
              className={cn(
                'w-full justify-start relative',
                isActive && 'bg-secondary'
              )}
              asChild
            >
              <Link href={item.href}>
                <item.icon className="mr-3 h-4 w-4" />
                <span className="flex-1 text-left">{item.name}</span>
                {badgeCount > 0 && (
                  <Badge variant="destructive" className="ml-auto h-5 min-w-[20px] px-1.5 text-xs">
                    {badgeCount}
                  </Badge>
                )}
              </Link>
            </Button>
          );
        })}
      </nav>

      <div className="p-4 border-t">
        <div className="flex items-center space-x-3">
          <div className="w-10 h-10 rounded-full bg-gradient-to-br from-primary to-primary/60 flex items-center justify-center text-primary-foreground font-medium">
            SC
          </div>
          <div className="flex-1 min-w-0">
            <p className="text-sm font-medium truncate">Sarah Chen</p>
            <p className="text-xs text-muted-foreground truncate">Project Manager</p>
          </div>
        </div>
      </div>
    </div>
  );
}
