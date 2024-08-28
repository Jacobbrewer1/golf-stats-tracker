create table round_hit_stats
(
    id             int         not null auto_increment,
    round_stats_id int         not null,
    type           enum ('GREEN', 'FAIRWAY') not null,
    miss           varchar(15) not null,
    count          int         not null,
    primary key (id),
    constraint round_hit_stats_round_stats_id_fk
        foreign key (round_stats_id) references round_stats (id)
);

