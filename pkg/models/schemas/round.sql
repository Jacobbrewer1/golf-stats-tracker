create table round
(
    id       int       not null auto_increment,
    user_id  int       not null,
    tee_time timestamp not null,
    primary key (id),
    constraint round_user_id_fk
        foreign key (user_id) references user (id)
);