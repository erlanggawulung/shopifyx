CREATE TABLE "bank_accounts" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id" UUID NOT NULL,
  "bank_name" VARCHAR NOT NULL,
  "bank_account_name" VARCHAR NOT NULL,
  "bank_account_number" VARCHAR NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT NOW()
);
