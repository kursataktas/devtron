/*
 * Copyright (c) 2024. Devtron Inc.
 */

CREATE SEQUENCE IF NOT EXISTS id_seq_timeout_window_resource_mappings;
CREATE TABLE IF NOT EXISTS "public"."timeout_window_resource_mappings"
(
    "id"                int          NOT NULL DEFAULT nextval('id_seq_timeout_window_resource_mappings'::regclass),
    "timeout_window_configuration_id"              int NOT NULL,
    "resource_id" int NOT NULL,
    "resource_type" int NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT timeout_window_configuration_id_fkey
        FOREIGN KEY("timeout_window_configuration_id")
            REFERENCES"public"."timeout_window_configuration" ("id")
            ON DELETE CASCADE
);

ALTER  TABLE global_policy_searchable_field ADD COLUMN IF NOT EXISTS field_name varchar;
ALTER  TABLE global_policy_searchable_field ADD COLUMN IF NOT EXISTS value_int integer;
ALTER  TABLE global_policy_searchable_field ADD COLUMN IF NOT EXISTS value_time_stamp timestamptz;
CREATE UNIQUE INDEX idx_unique_policy_name_policy_of
    ON global_policy (name,policy_of)
    WHERE deleted = false;
ALTER TABLE  resource_filter_evaluation_audit ADD COLUMN "filter_type" integer DEFAULT 1;


