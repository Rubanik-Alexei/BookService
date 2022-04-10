create database if not exists bookshop;
use bookshop;
create table if not exists books (
	book_id int auto_increment primary key,
    name varchar(255) not null
);
create table if not exists authors(
	author_id int auto_increment primary key,
    name varchar(255) not null
);

create table if not exists book_author(
	book_id int not null,
    author_id int not null,
    foreign key(book_id) references books (book_id) on delete restrict on update cascade,
    foreign key(author_id) references authors (author_id) on delete restrict on update cascade,
    primary key(book_id,author_id)
);

insert into authors(name) values("Dave Cheney");
insert into authors(name) values("Rob Pike");
insert into books(name) values ("Golang");
insert into book_author values(1,1);
insert into book_author values(1,2);
insert into authors(name) values("Martin Kleppman");
insert into books(name) values ("Designing Data-Intensive Applications");
insert into book_author values(2,3);
insert into books(name) values("Golang");
insert into authors(name) values("OtherGolangAuthor");
insert into book_author values(3,4);
insert into books(name) values ("Gopher's Guide");
insert into book_author values(4,1)