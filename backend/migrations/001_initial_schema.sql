-- Complete schema creation for micro-lending platform
-- Version: 001 (Complete)
-- Description: Creates all tables with all required columns

-- Clients table - Main entity storing borrower information
CREATE TABLE IF NOT EXISTS clients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    control_number VARCHAR(20) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    last_name VARCHAR(100) NOT NULL,
    nickname VARCHAR(50),
    date_of_birth DATETIME,
    gender VARCHAR(10),
    religion VARCHAR(50),
    civil_status VARCHAR(20),
    home_address TEXT,
    years_of_residence INTEGER,
    facebook_account VARCHAR(100),
    age INTEGER,
    contact_number VARCHAR(20),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL
);

-- Income information - Financial data for loan assessment
CREATE TABLE IF NOT EXISTS income_info (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    client_id INTEGER NOT NULL,
    family_income_daily DECIMAL(10,2),
    family_income_monthly DECIMAL(10,2),
    total_cost_daily DECIMAL(10,2),
    total_cost_monthly DECIMAL(10,2),
    net_income_daily DECIMAL(10,2),
    net_income_monthly DECIMAL(10,2),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE,
    UNIQUE(client_id)
);

-- Loans table - Core lending transactions
CREATE TABLE IF NOT EXISTS loans (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    client_id INTEGER NOT NULL,
    control_number VARCHAR(20) UNIQUE NOT NULL,
    date_of_release DATETIME NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL,
    ammortization DECIMAL(10,2) NOT NULL,
    terms INTEGER NOT NULL,
    mode VARCHAR(20) DEFAULT 'Weekly',
    outstanding_balance DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) DEFAULT 'Active',
    due_date VARCHAR(20),
    deductions VARCHAR(100),
    amount_release DECIMAL(10,2) NOT NULL,
    payment_period_weeks INTEGER,
    method_of_payment VARCHAR(50),
    credit_history VARCHAR(50),
    recommended_by VARCHAR(100),
    approved_by VARCHAR(100),
    loan_cycle INTEGER,
    recommended_loan_amount DECIMAL(10,2),
    approved_loan_amount DECIMAL(10,2),
    checked_by VARCHAR(100),
    name_ci VARCHAR(100),  -- ADDED: Name CI field
    noted_by VARCHAR(100),
    application_date DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE
);

-- Payments table - Tracks weekly loan payments
CREATE TABLE IF NOT EXISTS payments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    loan_id INTEGER NOT NULL,
    week_number INTEGER NOT NULL,
    payment_date DATETIME,
    amount_due DECIMAL(10,2) NOT NULL,
    amount_paid DECIMAL(10,2),
    status VARCHAR(20) DEFAULT 'Pending',
    payment_method VARCHAR(50),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    FOREIGN KEY (loan_id) REFERENCES loans(id) ON DELETE CASCADE,
    UNIQUE(loan_id, week_number)
);

-- Co-makers table - Loan guarantors
CREATE TABLE IF NOT EXISTS co_makers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    loan_id INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
    address TEXT,
    business VARCHAR(100),
    signature_path VARCHAR(255),
    id_copy_path VARCHAR(255),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    FOREIGN KEY (loan_id) REFERENCES loans(id) ON DELETE CASCADE
);

-- Family members table - Client's household information (WITH ALL COLUMNS)
CREATE TABLE IF NOT EXISTS family_members (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    client_id INTEGER NOT NULL,
    relationship VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    surname VARCHAR(100),
    age INTEGER,
    work VARCHAR(100),
    contact_number VARCHAR(20),
    address TEXT,
    birthday DATETIME,        -- ADDED: Birthday column for spouse
    nickname VARCHAR(50),     -- ADDED: Nickname column for spouse
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE
);

-- Documents table - Client document storage
CREATE TABLE IF NOT EXISTS documents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    client_id INTEGER NOT NULL,
    document_type VARCHAR(50) NOT NULL,
    file_path VARCHAR(255) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE
);

-- Users table - System user authentication
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    is_admin BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL
);

-- Performance indexes for frequently queried fields
CREATE INDEX IF NOT EXISTS idx_loans_client_id ON loans(client_id);
CREATE INDEX IF NOT EXISTS idx_loans_status ON loans(status);
CREATE INDEX IF NOT EXISTS idx_payments_loan_id ON payments(loan_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);
CREATE INDEX IF NOT EXISTS idx_clients_control_number ON clients(control_number);
CREATE INDEX IF NOT EXISTS idx_family_members_client_id ON family_members(client_id);



-- Add missing columns to family_members table
ALTER TABLE family_members ADD COLUMN birthday DATETIME;
ALTER TABLE family_members ADD COLUMN nickname VARCHAR(50);

-- Ensure name_ci column exists in loans table
ALTER TABLE loans ADD COLUMN name_ci VARCHAR(100);


-- Insert initial admin user (password: "admin123")
-- Note: Replace with actual bcrypt hash in production
INSERT OR IGNORE INTO users (username, password_hash, is_admin) 
VALUES ('admin', '$2a$10$r8R2Qk9U9W9p9V9X9Y9Z9A9B9C9D9E9F9G9H9I9J9K9L9M9N9O9P9Q', 1);
