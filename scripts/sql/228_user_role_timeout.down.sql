/*
 * Copyright (c) 2024. Devtron Inc.
 */

ALTER TABLE user_roles DROP CONSTRAINT user_roles_timeout_window_configuration_id_fkey;

ALTER TABLE user_roles DROP COLUMN timeout_window_configuration_id;