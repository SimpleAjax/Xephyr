'use client';

import { Bell, Search, Plus } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useAppStore } from '@/store/app-store';
import { CreateTaskDialog } from '@/components/tasks/create-task-dialog';
import { useState } from 'react';

export function Header() {
  const unreadNudges = useAppStore((state) => state.unreadNudgesCount);
  const [createTaskOpen, setCreateTaskOpen] = useState(false);

  return (
    <>
      <header className="h-16 border-b flex items-center justify-between px-6 bg-card">
        <div className="flex items-center flex-1">
          <div className="relative w-96">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input 
              placeholder="Search projects, tasks, people..." 
              className="pl-10"
            />
          </div>
        </div>

        <div className="flex items-center space-x-3">
          <Button onClick={() => setCreateTaskOpen(true)}>
            <Plus className="h-4 w-4 mr-2" />
            New Task
          </Button>
          
          <Button variant="ghost" size="icon" className="relative">
            <Bell className="h-5 w-5" />
            {unreadNudges > 0 && (
              <span className="absolute top-1 right-1 w-2 h-2 bg-destructive rounded-full" />
            )}
          </Button>
        </div>
      </header>
      
      <CreateTaskDialog open={createTaskOpen} onOpenChange={setCreateTaskOpen} />
    </>
  );
}
