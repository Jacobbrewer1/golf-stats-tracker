create table hole_stats
(
    id           int          not null auto_increment,
    hole_id      int          not null,
    score        int          not null,
    fairway_hit  enum ('HIT', 'LEFT', 'RIGHT', 'SHORT', 'LONG', 'NOT_APPLICABLE') not null,
    green_hit    enum ('HIT', 'LEFT', 'RIGHT', 'SHORT', 'LONG') not null,
    pin_location varchar(100) not null,
    putts        int          not null,
    penalties    int          not null,
    primary key (id),
    constraint hole_stats_hole_id_fk
        foreign key (hole_id) references hole (id)
);