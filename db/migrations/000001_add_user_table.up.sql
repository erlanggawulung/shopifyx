CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE "users" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "username" VARCHAR UNIQUE NOT NULL,
  "name" VARCHAR NOT NULL,
  "password" VARCHAR NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT NOW()
);