INSERT INTO public.event (id, event_type, description) VALUES (6, 'BLOCKED', '');
UPDATE "public"."notification_templates" SET template_payload = '{"from": "{{fromEmail}}", "to": "{{toEmail}}","subject": "▶️ Deployment pipeline triggered | Application: {{appName}} | Environment:  {{envName}}","html":"<table cellpadding=0 style=\"font-family:Arial,Verdana,Helvetica;width:600px;height:485px;border-collapse:inherit;border-spacing:0;border:1px solid #d0d4d9;border-radius:8px;padding:16px 20px;margin:20px auto;box-shadow:0 0 8px 0 rgba(0,0,0,.1)\"><tr><td colspan=3><div style=\"padding-bottom:16px;margin-bottom:20px;border-bottom:1px solid #edf1f5;max-width:600px\"><img src=https://devtron-public-asset.s3.us-east-2.amazonaws.com/images/devtron/devtron-logo.png style=max-width:122px alt=cd-triggered></div><tr><td colspan=3><div style=\"background-color:#e5f2ff;border-top-left-radius:8px;border-top-right-radius:8px;padding:20px 20px 16px 20px;display:flex;justify-content:space-between\"><div style=width:90%><div style=font-size:16px;line-height:24px;font-weight:600;margin-bottom:6px;color:#000a14>▶️ Deployment pipeline triggered</div><span style=font-size:14px;line-height:20px;color:#000a14>{{eventTime}}</span><br><div><span style=font-size:14px;line-height:20px;color:#000a14>by</span><span style=font-size:14px;line-height:20px;color:#06c;margin-left:4px>{{triggeredBy}}</span></div></div><div><img src=https://cdn.devtron.ai/images/image_deploy_notification.png style=height:72px;width:72px></div></div><tr><td colspan=3><div style=display:flex><div style=\"width:124px;background-color:#e5f2ff;border-bottom-left-radius:8px;padding:0 0 20px 20px\">{{#deploymentHistoryLink}}<a href={{&deploymentHistoryLink}} style=\"height:32px;padding:7px 12px;line-height:32px;font-size:12px;font-weight:600;border-radius:4px;text-decoration:none;outline:0;min-width:64px;text-transform:capitalize;text-align:center;background:#06c;color:#fff;border:1px solid transparent;cursor:pointer\">View Pipeline</a>{{/deploymentHistoryLink}}</div><div style=\"width:90%;background-color:#e5f2ff;border-bottom-right-radius:8px;padding:0 0 20px 20px\">{{#appDetailsLink}}<a href={{&appDetailsLink}} style=\"height:32px;padding:7px 12px;line-height:32px;font-size:12px;font-weight:600;border-radius:4px;text-decoration:none;outline:0;min-width:64px;text-transform:capitalize;text-align:center;border:1px solid #d0d4d9;color:#06c;cursor:pointer;width:90%\">App Details</a>{{/appDetailsLink}}</div></div><tr><td><br><tr><td colspan=3><div style=\"display:flex;font-weight:600;border-radius:4px;border:1px solid #fde7e7;background:#fff;padding:8px 16px;margin-bottom:16px;display:{{deploymentWindowCommentStyle}};\"><img src=https://cdn.devtron.ai/images/shield-image.png style=margin-right:4px;height:16px;width:16px>{{deploymentWindowComment}}</div><tr><td><div style=color:#3b444c;font-size:13px>Application</div><td colspan=2><div style=color:#3b444c;font-size:13px>Environment</div><tr><td><div style=color:#000a14;font-size:14px>{{appName}}</div><td><div style=color:#000a14;font-size:14px>{{envName}}</div><tr><td><div style=color:#3b444c;font-size:13px;margin-top:12px>Stage</div><tr><td><div style=color:#000a14;font-size:14px>{{stage}}</div><tr><tr><td colspan=3><div style=\"font-weight:600;margin-top:20px;width:100%;border-top:1px solid #edf1f5;padding:16px 0 12px;font-size:14px\">Source Code</div></tr>{{#ciMaterials}} {{^webhookType}}<tr><td><div style=color:#3b444c;font-size:13px>Branch</div><td colspan=2><div style=color:#3b444c;font-size:13px>Commit</div><tr><tr><td><div style=color:#000a14;font-size:14px>{{appName}}/{{branch}}</div><td><div style=color:#000a14;font-size:14px>{{commit}}</div></tr>{{/webhookType}} {{/ciMaterials}}<tr><td colspan=3><div style=\"font-weight:600;margin-top:20px;width:100%;border-top:1px solid #edf1f5;padding:16px 0 12px;font-size:14px\">Image Details</div><tr><td><div style=color:#3b444c;font-size:13px>Image tag</div><tr><td><div style=color:#000a14;font-size:14px>{{dockerImg}}</div><tr><td><br><tr><td colspan=3><div style=\"border-top:1px solid #edf1f5;margin:20px 0 16px 0;height:1px\"></div><tr><td colspan=2 style=display:flex><span><a href=https://twitter.com/DevtronL style=cursor:pointer;text-decoration:none;padding-right:12px;display:flex target=_blank><div><img src=https://cdn.devtron.ai/images/twitter_social_dark.png style=width:20px></div></a></span><span><a href=https://www.linkedin.com/company/devtron-labs/mycompany/ style=cursor:pointer;text-decoration:none;padding-right:12px;display:flex target=_blank><div><img src=https://cdn.devtron.ai/images/linkedin_social_dark.png style=width:20px></div></a></span><span><a href=https://devtron.ai/blog/ style=color:#000a14;font-size:13px;line-height:20px;cursor:pointer;text-decoration:underline;padding-right:12px target=_blank>Blog</a></span><span><a href=https://devtron.ai/ style=color:#000a14;font-size:13px;line-height:20px;cursor:pointer;text-decoration:underline target=_blank>Website</a></span><td colspan=2 style=text-align:right><div style=color:#767d84;font-size:13px;line-height:20px>© Devtron Labs 2024</div></table>"}' WHERE event_type_id=1 AND  node_type='CD' AND channel_type!='slack';
UPDATE "public"."notification_templates" SET template_payload = '{"from": "{{fromEmail}}", "to": "{{toEmail}}","subject": "🛎️ CD success for app: {{appName}} on environment: {{envName}}","html":"<table cellpadding=0 style=\"font-family:Arial,Verdana,Helvetica;width:600px;height:485px;border-collapse:inherit;border-spacing:0;border:1px solid #d0d4d9;border-radius:8px;padding:16px 20px;margin:20px auto;box-shadow:0 0 8px 0 rgba(0,0,0,.1)\"><tr><td colspan=3><div style=\"padding-bottom:16px;margin-bottom:20px;border-bottom:1px solid #edf1f5;max-width:600px\"><img src=https://devtron-public-asset.s3.us-east-2.amazonaws.com/images/devtron/devtron-logo.png style=max-width:122px alt=ci-triggered></div><tr><td colspan=3><div style=\"background-color:#E9FBF4;border-top-left-radius:8px;border-top-right-radius:8px;padding:20px 20px 16px 20px;display:flex;justify-content:space-between\"><div style=width:90%><div style=font-size:16px;line-height:24px;font-weight:600;margin-bottom:6px;color:#000a14>🎉  Deployment pipeline succeeded</div><span style=font-size:14px;line-height:20px;color:#000a14>{{eventTime}}</span><br><div><span style=font-size:14px;line-height:20px;color:#000a14>by</span><span style=font-size:14px;line-height:20px;color:#06c;margin-left:4px>{{triggeredBy}}</span></div></div><div><img src=https://cdn.devtron.ai/images/image_deploy_notification.png style=height:72px;width:72px></div></div><tr><td colspan=3><div style=display:flex><div style=\"width:124px;background-color:#E9FBF4;border-bottom-left-radius:8px;padding:0 0 20px 20px\">{{#deploymentHistoryLink}}<a href={{&deploymentHistoryLink}} style=\"height:32px;padding:7px 12px;line-height:32px;font-size:12px;font-weight:600;border-radius:4px;text-decoration:none;outline:0;min-width:64px;text-transform:capitalize;text-align:center;background:#06c;color:#fff;border:1px solid transparent;cursor:pointer\">View Pipeline</a>{{/deploymentHistoryLink}}</div><div style=\"width:90%;background-color:#E9FBF4;border-bottom-right-radius:8px;padding:0 0 20px 20px\">{{#appDetailsLink}}<a href={{&appDetailsLink}} style=\"height:32px;padding:7px 12px;line-height:32px;font-size:12px;font-weight:600;border-radius:4px;text-decoration:none;outline:0;min-width:64px;text-transform:capitalize;text-align:center;border:1px solid #d0d4d9;color:#06c;cursor:pointer;width:90%\">App Details</a>{{/appDetailsLink}}</div></div><tr><td><br><tr><td colspan=3><div style=\"display:flex;font-weight:600;border-radius:4px;border:1px solid #fde7e7;background:#fff;padding:8px 16px;margin-bottom:16px;display:{{deploymentWindowCommentStyle}};\"><img src=https://cdn.devtron.ai/images/shield-image.png style=margin-right:4px;height:16px;width:16px>{{deploymentWindowComment}}</div><tr><td><div style=color:#3b444c;font-size:13px>Application</div><td colspan=2><div style=color:#3b444c;font-size:13px>Environment</div><tr><td><div style=color:#000a14;font-size:14px>{{appName}}</div><td><div style=color:#000a14;font-size:14px>{{envName}}</div><tr><td><div style=color:#3b444c;font-size:13px;margin-top:12px>Stage</div><tr><td><div style=color:#000a14;font-size:14px>{{stage}}</div><tr><tr><td colspan=3><div style=\"font-weight:600;margin-top:20px;width:100%;border-top:1px solid #edf1f5;padding:16px 0 12px;font-size:14px\">Source Code</div></tr>{{#ciMaterials}} {{^webhookType}}<tr><td><div style=color:#3b444c;font-size:13px>Branch</div><td colspan=2><div style=color:#3b444c;font-size:13px>Commit</div><tr><tr><td><div style=color:#000a14;font-size:14px>{{appName}}/{{branch}}</div><td><div style=color:#000a14;font-size:14px>{{commit}}</div></tr>{{/webhookType}} {{/ciMaterials}}<tr><td colspan=3><div style=\"font-weight:600;margin-top:20px;width:100%;border-top:1px solid #edf1f5;padding:16px 0 12px;font-size:14px\">Image Details</div><tr><td><div style=color:#3b444c;font-size:13px>Image tag</div><tr><td><div style=color:#000a14;font-size:14px>{{dockerImg}}</div><tr><td><br><tr><td colspan=3><div style=\"border-top:1px solid #edf1f5;margin:20px 0 16px 0;height:1px\"></div><tr><td colspan=2 style=display:flex><span><a href=https://twitter.com/DevtronL style=cursor:pointer;text-decoration:none;padding-right:12px;display:flex target=_blank><div><img src=https://cdn.devtron.ai/images/twitter_social_dark.png style=width:20px></div></a></span><span><a href=https://www.linkedin.com/company/devtron-labs/mycompany/ style=cursor:pointer;text-decoration:none;padding-right:12px;display:flex target=_blank><div><img src=https://cdn.devtron.ai/images/linkedin_social_dark.png style=width:20px></div></a></span><span><a href=https://devtron.ai/blog/ style=color:#000a14;font-size:13px;line-height:20px;cursor:pointer;text-decoration:underline;padding-right:12px target=_blank>Blog</a></span><span><a href=https://devtron.ai/ style=color:#000a14;font-size:13px;line-height:20px;cursor:pointer;text-decoration:underline target=_blank>Website</a></span><td colspan=2 style=text-align:right><div style=color:#767d84;font-size:13px;line-height:20px>© Devtron Labs 2024</div></table>"}' WHERE event_type_id=2 AND  node_type='CD' AND channel_type!='slack';
UPDATE "public"."notification_templates" SET template_payload = '{"from": "{{fromEmail}}", "to": "{{toEmail}}","subject": "❌  Deployment pipeline failed |Application: {{appName}} | Environment:  {{envName}}","html":"<table cellpadding=0 style=\"font-family:Arial,Verdana,Helvetica;width:600px;height:485px;border-collapse:inherit;border-spacing:0;border:1px solid #d0d4d9;border-radius:8px;padding:16px 20px;margin:20px auto;box-shadow:0 0 8px 0 rgba(0,0,0,.1)\"><tr><td colspan=3><div style=\"padding-bottom:16px;margin-bottom:20px;border-bottom:1px solid #edf1f5;max-width:600px\"><img src=https://devtron-public-asset.s3.us-east-2.amazonaws.com/images/devtron/devtron-logo.png style=max-width:122px alt=ci-triggered></div><tr><td colspan=3><div style=\"background-color:#FDE7E7;border-top-left-radius:8px;border-top-right-radius:8px;padding:20px 20px 16px 20px;display:flex;justify-content:space-between\"><div style=width:90%><div style=font-size:16px;line-height:24px;font-weight:600;margin-bottom:6px;color:#000a14>❌  Deployment pipeline failed</div><span style=font-size:14px;line-height:20px;color:#000a14>{{eventTime}}</span><br><div><span style=font-size:14px;line-height:20px;color:#000a14>by</span><span style=font-size:14px;line-height:20px;color:#06c;margin-left:4px>{{triggeredBy}}</span></div></div><div><img src=https://cdn.devtron.ai/images/image_deploy_notification.png style=height:72px;width:72px></div></div><tr><td colspan=3><div style=display:flex><div style=\"width:124px;background-color:#FDE7E7;border-bottom-left-radius:8px;padding:0 0 20px 20px\">{{#deploymentHistoryLink}}<a href={{&deploymentHistoryLink}} style=\"height:32px;padding:7px 12px;line-height:32px;font-size:12px;font-weight:600;border-radius:4px;text-decoration:none;outline:0;min-width:64px;text-transform:capitalize;text-align:center;background:#06c;color:#fff;border:1px solid transparent;cursor:pointer\">View Pipeline</a>{{/deploymentHistoryLink}}</div><div style=\"width:90%;background-color:#FDE7E7;border-bottom-right-radius:8px;padding:0 0 20px 20px\">{{#appDetailsLink}}<a href={{&appDetailsLink}} style=\"height:32px;padding:7px 12px;line-height:32px;font-size:12px;font-weight:600;border-radius:4px;text-decoration:none;outline:0;min-width:64px;text-transform:capitalize;text-align:center;border:1px solid #d0d4d9;color:#06c;cursor:pointer;width:90%\">App Details</a>{{/appDetailsLink}}</div></div><tr><td><br><tr><td colspan=3><div style=\"display:flex;font-weight:600;border-radius:4px;border:1px solid #fde7e7;background:#fff;padding:8px 16px;margin-bottom:16px;display:{{deploymentWindowCommentStyle}};\"><img src=https://cdn.devtron.ai/images/shield-image.png style=margin-right:4px;height:16px;width:16px>{{deploymentWindowComment}}</div><tr><td><div style=color:#3b444c;font-size:13px>Application</div><td colspan=2><div style=color:#3b444c;font-size:13px>Environment</div><tr><td><div style=color:#000a14;font-size:14px>{{appName}}</div><td><div style=color:#000a14;font-size:14px>{{envName}}</div><tr><td><div style=color:#3b444c;font-size:13px;margin-top:12px>Stage</div><tr><td><div style=color:#000a14;font-size:14px>{{stage}}</div><tr><tr><td colspan=3><div style=\"font-weight:600;margin-top:20px;width:100%;border-top:1px solid #edf1f5;padding:16px 0 12px;font-size:14px\">Source Code</div></tr>{{#ciMaterials}} {{^webhookType}}<tr><td><div style=color:#3b444c;font-size:13px>Branch</div><td colspan=2><div style=color:#3b444c;font-size:13px>Commit</div><tr><tr><td><div style=color:#000a14;font-size:14px>{{appName}}/{{branch}}</div><td><div style=color:#000a14;font-size:14px>{{commit}}</div></tr>{{/webhookType}} {{/ciMaterials}}<tr><td colspan=3><div style=\"font-weight:600;margin-top:20px;width:100%;border-top:1px solid #edf1f5;padding:16px 0 12px;font-size:14px\">Image Details</div><tr><td><div style=color:#3b444c;font-size:13px>Image tag</div><tr><td><div style=color:#000a14;font-size:14px>{{dockerImg}}</div><tr><td><br><tr><td colspan=3><div style=\"border-top:1px solid #edf1f5;margin:20px 0 16px 0;height:1px\"></div><tr><td colspan=2 style=display:flex><span><a href=https://twitter.com/DevtronL style=cursor:pointer;text-decoration:none;padding-right:12px;display:flex target=_blank><div><img src=https://cdn.devtron.ai/images/twitter_social_dark.png style=width:20px></div></a></span><span><a href=https://www.linkedin.com/company/devtron-labs/mycompany/ style=cursor:pointer;text-decoration:none;padding-right:12px;display:flex target=_blank><div><img src=https://cdn.devtron.ai/images/linkedin_social_dark.png style=width:20px></div></a></span><span><a href=https://devtron.ai/blog/ style=color:#000a14;font-size:13px;line-height:20px;cursor:pointer;text-decoration:underline;padding-right:12px target=_blank>Blog</a></span><span><a href=https://devtron.ai/ style=color:#000a14;font-size:13px;line-height:20px;cursor:pointer;text-decoration:underline target=_blank>Website</a></span><td colspan=2 style=text-align:right><div style=color:#767d84;font-size:13px;line-height:20px>© Devtron Labs 2024</div></table>"}' WHERE event_type_id=3 AND  node_type='CD' AND channel_type!='slack';

