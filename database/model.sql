create database go_grpc;

create table todo (
	id uuid not null primary key,
	title varchar(128) not null,
	description varchar(1024) default null,
	completed bool default false
);