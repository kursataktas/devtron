/*
 * Copyright (c) 2024. Devtron Inc.
 */

UPDATE chart_ref_metadata set chart_description = 'Chart to deploy an advanced version of Deployment that supports blue-green and canary deployments. It requires a rollout controller to run inside the cluster to function.' WHERE chart_name = 'Rollout Deployment';
UPDATE chart_ref_metadata set chart_description = 'Chart to deploy a Job/CronJob. Job is a controller object that represents a finite task and CronJob can be used to schedule creation of Jobs.' WHERE chart_name = 'CronJob & Job';
UPDATE chart_ref_metadata set chart_description = 'Chart to deploy an Open-Source Enterprise-level solution to deploy Serverless apps.' WHERE chart_name = 'Knative';
UPDATE chart_ref_metadata set chart_description = 'Chart to deploy a Deployment that runs multiple replicas of your application and automatically replaces any instances that fail or become unresponsive.' WHERE chart_name = 'Deployment';