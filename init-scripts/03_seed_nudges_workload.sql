-- ============================================
-- NUDGES (AI Alerts)
-- ============================================
INSERT INTO nudges (id, organization_id, type, severity, status, title, description, ai_explanation, suggested_action, confidence_score, criticality_score, expires_at, related_project_id, related_task_id, related_user_id, metrics, created_at, updated_at)
VALUES 
  -- High Severity: Emma Overallocated
  ('dd0e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000', 'overload', 'high', 'unread', 
   'Emma Wilson is overallocated at 125%', 
   'Emma has 50 hours assigned this week across 2 projects',
   'Based on current assignments, Emma is scheduled for 50 hours this week (125% of her 40-hour capacity). This puts the Fitness App UI Design at risk of delay. Historical data shows that overallocated designers produce 30% more bugs.',
   'Reassign Marketing Website consultation to Rachel (75% compatibility)',
   0.92, 90, NOW() + INTERVAL '7 days',
   'aa0e8400-e29b-41d4-a716-446655440001', 'cc0e8400-e29b-41d4-a716-446655440110', '660e8400-e29b-41d4-a716-446655440002',
   '{"allocationPercentage": 125, "assignedTasks": 4, "totalHours": 50}'::jsonb,
   NOW(), NOW()),

  -- High Severity: Checkout Flow Blocked
  ('dd0e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440000', 'dependency_block', 'high', 'unread',
   'Checkout Flow blocked by Backend API delay',
   'Payment Integration dependency is at risk - Backend API is only 60% complete',
   'The critical path analysis shows that Backend API Development (task-ec-2) is 10 days behind schedule. This directly blocks Checkout Flow Implementation, which is on the critical path for MVP Launch. Delaying the MVP would cost approximately $12,000 per day in lost revenue.',
   'Assign additional backend developer (David Kim) to accelerate API completion',
   0.88, 85, NOW() + INTERVAL '5 days',
   'aa0e8400-e29b-41d4-a716-446655440000', 'cc0e8400-e29b-41d4-a716-446655440030', NULL,
   '{"blockedTasksCount": 2, "delayRisk": 10}'::jsonb,
   NOW(), NOW()),

  -- Medium Severity: Unassigned Critical Task
  ('dd0e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440000', 'unassigned', 'medium', 'unread',
   'Critical task unassigned: Checkout Flow Implementation',
   'High-priority task has no assignee and is at risk of missing deadline',
   'Checkout Flow Implementation is on the critical path with a due date of March 1st. The task requires React (Expert) and TypeScript (Expert) skills. Mike Rodriguez matches these requirements with 92% compatibility and has availability.',
   'Assign to Mike Rodriguez (92% skill match, 60% current allocation)',
   0.85, 75, NOW() + INTERVAL '7 days',
   'aa0e8400-e29b-41d4-a716-446655440000', 'cc0e8400-e29b-41d4-a716-446655440030', NULL,
   '{"daysUnassigned": 5, "skillMatchScore": 92}'::jsonb,
   NOW(), NOW()),

  -- Medium Severity: Delay Risk
  ('dd0e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440000', 'delay_risk', 'medium', 'unread',
   'Fitness App UI Design at risk of delay',
   'Progress is 44% vs expected 60% at this stage',
   'Based on velocity tracking, Emma''s current pace suggests a 5-day delay for Mobile App UI Design. The Fitness App launch date is at risk. Consider reducing scope or adding design resources.',
   'Reduce icon set scope by 30% to maintain timeline',
   0.78, 70, NOW() + INTERVAL '7 days',
   'aa0e8400-e29b-41d4-a716-446655440001', 'cc0e8400-e29b-41d4-a716-446655440110', '660e8400-e29b-41d4-a716-446655440002',
   '{"expectedProgress": 60, "actualProgress": 44, "variance": -16}'::jsonb,
   NOW(), NOW()),

  -- Low Severity: Skill Gap
  ('dd0e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440000', 'skill_gap', 'low', 'unread',
   'Kubernetes expertise needed for API Modernization',
   'Project requires k8s skills but no assigned member has this expertise',
   'The API Modernization project requires Kubernetes for deployment. Current team assignments show a gap in this area. Consider training existing team members or bringing in a contractor.',
   'Assign David Kim (DevOps) as advisor or schedule k8s training',
   0.72, 55, NOW() + INTERVAL '14 days',
   'aa0e8400-e29b-41d4-a716-446655440004', NULL, NULL,
   '{"requiredSkill": "Kubernetes", "coverage": 0}'::jsonb,
   NOW(), NOW()),

  -- High Severity: Alex Overallocated
  ('dd0e8400-e29b-41d4-a716-446655440005', '550e8400-e29b-41d4-a716-446655440000', 'overload', 'high', 'unread',
   'Alex Thompson is overallocated at 110%',
   'Alex has 44 hours assigned this week',
   'Alex is currently at 110% allocation with tasks from both Fitness App and E-Commerce projects. Adding new assignments would result in overallocation. Historical burnout risk is elevated.',
   'Delay non-critical tasks or redistribute to Rachel Green',
   0.90, 80, NOW() + INTERVAL '7 days',
   'aa0e8400-e29b-41d4-a716-446655440001', 'cc0e8400-e29b-41d4-a716-446655440112', '660e8400-e29b-41d4-a716-446655440003',
   '{"allocationPercentage": 110, "assignedTasks": 3, "totalHours": 44}'::jsonb,
   NOW(), NOW()),

  -- Low Severity: Blocked Task
  ('dd0e8400-e29b-41d4-a716-446655440006', '550e8400-e29b-41d4-a716-446655440000', 'blocked', 'low', 'read',
   'Product Catalog waiting on API endpoints',
   'Frontend work is blocked by incomplete backend',
   'Product Catalog frontend is ready to proceed but waiting on catalog API endpoints from Backend API Development. This is causing idle time for Mike Rodriguez.',
   'Review API contract and implement mock responses for parallel development',
   0.65, 45, NOW() + INTERVAL '3 days',
   'aa0e8400-e29b-41d4-a716-446655440000', 'cc0e8400-e29b-41d4-a716-446655440031', '660e8400-e29b-41d4-a716-446655440001',
   '{"blockedHours": 8, "waitingOn": "Backend API"}'::jsonb,
   NOW() - INTERVAL '1 day', NOW()),

  -- Medium Severity: Conflict
  ('dd0e8400-e29b-41d4-a716-446655440007', '550e8400-e29b-41d4-a716-446655440000', 'conflict', 'medium', 'unread',
   'Resource conflict: Rachel Green double-booked',
   'Two tasks scheduled for same time period',
   'Analysis shows Rachel has overlapping commitments on Marketing Website and E-Commerce Admin Dashboard for next week. Both tasks are medium priority and require her React skills.',
   'Reschedule Admin Dashboard work to following sprint',
   0.82, 65, NOW() + INTERVAL '5 days',
   'aa0e8400-e29b-41d4-a716-446655440002', 'cc0e8400-e29b-41d4-a716-446655440214', '660e8400-e29b-41d4-a716-446655440005',
   '{"conflictingTasks": 2, "overlapDays": 5}'::jsonb,
   NOW(), NOW());

