/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package adaptorinterface

import (
	"context"
	"errors"
	"log/slog"
	"os"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	hwmgmtv1alpha1 "github.com/openshift-kni/oran-o2ims/api/hardwaremanagement/v1alpha1"
)

type HwMgrAdaptorIntf interface {
	SetupAdaptor() error
	HandleNodePool(ctx context.Context, nodepool *hwmgmtv1alpha1.NodePool) (ctrl.Result, error)
	HandleNodePoolDeletion(ctx context.Context, nodepool *hwmgmtv1alpha1.NodePool) error
}

// Define the HwMgrAdaptor structures
type HwMgrAdaptorConfig struct {
	client.Client
	logger    *slog.Logger
	namespace string
}

type HwMgrAdaptor struct {
	config HwMgrAdaptorConfig
}

func NewHwMgrAdaptor(config *HwMgrAdaptorConfig) (hwmgr *HwMgrAdaptor, err error) {
	if config.logger == nil {
		err = errors.New("logger is mandatory")
		return
	}

	hwmgr = &HwMgrAdaptor{
		config: *config,
	}

	if hwmgr.config.namespace == "" {
		hwmgr.config.namespace = os.Getenv("MY_POD_NAMESPACE")
	}

	return
}
