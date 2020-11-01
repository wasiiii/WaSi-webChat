CREATE user tester with password 'postgres';

CREATE TABLE chat
(
    sender character varying COLLATE pg_catalog."default",
    receiver character varying COLLATE pg_catalog."default",
    date character varying COLLATE pg_catalog."default",
    msg character varying COLLATE pg_catalog."default",
    img character varying COLLATE pg_catalog."default"
);




CREATE TABLE userinfo
(
    name character varying COLLATE pg_catalog."default",
    email character varying COLLATE pg_catalog."default",
    phone character varying COLLATE pg_catalog."default",
    password character varying COLLATE pg_catalog."default"
);
