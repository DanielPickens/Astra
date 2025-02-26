package binding

import (
	"testing"

	"github.com/devfile/library/v2/pkg/devfile/parser"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"

	"github\.com/danielpickens/astra/pkg/api"
	"github\.com/danielpickens/astra/pkg/kclient"
	bindingApis "github.com/daniel-pickens/service-binding-operator/apis"
	"github.com/daniel-pickens/service-binding-operator/apis/binding/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var apiServiceBinding = api.ServiceBinding{
	Name: "my-nodejs-app-cluster-sample",
	Spec: api.ServiceBindingSpec{
		Application: api.ServiceBindingReference{
			Name:       "my-nodejs-app-app",
			APIVersion: deploymentApiVersion,
			Kind:       deploymentKind,
		},
		Services: []api.ServiceBindingReference{
			{
				Name:       "cluster-sample",
				APIVersion: ClusterAPIVersion,
				Kind:       clusterKind,
			},
		},
		BindAsFiles:            true,
		DetectBindingResources: true,
	},
}

var bindingServiceBinding = kclient.NewServiceBindingObject(
	"my-nodejs-app-cluster-sample",
	true,
	"my-nodejs-app-app",
	"",
	deploymentGVK,
	nil,
	[]v1alpha1.Service{
		{
			NamespacedRef: v1alpha1.NamespacedRef{
				Ref: v1alpha1.Ref{
					Group:   clusterGVK.Group,
					Version: clusterGVK.Version,
					Kind:    clusterKind,
					Name:    "cluster-sample",
				},
			},
		},
	},
	v1alpha1.ServiceBindingStatus{
		Conditions: []metav1.Condition{
			{
				Type:   bindingApis.InjectionReady,
				Status: metav1.ConditionTrue,
			},
		},
		Secret: "asecret",
	},
)

var sbSecret = corev1.Secret{
	Data: map[string][]byte{
		"akey": []byte("avalue"),
	},
}

func TestBindingClient_ListAllBindings(t *testing.T) {
	bindingServiceBinding.SetLabels(map[string]string{
		"astra.dev/mode": "Dev",
	})
	type fields struct {
		kubernetesClient func(ctrl *gomock.Controller) kclient.ClientInterface
	}
	type args struct {
		devfileObj *parser.DevfileObj
		context    string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          []api.ServiceBinding
		wantInDevfile []string
		wantErr       bool
	}{
		{
			name: "a servicebinding defined in Devfile, nothing in cluster",
			fields: fields{
				func(ctrl *gomock.Controller) kclient.ClientInterface {
					client := kclient.NewMockClientInterface(ctrl)
					client.EXPECT().ListServiceBindingsFromAllGroups().Return(nil, nil, nil)
					client.EXPECT().GetBindingServiceBinding(gomock.Any()).Return(
						v1alpha1.ServiceBinding{},
						errors.NewNotFound(
							schema.GroupResource{
								Group:    "dont",
								Resource: "care",
							},
							"my-nodejs-app-cluster-sample",
						),
					)
					return client
				},
			}, args: args{
				devfileObj: getDevfileObjWithServiceBinding("aname", "", true, ""),
				context:    "/apath",
			},
			want:          []api.ServiceBinding{apiServiceBinding},
			wantInDevfile: []string{"my-nodejs-app-cluster-sample"},
		},
		{
			name: "a servicebinding defined in Devfile, also in cluster",
			fields: fields{
				func(ctrl *gomock.Controller) kclient.ClientInterface {
					client := kclient.NewMockClientInterface(ctrl)
					client.EXPECT().ListServiceBindingsFromAllGroups().Return(nil, []v1alpha1.ServiceBinding{
						*bindingServiceBinding,
					}, nil)
					client.EXPECT().GetBindingServiceBinding(gomock.Any()).Return(
						*bindingServiceBinding,
						nil,
					).AnyTimes()
					client.EXPECT().GetCurrentNamespace().Return("anamespace").AnyTimes()
					client.EXPECT().GetSecret("asecret", "anamespace").Return(&sbSecret, nil).AnyTimes()
					return client
				},
			}, args: args{
				devfileObj: getDevfileObjWithServiceBinding("aname", "", true, ""),
				context:    "/apath",
			},
			want: []api.ServiceBinding{
				{
					Name: "my-nodejs-app-cluster-sample",
					Spec: api.ServiceBindingSpec{
						Application: api.ServiceBindingReference{
							Name:       "my-nodejs-app-app",
							APIVersion: deploymentApiVersion,
							Kind:       deploymentKind,
						},
						Services: []api.ServiceBindingReference{
							{
								Name:       "cluster-sample",
								APIVersion: ClusterAPIVersion,
								Kind:       clusterKind,
							},
						},
						BindAsFiles:            true,
						DetectBindingResources: true,
					},
					Status: &api.ServiceBindingStatus{
						BindingFiles: []string{"${SERVICE_BINDING_ROOT}/my-nodejs-app-cluster-sample/akey"},
						RunningIn:    api.RunningModes{"dev": true, "deploy": false},
					},
				},
			},
			wantInDevfile: []string{"my-nodejs-app-cluster-sample"},
		},
		{
			name: "a servicebinding defined in cluster",
			fields: fields{
				func(ctrl *gomock.Controller) kclient.ClientInterface {
					client := kclient.NewMockClientInterface(ctrl)
					client.EXPECT().ListServiceBindingsFromAllGroups().Return(nil, []v1alpha1.ServiceBinding{
						*bindingServiceBinding,
					}, nil)
					client.EXPECT().GetBindingServiceBinding(gomock.Any()).Return(
						*bindingServiceBinding,
						nil,
					).AnyTimes()
					client.EXPECT().GetCurrentNamespace().Return("anamespace").AnyTimes()
					client.EXPECT().GetSecret("asecret", "anamespace").Return(&sbSecret, nil).AnyTimes()
					return client
				},
			}, args: args{},
			want: []api.ServiceBinding{
				{
					Name: "my-nodejs-app-cluster-sample",
					Spec: api.ServiceBindingSpec{
						Application: api.ServiceBindingReference{
							Name:       "my-nodejs-app-app",
							APIVersion: deploymentApiVersion,
							Kind:       deploymentKind,
						},
						Services: []api.ServiceBindingReference{
							{
								Name:       "cluster-sample",
								APIVersion: ClusterAPIVersion,
								Kind:       clusterKind,
							},
						},
						BindAsFiles:            true,
						DetectBindingResources: true,
					},
					Status: &api.ServiceBindingStatus{
						BindingFiles: []string{"${SERVICE_BINDING_ROOT}/my-nodejs-app-cluster-sample/akey"},
						RunningIn:    api.RunningModes{"dev": true, "deploy": false},
					},
				},
			},
			wantInDevfile: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			o := &BindingClient{
				kubernetesClient: tt.fields.kubernetesClient(ctrl),
			}
			got, gotInDevfile, err := o.ListAllBindings(tt.args.devfileObj, tt.args.context)
			if (err != nil) != tt.wantErr {
				t.Errorf("BindingClient.ListAllBindings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("BindingClient.ListAllBindings() mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantInDevfile, gotInDevfile); diff != "" {
				t.Errorf("BindingClient.ListAllBindings() wantInDevfile mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