-- ============================================
-- WORKLOAD ENTRIES (Weekly Allocation)
-- ============================================
INSERT INTO workload_entries (id, organization_id, user_id, week_start, allocation_percentage, assigned_tasks, total_estimated_hours, available_hours, created_at, updated_at)
VALUES 
  -- Emma Wilson - Overallocated (125%)
  ('ee0e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440002', DATE_TRUNC('week', NOW()), 125, 4, 50, 0, NOW(), NOW()),
  
  -- Alex Thompson - Slightly overallocated (110%)
  ('ee0e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440003', DATE_TRUNC('week', NOW()), 110, 3, 44, 0, NOW(), NOW()),
  
  -- Mike Rodriguez - Good (75%)
  ('ee0e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440001', DATE_TRUNC('week', NOW()), 75, 3, 30, 10, NOW(), NOW()),
  
  -- James Liu - Optimal (90%)
  ('ee0e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440004', DATE_TRUNC('week', NOW()), 90, 4, 36, 4, NOW(), NOW()),
  
  -- Rachel Green - Optimal (85%)
  ('ee0e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440005', DATE_TRUNC('week', NOW()), 85, 3, 34, 6, NOW(), NOW()),
  
  -- David Kim - Available (40%)
  ('ee0e8400-e29b-41d4-a716-446655440005', '550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440006', DATE_TRUNC('week', NOW()), 40, 2, 16, 24, NOW(), NOW()),
  
  -- Lisa Patel - Optimal (80%)
  ('ee0e8400-e29b-41d4-a716-446655440006', '550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440007', DATE_TRUNC('week', NOW()), 80, 3, 32, 8, NOW(), NOW()),
  
  -- Sarah Chen - Optimal (70%)
  ('ee0e8400-e29b-41d4-a716-446655440007', '550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440000', DATE_TRUNC('week', NOW()), 70, 2, 28, 12, NOW(), NOW());

