drop table if exists "requests_for_delivery" CASCADE;
drop table if exists "rocket_flights" CASCADE;
drop table if exists "flights_requests" CASCADE;
drop table if exists "users" CASCADE;

create table "users"
(
    user_id         integer not null 
                    constraint user_pk primary key,
    login           varchar(30),
    password        varchar(30),
    is_admin        boolean
);

create table rocket_flights
(
    flight_id       integer not null 
                    constraint flight_pk primary key,
    user_id         integer 
                    constraint flight_creator_user_id_fk references "users",
    moderator_id    integer 
                    constraint flight_moderator_user_id_fk references "users", 
    status          varchar(20),
    created_at      timestamp,
    formed_at       timestamp,
    confirmed_at    timestamp,
    flight_date     timestamp,
    payload         integer,
    price           float,
    title           varchar(100),
    site_number     integer
);

create table requests_for_delivery
(
    request_id          integer not null 
                        constraint request_pk primary key,
    is_available        boolean,
    img_url             TEXT,
    title               varchar(100),
    load_capacity       float,
    description         TEXT,
    detailed_desc       TEXT,
    desired_price       float,
    flight_date_start   timestamp,
    flight_date_end     timestamp
);

create table flights_requests
(
    flight_id       integer 
                    constraint flight_request_flight_flight_id_fk references rocket_flights,
    request_id      integer 
                    constraint request_flight_request_request_id_fk references requests_for_delivery,
                primary key (flight_id, request_id)
);

alter table "users" owner to admin;
alter table rocket_flights owner to admin;
alter table requests_for_delivery owner to admin;
alter table flights_requests owner to admin;
