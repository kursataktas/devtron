/*
 * Copyright (c) 2024. Devtron Inc.
 */

ALTER TABLE public.ci_pipeline ADD CONSTRAINT  ci_pipeline_parent_ci_pipeline_fkey FOREIGN KEY (parent_ci_pipeline) REFERENCES public.ci_pipeline(id) ;