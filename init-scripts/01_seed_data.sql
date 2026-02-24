-- Xephyr Seed Data
-- Run this to populate the database with dummy data for testing

-- ============================================
-- ORGANIZATION
-- ============================================
INSERT INTO organizations (id, name, slug, plan, created_at, updated_at)
VALUES 
  ('550e8400-e29b-41d4-a716-446655440000', 'Acme Corp', 'acme-corp', 'pro', NOW(), NOW()),
  ('550e8400-e29b-41d4-a716-446655440001', 'Test Org', 'test-org', 'free', NOW(), NOW());

-- ============================================
-- USERS (Team Members)
-- ============================================
INSERT INTO users (id, email, name, avatar_url, hourly_rate, timezone, is_active, created_at, updated_at)
VALUES 
  -- Core Team
  ('660e8400-e29b-41d4-a716-446655440000', 'sarah.chen@acme.com', 'Sarah Chen', 'https://api.dicebear.com/7.x/avataaars/svg?seed=Sarah', 85.00, 'America/Los_Angeles', true, NOW(), NOW()),
  ('660e8400-e29b-41d4-a716-446655440001', 'mike.rodriguez@acme.com', 'Mike Rodriguez', 'https://api.dicebear.com/7.x/avataaars/svg?seed=Mike', 75.00, 'America/New_York', true, NOW(), NOW()),
  ('660e8400-e29b-41d4-a716-446655440002', 'emma.wilson@acme.com', 'Emma Wilson', 'https://api.dicebear.com/7.x/avataaars/svg?seed=Emma', 80.00, 'Europe/London', true, NOW(), NOW()),
  ('660e8400-e29b-41d4-a716-446655440003', 'alex.thompson@acme.com', 'Alex Thompson', 'https://api.dicebear.com/7.x/avataaars/svg?seed=Alex', 70.00, 'America/Chicago', true, NOW(), NOW()),
  ('660e8400-e29b-41d4-a716-446655440004', 'james.liu@acme.com', 'James Liu', 'https://api.dicebear.com/7.x/avataaars/svg?seed=James', 90.00, 'America/Los_Angeles', true, NOW(), NOW()),
  ('660e8400-e29b-41d4-a716-446655440005', 'rachel.green@acme.com', 'Rachel Green', 'https://api.dicebear.com/7.x/avataaars/svg?seed=Rachel', 72.00, 'America/New_York', true, NOW(), NOW()),
  ('660e8400-e29b-41d4-a716-446655440006', 'david.kim@acme.com', 'David Kim', 'https://api.dicebear.com/7.x/avataaars/svg?seed=David', 78.00, 'America/Los_Angeles', true, NOW(), NOW()),
  ('660e8400-e29b-41d4-a716-446655440007', 'lisa.patel@acme.com', 'Lisa Patel', 'https://api.dicebear.com/7.x/avataaars/svg?seed=Lisa', 82.00, 'America/New_York', true, NOW(), NOW());

-- ============================================
-- ORGANIZATION MEMBERS
-- ============================================
INSERT INTO organization_members (id, organization_id, user_id, role, joined_at, created_at, updated_at)
VALUES 
  ('770e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440000', 'admin', NOW(), NOW(), NOW()),
  ('770e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440001', 'member', NOW(), NOW(), NOW()),
  ('770e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440002', 'member', NOW(), NOW(), NOW()),
  ('770e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440003', 'member', NOW(), NOW(), NOW()),
  ('770e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440004', 'pm', NOW(), NOW(), NOW()),
  ('770e8400-e29b-41d4-a716-446655440005', '550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440005', 'member', NOW(), NOW(), NOW()),
  ('770e8400-e29b-41d4-a716-446655440006', '550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440006', 'member', NOW(), NOW(), NOW()),
  ('770e8400-e29b-41d4-a716-446655440007', '550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440007', 'member', NOW(), NOW(), NOW());

