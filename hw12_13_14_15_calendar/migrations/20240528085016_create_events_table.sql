-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
                id uuid primary key not null,
                title text not null,
                start_stamp timestamp not null,
                end_stamp timestamp not null,
                description text,
                author_id uuid not null,
                hours_before_to_notify integer not null);

CREATE INDEX start_date_idx ON events (start_stamp);
CREATE INDEX author_id_idx ON events (author_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