INSERT INTO "public"."notification_templates" (channel_type, node_type, event_type_id, template_name, template_payload)
VALUES ('ses', 'CD', 6, 'Config approval ses template', '{"from": "{{fromEmail}}", "to": "{{toEmail}}","subject": "🚫  Auto-deployment blocked:| Application: {{appName}} | Environment:  {{envName}}","html":"<table cellpadding=0 style=\"font-family:Arial,Verdana,Helvetica;width:600px;height:485px;border-collapse:inherit;border-spacing:0;border:1px solid #d0d4d9;border-radius:8px;padding:16px 20px;margin:20px auto;box-shadow:0 0 8px 0 rgba(0,0,0,.1)\"><tr><td colspan=3><div style=\"padding-bottom:16px;margin-bottom:20px;border-bottom:1px solid #edf1f5;max-width:600px\"><img src=https://devtron-public-asset.s3.us-east-2.amazonaws.com/images/devtron/devtron-logo.png style=max-width:122px alt=ci-triggered></div><tr><td colspan=3><div style=\"background-color:#FDE7E7;border-top-left-radius:8px;border-top-right-radius:8px;padding:20px 20px 16px 20px;display:flex;justify-content:space-between\"><div style=width:90%><div style=font-size:16px;line-height:24px;font-weight:600;margin-bottom:6px;color:#000a14>🚫  Auto-deployment blocked</div><span style=font-size:14px;line-height:20px;color:#000a14>{{eventTime}}</span><br><div><span style=font-size:14px;line-height:20px;color:#000a14>by</span><span style=font-size:14px;line-height:20px;color:#06c;margin-left:4px>{{triggeredBy}}</span></div></div><div><img src=https://cdn.devtron.ai/images/image_deploy_notification.png style=height:72px;width:72px></div></div><tr><td colspan=3><div style=display:flex><div style=\"width:124px;background-color:#FDE7E7;border-bottom-left-radius:8px;padding:0 0 20px 20px\">{{#deploymentHistoryLink}}<a href={{&deploymentHistoryLink}} style=\"height:32px;padding:7px 12px;line-height:32px;font-size:12px;font-weight:600;border-radius:4px;text-decoration:none;outline:0;min-width:64px;text-transform:capitalize;text-align:center;background:#06c;color:#fff;border:1px solid transparent;cursor:pointer\">View Pipeline</a>{{/deploymentHistoryLink}}</div><div style=\"width:90%;background-color:#FDE7E7;border-bottom-right-radius:8px;padding:0 0 20px 20px\">{{#appDetailsLink}}<a href={{&appDetailsLink}} style=\"height:32px;padding:7px 12px;line-height:32px;font-size:12px;font-weight:600;border-radius:4px;text-decoration:none;outline:0;min-width:64px;text-transform:capitalize;text-align:center;border:1px solid #d0d4d9;color:#06c;cursor:pointer;width:90%\">App Details</a>{{/appDetailsLink}}</div></div><tr><td><br><tr><td colspan=3><div style=\"display:flex;font-weight:600;border-radius:4px;border:1px solid #fde7e7;background:#fff;padding:8px 16px;margin-bottom:16px;display:{{deploymentWindowCommentStyle}};\"><img src=https://cdn.devtron.ai/images/shield-image.png style=margin-right:4px;height:16px;width:16px>{{deploymentWindowComment}}</div><tr><td><div style=color:#3b444c;font-size:13px>Application</div><td colspan=2><div style=color:#3b444c;font-size:13px>Environment</div><tr><td><div style=color:#000a14;font-size:14px>{{appName}}</div><td><div style=color:#000a14;font-size:14px>{{envName}}</div><tr><td><div style=color:#3b444c;font-size:13px;margin-top:12px>Stage</div><tr><td><div style=color:#000a14;font-size:14px>{{stage}}</div><tr><tr><td colspan=3><div style=\"font-weight:600;margin-top:20px;width:100%;border-top:1px solid #edf1f5;padding:16px 0 12px;font-size:14px\">Source Code</div></tr>{{#ciMaterials}} {{^webhookType}}<tr><td><div style=color:#3b444c;font-size:13px>Branch</div><td colspan=2><div style=color:#3b444c;font-size:13px>Commit</div><tr><tr><td><div style=color:#000a14;font-size:14px>{{appName}}/{{branch}}</div><td><div style=color:#000a14;font-size:14px>{{commit}}</div></tr>{{/webhookType}} {{/ciMaterials}}<tr><td colspan=3><div style=\"font-weight:600;margin-top:20px;width:100%;border-top:1px solid #edf1f5;padding:16px 0 12px;font-size:14px\">Image Details</div><tr><td><div style=color:#3b444c;font-size:13px>Image tag</div><tr><td><div style=color:#000a14;font-size:14px>{{dockerImg}}</div><tr><td><br><tr><td colspan=3><div style=\"border-top:1px solid #edf1f5;margin:20px 0 16px 0;height:1px\"></div><tr><td colspan=2 style=display:flex><span><a href=https://twitter.com/DevtronL style=cursor:pointer;text-decoration:none;padding-right:12px;display:flex target=_blank><div><img src=https://cdn.devtron.ai/images/twitter_social_dark.png style=width:20px></div></a></span><span><a href=https://www.linkedin.com/company/devtron-labs/mycompany/ style=cursor:pointer;text-decoration:none;padding-right:12px;display:flex target=_blank><div><img src=https://cdn.devtron.ai/images/linkedin_social_dark.png style=width:20px></div></a></span><span><a href=https://devtron.ai/blog/ style=color:#000a14;font-size:13px;line-height:20px;cursor:pointer;text-decoration:underline;padding-right:12px target=_blank>Blog</a></span><span><a href=https://devtron.ai/ style=color:#000a14;font-size:13px;line-height:20px;cursor:pointer;text-decoration:underline target=_blank>Website</a></span><td colspan=2 style=text-align:right><div style=color:#767d84;font-size:13px;line-height:20px>© Devtron Labs 2024</div></table>"}');
INSERT INTO "public"."notification_templates" (channel_type, node_type, event_type_id, template_name, template_payload)
VALUES ('smtp', 'CD', 6, 'Config approval smtp template', '{"from": "{{fromEmail}}", "to": "{{toEmail}}","subject": "🚫  Auto-deployment blocked:| Application: {{appName}} | Environment:  {{envName}}","html":"<table cellpadding=0 style=\"font-family:Arial,Verdana,Helvetica;width:600px;height:485px;border-collapse:inherit;border-spacing:0;border:1px solid #d0d4d9;border-radius:8px;padding:16px 20px;margin:20px auto;box-shadow:0 0 8px 0 rgba(0,0,0,.1)\"><tr><td colspan=3><div style=\"padding-bottom:16px;margin-bottom:20px;border-bottom:1px solid #edf1f5;max-width:600px\"><img src=https://devtron-public-asset.s3.us-east-2.amazonaws.com/images/devtron/devtron-logo.png style=max-width:122px alt=ci-triggered></div><tr><td colspan=3><div style=\"background-color:#FDE7E7;border-top-left-radius:8px;border-top-right-radius:8px;padding:20px 20px 16px 20px;display:flex;justify-content:space-between\"><div style=width:90%><div style=font-size:16px;line-height:24px;font-weight:600;margin-bottom:6px;color:#000a14>🚫  Auto-deployment blocked</div><span style=font-size:14px;line-height:20px;color:#000a14>{{eventTime}}</span><br><div><span style=font-size:14px;line-height:20px;color:#000a14>by</span><span style=font-size:14px;line-height:20px;color:#06c;margin-left:4px>{{triggeredBy}}</span></div></div><div><img src=https://cdn.devtron.ai/images/image_deploy_notification.png style=height:72px;width:72px></div></div><tr><td colspan=3><div style=display:flex><div style=\"width:124px;background-color:#FDE7E7;border-bottom-left-radius:8px;padding:0 0 20px 20px\">{{#deploymentHistoryLink}}<a href={{&deploymentHistoryLink}} style=\"height:32px;padding:7px 12px;line-height:32px;font-size:12px;font-weight:600;border-radius:4px;text-decoration:none;outline:0;min-width:64px;text-transform:capitalize;text-align:center;background:#06c;color:#fff;border:1px solid transparent;cursor:pointer\">View Pipeline</a>{{/deploymentHistoryLink}}</div><div style=\"width:90%;background-color:#FDE7E7;border-bottom-right-radius:8px;padding:0 0 20px 20px\">{{#appDetailsLink}}<a href={{&appDetailsLink}} style=\"height:32px;padding:7px 12px;line-height:32px;font-size:12px;font-weight:600;border-radius:4px;text-decoration:none;outline:0;min-width:64px;text-transform:capitalize;text-align:center;border:1px solid #d0d4d9;color:#06c;cursor:pointer;width:90%\">App Details</a>{{/appDetailsLink}}</div></div><tr><td><br><tr><td colspan=3><div style=\"display:flex;font-weight:600;border-radius:4px;border:1px solid #fde7e7;background:#fff;padding:8px 16px;margin-bottom:16px;display:{{deploymentWindowCommentStyle}};\"><img src=https://cdn.devtron.ai/images/shield-image.png style=margin-right:4px;height:16px;width:16px>{{deploymentWindowComment}}</div><tr><td><div style=color:#3b444c;font-size:13px>Application</div><td colspan=2><div style=color:#3b444c;font-size:13px>Environment</div><tr><td><div style=color:#000a14;font-size:14px>{{appName}}</div><td><div style=color:#000a14;font-size:14px>{{envName}}</div><tr><td><div style=color:#3b444c;font-size:13px;margin-top:12px>Stage</div><tr><td><div style=color:#000a14;font-size:14px>{{stage}}</div><tr><tr><td colspan=3><div style=\"font-weight:600;margin-top:20px;width:100%;border-top:1px solid #edf1f5;padding:16px 0 12px;font-size:14px\">Source Code</div></tr>{{#ciMaterials}} {{^webhookType}}<tr><td><div style=color:#3b444c;font-size:13px>Branch</div><td colspan=2><div style=color:#3b444c;font-size:13px>Commit</div><tr><tr><td><div style=color:#000a14;font-size:14px>{{appName}}/{{branch}}</div><td><div style=color:#000a14;font-size:14px>{{commit}}</div></tr>{{/webhookType}} {{/ciMaterials}}<tr><td colspan=3><div style=\"font-weight:600;margin-top:20px;width:100%;border-top:1px solid #edf1f5;padding:16px 0 12px;font-size:14px\">Image Details</div><tr><td><div style=color:#3b444c;font-size:13px>Image tag</div><tr><td><div style=color:#000a14;font-size:14px>{{dockerImg}}</div><tr><td><br><tr><td colspan=3><div style=\"border-top:1px solid #edf1f5;margin:20px 0 16px 0;height:1px\"></div><tr><td colspan=2 style=display:flex><span><a href=https://twitter.com/DevtronL style=cursor:pointer;text-decoration:none;padding-right:12px;display:flex target=_blank><div><img src=https://cdn.devtron.ai/images/twitter_social_dark.png style=width:20px></div></a></span><span><a href=https://www.linkedin.com/company/devtron-labs/mycompany/ style=cursor:pointer;text-decoration:none;padding-right:12px;display:flex target=_blank><div><img src=https://cdn.devtron.ai/images/linkedin_social_dark.png style=width:20px></div></a></span><span><a href=https://devtron.ai/blog/ style=color:#000a14;font-size:13px;line-height:20px;cursor:pointer;text-decoration:underline;padding-right:12px target=_blank>Blog</a></span><span><a href=https://devtron.ai/ style=color:#000a14;font-size:13px;line-height:20px;cursor:pointer;text-decoration:underline target=_blank>Website</a></span><td colspan=2 style=text-align:right><div style=color:#767d84;font-size:13px;line-height:20px>© Devtron Labs 2024</div></table>"}');