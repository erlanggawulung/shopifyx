CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "username" VARCHAR UNIQUE NOT NULL,
  "name" VARCHAR NOT NULL,
  "password" VARCHAR NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT NOW()
);
