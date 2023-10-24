-- +goose Up
-- +goose StatementBegin
create table reports
(
    id         integer primary key autoincrement,
    org_id     integer,
    created_at datetime default current_timestamp,
    updated_at datetime default current_timestamp,
    deleted_at datetime,
    state      varchar,
    dashboards text,
    recipients text,
    reply_to   varchar,
    message    varchar,
    options    text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table reports;
-- +goose StatementEnd
