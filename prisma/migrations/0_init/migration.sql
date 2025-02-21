-- CreateSchema
CREATE SCHEMA IF NOT EXISTS "account";

-- CreateSchema
CREATE SCHEMA IF NOT EXISTS "payment";

-- CreateSchema
CREATE SCHEMA IF NOT EXISTS "product";

-- CreateEnum
CREATE TYPE "account"."role" AS ENUM ('ADMIN', 'USER');

-- CreateEnum
CREATE TYPE "account"."gender" AS ENUM ('MALE', 'FEMALE', 'OTHER');

-- CreateEnum
CREATE TYPE "payment"."status" AS ENUM ('PENDING', 'SUCCESS', 'CANCELLED', 'FAILED');

-- CreateEnum
CREATE TYPE "payment"."payment_method" AS ENUM ('CASH', 'MOMO', 'VNPAY');

-- CreateTable
CREATE TABLE "account"."base" (
    "id" BIGINT NOT NULL,
    "username" VARCHAR(50) NOT NULL,
    "password" VARCHAR(100) NOT NULL,
    "role" "account"."role" NOT NULL,

    CONSTRAINT "base_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "account"."user" (
    "id" BIGINT NOT NULL,
    "email" VARCHAR(255) NOT NULL,
    "phone" TEXT NOT NULL,
    "gender" "account"."gender" NOT NULL,
    "full_name" VARCHAR(100) NOT NULL DEFAULT '',
    "default_address_id" BIGINT NOT NULL,

    CONSTRAINT "user_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "account"."admin" (
    "id" BIGINT NOT NULL,

    CONSTRAINT "admin_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "account"."address" (
    "id" BIGINT NOT NULL,
    "user_id" BIGINT NOT NULL,
    "address" TEXT NOT NULL,
    "city" TEXT NOT NULL,
    "province" TEXT NOT NULL,
    "country" TEXT NOT NULL,
    "postal_code" TEXT NOT NULL,

    CONSTRAINT "address_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "account"."cart" (
    "id" BIGINT NOT NULL,

    CONSTRAINT "cart_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "account"."item_on_cart" (
    "cart_id" BIGINT NOT NULL,
    "product_model_id" BIGINT NOT NULL,
    "quantity" BIGINT NOT NULL,

    CONSTRAINT "item_on_cart_pkey" PRIMARY KEY ("cart_id","product_model_id")
);

-- CreateTable
CREATE TABLE "product"."brand" (
    "id" BIGINT NOT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT NOT NULL,

    CONSTRAINT "brand_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "product"."model" (
    "id" BIGINT NOT NULL,
    "brand_id" BIGINT NOT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "list_price" BIGINT NOT NULL,
    "date_manufactured" TIMESTAMPTZ(3) NOT NULL,

    CONSTRAINT "model_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "product"."base" (
    "id" BIGSERIAL NOT NULL,
    "serial_id" TEXT NOT NULL,
    "product_model_id" BIGINT NOT NULL,
    "date_created" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "date_updated" TIMESTAMPTZ(3) NOT NULL,

    CONSTRAINT "base_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "product"."sale" (
    "id" BIGINT NOT NULL,
    "tag_name" TEXT,
    "product_model_id" BIGINT,
    "date_started" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "date_ended" TIMESTAMPTZ(3),
    "quantity" BIGINT NOT NULL,
    "used" BIGINT NOT NULL DEFAULT 0,
    "is_active" BOOLEAN NOT NULL DEFAULT true,
    "discount_percent" BIGINT,
    "discount_price" BIGINT,

    CONSTRAINT "sale_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "product"."tag_on_product" (
    "product_model_id" BIGINT NOT NULL,
    "tag_name" TEXT NOT NULL,

    CONSTRAINT "tag_on_product_pkey" PRIMARY KEY ("product_model_id","tag_name")
);

-- CreateTable
CREATE TABLE "product"."tag" (
    "tag_name" TEXT NOT NULL,
    "description" TEXT NOT NULL DEFAULT '',

    CONSTRAINT "tag_pkey" PRIMARY KEY ("tag_name")
);

