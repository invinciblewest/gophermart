-- +goose Up
-- +goose StatementBegin
CREATE TABLE "users" (
    "id" serial PRIMARY KEY,
    "login" varchar(255) UNIQUE NOT NULL,
    "password" varchar(255) NOT NULL,
    "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "orders" (
    "id" serial PRIMARY KEY,
    "user_id" int NOT NULL,
    "number" varchar(50) UNIQUE NOT NULL,
    "status" varchar(20) NOT NULL,
    "accrual" numeric(10,2),
    "uploaded_at" timestamp DEFAULT (now())
);

CREATE TABLE "withdrawals" (
    "id" serial PRIMARY KEY,
    "user_id" int NOT NULL,
    "order_number" varchar(50) NOT NULL,
    "amount" numeric(10,2) NOT NULL,
    "processed_at" timestamp DEFAULT (now())
);

ALTER TABLE "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "withdrawals" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "withdrawals";
DROP TABLE "orders";
DROP TABLE "users";
-- +goose StatementEnd
