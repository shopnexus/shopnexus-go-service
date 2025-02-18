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
    "id" BYTEA NOT NULL,
    "username" VARCHAR(50) NOT NULL,
    "password" VARCHAR(100) NOT NULL,
    "role" "account"."role" NOT NULL,

    CONSTRAINT "base_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "account"."user" (
    "id" BYTEA NOT NULL,
    "email" VARCHAR(255) NOT NULL,
    "phone" TEXT NOT NULL,
    "gender" "account"."gender" NOT NULL,
    "full_name" VARCHAR(100),
    "default_address_id" BYTEA NOT NULL,

    CONSTRAINT "user_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "account"."shop" (
    "id" BYTEA NOT NULL,

    CONSTRAINT "shop_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "account"."address" (
    "id" BYTEA NOT NULL,
    "user_id" BYTEA NOT NULL,
    "address" TEXT NOT NULL,
    "city" TEXT NOT NULL,
    "province" TEXT NOT NULL,
    "country" TEXT NOT NULL,
    "postal_code" TEXT NOT NULL,

    CONSTRAINT "address_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "account"."cart" (
    "user_id" BYTEA NOT NULL,

    CONSTRAINT "cart_pkey" PRIMARY KEY ("user_id")
);

-- CreateTable
CREATE TABLE "account"."item_on_cart" (
    "cart_id" BYTEA NOT NULL,
    "product_model_id" BYTEA NOT NULL,
    "quantity" BIGINT NOT NULL,

    CONSTRAINT "item_on_cart_pkey" PRIMARY KEY ("cart_id","product_model_id")
);

-- CreateTable
CREATE TABLE "product"."brand" (
    "id" BYTEA NOT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT NOT NULL,

    CONSTRAINT "brand_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "product"."model" (
    "id" BYTEA NOT NULL,
    "brand_id" BYTEA NOT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "list_price" DECIMAL(10,2) NOT NULL,
    "date_manufactured" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "model_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "product"."base" (
    "serial_id" BYTEA NOT NULL,
    "product_model_id" BYTEA NOT NULL,
    "date_created" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "date_update" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "base_pkey" PRIMARY KEY ("serial_id")
);

-- CreateTable
CREATE TABLE "product"."sale" (
    "id" BYTEA NOT NULL,
    "tag_name" TEXT,
    "product_model_id" BYTEA,
    "date_started" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "date_ended" TIMESTAMP(3),
    "quantity" BIGINT NOT NULL,
    "used" BIGINT NOT NULL DEFAULT 0,
    "is_active" BOOLEAN NOT NULL DEFAULT true,
    "discount_percent" INTEGER,
    "discount_price" DECIMAL(10,2),

    CONSTRAINT "sale_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "product"."tag_on_product" (
    "product_model_id" BYTEA NOT NULL,
    "tag_name" TEXT NOT NULL,

    CONSTRAINT "tag_on_product_pkey" PRIMARY KEY ("product_model_id","tag_name")
);

-- CreateTable
CREATE TABLE "product"."tag" (
    "tag_name" TEXT NOT NULL,
    "description" TEXT,

    CONSTRAINT "tag_pkey" PRIMARY KEY ("tag_name")
);

-- CreateTable
CREATE TABLE "payment"."invoice" (
    "id" BYTEA NOT NULL,
    "user_id" BYTEA NOT NULL,
    "address" TEXT NOT NULL,
    "total" DECIMAL(10,2) NOT NULL,
    "payment_method" "payment"."payment_method" NOT NULL,
    "date_created" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "invoice_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "payment"."product_on_invoice" (
    "id" BYTEA NOT NULL,
    "invoice_id" BYTEA NOT NULL,
    "product_serial_id" BYTEA NOT NULL,
    "quantity" BIGINT NOT NULL,
    "price" DECIMAL(10,2) NOT NULL,
    "totalPrice" DECIMAL(10,2) NOT NULL,

    CONSTRAINT "product_on_invoice_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "payment"."base" (
    "id" BYTEA NOT NULL,
    "status" "payment"."status" NOT NULL,
    "payment_method" "payment"."payment_method" NOT NULL,
    "invoice_id" BYTEA NOT NULL,
    "date_created" TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "date_expired" TIMESTAMP(0) NOT NULL,

    CONSTRAINT "base_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "product"."image" (
    "brand_id" BYTEA,
    "product_model_id" BYTEA,
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
CREATE UNIQUE INDEX "base_invoice_id_key" ON "payment"."base"("invoice_id");

-- AddForeignKey
ALTER TABLE "account"."user" ADD CONSTRAINT "user_id_fkey" FOREIGN KEY ("id") REFERENCES "account"."base"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "account"."shop" ADD CONSTRAINT "shop_id_fkey" FOREIGN KEY ("id") REFERENCES "account"."base"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "account"."address" ADD CONSTRAINT "address_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "account"."user"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "account"."cart" ADD CONSTRAINT "cart_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "account"."user"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "account"."item_on_cart" ADD CONSTRAINT "item_on_cart_cart_id_fkey" FOREIGN KEY ("cart_id") REFERENCES "account"."cart"("user_id") ON DELETE CASCADE ON UPDATE CASCADE;

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
ALTER TABLE "payment"."invoice" ADD CONSTRAINT "invoice_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "account"."user"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "payment"."product_on_invoice" ADD CONSTRAINT "product_on_invoice_invoice_id_fkey" FOREIGN KEY ("invoice_id") REFERENCES "payment"."invoice"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "payment"."base" ADD CONSTRAINT "base_invoice_id_fkey" FOREIGN KEY ("invoice_id") REFERENCES "payment"."invoice"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "product"."image" ADD CONSTRAINT "image_brand_id_fkey" FOREIGN KEY ("brand_id") REFERENCES "product"."brand"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "product"."image" ADD CONSTRAINT "image_product_model_id_fkey" FOREIGN KEY ("product_model_id") REFERENCES "product"."model"("id") ON DELETE CASCADE ON UPDATE CASCADE;

