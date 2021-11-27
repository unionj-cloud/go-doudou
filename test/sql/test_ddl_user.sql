create table ddl_user
(
    id         int auto_increment
        primary key,
    name       varchar(255) default 'jack'            not null,
    phone      varchar(255) default '13552053960'     not null comment 'mobile phone',
    age        int                                    not null,
    no         int                                    not null,
    school     varchar(255) default 'harvard'         null comment 'school',
    is_student tinyint                                not null,
    create_at  datetime     default CURRENT_TIMESTAMP null,
    update_at  datetime     default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    delete_at  datetime                               null,
    constraint no_idx
        unique (no)
);

create index age_idx
    on ddl_user (age);

create index name_phone_idx
    on ddl_user (phone, name);

create table ddl_book
(
    id           int auto_increment
        primary key,
    name         varchar(45) null,
    user_id      int         null,
    publisher_id int         null,
    constraint fk_user
        foreign key (user_id) references ddl_user (id)
        ON DELETE CASCADE ON UPDATE NO ACTION
);

create table ddl_publisher
(
    id   int auto_increment
        primary key,
    name varchar(45) null
);