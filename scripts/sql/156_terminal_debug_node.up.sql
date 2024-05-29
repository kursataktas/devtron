/*
 * Copyright (c) 2024. Devtron Inc.
 */

INSERT into terminal_access_templates(template_name, template_data, created_on, created_by, updated_on, updated_by) VALUES
('terminal-node-debug-pod','{"apiVersion":"v1","kind":"Pod","metadata":{"name":"${pod_name}"},"spec":{"serviceAccountName":"${pod_name}-sa","containers":[{"image":"${base_image}","imagePullPolicy":"IfNotPresent","name":"debugger","resources":{},"stdin":true,"terminationMessagePolicy":"File","tty":true,"volumeMounts":[{"mountPath":"/host","name":"host-root"}]}],"hostIPC":true,"hostNetwork":true,"hostPID":true,"nodeName":"${node_name}","restartPolicy":"Never","tolerations":[{"operator":"Exists"}],"volumes":[{"hostPath":{"path":"/"},"name":"host-root"}]}}', now(), 1, now(), 1);