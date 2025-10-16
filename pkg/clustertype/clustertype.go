// Copyright 2025 RedHat
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package clustertype

import (
	"context"

	configv1 "github.com/openshift/api/config/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ClusterType represents the type of OpenShift cluster
type ClusterType string

const (
	// ManagementCluster represents a management cluster (hosts HyperShift control planes)
	ManagementCluster ClusterType = "management"
	// ServiceCluster represents a service cluster (classic OSD/ROSA)
	ServiceCluster ClusterType = "service"
	// HyperShiftCluster represents a HyperShift cluster (control plane hosted externally)
	HyperShiftCluster ClusterType = "hypershift"
	// Unknown represents an unknown cluster type
	Unknown ClusterType = "unknown"
)

// Detector provides methods to detect cluster type
type Detector interface {
	GetClusterType(ctx context.Context) (ClusterType, error)
}

// detector implements the Detector interface
type detector struct {
	client client.Client
}

// NewDetector creates a new cluster type detector
func NewDetector(c client.Client) Detector {
	return &detector{client: c}
}

// GetClusterType determines the cluster type based on Infrastructure CR
func (d *detector) GetClusterType(ctx context.Context) (ClusterType, error) {
	infra := &configv1.Infrastructure{}
	err := d.client.Get(ctx, client.ObjectKey{Name: "cluster"}, infra)
	if err != nil {
		return Unknown, err
	}

	// Check control plane topology
	topology := infra.Status.ControlPlaneTopology
	switch topology {
	case configv1.ExternalTopologyMode:
		// External control plane = HyperShift cluster
		return HyperShiftCluster, nil
	case configv1.HighlyAvailableTopologyMode, configv1.SingleReplicaTopologyMode:
		// Classic topology - need to determine if it's MC or SC
		// Check for HyperShift operator presence or other MC indicators
		if d.isManagementCluster(ctx) {
			return ManagementCluster, nil
		}
		return ServiceCluster, nil
	default:
		return Unknown, nil
	}
}

// isManagementCluster checks if this is a management cluster by looking for HyperShift indicators
func (d *detector) isManagementCluster(ctx context.Context) bool {
	// TODO: Implement logic to detect management cluster
	// Options:
	// 1. Check for hypershift-operator namespace
	// 2. Check for specific labels on the Infrastructure CR
	// 3. Check for the presence of HostedCluster CRDs
	// 4. Use an explicit label or annotation
	//
	// For now, return false (assume service cluster)
	// This should be enhanced based on actual management cluster characteristics
	return false
}
