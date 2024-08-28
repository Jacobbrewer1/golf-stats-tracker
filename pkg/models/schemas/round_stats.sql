create table round_stats
(
    id               int           not null auto_increment,
    round_id         int           not null,
    avg_fairways_hit decimal(5, 2) not null,
    avg_greens_hit   decimal(5, 2) not null,
    avg_putts        decimal(5, 2) not null,
    penalties        int           not null,
    avg_par_3        decimal(5, 2) not null,
    avg_par_4        decimal(5, 2) not null,
    avg_par_5        decimal(5, 2) not null,
    primary key (id),
    constraint round_stats_round_id_fk
        foreign key (round_id) references round (id)
);

