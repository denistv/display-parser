-- +migrate Up
CREATE TABLE models
(
    id          bigserial primary key,
    external_id text unique not null,
    url         text,
    name        text        not null,
    brand_id    bigint
);

CREATE TABLE documents
(
    url  text primary key,
    body text not null
);

-- +migrate Down
