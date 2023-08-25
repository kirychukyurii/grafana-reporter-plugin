-- +goose Up
-- +goose StatementBegin
create table report_history
(
    id                 integer primary key autoincrement,
    created_at         datetime default current_timestamp,
    state              varchar,
    report_id          integer,
    report_schedule_id integer,
    recipients         varchar
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table report_history;
-- +goose StatementEnd
