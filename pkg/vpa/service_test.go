package vpa

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mercari/tortoise/api/v1beta1"
	autoscalingv1beta1 "github.com/mercari/tortoise/api/v1beta1"
)

func TestMakeAllVerticalContainerResourcePhaseRunning(t *testing.T) {
	type args struct {
		tortoise *autoscalingv1beta1.Tortoise
	}
	tests := []struct {
		name string
		args args
		want *autoscalingv1beta1.Tortoise
	}{
		{
			name: "test1",
			args: args{
				tortoise: &autoscalingv1beta1.Tortoise{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "tortoise",
						Namespace: "default",
					},
					Spec: autoscalingv1beta1.TortoiseSpec{
						ResourcePolicy: []autoscalingv1beta1.ContainerResourcePolicy{
							{
								ContainerName: "app",
								AutoscalingPolicy: map[v1.ResourceName]v1beta1.AutoscalingType{
									v1.ResourceMemory: v1beta1.AutoscalingTypeVertical,
									v1.ResourceCPU:    v1beta1.AutoscalingTypeHorizontal,
								},
							},
							{
								ContainerName: "istio-proxy",
								AutoscalingPolicy: map[v1.ResourceName]v1beta1.AutoscalingType{
									v1.ResourceMemory: v1beta1.AutoscalingTypeVertical,
									v1.ResourceCPU:    v1beta1.AutoscalingTypeHorizontal,
								},
							},
						},
					},
					Status: autoscalingv1beta1.TortoiseStatus{
						ContainerResourcePhases: []autoscalingv1beta1.ContainerResourcePhases{
							{
								ContainerName: "app",
								ResourcePhases: map[v1.ResourceName]autoscalingv1beta1.ContainerResourcePhase{
									v1.ResourceCPU:    autoscalingv1beta1.ContainerResourcePhaseGatheringData,
									v1.ResourceMemory: autoscalingv1beta1.ContainerResourcePhaseGatheringData,
								},
							},
							{
								ContainerName: "istio-proxy",
								ResourcePhases: map[v1.ResourceName]autoscalingv1beta1.ContainerResourcePhase{
									v1.ResourceCPU:    autoscalingv1beta1.ContainerResourcePhaseGatheringData,
									v1.ResourceMemory: autoscalingv1beta1.ContainerResourcePhaseGatheringData,
								},
							},
						},
					},
				},
			},
			want: &autoscalingv1beta1.Tortoise{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "tortoise",
					Namespace: "default",
				},
				Spec: autoscalingv1beta1.TortoiseSpec{
					ResourcePolicy: []autoscalingv1beta1.ContainerResourcePolicy{
						{
							ContainerName: "app",
							AutoscalingPolicy: map[v1.ResourceName]v1beta1.AutoscalingType{
								v1.ResourceMemory: v1beta1.AutoscalingTypeVertical,
								v1.ResourceCPU:    v1beta1.AutoscalingTypeHorizontal,
							},
						},
						{
							ContainerName: "istio-proxy",
							AutoscalingPolicy: map[v1.ResourceName]v1beta1.AutoscalingType{
								v1.ResourceMemory: v1beta1.AutoscalingTypeVertical,
								v1.ResourceCPU:    v1beta1.AutoscalingTypeHorizontal,
							},
						},
					},
				},
				Status: autoscalingv1beta1.TortoiseStatus{
					ContainerResourcePhases: []autoscalingv1beta1.ContainerResourcePhases{
						{
							ContainerName: "app",
							ResourcePhases: map[v1.ResourceName]autoscalingv1beta1.ContainerResourcePhase{
								v1.ResourceCPU:    autoscalingv1beta1.ContainerResourcePhaseGatheringData,
								v1.ResourceMemory: autoscalingv1beta1.ContainerResourcePhaseWorking,
							},
						},
						{
							ContainerName: "istio-proxy",
							ResourcePhases: map[v1.ResourceName]autoscalingv1beta1.ContainerResourcePhase{
								v1.ResourceCPU:    autoscalingv1beta1.ContainerResourcePhaseGatheringData,
								v1.ResourceMemory: autoscalingv1beta1.ContainerResourcePhaseWorking,
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MakeAllVerticalContainerResourcePhaseWorking(tt.args.tortoise)

			// use diff to compare
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Fatalf("MakeAllVerticalContainerResourcePhaseRunning() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}