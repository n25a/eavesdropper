create table test_one
(
    id int not null primary key auto_increment,
    data INT(11) not null,
    state varchar(256) not null,
    created_at timestamp default current_timestamp not null,
    updated_at timestamp null
);

create table test_two
(
    id int not null primary key auto_increment,
    var_one varchar(256) not null,
    status INT(11) not null,
    created_at timestamp default current_timestamp not null,
    updated_at timestamp null
);

create table test_three
(
    id int not null primary key auto_increment,
    var_one varchar(256) not null,
    var_three double not null,
    created_at timestamp default current_timestamp not null,
    updated_at timestamp null
);
