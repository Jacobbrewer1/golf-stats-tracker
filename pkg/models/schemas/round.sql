create table round
(
    id        int       not null auto_increment,
    course_id int       not null,
    tee_time  timestamp not null,
    primary key (id)
);