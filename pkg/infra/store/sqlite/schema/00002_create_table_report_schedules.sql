-- +goose Up
-- +goose StatementBegin
create table report_schedules
(
    id            integer primary key autoincrement,
    created_at    datetime default current_timestamp,
    updated_at    datetime default current_timestamp,
    deleted_at    datetime,
    name          varchar,
    active        boolean,
    report_id     integer,
    start_date    datetime,
    end_date      datetime,
    timezone      varchar,
    interval      varchar,
    workdays_only boolean
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table report_schedules;
-- +goose StatementEnd
