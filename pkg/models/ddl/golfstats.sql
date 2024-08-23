create table user
(
    id         int auto_increment
        primary key,
    name       varchar(50)  not null,
    username   varchar(100) not null,
    password   text         not null,
    last_login datetime     null,
    constraint user_username_uindex
        unique (username)
);

create table round
(
    id       int auto_increment
        primary key,
    user_id  int                                   not null,
    tee_time timestamp default current_timestamp() not null on update current_timestamp(),
    constraint round_user_id_fk
        foreign key (user_id) references user (id)
);

create table course
(
    id       int auto_increment
        primary key,
    round_id int          not null,
    name     varchar(255) not null,
    constraint course_round_id_fk
        foreign key (round_id) references round (id)
);

create table course_details
(
    id                int auto_increment
        primary key,
    course_id         int           not null,
    marker            varchar(255)  null,
    slope             int           not null,
    course_rating     decimal(4, 1) not null,
    front_nine_par    int           not null,
    back_nine_par     int           not null,
    total_par         int           not null,
    front_nine_yards  int           not null,
    back_nine_yards   int           not null,
    total_yards       int           not null,
    front_nine_meters int           not null,
    back_nine_meters  int           not null,
    total_meters      int           not null,
    constraint course_details_course_id_fk
        foreign key (course_id) references course (id)
);

create table hole
(
    id                int auto_increment
        primary key,
    course_details_id int not null,
    number            int not null,
    par               int not null,
    stroke            int not null,
    distance_yards    int not null,
    distance_meters   int not null,
    constraint hole_course_details_id_fk
        foreign key (course_details_id) references course_details (id)
);

create table hole_stats
(
    id           int                                            not null
        primary key,
    hole_id      int                                            not null,
    score        int                                            not null,
    fairway_hit  enum ('HIT', 'LEFT', 'RIGHT', 'SHORT', 'LONG') not null,
    green_hit    enum ('HIT', 'LEFT', 'RIGHT', 'SHORT', 'LONG') not null,
    pin_location varchar(100)                                   not null,
    putts        int                                            not null,
    penalties    int                                            not null,
    constraint hole_stats_hole_id_fk
        foreign key (hole_id) references hole (id)
);

