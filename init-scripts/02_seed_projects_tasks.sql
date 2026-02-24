-- ============================================
-- PROJECTS
-- ============================================
INSERT INTO projects (id, organization_id, name, description, status, priority, health_score, progress, start_date, target_end_date, budget, created_at, updated_at)
VALUES 
  -- E-Commerce Platform (Active, at risk)
  ('aa0e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000', 'E-Commerce Platform', 'Complete redesign of customer-facing e-commerce platform', 'active', 95, 72, 45, '2026-01-15', '2026-05-01', 250000.00, NOW(), NOW()),
  
  -- Fitness App (Active, at risk)
  ('aa0e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440000', 'Fitness App Mobile Launch', 'Native mobile app for iOS and Android', 'active', 88, 45, 25, '2026-01-20', '2026-04-15', 180000.00, NOW(), NOW()),
  
  -- Marketing Website (Active, healthy)
  ('aa0e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440000', 'Marketing Website Refresh', 'New marketing site with improved conversion', 'active', 75, 85, 70, '2026-02-01', '2026-03-30', 80000.00, NOW(), NOW()),
  
  -- Internal Dashboard (Active, healthy)
  ('aa0e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440000', 'Internal Analytics Dashboard', 'Executive dashboard for business metrics', 'active', 60, 90, 80, '2026-01-10', '2026-03-15', 120000.00, NOW(), NOW()),
  
  -- API Modernization (Paused)
  ('aa0e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440000', 'API Modernization', 'Migrate legacy APIs to GraphQL', 'paused', 70, 50, 30, '2025-12-01', '2026-06-01', 200000.00, NOW(), NOW());

-- ============================================
-- PROJECT MEMBERS
-- ============================================
INSERT INTO project_members (id, project_id, user_id, role, created_at, updated_at)
VALUES 
  -- E-Commerce Team
  ('bb0e8400-e29b-41d4-a716-446655440000', 'aa0e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440000', 'Tech Lead', NOW(), NOW()),
  ('bb0e8400-e29b-41d4-a716-446655440001', 'aa0e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440001', 'Frontend Dev', NOW(), NOW()),
  ('bb0e8400-e29b-41d4-a716-446655440002', 'aa0e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440004', 'Backend Lead', NOW(), NOW()),
  ('bb0e8400-e29b-41d4-a716-446655440003', 'aa0e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440002', 'UI/UX Designer', NOW(), NOW()),
  
  -- Fitness App Team
  ('bb0e8400-e29b-41d4-a716-446655440010', 'aa0e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440002', 'Lead Designer', NOW(), NOW()),
  ('bb0e8400-e29b-41d4-a716-446655440011', 'aa0e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440005', 'Full-stack Dev', NOW(), NOW()),
  ('bb0e8400-e29b-41d4-a716-446655440012', 'aa0e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440003', 'Backend Dev', NOW(), NOW()),
  
  -- Marketing Site Team
  ('bb0e8400-e29b-41d4-a716-446655440020', 'aa0e8400-e29b-41d4-a716-446655440002', '660e8400-e29b-41d4-a716-446655440005', 'Frontend Lead', NOW(), NOW()),
  ('bb0e8400-e29b-41d4-a716-446655440021', 'aa0e8400-e29b-41d4-a716-446655440002', '660e8400-e29b-41d4-a716-446655440002', 'Designer', NOW(), NOW()),
  
  -- Dashboard Team
  ('bb0e8400-e29b-41d4-a716-446655440030', 'aa0e8400-e29b-41d4-a716-446655440003', '660e8400-e29b-41d4-a716-446655440007', 'Data Engineer', NOW(), NOW()),
  ('bb0e8400-e29b-41d4-a716-446655440031', 'aa0e8400-e29b-41d4-a716-446655440003', '660e8400-e29b-41d4-a716-446655440001', 'Frontend Dev', NOW(), NOW());

