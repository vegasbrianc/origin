package v1_test

import (
	"reflect"
	"testing"

	kapi "k8s.io/kubernetes/pkg/api"
	// required to register defaulting functions for containers
	_ "k8s.io/kubernetes/pkg/api/install"
	kapiv1 "k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/util"
	"k8s.io/kubernetes/pkg/util/intstr"

	v1 "github.com/openshift/origin/pkg/api/v1"
	deployapi "github.com/openshift/origin/pkg/deploy/api"
	_ "github.com/openshift/origin/pkg/deploy/api/install"
	deployv1 "github.com/openshift/origin/pkg/deploy/api/v1"
)

func mkintp(i int64) *int64 {
	return &i
}

func TestDefaults(t *testing.T) {
	defaultIntOrString := intstr.FromString("25%")
	differentIntOrString := intstr.FromInt(5)
	tests := []struct {
		original *deployv1.DeploymentConfig
		expected *deployv1.DeploymentConfig
	}{
		{
			original: &deployv1.DeploymentConfig{},
			expected: &deployv1.DeploymentConfig{
				Spec: deployv1.DeploymentConfigSpec{
					Strategy: deployv1.DeploymentStrategy{
						Type: deployv1.DeploymentStrategyTypeRolling,
						RollingParams: &deployv1.RollingDeploymentStrategyParams{
							UpdatePeriodSeconds: newInt64(deployapi.DefaultRollingUpdatePeriodSeconds),
							IntervalSeconds:     newInt64(deployapi.DefaultRollingIntervalSeconds),
							TimeoutSeconds:      newInt64(deployapi.DefaultRollingTimeoutSeconds),
							MaxSurge:            &defaultIntOrString,
							MaxUnavailable:      &defaultIntOrString,
						},
					},
					Triggers: []deployv1.DeploymentTriggerPolicy{
						{
							Type: deployv1.DeploymentTriggerOnConfigChange,
						},
					},
				},
			},
		},
		{
			original: &deployv1.DeploymentConfig{
				Spec: deployv1.DeploymentConfigSpec{
					Strategy: deployv1.DeploymentStrategy{
						Type: deployv1.DeploymentStrategyTypeRecreate,
						RecreateParams: &deployv1.RecreateDeploymentStrategyParams{
							TimeoutSeconds: newInt64(deployapi.DefaultRollingTimeoutSeconds),
							Pre: &deployv1.LifecycleHook{
								TagImages: []deployv1.TagImageHook{{}, {}},
							},
							Mid: &deployv1.LifecycleHook{
								TagImages: []deployv1.TagImageHook{{}, {}},
							},
							Post: &deployv1.LifecycleHook{
								TagImages: []deployv1.TagImageHook{{}, {}},
							},
						},
						RollingParams: &deployv1.RollingDeploymentStrategyParams{
							Pre: &deployv1.LifecycleHook{
								TagImages: []deployv1.TagImageHook{{}, {}},
							},
							Post: &deployv1.LifecycleHook{
								TagImages: []deployv1.TagImageHook{{}, {}},
							},
							UpdatePeriodSeconds: newInt64(5),
							IntervalSeconds:     newInt64(6),
							TimeoutSeconds:      newInt64(7),
							MaxSurge:            &differentIntOrString,
							MaxUnavailable:      &differentIntOrString,
						},
					},
					Triggers: []deployv1.DeploymentTriggerPolicy{
						{
							Type: deployv1.DeploymentTriggerOnImageChange,
						},
					},
					Template: &kapiv1.PodTemplateSpec{
						Spec: kapiv1.PodSpec{
							Containers: []kapiv1.Container{
								{
									Name: "test",
								},
							},
						},
					},
				},
			},
			expected: &deployv1.DeploymentConfig{
				Spec: deployv1.DeploymentConfigSpec{
					Strategy: deployv1.DeploymentStrategy{
						Type: deployv1.DeploymentStrategyTypeRecreate,
						RecreateParams: &deployv1.RecreateDeploymentStrategyParams{
							TimeoutSeconds: newInt64(deployapi.DefaultRollingTimeoutSeconds),
							Pre: &deployv1.LifecycleHook{
								TagImages: []deployv1.TagImageHook{{ContainerName: "test"}, {ContainerName: "test"}},
							},
							Mid: &deployv1.LifecycleHook{
								TagImages: []deployv1.TagImageHook{{ContainerName: "test"}, {ContainerName: "test"}},
							},
							Post: &deployv1.LifecycleHook{
								TagImages: []deployv1.TagImageHook{{ContainerName: "test"}, {ContainerName: "test"}},
							},
						},
						RollingParams: &deployv1.RollingDeploymentStrategyParams{
							Pre: &deployv1.LifecycleHook{
								TagImages: []deployv1.TagImageHook{{ContainerName: "test"}, {ContainerName: "test"}},
							},
							Post: &deployv1.LifecycleHook{
								TagImages: []deployv1.TagImageHook{{ContainerName: "test"}, {ContainerName: "test"}},
							},
							UpdatePeriodSeconds: newInt64(5),
							IntervalSeconds:     newInt64(6),
							TimeoutSeconds:      newInt64(7),
							MaxSurge:            &differentIntOrString,
							MaxUnavailable:      &differentIntOrString,
						},
					},
					Triggers: []deployv1.DeploymentTriggerPolicy{
						{
							Type: deployv1.DeploymentTriggerOnImageChange,
						},
					},
					Template: &kapiv1.PodTemplateSpec{
						Spec: kapiv1.PodSpec{
							SecurityContext:               &kapiv1.PodSecurityContext{},
							RestartPolicy:                 kapiv1.RestartPolicyAlways,
							TerminationGracePeriodSeconds: mkintp(30),
							DNSPolicy:                     kapiv1.DNSClusterFirst,
							Containers: []kapiv1.Container{
								{
									Name: "test",
									TerminationMessagePath: "/dev/termination-log",
									ImagePullPolicy:        kapiv1.PullAlways,
								},
							},
						},
					},
				},
			},
		},
		{
			original: &deployv1.DeploymentConfig{
				Spec: deployv1.DeploymentConfigSpec{
					Strategy: deployv1.DeploymentStrategy{
						Type: deployv1.DeploymentStrategyTypeRolling,
						RollingParams: &deployv1.RollingDeploymentStrategyParams{
							UpdatePeriodSeconds: newInt64(5),
							IntervalSeconds:     newInt64(6),
							TimeoutSeconds:      newInt64(7),
							UpdatePercent:       newInt(50),
						},
					},
					Triggers: []deployv1.DeploymentTriggerPolicy{
						{
							Type: deployv1.DeploymentTriggerOnImageChange,
						},
					},
				},
			},
			expected: &deployv1.DeploymentConfig{
				Spec: deployv1.DeploymentConfigSpec{
					Strategy: deployv1.DeploymentStrategy{
						Type: deployv1.DeploymentStrategyTypeRolling,
						RollingParams: &deployv1.RollingDeploymentStrategyParams{
							UpdatePeriodSeconds: newInt64(5),
							IntervalSeconds:     newInt64(6),
							TimeoutSeconds:      newInt64(7),
							UpdatePercent:       newInt(50),
							MaxSurge:            newIntOrString(intstr.FromString("50%")),
							MaxUnavailable:      newIntOrString(intstr.FromInt(0)),
						},
					},
					Triggers: []deployv1.DeploymentTriggerPolicy{
						{
							Type: deployv1.DeploymentTriggerOnImageChange,
						},
					},
				},
			},
		},
		{
			original: &deployv1.DeploymentConfig{
				Spec: deployv1.DeploymentConfigSpec{
					Strategy: deployv1.DeploymentStrategy{
						Type: deployv1.DeploymentStrategyTypeRolling,
						RollingParams: &deployv1.RollingDeploymentStrategyParams{
							UpdatePeriodSeconds: newInt64(5),
							IntervalSeconds:     newInt64(6),
							TimeoutSeconds:      newInt64(7),
							UpdatePercent:       newInt(-25),
						},
					},
					Triggers: []deployv1.DeploymentTriggerPolicy{
						{
							Type: deployv1.DeploymentTriggerOnImageChange,
						},
					},
				},
			},
			expected: &deployv1.DeploymentConfig{
				Spec: deployv1.DeploymentConfigSpec{
					Strategy: deployv1.DeploymentStrategy{
						Type: deployv1.DeploymentStrategyTypeRolling,
						RollingParams: &deployv1.RollingDeploymentStrategyParams{
							UpdatePeriodSeconds: newInt64(5),
							IntervalSeconds:     newInt64(6),
							TimeoutSeconds:      newInt64(7),
							UpdatePercent:       newInt(-25),
							MaxSurge:            newIntOrString(intstr.FromInt(0)),
							MaxUnavailable:      newIntOrString(intstr.FromString("25%")),
						},
					},
					Triggers: []deployv1.DeploymentTriggerPolicy{
						{
							Type: deployv1.DeploymentTriggerOnImageChange,
						},
					},
				},
			},
		},
	}

	for i, test := range tests {
		t.Logf("test %d", i)
		original := test.original
		expected := test.expected
		obj2 := roundTrip(t, runtime.Object(original))
		got, ok := obj2.(*deployv1.DeploymentConfig)
		if !ok {
			t.Errorf("unexpected object: %v", got)
			t.FailNow()
		}
		if !reflect.DeepEqual(got.Spec, expected.Spec) {
			t.Errorf("got different than expected:\nA:\t%#v\nB:\t%#v\n\nDiff:\n%s\n\n%s", got, expected, util.ObjectDiff(expected, got), util.ObjectGoPrintSideBySide(expected, got))
		}
	}
}

func roundTrip(t *testing.T, obj runtime.Object) runtime.Object {
	data, err := runtime.Encode(kapi.Codecs.LegacyCodec(v1.SchemeGroupVersion), obj)
	if err != nil {
		t.Errorf("%v\n %#v", err, obj)
		return nil
	}
	obj2, err := runtime.Decode(kapi.Codecs.UniversalDecoder(), data)
	if err != nil {
		t.Errorf("%v\nData: %s\nSource: %#v", err, string(data), obj)
		return nil
	}
	obj3 := reflect.New(reflect.TypeOf(obj).Elem()).Interface().(runtime.Object)
	err = kapi.Scheme.Convert(obj2, obj3)
	if err != nil {
		t.Errorf("%v\nSource: %#v", err, obj2)
		return nil
	}
	return obj3
}

func newInt64(val int64) *int64 {
	return &val
}

func newInt(val int) *int {
	return &val
}

func newIntOrString(ios intstr.IntOrString) *intstr.IntOrString {
	return &ios
}