-- ============================================
-- TASK SKILLS (Required skills for tasks)
-- ============================================
INSERT INTO task_skills (id, task_id, skill_id, proficiency_required, is_required, created_at, updated_at)
VALUES 
  -- E-Commerce: Backend API Development
  ('ff0e8400-e29b-41d4-a716-446655440000', 'cc0e8400-e29b-41d4-a716-446655440020', '880e8400-e29b-41d4-a716-446655440010', 3, true, NOW(), NOW()),
  ('ff0e8400-e29b-41d4-a716-446655440001', 'cc0e8400-e29b-41d4-a716-446655440020', '880e8400-e29b-41d4-a716-446655440013', 3, true, NOW(), NOW()),
  
  -- E-Commerce: Checkout Flow
  ('ff0e8400-e29b-41d4-a716-446655440010', 'cc0e8400-e29b-41d4-a716-446655440030', '880e8400-e29b-41d4-a716-446655440000', 4, true, NOW(), NOW()),
  ('ff0e8400-e29b-41d4-a716-446655440011', 'cc0e8400-e29b-41d4-a716-446655440030', '880e8400-e29b-41d4-a716-446655440001', 4, true, NOW(), NOW()),
  
  -- E-Commerce: Product Catalog
  ('ff0e8400-e29b-41d4-a716-446655440020', 'cc0e8400-e29b-41d4-a716-446655440031', '880e8400-e29b-41d4-a716-446655440000', 3, true, NOW(), NOW()),
  ('ff0e8400-e29b-41d4-a716-446655440021', 'cc0e8400-e29b-41d4-a716-446655440031', '880e8400-e29b-41d4-a716-446655440003', 2, false, NOW(), NOW()),
  
  -- Fitness App: UI Design
  ('ff0e8400-e29b-41d4-a716-446655440030', 'cc0e8400-e29b-41d4-a716-446655440110', '880e8400-e29b-41d4-a716-446655440020', 4, true, NOW(), NOW()),
  ('ff0e8400-e29b-41d4-a716-446655440031', 'cc0e8400-e29b-41d4-a716-446655440110', '880e8400-e29b-41d4-a716-446655440021', 4, true, NOW(), NOW()),
  
  -- Fitness App: Workout Tracking
  ('ff0e8400-e29b-41d4-a716-446655440040', 'cc0e8400-e29b-41d4-a716-446655440111', '880e8400-e29b-41d4-a716-446655440000', 3, true, NOW(), NOW()),
  ('ff0e8400-e29b-41d4-a716-446655440041', 'cc0e8400-e29b-41d4-a716-446655440111', '880e8400-e29b-41d4-a716-446655440011', 3, true, NOW(), NOW());

-- ============================================
-- TASK DEPENDENCIES
-- ============================================
INSERT INTO task_dependencies (id, task_id, depends_on_task_id, dependency_type, lag_hours, created_at, updated_at)
VALUES 
  -- Checkout Flow depends on Backend API
  ('gg0e8400-e29b-41d4-a716-446655440000', 'cc0e8400-e29b-41d4-a716-446655440030', 'cc0e8400-e29b-41d4-a716-446655440020', 'finish_to_start', 0, NOW(), NOW()),
  
  -- Payment Integration depends on Backend API
  ('gg0e8400-e29b-41d4-a716-446655440001', 'cc0e8400-e29b-41d4-a716-446655440021', 'cc0e8400-e29b-41d4-a716-446655440020', 'finish_to_start', 8, NOW(), NOW()),
  
  -- Product Catalog depends on Design System
  ('gg0e8400-e29b-41d4-a716-446655440002', 'cc0e8400-e29b-41d4-a716-446655440031', 'cc0e8400-e29b-41d4-a716-446655440010', 'finish_to_start', 0, NOW(), NOW()),
  
  -- User Account Pages depends on Backend API
  ('gg0e8400-e29b-41d4-a716-446655440003', 'cc0e8400-e29b-41d4-a716-446655440032', 'cc0e8400-e29b-41d4-a716-446655440020', 'finish_to_start', 0, NOW(), NOW()),
  
  -- Full Release depends on MVP Launch
  ('gg0e8400-e29b-41d4-a716-446655440004', 'cc0e8400-e29b-41d4-a716-446655440001', 'cc0e8400-e29b-41d4-a716-446655440000', 'finish_to_start', 168, NOW(), NOW()),
  
  -- Fitness App: Workout Tracking depends on UI Design
  ('gg0e8400-e29b-41d4-a716-446655440005', 'cc0e8400-e29b-41d4-a716-446655440111', 'cc0e8400-e29b-41d4-a716-446655440110', 'start_to_start', 40, NOW(), NOW()),
  
  -- Fitness App: Social Features depends on Workout Tracking
  ('gg0e8400-e29b-41d4-a716-446655440006', 'cc0e8400-e29b-41d4-a716-446655440112', 'cc0e8400-e29b-41d4-a716-446655440111', 'finish_to_start', 0, NOW(), NOW()),
  
  -- Marketing: Product Pages depends on Homepage
  ('gg0e8400-e29b-41d4-a716-446655440007', 'cc0e8400-e29b-41d4-a716-446655440211', 'cc0e8400-e29b-41d4-a716-446655440210', 'start_to_start', 0, NOW(), NOW());

