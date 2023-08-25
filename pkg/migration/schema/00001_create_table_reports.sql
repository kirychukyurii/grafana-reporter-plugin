-- +goose Up
-- +goose StatementBegin
create table reports
(
    id                   integer primary key autoincrement,
    created_at           datetime default current_timestamp,
    updated_at           datetime default current_timestamp,
    deleted_at           datetime,
    state                varchar,
    recipients           varchar,
    reply_to             varchar,
    message              varchar,
    orientation          varchar,
    layout               varchar,
    enable_dashboard_url boolean  default 0,
    enable_csv           boolean  default 0,
    scale_factor         integer  default 1
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table reports;
-- +goose StatementEnd
