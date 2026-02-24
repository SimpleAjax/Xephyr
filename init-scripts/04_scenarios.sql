-- Fix scenarios insert with valid UUIDs
INSERT INTO scenarios (id, organization_id, title, description, change_type, status, proposed_changes, created_by_id, created_at, updated_at)
VALUES 
  ('aaaaaaaa-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000', 
   'Emma takes 1-week vacation', 
   'Simulate impact of Emma taking vacation next week',
   'employee_leave', 'pending',
   '{"personId": "660e8400-e29b-41d4-a716-446655440002", "leaveStartDate": "2026-02-24", "leaveEndDate": "2026-02-28"}'::jsonb,
   '660e8400-e29b-41d4-a716-446655440000', NOW(), NOW()),
   
  ('aaaaaaaa-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440000',
   'Add 2 sprints to Fitness App',
   'Scope increase to add social features earlier',
   'scope_change', 'pending',
   '{"projectId": "aa0e8400-e29b-41d4-a716-446655440001", "additionalSprints": 2}'::jsonb,
   '660e8400-e29b-41d4-a716-446655440000', NOW(), NOW()),
   
  ('aaaaaaaa-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440000',
   'Reassign E-Commerce Backend Tasks',
   'Move some backend tasks from James to David',
   'reallocation', 'approved',
   '{"taskId": "cc0e8400-e29b-41d4-a716-446655440022"}'::jsonb,
   '660e8400-e29b-41d4-a716-446655440004', NOW(), NOW());
