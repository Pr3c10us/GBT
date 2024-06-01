CREATE TABLE IF NOT EXISTS debtors (
    id UUID NOT NULL DEFAULT (uuid_generate_v4()) PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    phone_number VARCHAR(32),
    debt INT NOT NULL DEFAULT (0),
    user_id UUID NOT NULL,
    FOREIGN KEY (user_id)
    REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
)