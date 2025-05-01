-- +goose Up
CREATE TABLE users(
    username VARCHAR(256) NOT NULL,
    chat_id BIGINT NOT NULL
);

CREATE TABLE cities(
    chat_id BIGINT NOT NULL,
    main_city TEXT,
    other_cities TEXT[]
);

CREATE TABLE geodata(
    main_city TEXT,
    latitude TEXT,
    Longitude TEXT
);
-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
DROP TABLE geodata;
DROP TABLE cities;
DROP TABLE users;
-- +goose StatementBegin
-- +goose StatementEnd
