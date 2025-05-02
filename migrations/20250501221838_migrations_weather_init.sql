-- +goose Up
CREATE TABLE users(
    username VARCHAR(256) NOT NULL,
    chat_id BIGINT NOT NULL UNIQUE,
    language VARCHAR(50) NOT NULL,
    PRIMARY KEY (chat_id)
);

CREATE TABLE cities(
    chat_id BIGINT NOT NULL UNIQUE,
    main_city TEXT NOT NULL,
    other_cities TEXT[] DEFAULT '{}',
    FOREIGN KEY (chat_id) REFERENCES users(chat_id)
        ON DELETE CASCADE
);

CREATE TABLE geodata(
    main_city TEXT UNIQUE,
    latitude DOUBLE PRECISION NOT NULL,
    Longitude DOUBLE PRECISION NOT NULL
);
-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
DROP TABLE geodata;
DROP TABLE cities;
DROP TABLE users;
-- +goose StatementBegin
-- +goose StatementEnd
