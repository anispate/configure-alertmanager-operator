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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	configv1 "github.com/openshift/api/config/v1"
)

// EXAMPLE IMPLEMENTATIONS FOR MANAGEMENT CLUSTER DETECTION
// Choose one or combine multiple approaches based on your requirements

// Example 1: Check for explicit label on Infrastructure CR
func (d *detector) isManagementClusterByLabel(ctx context.Context) bool {
	infra := &configv1.Infrastructure{}
	err := d.client.Get(ctx, client.ObjectKey{Name: "cluster"}, infra)
	if err != nil {
		return false
	}

	// Check for explicit cluster-type label
	// This label would be set during cluster provisioning
	if infra.Labels["cluster-type"] == "management" {
		return true
	}

	// Alternative: Check for annotation
	if infra.Annotations["openshift.io/cluster-role"] == "management" {
		return true
	}

	return false
}

// Example 2: Check for HyperShift operator namespace
func (d *detector) isManagementClusterByHyperShiftNamespace(ctx context.Context) bool {
	// Management clusters run the HyperShift operator
	ns := &corev1.Namespace{}
	err := d.client.Get(ctx, client.ObjectKey{Name: "hypershift"}, ns)
	if err != nil {
		if errors.IsNotFound(err) {
			return false
		}
		// If we can't determine, default to false
		return false
	}

	// If hypershift namespace exists, this is likely an MC
	return true
}

// Example 3: Check for multiple indicators (recommended)
func (d *detector) isManagementClusterMultiCheck(ctx context.Context) bool {
	// Priority 1: Check explicit label (most reliable)
	infra := &configv1.Infrastructure{}
	err := d.client.Get(ctx, client.ObjectKey{Name: "cluster"}, infra)
	if err == nil {
		if infra.Labels["cluster-type"] == "management" {
			return true
		}
		if infra.Annotations["openshift.io/cluster-role"] == "management" {
			return true
		}
	}

	// Priority 2: Check for HyperShift operator namespace
	ns := &corev1.Namespace{}
	err = d.client.Get(ctx, client.ObjectKey{Name: "hypershift"}, ns)
	if err == nil {
		// Namespace exists, likely an MC
		return true
	}

	// Priority 3: Could check for HostedCluster CRDs or other indicators
	// For now, default to false (service cluster)
	return false
}

// Example 4: Check for environment variable (for testing/override)
// This would be set in the operator deployment
/*
func isManagementClusterByEnv() bool {
	clusterType := os.Getenv("CLUSTER_TYPE")
	return clusterType == "management" || clusterType == "mc"
}
*/

// RECOMMENDED IMPLEMENTATION:
// Replace the isManagementCluster() function in clustertype.go with:
//
// func (d *detector) isManagementCluster(ctx context.Context) bool {
//     return d.isManagementClusterMultiCheck(ctx)
// }
