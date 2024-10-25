CREATE TABLE
    users (
        id smallint not null auto_increment,
        email varchar(255) not null,
        username varchar(255) not null,
        password varchar(255) not null,
        name varchar(255),
        avatar varchar(255),
        description varchar(255),
        primary key (id)
    );

CREATE TABLE
    posts (
        id smallint not null auto_increment,
        user_id smallint not null,
        category_id smallint,
        title varchar(255) not null default '',
        preview varchar(255),
        date datetime default current_timestamp,
        content text not null default (''),
        public bool default false not null,
        foreign key (user_id) references users (id) on delete cascade,
        primary key (id)
    );

create table relations (
	id_1 smallint not null,
  id_2 smallint not null,
  foreign key (id_1) references users(id) on delete cascade,
  foreign key (id_2) references users(id) on delete cascade,
  primary key (id_1, id_2)
);

create table tags (
    post_id smallint not null,
    tag varchar(255) not null,
    foreign key (post_id) references posts(id) on delete cascade
);

CREATE TABLE
    comments (
        id smallint not null auto_increment,
        post_id smallint not null,
        user_id smallint not null,
        text varchar(255) not null,
        foreign key (user_id) references users (id) on delete no action,
        foreign key (post_id) references posts (id) on delete cascade,
        primary key (id)
    );

-- CREATE TABLE
--     likes (
--         post_id smallint not null,
--         user_id smallint not null,
--         foreign key (user_id) references users (id) on delete no action,
--         foreign key (post_id) references posts (id) on delete cascade
--     );

-- CREATE TABLE
--     views (
--         post_id smallint not null,
--         user_id smallint not null,
--         foreign key (user_id) references users (id) on delete no action,
--         foreign key (post_id) references posts (id) on delete cascade,
--     );
