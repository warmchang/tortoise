package pod

import (
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	"github.com/mercari/tortoise/api/v1beta3"
	"github.com/mercari/tortoise/pkg/annotation"
	"github.com/mercari/tortoise/pkg/features"
)

func TestService_ModifyPodTemplateResource(t *testing.T) {
	type args struct {
		podTemplate *v1.PodTemplateSpec
		tortoise    *v1beta3.Tortoise
		opts        []ModifyPodSpecResourceOption
	}
	tests := []struct {
		name string
		args args
		want *v1.PodTemplateSpec
	}{
		{
			name: "Tortoise is Auto; NoScaleDown option; istio CPU is changed",
			args: args{
				opts: []ModifyPodSpecResourceOption{NoScaleDown},
				podTemplate: &v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							annotation.IstioSidecarInjectionAnnotation:        "true",
							annotation.IstioSidecarProxyCPUAnnotation:         "100m",
							annotation.IstioSidecarProxyCPULimitAnnotation:    "300m",
							annotation.IstioSidecarProxyMemoryAnnotation:      "100Mi",
							annotation.IstioSidecarProxyMemoryLimitAnnotation: "100Mi",
						},
					},
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name: "container",
								Resources: v1.ResourceRequirements{
									Requests: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("100m"),
										v1.ResourceMemory: resource.MustParse("100Mi"),
									},
									Limits: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("300m"),
										v1.ResourceMemory: resource.MustParse("100Mi"),
									},
								},
							},
						},
					},
				},
				tortoise: &v1beta3.Tortoise{
					Spec: v1beta3.TortoiseSpec{
						UpdateMode: v1beta3.UpdateModeAuto,
					},
					Status: v1beta3.TortoiseStatus{
						TortoisePhase: v1beta3.TortoisePhaseWorking,
						Conditions: v1beta3.Conditions{
							ContainerResourceRequests: []v1beta3.ContainerResourceRequests{
								{
									ContainerName: "container",
									Resource: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("200m"), // scale up
										v1.ResourceMemory: resource.MustParse("1Mi"),  // scale down
									},
								},
								{
									ContainerName: "istio-proxy",
									Resource: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("300m"), // scale up
										v1.ResourceMemory: resource.MustParse("1Mi"),  // scale down
									},
								},
							},
						},
					},
				},
			},
			want: &v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						annotation.IstioSidecarInjectionAnnotation:        "true",
						annotation.IstioSidecarProxyCPUAnnotation:         "300m",
						annotation.IstioSidecarProxyCPULimitAnnotation:    "900m",
						annotation.IstioSidecarProxyMemoryAnnotation:      "100Mi",
						annotation.IstioSidecarProxyMemoryLimitAnnotation: "100Mi",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "container",
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("200m"),
									v1.ResourceMemory: resource.MustParse("100Mi"),
								},
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("600m"),
									v1.ResourceMemory: resource.MustParse("100Mi"),
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Tortoise is Auto; NoScaleDown option; istio Memory is changed",
			args: args{
				opts: []ModifyPodSpecResourceOption{NoScaleDown},
				podTemplate: &v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							annotation.IstioSidecarInjectionAnnotation:        "true",
							annotation.IstioSidecarProxyCPUAnnotation:         "100m",
							annotation.IstioSidecarProxyCPULimitAnnotation:    "300m",
							annotation.IstioSidecarProxyMemoryAnnotation:      "100Mi",
							annotation.IstioSidecarProxyMemoryLimitAnnotation: "100Mi",
						},
					},
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name: "container",
								Resources: v1.ResourceRequirements{
									Requests: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("100m"),
										v1.ResourceMemory: resource.MustParse("100Mi"),
									},
									Limits: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("300m"),
										v1.ResourceMemory: resource.MustParse("100Mi"),
									},
								},
							},
						},
					},
				},
				tortoise: &v1beta3.Tortoise{
					Spec: v1beta3.TortoiseSpec{
						UpdateMode: v1beta3.UpdateModeAuto,
					},
					Status: v1beta3.TortoiseStatus{
						TortoisePhase: v1beta3.TortoisePhaseWorking,
						Conditions: v1beta3.Conditions{
							ContainerResourceRequests: []v1beta3.ContainerResourceRequests{
								{
									ContainerName: "container",
									Resource: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("200m"), // scale up
										v1.ResourceMemory: resource.MustParse("1Mi"),  // scale down
									},
								},
								{
									ContainerName: "istio-proxy",
									Resource: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("20m"),   // scale down
										v1.ResourceMemory: resource.MustParse("300Mi"), // scale up
									},
								},
							},
						},
					},
				},
			},
			want: &v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						annotation.IstioSidecarInjectionAnnotation:        "true",
						annotation.IstioSidecarProxyCPUAnnotation:         "100m",
						annotation.IstioSidecarProxyCPULimitAnnotation:    "300m",
						annotation.IstioSidecarProxyMemoryAnnotation:      "300Mi",
						annotation.IstioSidecarProxyMemoryLimitAnnotation: "300Mi",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "container",
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("200m"),
									v1.ResourceMemory: resource.MustParse("100Mi"),
								},
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("600m"),
									v1.ResourceMemory: resource.MustParse("100Mi"),
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := New(nil, "", nil, nil)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			got := tt.args.podTemplate.DeepCopy()
			s.ModifyPodTemplateResource(got, tt.args.tortoise, tt.args.opts...)
			if d := cmp.Diff(got, tt.want); d != "" {
				t.Errorf("ModifyPodResource() mismatch (-want +got):\n%s", d)
			}
		})
	}
}

