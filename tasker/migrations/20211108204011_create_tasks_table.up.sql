CREATE TABLE tasks
(
    id          serial primary key,
    public_id   text,
    assignee_id text,
    description text,
    fee         integer,
    reward      integer,
    status      text
);
