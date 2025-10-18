-- Add columns for partial payment support
ALTER TABLE payments ADD COLUMN is_partial BOOLEAN DEFAULT FALSE;
ALTER TABLE payments ADD COLUMN completes_week BOOLEAN DEFAULT FALSE;

-- Update existing payments to mark them as non-partial
UPDATE payments SET is_partial = FALSE;

-- Create index for better performance on partial payment queries
CREATE INDEX IF NOT EXISTS idx_payments_partial ON payments(loan_id, week_number, status) WHERE status = 'Partial';
