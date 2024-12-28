CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    phone_number TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS bills (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    amount DECIMAL(10, 2) NOT NULL,
    due_date TIMESTAMP NOT NULL,
    status TEXT NOT NULL DEFAULT 'unpaid'
);

CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    bill_id INTEGER REFERENCES bills(id),
    amount DECIMAL(10, 2) NOT NULL,
    payment_date TIMESTAMP NOT NULL,
    status TEXT NOT NULL,
    reference TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    message TEXT NOT NULL,
    sent_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS settings (
    id SERIAL PRIMARY KEY,
    payment_due_days INTEGER NOT NULL,
    late_fee_percentage DECIMAL(5, 2) NOT NULL
);

-- Insert default settings
INSERT INTO settings (id, payment_due_days, late_fee_percentage)
VALUES (1, 30, 5.00)
ON CONFLICT (id) DO NOTHING;