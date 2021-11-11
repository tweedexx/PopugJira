CREATE TABLE users
(
    id        serial primary key,
    public_id text,
    role      text,
    name      text,
    password  text,
    email     text
);
