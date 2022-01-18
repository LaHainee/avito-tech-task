--|------------------Balance------------------|--
create table balance
(
    id      serial
        constraint balance_pk
            primary key,
    user_id bigint           not null,
    balance double precision not null
);

create unique index balance_id_uindex
    on balance (id);

create unique index balance_user_id_uindex
    on balance (user_id);

create index balance_user_id on balance using hash (user_id);
--|------------------Balance------------------|--


--|------------------Transactions------------------|--
create type operation_type as
    enum ('write_off', 'add', 'transfer');

create table transactions
(
    id             serial
        constraint transactions_pk
            primary key,
    operation_type operation_type not null,
    sender         bigint
        constraint transactions_balance_user_id_fk_2
            references balance (user_id)
            on delete cascade,
    receiver       bigint
        constraint transactions_balance_user_id_fk
            references balance (user_id)
            on delete cascade,
    amount         double precision,
    created        timestamp with time zone default now()
);

create unique index transactions_id_uindex
    on transactions (id);

create index transactions_sender_operation on transactions (sender, operation_type);
create index transactions_sender_created on transactions (sender, created);
--|------------------Transactions------------------|--