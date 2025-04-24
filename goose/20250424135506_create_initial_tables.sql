-- +goose Up
-- +goose StatementBegin

-- --------------------------------------------------------
--
-- Table structure for table `migration`
--
CREATE TABLE "public"."migration" (
  "version" character varying(255) NOT NULL,
  "apply_time" int DEFAULT NULL
);

-- --------------------------------------------------------
--
-- Table structure for table `users`
--
CREATE TABLE "public"."users" (
    "id" SERIAL PRIMARY KEY,
    "email" character varying(255) NOT NULL,
    "employee_id" character varying(255) NOT NULL,
    "first_name" character varying(255) NOT NULL,
    "last_name" character varying(255) NOT NULL,
    "display_name" character varying(255) NOT NULL,
    "username" character varying(255) NOT NULL,
    "active" boolean NOT NULL,
    "locked" boolean NOT NULL,
    "last_login_at" timestamp NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL
);

-- --------------------------------------------------------
--
-- Table structure for table `tokens`
--
CREATE TABLE tokens (
    id SERIAL PRIMARY KEY,
    user_id int NOT NULL,
    hash character varying(255) NOT NULL UNIQUE,
    expires_at timestamp NOT NULL,
    last_used_at timestamp DEFAULT NULL,
    created_utc timestamp NOT NULL,
    updated_utc timestamp NOT NULL,
    CONSTRAINT tokens_user_id FOREIGN KEY (user_id)
        REFERENCES users(id) ON DELETE CASCADE ON UPDATE NO ACTION
);

-- --------------------------------------------------------
--
-- Table structure for table `email_logs`
--
CREATE TABLE "public"."email_logs" (
    "id" SERIAL PRIMARY KEY,
    "user_id" int NOT NULL,
    "message_type" character varying(255) NOT NULL,
    "created_at" timestamp NOT NULL,
    CONSTRAINT email_logs_users_id_fk FOREIGN KEY (user_id)
        REFERENCES users(id) ON DELETE SET NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "email_logs";
DROP TABLE "tokens";
DROP TABLE "user";
DROP TABLE "migration";
-- +goose StatementEnd
