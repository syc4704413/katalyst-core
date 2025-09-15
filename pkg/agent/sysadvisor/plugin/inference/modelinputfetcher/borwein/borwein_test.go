/*
Copyright 2022 The Katalyst Authors.

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

package borwein

import (
	"context"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes/fake"

	internalfake "github.com/kubewharf/katalyst-api/pkg/client/clientset/versioned/fake"
	"github.com/kubewharf/katalyst-core/pkg/agent/sysadvisor/metacache"
	"github.com/kubewharf/katalyst-core/pkg/client"
	"github.com/kubewharf/katalyst-core/pkg/config"
	"github.com/kubewharf/katalyst-core/pkg/config/agent/global"
	metaconfig "github.com/kubewharf/katalyst-core/pkg/config/agent/metaserver"
	"github.com/kubewharf/katalyst-core/pkg/consts"
	"github.com/kubewharf/katalyst-core/pkg/metaserver"
	"github.com/kubewharf/katalyst-core/pkg/metaserver/agent"
	"github.com/kubewharf/katalyst-core/pkg/metaserver/agent/metric"
	"github.com/kubewharf/katalyst-core/pkg/metaserver/agent/node"
	dynamicconfig "github.com/kubewharf/katalyst-core/pkg/metaserver/kcc"
	"github.com/kubewharf/katalyst-core/pkg/metrics"
	metricspool "github.com/kubewharf/katalyst-core/pkg/metrics/metrics-pool"
	metricutil "github.com/kubewharf/katalyst-core/pkg/util/metric"
)

func generateTestMetaServer(clientSet *client.GenericClientSet) *metaserver.MetaServer {
	return &metaserver.MetaServer{
		MetaAgent: &agent.MetaAgent{
			MetricsFetcher: metric.NewFakeMetricsFetcher(metrics.DummyMetrics{}),
		},
		ConfigurationManager: &dynamicconfig.DummyConfigurationManager{},
	}
}

func generateTestGenericClientSet(kubeObjects, internalObjects []runtime.Object) *client.GenericClientSet {
	return &client.GenericClientSet{
		KubeClient:     fake.NewSimpleClientset(kubeObjects...),
		InternalClient: internalfake.NewSimpleClientset(internalObjects...),
		DynamicClient:  dynamicfake.NewSimpleDynamicClient(runtime.NewScheme(), internalObjects...),
	}
}

func TestNativeGetNodeFeatureValue(t *testing.T) {
	t.Parallel()
	nodeName := "node1"
	type args struct {
		featureName string
		conf        *config.Configuration
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test with fake node",
			args: args{
				featureName: NodeFeatureNodeName,
				conf:        config.NewConfiguration(),
			},
			want:    "node1",
			wantErr: false,
		},
	}
	nowTimestamp := time.Now().Unix()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clientSet := generateTestGenericClientSet([]runtime.Object{&v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: nodeName,
				},
			}}, nil)
			metaServer := generateTestMetaServer(clientSet)
			metaServer.NodeFetcher = node.NewRemoteNodeFetcher(&global.BaseConfiguration{NodeName: nodeName}, &metaconfig.NodeConfiguration{}, clientSet.KubeClient.CoreV1().Nodes())
			_, err := nativeGetNodeFeatureValue(nowTimestamp, metaServer, nil, tt.args.conf, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("NativeGetNodeFeatureValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
func TestNativeGetContainerFeatureValue(t *testing.T) {
	t.Parallel()
	podUID := "test-pod-uid"
	containerName := "test-container"
	fakeCPUUsage := 20.0

	clientSet := generateTestGenericClientSet(nil, nil)
	metaServer := generateTestMetaServer(clientSet)
	metaServer.MetricsFetcher.RegisterExternalMetric(func(store *metricutil.MetricStore) {
		store.SetContainerMetric(podUID, containerName, consts.MetricCPUUsageContainer, metricutil.MetricData{
			Value: fakeCPUUsage,
		})
	})
	metaServer.MetricsFetcher.Run(context.Background())

	type args struct {
		podUID        string
		containerName string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "test with fake container",
			args: args{
				podUID:        podUID,
				containerName: containerName,
			},
			want: map[string]interface{}{
				containerName: map[string]interface{}{
					consts.MetricCPUUsageContainer: fakeCPUUsage,
				},
			},
			wantErr: false,
		},
		{
			name: "test with not exist container",
			args: args{
				podUID:        podUID,
				containerName: "not-exist-container",
			},
			want:    nil,
			wantErr: true,
		},
	}
	nowTimestamp := time.Now().Unix()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := nativeGetContainerFeatureValue(nowTimestamp, tt.args.podUID, tt.args.containerName, metaServer,
				&metacache.MetaCacheImp{MetricsReader: metaServer.MetricsFetcher}, config.NewConfiguration(), nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("NativeGetContainerFeatureValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NativeGetContainerFeatureValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNativeGetNumaFeatureValue(t *testing.T) {
	t.Parallel()
	numaID := 0
	fakeMemBandwidth := 1000.0

	stateFileDir, err := ioutil.TempDir("", "metacache-TestNativeGetNumaFeatureValue")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(stateFileDir)
	conf := config.NewConfiguration()
	conf.GenericSysAdvisorConfiguration.StateFileDirectory = stateFileDir

	clientSet := generateTestGenericClientSet(nil, nil)
	metaServer := generateTestMetaServer(clientSet)
	metaServer.MetricsFetcher.RegisterExternalMetric(func(store *metricutil.MetricStore) {
		store.SetNumaMetric(numaID, consts.MetricMemBandwidthNuma, metricutil.MetricData{
			Value: fakeMemBandwidth,
		})
	})
	metaServer.MetricsFetcher.Run(context.Background())
	metaCache, err := metacache.NewMetaCacheImp(conf, metricspool.DummyMetricsEmitterPool{}, metaServer.MetricsFetcher)
	if err != nil {
		t.Fatalf("failed to create metacache: %v", err)
	}

	type args struct {
		numaID int
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "test successful fetch",
			args: args{
				numaID: numaID,
			},
			want: map[string]interface{}{
				consts.MetricMemBandwidthNuma: fakeMemBandwidth,
			},
			wantErr: false,
		},
		{
			name: "test metric not found",
			args: args{
				numaID: 1, // different numa id
			},
			want:    nil,
			wantErr: true,
		},
	}

	nowTimestamp := time.Now().Unix()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := nativeGetNumaFeatureValue(nowTimestamp, tt.args.numaID, metaServer, metaCache, config.NewConfiguration(), nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("nativeGetNumaFeatureValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nativeGetNumaFeatureValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