-- ============================================
-- ASSIGNMENT SUGGESTIONS (AI Recommendations)
-- ============================================
INSERT INTO assignment_suggestions (id, task_id, suggested_user_id, total_score, skill_match_score, availability_score, workload_score, performance_score, reasons, warnings, ai_explanation, status, created_at, updated_at)
VALUES 
  ('hh0e8400-e29b-41d4-a716-446655440000', 'cc0e8400-e29b-41d4-a716-446655440030', '660e8400-e29b-41d4-a716-446655440001', 92, 38, 30, 18, 6, 
   '["Expert-level React skills", "TypeScript proficiency", "Available capacity"]'::jsonb,
   '[]'::jsonb,
   'Mike is the best match with 92% compatibility. He has all required skills at high proficiency (React: Expert, TypeScript: Expert). Currently at 60% capacity with availability this week.',
   'pending', NOW(), NOW()),
   
  ('hh0e8400-e29b-41d4-a716-446655440001', 'cc0e8400-e29b-41d4-a716-446655440030', '660e8400-e29b-41d4-a716-446655440005', 78, 30, 10, 10, 8,
   '["Good React skills", "Available next week"]'::jsonb,
   '["Currently at 85% allocation", "Would require task reshuffling"]'::jsonb,
   'Rachel has strong frontend skills but is currently well-utilized. Consider assigning after current Marketing Website tasks wrap up.',
   'pending', NOW(), NOW()),
   
  ('hh0e8400-e29b-41d4-a716-446655440002', 'cc0e8400-e29b-41d4-a716-446655440040', '660e8400-e29b-41d4-a716-446655440005', 85, 35, 20, 20, 10,
   '["Full-stack experience", "React proficiency", "Good availability"]'::jsonb,
   '[]'::jsonb,
   'Rachel is a good fit for the Admin Dashboard with her full-stack background. She can handle both frontend and API integration.',
   'pending', NOW(), NOW());

-- ============================================
-- SCENARIOS (What-if Analysis)
-- ============================================
INSERT INTO scenarios (id, organization_id, title, description, change_type, status, proposed_changes, created_by_id, created_at, updated_at)
VALUES 
  ('ii0e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000', 
   'Emma takes 1-week vacation', 
   'Simulate impact of Emma taking vacation next week',
   'employee_leave', 'pending',
   '{"personId": "660e8400-e29b-41d4-a716-446655440002", "leaveStartDate": "2026-02-24", "leaveEndDate": "2026-02-28", "coverageStrategy": "reassign"}'::jsonb,
   '660e8400-e29b-41d4-a716-446655440000', NOW(), NOW()),
   
  ('ii0e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440000',
   'Add 2 sprints to Fitness App',
   'Scope increase to add social features earlier',
   'scope_change', 'pending',
   '{"projectId": "aa0e8400-e29b-41d4-a716-446655440001", "additionalSprints": 2, "newDeadline": "2026-05-01"}'::jsonb,
   '660e8400-e29b-41d4-a716-446655440000', NOW(), NOW()),
   
  ('ii0e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440000',
   'Reassign E-Commerce Backend Tasks',
   'Move some backend tasks from James to David',
   'reallocation', 'applied',
   '{"reassignments": [{"taskId": "cc0e8400-e29b-41d4-a716-446655440022", "fromUserId": "660e8400-e29b-41d4-a716-446655440004", "toUserId": "660e8400-e29b-41d4-a716-446655440006"}]}'::jsonb,
   '660e8400-e29b-41d4-a716-446655440004', NOW(), NOW());

-- Insert scenario impact analysis for applied scenario
INSERT INTO scenario_impact_analyses (id, scenario_id, delay_hours_total, cost_impact, affected_project_ids, affected_task_ids, recommendations, timeline_comparison, created_at, updated_at)
VALUES 
  ('jj0e8400-e29b-41d4-a716-446655440000', 'ii0e8400-e29b-41d4-a716-446655440002', 0, 0,
   '["aa0e8400-e29b-41d4-a716-446655440000"]'::jsonb,
   '["cc0e8400-e29b-41d4-a716-446655440022"]'::jsonb,
   '["Better workload distribution", "David gains domain knowledge"]'::jsonb,
   '{"originalEndDate": "2026-03-05", "newEndDate": "2026-03-05"}'::jsonb,
   NOW(), NOW());
