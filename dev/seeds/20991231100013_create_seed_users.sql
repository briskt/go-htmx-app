-- +goose Up
-- begin creating seeds with ID 1001 in case records were created already
INSERT INTO users (email, employee_id, first_name, last_name, display_name, username, active, locked,
                  last_login_at, created_at, updated_at) VALUES
  ('10001@example.com', '10001', 'John', 'Doe', 'John Doe', 'john_doe', true, false,
   '2024-08-31 08:39:52', '2024-08-31 06:04:47', '2024-08-31 08:39:52');
-- +goose Down
DELETE FROM `user`;
