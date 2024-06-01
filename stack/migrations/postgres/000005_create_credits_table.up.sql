CREATE TABLE IF NOT EXISTS credits (
    amount INT NOT NULL,
    reason VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    creditor_id UUID NOT NULL ,
    FOREIGN KEY (creditor_id)
        REFERENCES creditors (id) ON DELETE CASCADE ON UPDATE CASCADE
)