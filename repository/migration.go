package repository

import "fmt"

// MigrateSchema removes user_id columns and users table from existing database
// This should be run once to migrate from the old schema to the new schema
func MigrateSchema() error {
	queries := []string{
		// Drop foreign key constraints first
		`ALTER TABLE IF EXISTS expenses DROP CONSTRAINT IF EXISTS expenses_user_id_fkey`,
		`ALTER TABLE IF EXISTS categories DROP CONSTRAINT IF EXISTS categories_user_id_fkey`,
		`ALTER TABLE IF EXISTS budgets DROP CONSTRAINT IF EXISTS budgets_user_id_fkey`,

		// Drop user_id columns
		`ALTER TABLE IF EXISTS expenses DROP COLUMN IF EXISTS user_id`,
		`ALTER TABLE IF EXISTS categories DROP COLUMN IF EXISTS user_id`,
		`ALTER TABLE IF EXISTS budgets DROP COLUMN IF EXISTS user_id`,

		// Drop unique constraints that included user_id
		`ALTER TABLE IF EXISTS categories DROP CONSTRAINT IF EXISTS categories_user_id_name_key`,
		`ALTER TABLE IF EXISTS budgets DROP CONSTRAINT IF EXISTS budgets_user_id_month_year_key`,

		// Add new unique constraints without user_id (only if they don't exist)
		`DO $$ 
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'categories_name_key') THEN
				ALTER TABLE categories ADD CONSTRAINT categories_name_key UNIQUE (name);
			END IF;
		END $$`,
		`DO $$ 
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'budgets_month_year_key') THEN
				ALTER TABLE budgets ADD CONSTRAINT budgets_month_year_key UNIQUE (month, year);
			END IF;
		END $$`,

		// Drop users table if it exists
		`DROP TABLE IF EXISTS users CASCADE`,

		// Drop indexes that included user_id
		`DROP INDEX IF EXISTS idx_expenses_user_id`,
		`DROP INDEX IF EXISTS idx_categories_user_id`,
		`DROP INDEX IF EXISTS idx_budgets_user_id`,
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			// Log error but continue - some constraints/columns might not exist
			fmt.Printf("Warning: Migration query failed (may not exist): %v\n", err)
		}
	}

	return nil
}
