CREATE TABLE "products" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "name" VARCHAR NOT NULL,
  "price" INT NOT NULL,
  "image_url" VARCHAR NOT NULL,
  "stock" INT NOT NULL,
  "condition" VARCHAR NOT NULL,
  "tags" TEXT NOT NULL,
  "is_purchaseable" BOOLEAN NOT NULL,
  "created_by" UUID NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT NOW()
);
