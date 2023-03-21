CREATE TABLE accounts (
    id bigserial PRIMARY KEY,
    owner varchar(100) NOT NULL,
    balance bigint NOT NULL,
    currency varchar(50) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE entries (
    id bigserial PRIMARY KEY,
    account_id bigint NOT NULL,
    amount bigint NOT NULL,
    created_at timestamptz NOT NULL DEFAULT (now()),
    CONSTRAINT acc_index FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE NO ACTION
);

CREATE TABLE transfers (
    id bigserial PRIMARY KEY,
    from_account_id bigint NOT NULL,
    to_account_id bigint NOT NULL,
    amount bigint NOT NULL,
    created_at timestamptz NOT NULL DEFAULT (now()),
    CONSTRAINT from_acc_index FOREIGN KEY (from_account_id) REFERENCES accounts(id) ON DELETE NO ACTION,
    CONSTRAINT to_acc_index FOREIGN KEY (to_account_id) REFERENCES accounts (id) ON DELETE NO ACTION
);