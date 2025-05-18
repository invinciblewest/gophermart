-- +goose Up
-- +goose StatementBegin
CREATE TABLE "users" (
    "id" serial PRIMARY KEY,
    "login" varchar(255) UNIQUE NOT NULL,
    "password" varchar(255) NOT NULL,
    "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "orders" (
    "id" serial PRIMARY KEY,
    "user_id" int NOT NULL REFERENCES "users" ("id"),
    "number" varchar(50) UNIQUE NOT NULL,
    "status" varchar(20) NOT NULL,
    "accrual" int,
    "uploaded_at" timestamptz DEFAULT (now())
);

CREATE TABLE "withdrawals" (
    "id" serial PRIMARY KEY,
    "user_id" int NOT NULL REFERENCES "users" ("id"),
    "order_number" varchar(50) NOT NULL,
    "amount" int NOT NULL,
    "processed_at" timestamptz DEFAULT (now())
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "withdrawals";
DROP TABLE "orders";
DROP TABLE "users";
-- +goose StatementEnd
