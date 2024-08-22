create table club
(
    id          int auto_increment
        primary key,
    name        varchar(255) not null,
    address1    varchar(255) not null,
    address2    varchar(255) not null,
    address3    varchar(255) not null,
    address4    varchar(255) not null,
    postal_code varchar(8)   not null,
    email       varchar(255) null,
    phone       varchar(30) null
);

create table course
(
    id      int auto_increment
        primary key,
    club_id int          not null,
    name    varchar(255) not null,
    constraint course_club_id_fk
        foreign key (club_id) references club (id)
);

create table course_details
(
    id                int auto_increment
        primary key,
    course_id         int           not null,
    marker            varchar(255) null,
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

create table round
(
    id        int auto_increment
        primary key,
    course_id int                                   not null,
    tee_time  timestamp default current_timestamp() not null on update current_timestamp ()
);

create table user
(
    id         int auto_increment
        primary key,
    name       varchar(50)  not null,
    username   varchar(100) not null,
    password   text         not null,
    last_login datetime null,
    constraint user_username_uindex
        unique (username)
);

