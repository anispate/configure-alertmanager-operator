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
	"testing"

	configv1 "github.com/openshift/api/config/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestGetClusterType(t *testing.T) {
	tests := []struct {
		name             string
		topology         configv1.TopologyMode
		expectedType     ClusterType
		expectError      bool
	}{
		{
			name:         "External topology returns HyperShift",
			topology:     configv1.ExternalTopologyMode,
			expectedType: HyperShiftCluster,
			expectError:  false,
		},
		{
			name:         "HighlyAvailable topology returns ServiceCluster",
			topology:     configv1.HighlyAvailableTopologyMode,
			expectedType: ServiceCluster,
			expectError:  false,
		},
		{
			name:         "SingleReplica topology returns ServiceCluster",
			topology:     configv1.SingleReplicaTopologyMode,
			expectedType: ServiceCluster,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fake infrastructure object
			infra := &configv1.Infrastructure{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cluster",
				},
				Status: configv1.InfrastructureStatus{
					ControlPlaneTopology: tt.topology,
				},
			}

			// Create a scheme and add the necessary types
			scheme := runtime.NewScheme()
			_ = configv1.AddToScheme(scheme)

			// Create a fake client with the infrastructure object
			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(infra).
				Build()

			// Create the detector
			detector := NewDetector(fakeClient)

			// Get the cluster type
			clusterType, err := detector.GetClusterType(context.TODO())

			// Check for errors
			if (err != nil) != tt.expectError {
				t.Errorf("GetClusterType() error = %v, expectError %v", err, tt.expectError)
				return
			}

			// Check the cluster type
			if clusterType != tt.expectedType {
				t.Errorf("GetClusterType() = %v, want %v", clusterType, tt.expectedType)
			}
		})
	}
}