-- ============================================
-- TASKS - E-Commerce Platform
-- ============================================
INSERT INTO tasks (id, project_id, parent_task_id, hierarchy_level, title, description, status, priority, priority_score, business_value, estimated_hours, actual_hours, start_date, due_date, assignee_id, is_milestone, is_critical_path, risk_score, created_at, updated_at)
VALUES 
  -- Milestones
  ('cc0e8400-e29b-41d4-a716-446655440000', 'aa0e8400-e29b-41d4-a716-446655440000', NULL, 1, 'MVP Launch', 'Initial MVP release with core features', 'in_progress', 'high', 95, 100, 400, 180, '2026-01-15', '2026-03-15', NULL, true, true, 20, NOW(), NOW()),
  ('cc0e8400-e29b-41d4-a716-446655440001', 'aa0e8400-e29b-41d4-a716-446655440000', NULL, 1, 'Full Release', 'Complete platform launch', 'backlog', 'critical', 90, 100, 300, 0, NULL, '2026-05-01', NULL, true, true, 40, NOW(), NOW()),

  -- Design Phase (Done)
  ('cc0e8400-e29b-41d4-a716-446655440010', 'aa0e8400-e29b-41d4-a716-446655440000', 'cc0e8400-e29b-41d4-a716-446655440000', 2, 'Design System Architecture', 'Create component library and design tokens', 'done', 'high', 85, 90, 40, 42, '2026-01-15', '2026-01-30', '660e8400-e29b-41d4-a716-446655440002', false, false, 10, NOW(), NOW()),
  ('cc0e8400-e29b-41d4-a716-446655440011', 'aa0e8400-e29b-41d4-a716-446655440000', 'cc0e8400-e29b-41d4-a716-446655440000', 2, 'High-fidelity Mockups', 'Complete UI designs for all pages', 'done', 'high', 80, 85, 60, 65, '2026-01-20', '2026-02-10', '660e8400-e29b-41d4-a716-446655440002', false, false, 15, NOW(), NOW()),
  
  -- Backend (In Progress)
  ('cc0e8400-e29b-41d4-a716-446655440020', 'aa0e8400-e29b-41d4-a716-446655440000', 'cc0e8400-e29b-41d4-a716-446655440000', 2, 'Backend API Development', 'RESTful APIs for product catalog, cart, orders', 'in_progress', 'high', 92, 95, 80, 48, '2026-01-25', '2026-03-01', '660e8400-e29b-41d4-a716-446655440004', false, true, 25, NOW(), NOW()),
  ('cc0e8400-e29b-41d4-a716-446655440021', 'aa0e8400-e29b-41d4-a716-446655440000', 'cc0e8400-e29b-41d4-a716-446655440000', 2, 'Payment Integration', 'Stripe integration for payments', 'ready', 'critical', 88, 100, 40, 0, NULL, '2026-03-10', '660e8400-e29b-41d4-a716-446655440004', false, true, 30, NOW(), NOW()),
  ('cc0e8400-e29b-41d4-a716-446655440022', 'aa0e8400-e29b-41d4-a716-446655440000', 'cc0e8400-e29b-41d4-a716-446655440000', 2, 'Inventory Management', 'Real-time inventory sync', 'in_progress', 'medium', 75, 80, 50, 20, '2026-02-05', '2026-03-05', '660e8400-e29b-41d4-a716-446655440003', false, false, 20, NOW(), NOW()),
  
  -- Frontend (In Progress)
  ('cc0e8400-e29b-41d4-a716-446655440030', 'aa0e8400-e29b-41d4-a716-446655440000', 'cc0e8400-e29b-41d4-a716-446655440000', 2, 'Checkout Flow Implementation', 'Complete checkout with validation', 'ready', 'critical', 95, 100, 50, 0, NULL, '2026-03-01', NULL, false, true, 35, NOW(), NOW()),
  ('cc0e8400-e29b-41d4-a716-446655440031', 'aa0e8400-e29b-41d4-a716-446655440000', 'cc0e8400-e29b-41d4-a716-446655440000', 2, 'Product Catalog', 'Browse and search products', 'in_progress', 'high', 85, 90, 60, 35, '2026-02-01', '2026-02-28', '660e8400-e29b-41d4-a716-446655440001', false, true, 25, NOW(), NOW()),
  ('cc0e8400-e29b-41d4-a716-446655440032', 'aa0e8400-e29b-41d4-a716-446655440000', 'cc0e8400-e29b-41d4-a716-446655440000', 2, 'User Account Pages', 'Login, profile, order history', 'in_progress', 'medium', 70, 75, 40, 15, '2026-02-10', '2026-03-01', '660e8400-e29b-41d4-a716-446655440001', false, false, 15, NOW(), NOW()),
  
  -- Admin Dashboard
  ('cc0e8400-e29b-41d4-a716-446655440040', 'aa0e8400-e29b-41d4-a716-446655440000', 'cc0e8400-e29b-41d4-a716-446655440000', 2, 'Admin Dashboard', 'Internal admin for product management', 'backlog', 'medium', 65, 70, 80, 0, NULL, '2026-04-01', '660e8400-e29b-41d4-a716-446655440005', false, false, 10, NOW(), NOW()),
  ('cc0e8400-e29b-41d4-a716-446655440041', 'aa0e8400-e29b-41d4-a716-446655440000', 'cc0e8400-e29b-41d4-a716-446655440000', 2, 'Analytics Integration', 'Track user behavior and sales', 'backlog', 'low', 55, 60, 40, 0, NULL, '2026-04-15', NULL, false, false, 15, NOW(), NOW());

