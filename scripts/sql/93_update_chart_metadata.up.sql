/*
 * Copyright (c) 2024. Devtron Inc.
 */

UPDATE chart_ref_metadata set chart_description = 'This chart deploys an advanced version of deployment that supports Blue/Green and Canary deployments. For functioning, it requires a rollout controller to run inside the cluster.' WHERE chart_name = 'Rollout Deployment';
UPDATE chart_ref_metadata set chart_description = 'This chart deploys Job & CronJob.  A Job is a controller object that represents a finite task and CronJob is used to schedule creation of Jobs.' WHERE chart_name = 'CronJob & Job';
UPDATE chart_ref_metadata set chart_description = 'This chart deploys Knative which is an Open-Source Enterprise-level solution to deploy Serverless apps.' WHERE chart_name = 'Knative';
UPDATE chart_ref_metadata set chart_description = 'Creates a deployment that runs multiple replicas of your application and automatically replaces any instances that fail or become unresponsive.' WHERE chart_name = 'Deployment';