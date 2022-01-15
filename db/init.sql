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

create table transactions
(
    id          serial
        constraint transactions_pk
            primary key,
    description varchar                                not null,
    created     timestamp with time zone default now() not null,
    amount      double precision                       not null,
    user_id     bigint                                 not null
        constraint transactions_balance_user_id_fk
            references balance (user_id)
            on delete cascade
);

create unique index transactions_id_uindex
    on transactions (id);