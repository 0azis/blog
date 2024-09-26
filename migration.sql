CREATE TABLE
    users (
        id smallint not null auto_increment,
        firstName varchar(255) not null,
        lastName varchar(255) not null,
        nick varchar(255) not null,
        password varchar(255) not null,
        avatar varchar(255) not null,
        description varchar(255) not null,
        primary key (id)
    );

CREATE TABLE
    posts (
        id smallint not null auto_increment,
        primary key (id)
    );

CREATE TABLE
    subscribers (
        user_id smallint not null,
        subscriber_id smallint not null,
        foreign key (user_id) references users (id) on delete cascade,
        foreign key (subscriber_id) references users (id) on delete cascade
    );

CREATE TABLE
    followers (
        user_id smallint not null,
        follower_id smallint not null,
        foreign key (user_id) references users (id) on delete cascade,
        foreign key (follower_id) references users (id) on delete cascade
    );

CREATE TABLE
    comments (
        id smallint not null auto_increment,
        post_id smallint not null,
        user_id smallint not null,
        comment_text varchar(255) not null,
        foreign key (user_id) references users (id) on delete no action,
        foreign key (post_id) references posts (id) on delete cascade,
        primary key (id)
    );

CREATE TABLE
    likes (
        post_id smallint not null,
        user_id smallint not null,
        foreign key (user_id) references users (id) on delete no action,
        foreign key (post_id) references posts (id) on delete cascade
    );

CREATE TABLE
    views (
        post_id smallint not null,
        user_id smallint not null,
        foreign key (user_id) references users (id) on delete no action,
        foreign key (post_id) references posts (id) on delete cascade,
    );