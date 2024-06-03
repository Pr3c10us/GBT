ALTER TABLE debtors
    ADD CONSTRAINT unique_debtor_name UNIQUE (name);

ALTER TABLE creditors
    ADD CONSTRAINT unique_creditor_name UNIQUE (name);