func TestService_ModifyPodSpecResource(t *testing.T) {
	type fields struct {
		resourceLimitMultiplier map[string]int64
		minimumCPULimit         string
		featureFlags            []features.FeatureFlag
	}
	type args struct {
		pod      *v1.Pod
		tortoise *v1beta3.Tortoise
		opts     []ModifyPodSpecResourceOption
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *v1.Pod
	}{
		{
			name: "Tortoise is Off",
			args: args{
				pod: &v1.Pod{
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name: "container",
								Resources: v1.ResourceRequirements{
									Requests: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("100m"),
										v1.ResourceMemory: resource.MustParse("100Mi"),
									},
									Limits: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("200m"),
										v1.ResourceMemory: resource.MustParse("200Mi"),
									},
								},
							},
						},
					},
				},
				tortoise: &v1beta3.Tortoise{
					Spec: v1beta3.TortoiseSpec{
						UpdateMode: v1beta3.UpdateModeOff,
					},
					Status: v1beta3.TortoiseStatus{
						TortoisePhase: v1beta3.TortoisePhaseWorking,
					},
				},
			},
			want: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "container",
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("100m"),
									v1.ResourceMemory: resource.MustParse("100Mi"),
								},
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("200m"),
									v1.ResourceMemory: resource.MustParse("200Mi"),
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Tortoise is just created",
			args: args{
				pod: &v1.Pod{
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name: "container",
								Resources: v1.ResourceRequirements{
									Requests: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("100m"),
										v1.ResourceMemory: resource.MustParse("100Mi"),
									},
									Limits: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("200m"),
										v1.ResourceMemory: resource.MustParse("200Mi"),
									},
								},
							},
						},
					},
				},
				tortoise: &v1beta3.Tortoise{
					Spec: v1beta3.TortoiseSpec{
						UpdateMode: v1beta3.UpdateModeAuto,
					},
					Status: v1beta3.TortoiseStatus{
						// TortoisePhase is empty
					},
				},
			},
			want: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "container",
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("100m"),
									v1.ResourceMemory: resource.MustParse("100Mi"),
								},
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("200m"),
									v1.ResourceMemory: resource.MustParse("200Mi"),
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Tortoise is Initializing",
			args: args{
				pod: &v1.Pod{
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name: "container",
								Resources: v1.ResourceRequirements{
									Requests: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("100m"),
										v1.ResourceMemory: resource.MustParse("100Mi"),
									},
									Limits: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("200m"),
										v1.ResourceMemory: resource.MustParse("200Mi"),
									},
								},
							},
						},
					},
				},
				tortoise: &v1beta3.Tortoise{
					Spec: v1beta3.TortoiseSpec{
						UpdateMode: v1beta3.UpdateModeAuto,
					},
					Status: v1beta3.TortoiseStatus{
						TortoisePhase: v1beta3.TortoisePhaseInitializing,
					},
				},
			},
			want: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "container",
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("100m"),
									v1.ResourceMemory: resource.MustParse("100Mi"),
								},
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("200m"),
									v1.ResourceMemory: resource.MustParse("200Mi"),
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Tortoise is Gatheringdata",
			args: args{
				pod: &v1.Pod{
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name: "container",
								Resources: v1.ResourceRequirements{
									Requests: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("100m"),
										v1.ResourceMemory: resource.MustParse("100Mi"),
									},
									Limits: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("200m"),
										v1.ResourceMemory: resource.MustParse("200Mi"),
									},
								},
							},
						},
					},
				},
				tortoise: &v1beta3.Tortoise{
					Spec: v1beta3.TortoiseSpec{
						UpdateMode: v1beta3.UpdateModeAuto,
					},
					Status: v1beta3.TortoiseStatus{
						TortoisePhase: v1beta3.TortoisePhaseGatheringData,
					},
				},
			},
			want: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "container",
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("100m"),
									v1.ResourceMemory: resource.MustParse("100Mi"),
								},
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("200m"),
									v1.ResourceMemory: resource.MustParse("200Mi"),
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Tortoise is Auto; Resource Request and Limit are updated based on the recommendation",
			args: args{
				pod: &v1.Pod{
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name: "container",
								Resources: v1.ResourceRequirements{
									Requests: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("100m"),
										v1.ResourceMemory: resource.MustParse("100Mi"),
									},
									Limits: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("300m"),
										v1.ResourceMemory: resource.MustParse("100Mi"),
									},
								},
							},
							{
								Name: "istio-proxy",
								Resources: v1.ResourceRequirements{
									Requests: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("100m"),
										v1.ResourceMemory: resource.MustParse("100Mi"),
									},
									Limits: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("400m"),
										v1.ResourceMemory: resource.MustParse("200Mi"),
									},
								},
							},
						},
					},
				},
				tortoise: &v1beta3.Tortoise{
					Spec: v1beta3.TortoiseSpec{
						UpdateMode: v1beta3.UpdateModeAuto,
					},
					Status: v1beta3.TortoiseStatus{
						TortoisePhase: v1beta3.TortoisePhaseWorking,
						Conditions: v1beta3.Conditions{
							ContainerResourceRequests: []v1beta3.ContainerResourceRequests{
								{
									ContainerName: "container",
									Resource: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("200m"),
										v1.ResourceMemory: resource.MustParse("200Mi"),
									},
								},
								{
									ContainerName: "istio-proxy",
									Resource: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("200m"),
										v1.ResourceMemory: resource.MustParse("200Mi"),
									},
								},
							},
						},
					},
				},
			},
			want: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "container",
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("200m"),
									v1.ResourceMemory: resource.MustParse("200Mi"),
								},
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("600m"),
									v1.ResourceMemory: resource.MustParse("200Mi"),
								},
							},
						},
						{
							Name: "istio-proxy",
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("200m"),
									v1.ResourceMemory: resource.MustParse("200Mi"),
								},
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("800m"),
									v1.ResourceMemory: resource.MustParse("400Mi"),
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Tortoise is Auto; some recommendation isn't found",
			args: args{
				pod: &v1.Pod{
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name: "container",
								Resources: v1.ResourceRequirements{
									Requests: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("100m"),
										v1.ResourceMemory: resource.MustParse("100Mi"),
									},
									Limits: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("300m"),
										v1.ResourceMemory: resource.MustParse("100Mi"),
									},
								},
							},
							{
								Name: "istio-proxy",
								Resources: v1.ResourceRequirements{
									Requests: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("100m"),
										v1.ResourceMemory: resource.MustParse("100Mi"),
									},
									Limits: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("400m"),
										v1.ResourceMemory: resource.MustParse("200Mi"),
									},
								},
							},
						},
					},
				},
				tortoise: &v1beta3.Tortoise{
					Spec: v1beta3.TortoiseSpec{
						UpdateMode: v1beta3.UpdateModeAuto,
					},
					Status: v1beta3.TortoiseStatus{
						TortoisePhase: v1beta3.TortoisePhaseWorking,
						Conditions: v1beta3.Conditions{
							ContainerResourceRequests: []v1beta3.ContainerResourceRequests{
								{
									ContainerName: "istio-proxy",
									Resource: v1.ResourceList{
										v1.ResourceCPU: resource.MustParse("200m"),
									},
								},
							},
						},
					},
				},
			},
			want: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "container",
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("100m"),
									v1.ResourceMemory: resource.MustParse("100Mi"),
								},
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("300m"),
									v1.ResourceMemory: resource.MustParse("100Mi"),
								},
							},
						},
						{
							Name: "istio-proxy",
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("200m"),
									v1.ResourceMemory: resource.MustParse("100Mi"),
								},
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("800m"),
									v1.ResourceMemory: resource.MustParse("200Mi"),
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Tortoise is Auto; NoScaleDown option",
			args: args{
				opts: []ModifyPodSpecResourceOption{NoScaleDown},
				pod: &v1.Pod{
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name: "container",
								Resources: v1.ResourceRequirements{
									Requests: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("100m"),
										v1.ResourceMemory: resource.MustParse("100Mi"),
									},
									Limits: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("300m"),
										v1.ResourceMemory: resource.MustParse("100Mi"),
									},
								},
							},
						},
					},
				},
				tortoise: &v1beta3.Tortoise{
					Spec: v1beta3.TortoiseSpec{
						UpdateMode: v1beta3.UpdateModeAuto,
					},
					Status: v1beta3.TortoiseStatus{
						TortoisePhase: v1beta3.TortoisePhaseWorking,
						Conditions: v1beta3.Conditions{
							ContainerResourceRequests: []v1beta3.ContainerResourceRequests{
								{
									ContainerName: "container",
									Resource: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("200m"), // scale up
										v1.ResourceMemory: resource.MustParse("1Mi"),  // scale down
									},
								},
							},
						},
					},
				},
			},
			want: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "container",
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("200m"),  // scale up
									v1.ResourceMemory: resource.MustParse("100Mi"), // scale down is ignored
								},
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("600m"),
									v1.ResourceMemory: resource.MustParse("100Mi"),
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Tortoise is Auto; hits resourceLimitMultiplier",
			fields: fields{
				resourceLimitMultiplier: map[string]int64{
					v1.ResourceCPU.String():    3,
					v1.ResourceMemory.String(): 1,
				},
			},
			args: args{
				pod: &v1.Pod{
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name: "container",
								Resources: v1.ResourceRequirements{
									Requests: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("100m"),
										v1.ResourceMemory: resource.MustParse("100Mi"),
									},
									Limits: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("200m"), // 1:2 -> hit resourceLimitMultiplier
										v1.ResourceMemory: resource.MustParse("100Mi"),
									},
								},
							},
							{
								Name: "istio-proxy",
								Resources: v1.ResourceRequirements{
									Requests: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("100m"),
										v1.ResourceMemory: resource.MustParse("100Mi"),
									},
									Limits: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("400m"),
										v1.ResourceMemory: resource.MustParse("200Mi"),
									},
								},
							},
						},
					},
				},
				tortoise: &v1beta3.Tortoise{
					Spec: v1beta3.TortoiseSpec{
						UpdateMode: v1beta3.UpdateModeAuto,
					},
					Status: v1beta3.TortoiseStatus{
						TortoisePhase: v1beta3.TortoisePhaseWorking,
						Conditions: v1beta3.Conditions{
							ContainerResourceRequests: []v1beta3.ContainerResourceRequests{
								{
									ContainerName: "container",
									Resource: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("200m"),
										v1.ResourceMemory: resource.MustParse("200Mi"),
									},
								},
								{
									ContainerName: "istio-proxy",
									Resource: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("200m"),
										v1.ResourceMemory: resource.MustParse("200Mi"),
									},
								},
							},
						},
					},
				},
			},
			want: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "container",
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("200m"),
									v1.ResourceMemory: resource.MustParse("200Mi"),
								},
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("600m"), // Changed to 1:3
									v1.ResourceMemory: resource.MustParse("200Mi"),
								},
							},
						},
						{
							Name: "istio-proxy",
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("200m"),
									v1.ResourceMemory: resource.MustParse("200Mi"),
								},
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("800m"),
									v1.ResourceMemory: resource.MustParse("400Mi"),
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Tortoise is Auto; hit minimumCPULimit",
			fields: fields{
				minimumCPULimit: "700m",
			},
			args: args{
				pod: &v1.Pod{
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name: "container",
								Resources: v1.ResourceRequirements{
									Requests: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("100m"),
										v1.ResourceMemory: resource.MustParse("100Mi"),
									},
									Limits: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("300m"),
										v1.ResourceMemory: resource.MustParse("100Mi"),
									},
								},
							},
							{
								Name: "istio-proxy",
								Resources: v1.ResourceRequirements{
									Requests: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("100m"),
										v1.ResourceMemory: resource.MustParse("100Mi"),
									},
									Limits: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("400m"),
										v1.ResourceMemory: resource.MustParse("200Mi"),
									},
								},
							},
						},
					},
				},
				tortoise: &v1beta3.Tortoise{
					Spec: v1beta3.TortoiseSpec{
						UpdateMode: v1beta3.UpdateModeAuto,
					},
					Status: v1beta3.TortoiseStatus{
						TortoisePhase: v1beta3.TortoisePhaseWorking,
						Conditions: v1beta3.Conditions{
							ContainerResourceRequests: []v1beta3.ContainerResourceRequests{
								{
									ContainerName: "container",
									Resource: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("200m"),
										v1.ResourceMemory: resource.MustParse("200Mi"),
									},
								},
								{
									ContainerName: "istio-proxy",
									Resource: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("200m"),
										v1.ResourceMemory: resource.MustParse("200Mi"),
									},
								},
							},
						},
					},
				},
			},
			want: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "container",
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("200m"),
									v1.ResourceMemory: resource.MustParse("200Mi"),
								},
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("700m"), // 600m hits the minimumCPULimit
									v1.ResourceMemory: resource.MustParse("200Mi"),
								},
							},
						},
						{
							Name: "istio-proxy",
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("200m"),
									v1.ResourceMemory: resource.MustParse("200Mi"),
								},
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("800m"),
									v1.ResourceMemory: resource.MustParse("400Mi"),
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Tortoise is Auto; GOMEMLIMIT and GOMAXPROCS are updated based on the recommendation",
			fields: fields{
				featureFlags: []features.FeatureFlag{features.GoMemLimitModificationEnabled},
			},
			args: args{
				pod: &v1.Pod{
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name: "container",
								Env: []v1.EnvVar{
									{
										Name:  "GOMAXPROCS",
										Value: "1",
									},
									{
										Name:  "GOMEMLIMIT",
										Value: "100MiB",
									},
								},
								Resources: v1.ResourceRequirements{
									Requests: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("100m"),
										v1.ResourceMemory: resource.MustParse("100Mi"),
									},
									Limits: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("300m"),
										v1.ResourceMemory: resource.MustParse("100Mi"),
									},
								},
							},
							{
								Name: "istio-proxy",
								Env: []v1.EnvVar{
									{
										Name:  "GOMAXPROCS",
										Value: "1",
									},
									{
										Name:  "GOMEMLIMIT",
										Value: "100MiB",
									},
								},
								Resources: v1.ResourceRequirements{
									Requests: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("100m"),
										v1.ResourceMemory: resource.MustParse("100Mi"),
									},
									Limits: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("400m"),
										v1.ResourceMemory: resource.MustParse("200Mi"),
									},
								},
							},
						},
					},
				},
				tortoise: &v1beta3.Tortoise{
					Spec: v1beta3.TortoiseSpec{
						UpdateMode: v1beta3.UpdateModeAuto,
					},
					Status: v1beta3.TortoiseStatus{
						TortoisePhase: v1beta3.TortoisePhaseWorking,
						Conditions: v1beta3.Conditions{
							ContainerResourceRequests: []v1beta3.ContainerResourceRequests{
								{
									ContainerName: "container",
									Resource: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("200m"),
										v1.ResourceMemory: resource.MustParse("200Mi"),
									},
								},
								{
									ContainerName: "istio-proxy",
									Resource: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("50m"),
										v1.ResourceMemory: resource.MustParse("2000Mi"),
									},
								},
							},
						},
					},
				},
			},
			want: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "container",
							Env: []v1.EnvVar{
								{
									Name:  "GOMAXPROCS",
									Value: "2",
								},
								{
									Name:  "GOMEMLIMIT",
									Value: strconv.Itoa(int(ptr.To(resource.MustParse("200Mi")).Value())),
								},
							},
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("200m"),
									v1.ResourceMemory: resource.MustParse("200Mi"),
								},
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("600m"),
									v1.ResourceMemory: resource.MustParse("200Mi"),
								},
							},
						},
						{
							Name: "istio-proxy",
							Env: []v1.EnvVar{
								{
									Name:  "GOMAXPROCS",
									Value: "1", // It wants to be 0.5, but GOMAXPROCS should be an integer. So, we ceil it to 1.
								},
								{
									Name:  "GOMEMLIMIT",
									Value: strconv.Itoa(int(ptr.To(resource.MustParse("2000Mi")).Value())),
								},
							},
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("50m"),
									v1.ResourceMemory: resource.MustParse("2000Mi"),
								},
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("200m"),
									v1.ResourceMemory: resource.MustParse("4000Mi"),
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Tortoise is Auto; GOMEMLIMIT is ignored if no feature flag",
			args: args{
				pod: &v1.Pod{
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name: "container",
								Env: []v1.EnvVar{
									{
										Name:  "GOMAXPROCS",
										Value: "1",
									},
									{
										Name:  "GOMEMLIMIT",
										Value: "100MiB",
									},
								},
								Resources: v1.ResourceRequirements{
									Requests: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("100m"),
										v1.ResourceMemory: resource.MustParse("100Mi"),
									},
									Limits: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("300m"),
										v1.ResourceMemory: resource.MustParse("100Mi"),
									},
								},
							},
							{
								Name: "istio-proxy",
								Env: []v1.EnvVar{
									{
										Name:  "GOMAXPROCS",
										Value: "1",
									},
									{
										Name:  "GOMEMLIMIT",
										Value: "100MiB",
									},
								},
								Resources: v1.ResourceRequirements{
									Requests: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("100m"),
										v1.ResourceMemory: resource.MustParse("100Mi"),
									},
									Limits: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("400m"),
										v1.ResourceMemory: resource.MustParse("200Mi"),
									},
								},
							},
						},
					},
				},
				tortoise: &v1beta3.Tortoise{
					Spec: v1beta3.TortoiseSpec{
						UpdateMode: v1beta3.UpdateModeAuto,
					},
					Status: v1beta3.TortoiseStatus{
						TortoisePhase: v1beta3.TortoisePhaseWorking,
						Conditions: v1beta3.Conditions{
							ContainerResourceRequests: []v1beta3.ContainerResourceRequests{
								{
									ContainerName: "container",
									Resource: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("200m"),
										v1.ResourceMemory: resource.MustParse("200Mi"),
									},
								},
								{
									ContainerName: "istio-proxy",
									Resource: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("50m"),
										v1.ResourceMemory: resource.MustParse("2000Mi"),
									},
								},
							},
						},
					},
				},
			},
			want: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "container",
							Env: []v1.EnvVar{
								{
									Name:  "GOMAXPROCS",
									Value: "2",
								},
								{
									Name:  "GOMEMLIMIT",
									Value: "100MiB",
								},
							},
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("200m"),
									v1.ResourceMemory: resource.MustParse("200Mi"),
								},
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("600m"),
									v1.ResourceMemory: resource.MustParse("200Mi"),
								},
							},
						},
						{
							Name: "istio-proxy",
							Env: []v1.EnvVar{
								{
									Name:  "GOMAXPROCS",
									Value: "1", // It wants to be 0.5, but GOMAXPROCS should be an integer. So, we ceil it to 1.
								},
								{
									Name:  "GOMEMLIMIT",
									Value: "100MiB",
								},
							},
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("50m"),
									v1.ResourceMemory: resource.MustParse("2000Mi"),
								},
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("200m"),
									v1.ResourceMemory: resource.MustParse("4000Mi"),
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := New(tt.fields.resourceLimitMultiplier, tt.fields.minimumCPULimit, nil, tt.fields.featureFlags)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			got := tt.args.pod.DeepCopy()
			s.ModifyPodSpecResource(&got.Spec, tt.args.tortoise, tt.args.opts...)
			if d := cmp.Diff(got.Spec, tt.want.Spec); d != "" {
				t.Errorf("ModifyPodResource() mismatch (-want +got):\n%s", d)
			}
		})
	}
}
