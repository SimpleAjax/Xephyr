-- Fix task dependencies insert with valid UUIDs
INSERT INTO task_dependencies (id, task_id, depends_on_task_id, dependency_type, lag_hours, created_at, updated_at)
VALUES 
  ('dddddddd-e29b-41d4-a716-446655440000', 'cc0e8400-e29b-41d4-a716-446655440030', 'cc0e8400-e29b-41d4-a716-446655440020', 'finish_to_start', 0, NOW(), NOW()),
  ('dddddddd-e29b-41d4-a716-446655440001', 'cc0e8400-e29b-41d4-a716-446655440021', 'cc0e8400-e29b-41d4-a716-446655440020', 'finish_to_start', 8, NOW(), NOW()),
  ('dddddddd-e29b-41d4-a716-446655440002', 'cc0e8400-e29b-41d4-a716-446655440031', 'cc0e8400-e29b-41d4-a716-446655440010', 'finish_to_start', 0, NOW(), NOW()),
  ('dddddddd-e29b-41d4-a716-446655440003', 'cc0e8400-e29b-41d4-a716-446655440032', 'cc0e8400-e29b-41d4-a716-446655440020', 'finish_to_start', 0, NOW(), NOW()),
  ('dddddddd-e29b-41d4-a716-446655440004', 'cc0e8400-e29b-41d4-a716-446655440001', 'cc0e8400-e29b-41d4-a716-446655440000', 'finish_to_start', 168, NOW(), NOW()),
  ('dddddddd-e29b-41d4-a716-446655440005', 'cc0e8400-e29b-41d4-a716-446655440111', 'cc0e8400-e29b-41d4-a716-446655440110', 'start_to_start', 40, NOW(), NOW()),
  ('dddddddd-e29b-41d4-a716-446655440006', 'cc0e8400-e29b-41d4-a716-446655440112', 'cc0e8400-e29b-41d4-a716-446655440111', 'finish_to_start', 0, NOW(), NOW()),
  ('dddddddd-e29b-41d4-a716-446655440007', 'cc0e8400-e29b-41d4-a716-446655440211', 'cc0e8400-e29b-41d4-a716-446655440210', 'start_to_start', 0, NOW(), NOW());
