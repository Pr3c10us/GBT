ALTER TABLE debtors
    ADD COLUMN  debt INT NOT NULL DEFAULT (0);

ALTER TABLE creditors
    ADD COLUMN  credit INT NOT NULL DEFAULT (0);