-- ============================================
-- TASKS - Fitness App
-- ============================================
INSERT INTO tasks (id, project_id, parent_task_id, hierarchy_level, title, description, status, priority, priority_score, business_value, estimated_hours, actual_hours, start_date, due_date, assignee_id, is_milestone, is_critical_path, risk_score, created_at, updated_at)
VALUES 
  ('cc0e8400-e29b-41d4-a716-446655440100', 'aa0e8400-e29b-41d4-a716-446655440001', NULL, 1, 'iOS App Launch', 'Release on App Store', 'backlog', 'critical', 90, 100, 350, 0, NULL, '2026-04-15', NULL, true, true, 45, NOW(), NOW()),
  
  ('cc0e8400-e29b-41d4-a716-446655440110', 'aa0e8400-e29b-41d4-a716-446655440001', 'cc0e8400-e29b-41d4-a716-446655440100', 2, 'Mobile App UI Design', 'Complete iOS/Android designs', 'in_progress', 'high', 88, 95, 80, 35, '2026-01-20', '2026-03-01', '660e8400-e29b-41d4-a716-446655440002', false, true, 35, NOW(), NOW()),
  ('cc0e8400-e29b-41d4-a716-446655440111', 'aa0e8400-e29b-41d4-a716-446655440001', 'cc0e8400-e29b-41d4-a716-446655440100', 2, 'Workout Tracking Module', 'Core exercise tracking functionality', 'in_progress', 'high', 85, 90, 100, 25, '2026-02-01', '2026-03-15', '660e8400-e29b-41d4-a716-446655440005', false, true, 40, NOW(), NOW()),
  ('cc0e8400-e29b-41d4-a716-446655440112', 'aa0e8400-e29b-41d4-a716-446655440001', 'cc0e8400-e29b-41d4-a716-446655440100', 2, 'Social Features', 'Share workouts, challenges', 'backlog', 'medium', 60, 70, 80, 0, NULL, '2026-04-01', '660e8400-e29b-41d4-a716-446655440003', false, false, 25, NOW(), NOW()),
  ('cc0e8400-e29b-41d4-a716-446655440113', 'aa0e8400-e29b-41d4-a716-446655440001', 'cc0e8400-e29b-41d4-a716-446655440100', 2, 'Apple Health Integration', 'Sync with Apple HealthKit', 'ready', 'medium', 70, 75, 40, 0, NULL, '2026-03-10', NULL, false, false, 20, NOW(), NOW());

