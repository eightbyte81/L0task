CREATE TABLE "order" (
                         "order_id" serial NOT NULL,
                         "order_uid" TEXT NOT NULL,
                         "track_number" TEXT NOT NULL,
                         "entry" TEXT NOT NULL,
                         "delivery_id" integer NOT NULL,
                         "payment_id" integer NOT NULL,
                         "locale" TEXT NOT NULL,
                         "internal_signature" TEXT NOT NULL,
                         "customer_id" TEXT NOT NULL,
                         "delivery_service" TEXT NOT NULL,
                         "shardkey" TEXT NOT NULL,
                         "sm_id" integer NOT NULL,
                         "date_created" TIMESTAMP NOT NULL,
                         "oof_shard" TEXT NOT NULL,
                         CONSTRAINT "order_pk" PRIMARY KEY ("order_id")
) WITH (
      OIDS=FALSE
      );



CREATE TABLE "delivery" (
                            "delivery_id" serial NOT NULL,
                            "name" TEXT NOT NULL,
                            "phone" TEXT NOT NULL,
                            "zip" TEXT NOT NULL,
                            "city" TEXT NOT NULL,
                            "address" TEXT NOT NULL,
                            "region" TEXT NOT NULL,
                            "email" TEXT NOT NULL,
                            CONSTRAINT "delivery_pk" PRIMARY KEY ("delivery_id")
) WITH (
      OIDS=FALSE
      );



CREATE TABLE "payment" (
                           "payment_id" serial NOT NULL,
                           "transaction" TEXT NOT NULL,
                           "request_id" TEXT NOT NULL,
                           "currency" TEXT NOT NULL,
                           "provider" TEXT NOT NULL,
                           "amount" integer NOT NULL,
                           "payment_dt" integer NOT NULL,
                           "bank" TEXT NOT NULL,
                           "delivery_cost" integer NOT NULL,
                           "goods_total" integer NOT NULL,
                           "custom_fee" integer NOT NULL,
                           CONSTRAINT "payment_pk" PRIMARY KEY ("payment_id")
) WITH (
      OIDS=FALSE
      );



CREATE TABLE "order_item" (
                              "order_items_id" serial NOT NULL,
                              "order_id" serial NOT NULL,
                              "chrt_id" integer NOT NULL,
                              CONSTRAINT "order_item_pk" PRIMARY KEY ("order_items_id")
) WITH (
      OIDS=FALSE
      );



CREATE TABLE "item" (
                        "chrt_id" integer NOT NULL,
                        "track_number" TEXT NOT NULL,
                        "price" integer NOT NULL,
                        "rid" TEXT NOT NULL,
                        "name" TEXT NOT NULL,
                        "sale" integer NOT NULL,
                        "size" TEXT NOT NULL,
                        "total_price" integer NOT NULL,
                        "nm_id" integer NOT NULL,
                        "brand" TEXT NOT NULL,
                        "status" integer NOT NULL,
                        CONSTRAINT "item_pk" PRIMARY KEY ("chrt_id")
) WITH (
      OIDS=FALSE
      );



ALTER TABLE "order" ADD CONSTRAINT "order_fk0" FOREIGN KEY ("delivery_id") REFERENCES "delivery"("delivery_id");
ALTER TABLE "order" ADD CONSTRAINT "order_fk1" FOREIGN KEY ("payment_id") REFERENCES "payment"("payment_id");



ALTER TABLE "order_item" ADD CONSTRAINT "order_item_fk0" FOREIGN KEY ("order_id") REFERENCES "order"("order_id");
ALTER TABLE "order_item" ADD CONSTRAINT "order_item_fk1" FOREIGN KEY ("chrt_id") REFERENCES "item"("chrt_id");






