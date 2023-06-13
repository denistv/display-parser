-- +migrate Up
CREATE TABLE models
(
    id          bigserial primary key,
    url         text,
    brand    text,
    series text,
    name text,
    year bigint,
    size bigint,
    ppi bigint,
    created_at timestamptz not null
);
COMMENT ON TABLE models IS 'Содержит распаршенные модели мониторов';

CREATE TABLE pages
(
    url  text primary key,
    body text not null,
    created_at timestamptz not null
);
COMMENT ON TABLE pages IS 'Содержит сырые html-страницы с описанием мониторов';

-- +migrate Down

DROP TABLE models;
DROP TABLE pages;