-- CreateTable
CREATE TABLE "payment"."product_on_payment" (
    "id" BIGINT NOT NULL,
    "payment_id" BIGINT NOT NULL,
    "product_serial_id" TEXT NOT NULL,
    "quantity" BIGINT NOT NULL,
    "price" BIGINT NOT NULL,
    "total_price" BIGINT NOT NULL,

    CONSTRAINT "product_on_payment_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "payment"."base" (
    "id" BIGINT NOT NULL,
    "user_id" BIGINT NOT NULL,
    "address" TEXT NOT NULL,
    "payment_method" "payment"."payment_method" NOT NULL,
    "total" BIGINT NOT NULL,
    "status" "payment"."status" NOT NULL,
    "date_created" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "base_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "product"."image" (
    "brand_id" BIGINT,
    "product_model_id" BIGINT,
    "url" TEXT NOT NULL,

    CONSTRAINT "image_pkey" PRIMARY KEY ("url")
);

-- CreateIndex
CREATE UNIQUE INDEX "user_email_key" ON "account"."user"("email");

-- CreateIndex
CREATE UNIQUE INDEX "user_phone_key" ON "account"."user"("phone");

-- CreateIndex
CREATE UNIQUE INDEX "user_default_address_id_key" ON "account"."user"("default_address_id");

-- CreateIndex
CREATE UNIQUE INDEX "base_serial_id_key" ON "product"."base"("serial_id");

-- AddForeignKey
ALTER TABLE "account"."user" ADD CONSTRAINT "user_id_fkey" FOREIGN KEY ("id") REFERENCES "account"."base"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "account"."admin" ADD CONSTRAINT "admin_id_fkey" FOREIGN KEY ("id") REFERENCES "account"."base"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "account"."address" ADD CONSTRAINT "address_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "account"."user"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "account"."cart" ADD CONSTRAINT "cart_id_fkey" FOREIGN KEY ("id") REFERENCES "account"."user"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "account"."item_on_cart" ADD CONSTRAINT "item_on_cart_cart_id_fkey" FOREIGN KEY ("cart_id") REFERENCES "account"."cart"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "account"."item_on_cart" ADD CONSTRAINT "item_on_cart_product_model_id_fkey" FOREIGN KEY ("product_model_id") REFERENCES "product"."model"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "product"."model" ADD CONSTRAINT "model_brand_id_fkey" FOREIGN KEY ("brand_id") REFERENCES "product"."brand"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "product"."base" ADD CONSTRAINT "base_product_model_id_fkey" FOREIGN KEY ("product_model_id") REFERENCES "product"."model"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "product"."sale" ADD CONSTRAINT "sale_tag_name_fkey" FOREIGN KEY ("tag_name") REFERENCES "product"."tag"("tag_name") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "product"."sale" ADD CONSTRAINT "sale_product_model_id_fkey" FOREIGN KEY ("product_model_id") REFERENCES "product"."model"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "product"."tag_on_product" ADD CONSTRAINT "tag_on_product_product_model_id_fkey" FOREIGN KEY ("product_model_id") REFERENCES "product"."model"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "product"."tag_on_product" ADD CONSTRAINT "tag_on_product_tag_name_fkey" FOREIGN KEY ("tag_name") REFERENCES "product"."tag"("tag_name") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "payment"."product_on_payment" ADD CONSTRAINT "product_on_payment_payment_id_fkey" FOREIGN KEY ("payment_id") REFERENCES "payment"."base"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "payment"."base" ADD CONSTRAINT "base_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "account"."user"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "product"."image" ADD CONSTRAINT "image_brand_id_fkey" FOREIGN KEY ("brand_id") REFERENCES "product"."brand"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "product"."image" ADD CONSTRAINT "image_product_model_id_fkey" FOREIGN KEY ("product_model_id") REFERENCES "product"."model"("id") ON DELETE CASCADE ON UPDATE CASCADE;

