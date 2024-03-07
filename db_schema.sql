create table products(
	id serial not null primary key,
	name varchar(50) not null,
	price float8 not null,
	category_id bigint not null
)

alter table products
	add column stock integer not null
	
insert into products(name,price,stock,category_id) values
(
	'adidek dudi',
	20.00,
	10,
	1
),
(
	'nikes kurmili',
	80.00,
	100,
	1
)
select * from products


create table categories(
	id serial8 not null primary key,
	name varchar(50) not null
)


insert into categories(name) values
(
	'cepatu'
)

select * from categories


