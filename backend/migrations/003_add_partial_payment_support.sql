-- Add columns for partial payment support and remaining balance tracking
ALTER TABLE payments ADD COLUMN remaining_balance DECIMAL(10,2) DEFAULT 0;

-- Remove the unique constraint to allow multiple partial payments per week
DROP INDEX IF EXISTS idx_payments_loan_week_unique;

-- Create new indexes for better performance
CREATE INDEX IF NOT EXISTS idx_payments_partial ON payments(loan_id, week_number, is_partial);
CREATE INDEX IF NOT EXISTS idx_payments_remaining_balance ON payments(loan_id, week_number, remaining_balance);

-- Update existing payments to mark them as non-partial and set remaining balance
UPDATE payments SET is_partial = FALSE, remaining_balance = 0;