-- ============================================
-- SKILLS
-- ============================================
INSERT INTO skills (id, organization_id, name, category, description, created_at, updated_at)
VALUES 
  -- Frontend Skills
  ('880e8400-e29b-41d4-a716-446655440000', NULL, 'React', 'Frontend', 'React.js framework', NOW(), NOW()),
  ('880e8400-e29b-41d4-a716-446655440001', NULL, 'TypeScript', 'Frontend', 'TypeScript language', NOW(), NOW()),
  ('880e8400-e29b-41d4-a716-446655440002', NULL, 'Vue.js', 'Frontend', 'Vue.js framework', NOW(), NOW()),
  ('880e8400-e29b-41d4-a716-446655440003', NULL, 'Next.js', 'Frontend', 'Next.js framework', NOW(), NOW()),
  
  -- Backend Skills
  ('880e8400-e29b-41d4-a716-446655440010', NULL, 'Go', 'Backend', 'Go programming language', NOW(), NOW()),
  ('880e8400-e29b-41d4-a716-446655440011', NULL, 'Node.js', 'Backend', 'Node.js runtime', NOW(), NOW()),
  ('880e8400-e29b-41d4-a716-446655440012', NULL, 'Python', 'Backend', 'Python programming', NOW(), NOW()),
  ('880e8400-e29b-41d4-a716-446655440013', NULL, 'PostgreSQL', 'Database', 'PostgreSQL database', NOW(), NOW()),
  
  -- Design Skills
  ('880e8400-e29b-41d4-a716-446655440020', NULL, 'Figma', 'Design', 'Figma design tool', NOW(), NOW()),
  ('880e8400-e29b-41d4-a716-446655440021', NULL, 'UI/UX', 'Design', 'UI/UX design', NOW(), NOW()),
  
  -- DevOps/Other
  ('880e8400-e29b-41d4-a716-446655440030', NULL, 'Docker', 'DevOps', 'Docker containerization', NOW(), NOW()),
  ('880e8400-e29b-41d4-a716-446655440031', NULL, 'AWS', 'DevOps', 'Amazon Web Services', NOW(), NOW()),
  ('880e8400-e29b-41d4-a716-446655440032', NULL, 'Kubernetes', 'DevOps', 'Kubernetes orchestration', NOW(), NOW());

