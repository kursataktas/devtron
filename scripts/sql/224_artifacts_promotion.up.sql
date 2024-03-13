-- 1 for  deployment approval request, 2 for artifact promotion approval request
ALTER TABLE deployment_approval_user_data ADD COLUMN "request_type" integer NOT NULL DEFAULT 1;

--  drop the constraint as this is no longer valid
ALTER TABLE deployment_approval_user_data DROP CONSTRAINT deployment_approval_user_data_approval_request_id_fkey;
DROP INDEX deployment_approval_user_data_approval_request_id_user_id_key;

-- rename deployment_approval_user_data table to request_approval_user_data
ALTER TABLE deployment_approval_user_data RENAME TO request_approval_user_data;

-- user can take action only once on any approval_request
CREATE UNIQUE INDEX "unique_user_request_action"
    ON request_approval_user_data(user_id,approval_request_id,request_type);
-- 1 for  resource_filter, 2 for artifact promotion policy filter evaluation
ALTER TABLE  resource_filter_evaluation_audit ADD COLUMN "filter_type" integer DEFAULT 1;



-- create artifact promotion approval request table
CREATE SEQUENCE IF NOT EXISTS id_artifact_promotion_approval_request;
CREATE TABLE IF NOT EXISTS public.artifact_promotion_approval_request
(
    "created_by"                   int4         NOT NULL,
    "updated_by"                   int4         NOT NULL,
    "id"                           int          NOT NULL DEFAULT nextval('id_artifact_promotion_approval_request'::regclass),
--     foreign key to artifact_promotion_policy
    "policy_id"                    int          NOT NULL,
--     foreign key to filter_evaluation_audit
    "policy_evaluation_audit_id"   int          NOT NULL,
--     foreign key to ci_artifact
    "artifact_id"                  int          NOT NULL,
    "source_pipeline_id"           int          NOT NULL,
--     CI_PIPELINE(0) or WEBHOOK(1) or CD_PIPELINE(2)
    "source_type"                  int          NOT NULL,
    "destination_pipeline_id"      int          NOT NULL,
--  CD_PIPELINE(2) , currently not defining this column as destination is always CD_PIPELINE
--  "destination_type"             int          NOT NULL,
    "status"                       int          NOT NULL,
    "created_on"                   timestamptz  NOT NULL,
    "updated_on"                   timestamptz  NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "artifact_promotion_approval_request_policy_id_fkey" FOREIGN KEY ("policy_id") REFERENCES "public"."global_policy" ("id"),
    CONSTRAINT "artifact_promotion_approval_request_artifact_id_fkey" FOREIGN KEY ("artifact_id") REFERENCES "public"."ci_artifact" ("id"),
    CONSTRAINT "artifact_promotion_approval_request_policy_evaluation_audit_id_fkey" FOREIGN KEY ("policy_evaluation_audit_id") REFERENCES "public"."resource_filter_evaluation_audit" ("id")
    );


CREATE UNIQUE INDEX "idx_unique_artifact_promoted_to_destination"
    ON artifact_promotion_approval_request(artifact_id,destination_pipeline_id)
    WHERE status = 3;

-- custom role queries
insert into rbac_policy_resource_detail
(resource,
 policy_resource_value,
 allowed_actions,
 resource_object,
 eligible_entity_access_types,
 deleted,created_on,
 created_by,
 updated_on,
 updated_by)
values ('approvalPolicy',
        '{"value": "approvalPolicy", "indexKeyMap": {}}','{promoter/artifact}','{"value": "%/%/%", "indexKeyMap": {"0": "TeamObj", "2": "EnvObj", "4": "AppObj"}}','{apps/devtron-app}',
        false,
        now(),
        1,
        now(),
        1);



insert into default_rbac_role_data (role,
                                    default_role_data,
                                    created_on,
                                    created_by,
                                    updated_on,
                                    updated_by,
                                    enabled)
VALUES ('artifactPromoter',
        '{"entity": "apps", "roleName": "artifactPromoter", "accessType": "devtron-app", "roleDescription": "can promote artifact for a particular CD Pipeline", "roleDisplayName": "Artifact Promoter", "policyResourceList": [{"actions": ["promoter/artifact"],
"resource": "approvalPolicy"}], "updatePoliciesForExistingProvidedRoles": false}',
        now(),
        1,
        now(),
        1,
        true);