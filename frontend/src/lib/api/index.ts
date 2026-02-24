// Xephyr API Client
// Main export for all API modules

// Client and types
export * from './client';
export * from './types';

// API modules
export * from './priority';
export * from './health';
export * from './nudge';
export * from './progress';
export * from './dependency';
export * from './assignment';
export * from './scenario';
export * from './workload';
export * from './projects';
export * from './tasks';
export * from './users';
export * from './skills';

// Default export
import { apiClient } from './client';
export default apiClient;
