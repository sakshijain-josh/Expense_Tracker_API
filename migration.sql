-- Migration script to remove user_id columns and users table
-- Run this script if you have an existing database with the old schema

-- Drop foreign key constraints first
ALTER TABLE IF EXISTS expenses DROP CONSTRAINT IF EXISTS expenses_user_id_fkey;
ALTER TABLE IF EXISTS categories DROP CONSTRAINT IF EXISTS categories_user_id_fkey;
ALTER TABLE IF EXISTS budgets DROP CONSTRAINT IF EXISTS budgets_user_id_fkey;

-- Drop user_id columns
ALTER TABLE IF EXISTS expenses DROP COLUMN IF EXISTS user_id;
ALTER TABLE IF EXISTS categories DROP COLUMN IF EXISTS user_id;
ALTER TABLE IF EXISTS budgets DROP COLUMN IF EXISTS user_id;

-- Drop unique constraints that included user_id
ALTER TABLE IF EXISTS categories DROP CONSTRAINT IF EXISTS categories_user_id_name_key;
ALTER TABLE IF EXISTS budgets DROP CONSTRAINT IF EXISTS budgets_user_id_month_year_key;

-- Add new unique constraints without user_id
ALTER TABLE IF EXISTS categories ADD CONSTRAINT categories_name_key UNIQUE (name);
ALTER TABLE IF EXISTS budgets ADD CONSTRAINT budgets_month_year_key UNIQUE (month, year);

-- Drop users table if it exists
DROP TABLE IF EXISTS users CASCADE;

-- Drop indexes that included user_id
DROP INDEX IF EXISTS idx_expenses_user_id;
DROP INDEX IF EXISTS idx_categories_user_id;
DROP INDEX IF EXISTS idx_budgets_user_id;
