alter table categories
    add deleted_at timestamp;

alter table transactions
    add deleted_at timestamp;

alter table user_tokens
    add deleted_at timestamp;

alter table users
    add deleted_at timestamp;

alter table wallets
    add deleted_at timestamp;