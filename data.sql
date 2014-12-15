
-- ----------------------------
--  Sequence structure for reminders_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."reminders_id_seq";
CREATE SEQUENCE "public"."reminders_id_seq" INCREMENT 1 START 5 MAXVALUE 9223372036854775807 MINVALUE 1 CACHE 1;

-- ----------------------------
--  Table structure for books
-- ----------------------------
DROP TABLE IF EXISTS "public"."reminders";
CREATE TABLE "public"."reminders" (
	"id" int4 NOT NULL DEFAULT nextval('reminders_id_seq'::regclass),
	"reminder" varchar(255) NOT NULL COLLATE "default",
	"sent" int NOT NULL
)
WITH (OIDS=FALSE);

-- ----------------------------
--  Records of books
-- ----------------------------
BEGIN;
INSERT INTO "public"."reminders" VALUES ('1', 'Crypto Course Week 1', 1);
COMMIT;


-- ----------------------------
--  Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."reminders_id_seq" RESTART 4 OWNED BY "reminders"."id";

-- ----------------------------
--  Primary key structure for table books
-- ----------------------------
ALTER TABLE "public"."reminders" ADD PRIMARY KEY ("id") NOT DEFERRABLE INITIALLY IMMEDIATE;