-- ============================================
-- TASKS - Marketing Website
-- ============================================
INSERT INTO tasks (id, project_id, parent_task_id, hierarchy_level, title, description, status, priority, priority_score, business_value, estimated_hours, actual_hours, start_date, due_date, assignee_id, is_milestone, is_critical_path, risk_score, created_at, updated_at)
VALUES 
  ('cc0e8400-e29b-41d4-a716-446655440200', 'aa0e8400-e29b-41d4-a716-446655440002', NULL, 1, 'Website Launch', 'Go live with new site', 'review', 'high', 80, 85, 200, 140, '2026-02-01', '2026-03-30', NULL, true, true, 15, NOW(), NOW()),
  
  ('cc0e8400-e29b-41d4-a716-446655440210', 'aa0e8400-e29b-41d4-a716-446655440002', 'cc0e8400-e29b-41d4-a716-446655440200', 2, 'Homepage Redesign', 'New hero and value props', 'done', 'high', 85, 90, 40, 42, '2026-02-01', '2026-02-15', '660e8400-e29b-41d4-a716-446655440002', false, true, 10, NOW(), NOW()),
  ('cc0e8400-e29b-41d4-a716-446655440211', 'aa0e8400-e29b-41d4-a716-446655440002', 'cc0e8400-e29b-41d4-a716-446655440200', 2, 'Product Pages', 'Individual product detail pages', 'done', 'high', 80, 85, 50, 48, '2026-02-10', '2026-02-28', '660e8400-e29b-41d4-a716-446655440005', false, true, 15, NOW(), NOW()),
  ('cc0e8400-e29b-41d4-a716-446655440212', 'aa0e8400-e29b-41d4-a716-446655440002', 'cc0e8400-e29b-41d4-a716-446655440200', 2, 'Blog & Content', 'CMS integration for blog', 'in_progress', 'medium', 65, 70, 40, 30, '2026-02-20', '2026-03-15', '660e8400-e29b-41d4-a716-446655440005', false, false, 10, NOW(), NOW()),
  ('cc0e8400-e29b-41d4-a716-446655440213', 'aa0e8400-e29b-41d4-a716-446655440002', 'cc0e8400-e29b-41d4-a716-446655440200', 2, 'SEO Optimization', 'Meta tags, structured data', 'review', 'medium', 60, 65, 30, 25, '2026-03-01', '2026-03-20', NULL, false, false, 15, NOW(), NOW()),
  ('cc0e8400-e29b-41d4-a716-446655440214', 'aa0e8400-e29b-41d4-a716-446655440002', 'cc0e8400-e29b-41d4-a716-446655440200', 2, 'Performance Tuning', 'Core Web Vitals optimization', 'in_progress', 'high', 75, 80, 40, 20, '2026-03-10', '2026-03-30', '660e8400-e29b-41d4-a716-446655440005', false, false, 20, NOW(), NOW());

-- ============================================
-- TASKS - Internal Dashboard
-- ============================================
INSERT INTO tasks (id, project_id, parent_task_id, hierarchy_level, title, description, status, priority, priority_score, business_value, estimated_hours, actual_hours, start_date, due_date, assignee_id, is_milestone, is_critical_path, risk_score, created_at, updated_at)
VALUES 
  ('cc0e8400-e29b-41d4-a716-446655440300', 'aa0e8400-e29b-41d4-a716-446655440003', NULL, 1, 'Dashboard v1.0', 'Initial release with core metrics', 'in_progress', 'medium', 70, 75, 250, 200, '2026-01-10', '2026-03-15', NULL, true, true, 10, NOW(), NOW()),
  
  ('cc0e8400-e29b-41d4-a716-446655440310', 'aa0e8400-e29b-41d4-a716-446655440003', 'cc0e8400-e29b-41d4-a716-446655440300', 2, 'Data Pipeline Setup', 'ETL from production DB', 'done', 'high', 80, 85, 60, 65, '2026-01-10', '2026-01-30', '660e8400-e29b-41d4-a716-446655440007', false, true, 15, NOW(), NOW()),
  ('cc0e8400-e29b-41d4-a716-446655440311', 'aa0e8400-e29b-41d4-a716-446655440003', 'cc0e8400-e29b-41d4-a716-446655440300', 2, 'Executive Summary View', 'High-level KPIs for leadership', 'done', 'high', 75, 80, 50, 52, '2026-01-25', '2026-02-15', '660e8400-e29b-41d4-a716-446655440001', false, true, 10, NOW(), NOW()),
  ('cc0e8400-e29b-41d4-a716-446655440312', 'aa0e8400-e29b-41d4-a716-446655440003', 'cc0e8400-e29b-41d4-a716-446655440300', 2, 'Sales Analytics', 'Revenue and conversion metrics', 'in_progress', 'medium', 65, 70, 60, 48, '2026-02-01', '2026-03-01', '660e8400-e29b-41d4-a716-446655440007', false, false, 15, NOW(), NOW()),
  ('cc0e8400-e29b-41d4-a716-446655440313', 'aa0e8400-e29b-41d4-a716-446655440003', 'cc0e8400-e29b-41d4-a716-446655440300', 2, 'User Behavior Tracking', 'Engagement and retention metrics', 'in_progress', 'medium', 60, 65, 50, 35, '2026-02-15', '2026-03-10', '660e8400-e29b-41d4-a716-446655440001', false, false, 10, NOW(), NOW()),
  ('cc0e8400-e29b-41d4-a716-446655440314', 'aa0e8400-e29b-41d4-a716-446655440003', 'cc0e8400-e29b-41d4-a716-446655440300', 2, 'Custom Reports', 'User-defined report builder', 'backlog', 'low', 50, 55, 80, 0, NULL, '2026-04-01', NULL, false, false, 20, NOW(), NOW());
