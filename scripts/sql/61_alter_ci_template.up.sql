/*
 * Copyright (c) 2024. Devtron Inc.
 */

ALTER TABLE ci_template
    ADD COLUMN target_platform VARCHAR(1000) NOT NULL DEFAULT '';
