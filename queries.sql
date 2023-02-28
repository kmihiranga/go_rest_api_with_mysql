use springboot_blog_rest_api;

create table tbl_users
(
	id varchar(255) not null primary key,
    email varchar(255) not null,
    password varchar(255) not null,
    first_name varchar(255) not null,
    last_name varchar(255) not null,
    created_at datetime null,
    updated_at datetime null
);

show columns from tbl_users;

select * from tbl_users; 
