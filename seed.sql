-- Seed Categories
INSERT INTO categories (name, created_at) VALUES 
('Food', NOW()),
('Transport', NOW()),
('Rent', NOW()),
('Entertainment', NOW())
ON CONFLICT DO NOTHING;

-- Seed Budget for Feb 2026
INSERT INTO budgets (month, year, budget_amount, created_at, updated_at) VALUES 
(2, 2026, 12000.00, NOW(), NOW())
ON CONFLICT (month, year) DO UPDATE SET budget_amount = 12000.00, updated_at = NOW();

-- Seed some initial expenses for Feb 2026 (Total: 8000, within 12000 budget)
INSERT INTO expenses (category_id, amount, description, payment_mode, expense_date, created_at) VALUES 
((SELECT id FROM categories WHERE name = 'Rent' LIMIT 1), 5000.00, 'Monthly Rent', 'Bank Transfer', '2026-02-01', NOW()),
((SELECT id FROM categories WHERE name = 'Food' LIMIT 1), 2000.00, 'Grocery shopping', 'UPI', '2026-02-05', NOW()),
((SELECT id FROM categories WHERE name = 'Transport' LIMIT 1), 1000.00, 'Fuel', 'Cash', '2026-02-10', NOW());
