create table user
(
    id         int          not null auto_increment,
    name       varchar(50)  not null,
    username   varchar(100) not null,
    password   text         not null,
    last_login datetime null,
    primary key (id)
);