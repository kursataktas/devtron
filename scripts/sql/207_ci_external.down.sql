/*
 * Copyright (c) 2024. Devtron Inc.
 */

UPDATE "public"."ci_pipeline" SET ci_pipeline_type = NULL WHERE parent_ci_pipeline > 0 AND ci_pipeline_type = 'LINKED';