-- ============================================
-- USER SKILLS
-- ============================================
INSERT INTO user_skills (id, user_id, skill_id, proficiency, years_of_experience, created_at, updated_at)
VALUES 
  -- Sarah (Full-stack lead)
  ('990e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440000', '880e8400-e29b-41d4-a716-446655440000', 4, 5.0, NOW(), NOW()),
  ('990e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440000', '880e8400-e29b-41d4-a716-446655440001', 4, 4.0, NOW(), NOW()),
  ('990e8400-e29b-41d4-a716-446655440002', '660e8400-e29b-41d4-a716-446655440000', '880e8400-e29b-41d4-a716-446655440010', 3, 3.0, NOW(), NOW()),
  
  -- Mike (Frontend specialist)
  ('990e8400-e29b-41d4-a716-446655440010', '660e8400-e29b-41d4-a716-446655440001', '880e8400-e29b-41d4-a716-446655440000', 4, 6.0, NOW(), NOW()),
  ('990e8400-e29b-41d4-a716-446655440011', '660e8400-e29b-41d4-a716-446655440001', '880e8400-e29b-41d4-a716-446655440001', 4, 5.0, NOW(), NOW()),
  ('990e8400-e29b-41d4-a716-446655440012', '660e8400-e29b-41d4-a716-446655440001', '880e8400-e29b-41d4-a716-446655440003', 3, 3.0, NOW(), NOW()),
  
  -- Emma (Designer)
  ('990e8400-e29b-41d4-a716-446655440020', '660e8400-e29b-41d4-a716-446655440002', '880e8400-e29b-41d4-a716-446655440020', 4, 5.0, NOW(), NOW()),
  ('990e8400-e29b-41d4-a716-446655440021', '660e8400-e29b-41d4-a716-446655440002', '880e8400-e29b-41d4-a716-446655440021', 4, 6.0, NOW(), NOW()),
  ('990e8400-e29b-41d4-a716-446655440022', '660e8400-e29b-41d4-a716-446655440002', '880e8400-e29b-41d4-a716-446655440000', 2, 2.0, NOW(), NOW()),
  
  -- Alex (Backend)
  ('990e8400-e29b-41d4-a716-446655440030', '660e8400-e29b-41d4-a716-446655440003', '880e8400-e29b-41d4-a716-446655440010', 4, 4.0, NOW(), NOW()),
  ('990e8400-e29b-41d4-a716-446655440031', '660e8400-e29b-41d4-a716-446655440003', '880e8400-e29b-41d4-a716-446655440013', 3, 3.0, NOW(), NOW()),
  ('990e8400-e29b-41d4-a716-446655440032', '660e8400-e29b-41d4-a716-446655440003', '880e8400-e29b-41d4-a716-446655440011', 3, 4.0, NOW(), NOW()),
  
  -- James (Senior Backend)
  ('990e8400-e29b-41d4-a716-446655440040', '660e8400-e29b-41d4-a716-446655440004', '880e8400-e29b-41d4-a716-446655440010', 4, 7.0, NOW(), NOW()),
  ('990e8400-e29b-41d4-a716-446655440041', '660e8400-e29b-41d4-a716-446655440004', '880e8400-e29b-41d4-a716-446655440013', 4, 6.0, NOW(), NOW()),
  ('990e8400-e29b-41d4-a716-446655440042', '660e8400-e29b-41d4-a716-446655440004', '880e8400-e29b-41d4-a716-446655440030', 3, 4.0, NOW(), NOW()),
  
  -- Rachel (Full-stack)
  ('990e8400-e29b-41d4-a716-446655440050', '660e8400-e29b-41d4-a716-446655440005', '880e8400-e29b-41d4-a716-446655440000', 3, 3.0, NOW(), NOW()),
  ('990e8400-e29b-41d4-a716-446655440051', '660e8400-e29b-41d4-a716-446655440005', '880e8400-e29b-41d4-a716-446655440011', 3, 4.0, NOW(), NOW()),
  ('990e8400-e29b-41d4-a716-446655440052', '660e8400-e29b-41d4-a716-446655440005', '880e8400-e29b-41d4-a716-446655440020', 2, 2.0, NOW(), NOW()),
  
  -- David (DevOps)
  ('990e8400-e29b-41d4-a716-446655440060', '660e8400-e29b-41d4-a716-446655440006', '880e8400-e29b-41d4-a716-446655440030', 4, 6.0, NOW(), NOW()),
  ('990e8400-e29b-41d4-a716-446655440061', '660e8400-e29b-41d4-a716-446655440006', '880e8400-e29b-41d4-a716-446655440031', 4, 5.0, NOW(), NOW()),
  ('990e8400-e29b-41d4-a716-446655440062', '660e8400-e29b-41d4-a716-446655440006', '880e8400-e29b-41d4-a716-446655440032', 3, 3.0, NOW(), NOW()),
  
  -- Lisa (Data/Python)
  ('990e8400-e29b-41d4-a716-446655440070', '660e8400-e29b-41d4-a716-446655440007', '880e8400-e29b-41d4-a716-446655440012', 4, 5.0, NOW(), NOW()),
  ('990e8400-e29b-41d4-a716-446655440071', '660e8400-e29b-41d4-a716-446655440007', '880e8400-e29b-41d4-a716-446655440013', 4, 6.0, NOW(), NOW()),
  ('990e8400-e29b-41d4-a716-446655440072', '660e8400-e29b-41d4-a716-446655440007', '880e8400-e29b-41d4-a716-446655440011', 3, 4.0, NOW(), NOW());
