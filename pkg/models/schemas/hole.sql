create table hole
(
    id                int not null auto_increment,
    course_details_id int not null,
    number            int not null,
    par               int not null,
    stroke            int not null,
    distance_yards    int not null,
    distance_meters   int not null,
    primary key (id),
    constraint hole_course_details_id_fk
        foreign key (course_details_id) references course_details (id)
);