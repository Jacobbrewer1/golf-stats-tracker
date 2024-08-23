create table course
(
    id       int          not null auto_increment,
    round_id int          not null,
    name     varchar(255) not null,
    primary key (id)
);
