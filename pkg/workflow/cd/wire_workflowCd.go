/*
 * Copyright (c) 2024. Devtron Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cd

import (
	"github.com/devtron-labs/devtron/pkg/workflow/cd/configHistory"
	"github.com/google/wire"
)

var CdWorkflowWireSet = wire.NewSet(
	configHistory.NewPipelineConfigOverrideReadServiceImpl,
	wire.Bind(new(configHistory.PipelineConfigOverrideReadService), new(*configHistory.PipelineConfigOverrideReadServiceImpl)),
	NewCdWorkflowCommonServiceImpl,
	wire.Bind(new(CdWorkflowCommonService), new(*CdWorkflowCommonServiceImpl)),
	NewCdWorkflowServiceImpl,
	wire.Bind(new(CdWorkflowService), new(*CdWorkflowServiceImpl)),
	NewCdWorkflowRunnerServiceImpl,
	wire.Bind(new(CdWorkflowRunnerService), new(*CdWorkflowRunnerServiceImpl)),
)
