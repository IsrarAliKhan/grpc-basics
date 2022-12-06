CREATE TABLE item
(
    id          SERIAL       NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP,
    deleted_at  TIMESTAMP,

    name        VARCHAR(100) NOT NULL,
    price       INT NOT NULL,
    quantity    INT NOT NULL
);

CREATE UNIQUE INDEX item_id_uindex
    ON item (id);

ALTER TABLE item
    ADD CONSTRAINT item_pk
        PRIMARY KEY (id);


CREATE TABLE users
(
    id          SERIAL       NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP,
    deleted_at  TIMESTAMP,  

    username    VARCHAR(100) NOT NULL,
    password    VARCHAR(100) NOT NULL,
    role        VARCHAR(100) NOT NULL
);

CREATE UNIQUE INDEX users_id_uindex
    ON users (id);

ALTER TABLE users
    ADD CONSTRAINT users_pk
        PRIMARY KEY (id);
