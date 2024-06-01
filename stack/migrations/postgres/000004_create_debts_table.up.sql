CREATE TABLE IF NOT EXISTS debts (
    amount INT NOT NULL,
    reason VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    debtor_id UUID NOT NULL ,
    FOREIGN KEY (debtor_id)
        REFERENCES debtors (id) ON DELETE CASCADE ON UPDATE CASCADE
)