alter table categories
    drop column deleted_at;

alter table transactions
    drop column deleted_at;

alter table user_tokens
    drop column deleted_at;

alter table users
    drop column deleted_at;

alter table wallets
    drop column deleted_at;