package storage

import astralabels "github\.com/danielpickens/astra/pkg/labels"

func getStorageLabels(storageName, componentName, applicationName string) map[string]string {
	labels := astralabels.GetLabels(componentName, applicationName, "", astralabels.ComponentDevMode, false)
	astralabels.AddStorageInfo(labels, storageName, false)
	return labels
}

/*
func TestPush(t *testing.T) {
	localStorage0 := localConfigProvider.LocalStorage{
		Name:      "storage-0",
		Size:      "1Gi",
		Path:      "/data",
		Container: "runtime-0",
		Ephemeral: util.GetBool(false),
	}
	localStorage1 := localConfigProvider.LocalStorage{
		Name:      "storage-1",
		Size:      "5Gi",
		Path:      "/path",
		Container: "runtime-1",
		Ephemeral: util.GetBool(false),
	}
	localEphemeralStorage0 := localConfigProvider.LocalStorage{
		Name:      "ephemeral-storage-0",
		Size:      "5Gi",
		Path:      "/path",
		Container: "runtime-1",
		Ephemeral: util.GetBool(true),
	}

	clusterStorage0 := NewStorageWithContainer("storage-0", "1Gi", "/data", "runtime-0", util.GetBool(false))
	clusterStorage1 := NewStorageWithContainer("storage-1", "5Gi", "/path", "runtime-1", util.GetBool(false))

	tests := []struct {
		name                string
		returnedFromLocal   []localConfigProvider.LocalStorage
		returnedFromCluster StorageList
		createdItems        []localConfigProvider.LocalStorage
		deletedItems        []string
		wantErr             bool
		wantEphemeralNames  []string
	}{
		{
			name:                "case 1: no storage in both local and cluster",
			returnedFromLocal:   []localConfigProvider.LocalStorage{},
			returnedFromCluster: StorageList{},
			wantEphemeralNames:  []string{},
		},
		{
			name:                "case 2: two persistent storage in local and no on cluster",
			returnedFromLocal:   []localConfigProvider.LocalStorage{localStorage0, localStorage1},
			returnedFromCluster: StorageList{},
			createdItems: []localConfigProvider.LocalStorage{
				{
					Name:      "storage-0",
					Size:      "1Gi",
					Path:      "/data",
					Container: "runtime-0",
					Ephemeral: util.GetBool(false),
				},
				{
					Name:      "storage-1",
					Size:      "5Gi",
					Path:      "/path",
					Container: "runtime-1",
					Ephemeral: util.GetBool(false),
				},
			},
			wantEphemeralNames: []string{},
		},
		{
			name:              "case 3: 0 persistent storage in local and two on cluster",
			returnedFromLocal: []localConfigProvider.LocalStorage{},
			returnedFromCluster: StorageList{
				Items: []Storage{clusterStorage0, clusterStorage1},
			},
			createdItems:       []localConfigProvider.LocalStorage{},
			deletedItems:       []string{"storage-0", "storage-1"},
			wantEphemeralNames: []string{},
		},
		{
			name:              "case 4: same two persistent storage in local and cluster",
			returnedFromLocal: []localConfigProvider.LocalStorage{localStorage0, localStorage1},
			returnedFromCluster: StorageList{
				Items: []Storage{clusterStorage0, clusterStorage1},
			},
			createdItems:       []localConfigProvider.LocalStorage{},
			deletedItems:       []string{},
			wantEphemeralNames: []string{},
		},
		{
			name: "case 5: two persistent storage in both local and cluster but two of them are different and the other two are same",
			returnedFromLocal: []localConfigProvider.LocalStorage{localStorage0,
				{
					Name:      "storage-1-1",
					Size:      "5Gi",
					Path:      "/path",
					Container: "runtime-1",
					Ephemeral: util.GetBool(false),
				},
			},
			returnedFromCluster: StorageList{
				Items: []Storage{
					clusterStorage0,
					clusterStorage1,
				},
			},
			createdItems: []localConfigProvider.LocalStorage{
				{
					Name:      "storage-1-1",
					Size:      "5Gi",
					Path:      "/path",
					Container: "runtime-1",
					Ephemeral: util.GetBool(false),
				},
			},
			deletedItems:       []string{clusterStorage1.Name},
			wantEphemeralNames: []string{},
		},
		{
			name: "case 6: spec mismatch",
			returnedFromLocal: []localConfigProvider.LocalStorage{
				{
					Name:      "storage-1",
					Size:      "3Gi",
					Path:      "/path",
					Container: "runtime-1",
					Ephemeral: util.GetBool(false),
				},
			},
			returnedFromCluster: StorageList{
				Items: []Storage{
					clusterStorage1,
				},
			},
			createdItems:       []localConfigProvider.LocalStorage{},
			deletedItems:       []string{},
			wantErr:            true,
			wantEphemeralNames: []string{},
		},
		{
			name: "case 7: only one PVC created for two storage with same name but on different containers",
			returnedFromLocal: []localConfigProvider.LocalStorage{
				{
					Name:      "storage-0",
					Size:      "1Gi",
					Path:      "/data",
					Container: "runtime-0",
					Ephemeral: util.GetBool(false),
				},
				{
					Name:      "storage-0",
					Size:      "1Gi",
					Path:      "/path",
					Container: "runtime-1",
					Ephemeral: util.GetBool(false),
				},
			},
			returnedFromCluster: StorageList{},
			createdItems: []localConfigProvider.LocalStorage{
				{
					Name:      "storage-0",
					Size:      "1Gi",
					Path:      "/path",
					Container: "runtime-1",
					Ephemeral: util.GetBool(false),
				},
			},
			wantEphemeralNames: []string{},
		},
		{
			name: "case 8: only path spec mismatch",
			returnedFromLocal: []localConfigProvider.LocalStorage{
				{
					Name:      "storage-1",
					Size:      "5Gi",
					Path:      "/data",
					Container: "runtime-1",
					Ephemeral: util.GetBool(false),
				},
			},
			returnedFromCluster: StorageList{
				Items: []Storage{
					clusterStorage1,
				},
			},
			wantEphemeralNames: []string{},
		},
		{
			name:              "case 9: only one PVC deleted for two storage with same name but on different containers",
			returnedFromLocal: []localConfigProvider.LocalStorage{},
			returnedFromCluster: StorageList{
				Items: []Storage{
					NewStorageWithContainer("storage-0", "1Gi", "/data", "runtime-0", util.GetBool(false)),
					NewStorageWithContainer("storage-0", "1Gi", "/data", "runtime-1", util.GetBool(false)),
				},
			},
			deletedItems:       []string{"storage-0"},
			wantEphemeralNames: []string{},
		},
		{
			name:                "case 10: one ephemeral storage in local, none in cluster",
			returnedFromLocal:   []localConfigProvider.LocalStorage{localEphemeralStorage0},
			returnedFromCluster: StorageList{},
			wantEphemeralNames:  []string{"ephemeral-storage-0"},
		},
		{
			name:                "case 11: one persistent + one ephemeral storage in local and no on cluster",
			returnedFromLocal:   []localConfigProvider.LocalStorage{localStorage0, localEphemeralStorage0},
			returnedFromCluster: StorageList{},
			createdItems: []localConfigProvider.LocalStorage{
				{
					Name:      "storage-0",
					Size:      "1Gi",
					Path:      "/data",
					Container: "runtime-0",
					Ephemeral: util.GetBool(false),
				},
			},
			wantEphemeralNames: []string{"ephemeral-storage-0"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fakeStorageClient := NewMockClient(ctrl)
			fakeLocalConfig := localConfigProvider.NewMockLocalConfigProvider(ctrl)

			fakeStorageClient.EXPECT().List().Return(tt.returnedFromCluster, nil).AnyTimes()
			fakeLocalConfig.EXPECT().ListStorage(gomock.Any()).Return(tt.returnedFromLocal, nil).AnyTimes()

			convert := ConvertListLocalToMachine(tt.createdItems)
			for i := range convert.Items {
				fakeStorageClient.EXPECT().Create(convert.Items[i]).Return(nil).Times(1)
			}

			for i := range tt.deletedItems {
				fakeStorageClient.EXPECT().Delete(tt.deletedItems[i]).Return(nil).Times(1)
			}

			ephemerals, err := Push(fakeStorageClient, parser.DevfileObj{}) // feloy:Tastra
			if (err != nil) != tt.wantErr {
				t.Errorf("Push() error = %v, wantErr %v", err, tt.wantErr)
			}
			ephemeralKeys := make([]string, 0, len(ephemerals))
			for k := range ephemerals {
				ephemeralKeys = append(ephemeralKeys, k)
			}
			if diff := cmp.Diff(tt.wantEphemeralNames, ephemeralKeys); diff != "" {
				t.Errorf("Push() ephemeral names mismatch (-want +got):\n%s", diff)
			}
		})
	}
}*/
