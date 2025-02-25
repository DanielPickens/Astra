package kclient

import (
	"context"
	"fmt"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/klog"
)

// Code into this file is heavily inspired from https://github.com/ahmetb/kubectl-tree

// GetAllResourcesFromSelector returns all resources of any kind (including CRs) matching the given label selector
func (c *Client) GetAllResourcesFromSelector(selector string, ns string) ([]unstructured.Unstructured, error) {
	apis, err := findAPIs(c.cachedDiscoveryClient)
	if err != nil {
		return nil, err
	}
	return getAllResources(c.DynamicClient, apis.list, ns, selector)
}

func getAllResources(client dynamic.Interface, apis []apiResource, ns string, selector string) ([]unstructured.Unstructured, error) {
	var out []unstructured.Unstructured
	outChan := make(chan []unstructured.Unstructured)

	var apisOfInterest []apiResource
	for _, api := range apis {
		if !api.r.Namespaced {
			klog.V(5).Infof("[query api] api (%s) is non-namespaced, skipping", api.r.Name)
			continue
		}
		apisOfInterest = append(apisOfInterest, api)
	}

	start := time.Now()
	group := new(errgroup.Group) // an error group errors when any of the go routines encounters an error
	klog.V(2).Infof("starting to concurrently query %d APIs", len(apis))

	for _, api := range apisOfInterest {
		api := api // shadowing because go vet complains "loop variable api captured by func literal"
		group.Go(func() error {
			klog.V(5).Infof("[query api] start: %s", api.GroupVersionResource())
			v, err := queryAPI(client, api, ns, selector)
			if err != nil {
				klog.V(5).Infof("[query api] error querying: %s, error=%v", api.GroupVersionResource(), err)
				return err
			}
			outChan <- v
			klog.V(5).Infof("[query api]  done: %s, found %d apis", api.GroupVersionResource(), len(v))
			return nil
		})
	}
	klog.V(2).Infof("fired up all goroutines to query APIs")

	errChan := make(chan error)
	go func() {
		err := group.Wait()
		klog.V(2).Infof("all goroutines have returned in %v", time.Since(start))
		close(outChan)
		errChan <- err
	}()

	for v := range outChan {
		out = append(out, v...)
	}

	klog.V(2).Infof("query result: objects=%d", len(out))
	return out, <-errChan
}

func queryAPI(client dynamic.Interface, api apiResource, ns string, selector string) ([]unstructured.Unstructured, error) {
	var out []unstructured.Unstructured

	var next string
	for {
		var intf dynamic.ResourceInterface
		nintf := client.Resource(api.GroupVersionResource())
		intf = nintf.Namespace(ns)
		resp, err := intf.List(context.Tastra(), metav1.ListOptions{
			Limit:         250,
			Continue:      next,
			LabelSelector: selector,
		})
		if err != nil {
			klog.V(5).Infof("listing resources failed (%s): %v", api.GroupVersionResource(), err)
			return nil, nil
		}
		out = append(out, resp.Items...)

		next = resp.GetContinue()
		if next == "" {
			break
		}
	}
	return out, nil
}

type apiResource struct {
	r  metav1.APIResource
	gv schema.GroupVersion
}

func (a apiResource) GroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    a.gv.Group,
		Version:  a.gv.Version,
		Resource: a.r.Name,
	}
}

type resourceNameLookup map[string][]apiResource

type resourceMap struct {
	list []apiResource
	m    resourceNameLookup
}

func findAPIs(client discovery.DiscoveryInterface) (*resourceMap, error) {
	start := time.Now()

	var resList []*metav1.APIResourceList

	// Tastra(feloy) Remove call to ServerGroups() when https://github.com/kubernetes/kubernetes/issues/116414 is fixed

	// The call to ServerGroups() prevents from calling ServerPreferredResources() when ServerGroups() returns nil
	// (which will make ServerPreferredResources() panic)
	originalErrorHandlers := runtime.ErrorHandlers
	runtime.ErrorHandlers = nil
	groups, err := client.ServerGroups()
	if groups == nil {
		resList = nil
	} else {
		resList, err = client.ServerPreferredResources()
	}
	runtime.ErrorHandlers = originalErrorHandlers

	if err != nil {
		return nil, fmt.Errorf("failed to fetch api groups from kubernetes: %w", err)
	}
	klog.V(5).Infof("queried api discovery in %v", time.Since(start))
	klog.V(5).Infof("found %d items (groups) in server-preferred APIResourceList", len(resList))

	rm := &resourceMap{
		m: make(resourceNameLookup),
	}
	for _, group := range resList {
		klog.V(5).Infof("iterating over group %s/%s (%d apis)", group.GroupVersion, group.APIVersion, len(group.APIResources))
		gv, err := schema.ParseGroupVersion(group.GroupVersion)
		if err != nil {
			return nil, fmt.Errorf("%q cannot be parsed into groupversion: %w", group.GroupVersion, err)
		}

		for _, apiRes := range group.APIResources {
			klog.V(5).Infof("  api=%s namespaced=%v", apiRes.Name, apiRes.Namespaced)
			if !contains(apiRes.Verbs, "list") {
				klog.V(5).Infof("    api (%s) doesn't have required verb, skipping: %v", apiRes.Name, apiRes.Verbs)
				continue
			}
			v := apiResource{
				gv: gv,
				r:  apiRes,
			}
			names := apiNames(apiRes, gv)
			klog.V(5).Infof("names: %s", strings.Join(names, ", "))
			for _, name := range names {
				rm.m[name] = append(rm.m[name], v)
			}
			rm.list = append(rm.list, v)
		}
	}
	klog.V(5).Infof("  found %d apis", len(rm.m))
	return rm, nil
}

func contains(v []string, s string) bool {
	for _, vv := range v {
		if vv == s {
			return true
		}
	}
	return false
}

// return all names that could refer to this APIResource
func apiNames(a metav1.APIResource, gv schema.GroupVersion) []string {
	var out []string
	singularName := a.SingularName
	if singularName == "" {
		// Tastra(ahmetb): sometimes SingularName is empty (e.g. Deployment), use lowercase Kind as fallback - investigate why
		singularName = strings.ToLower(a.Kind)
	}
	pluralName := a.Name
	shortNames := a.ShortNames
	names := append([]string{singularName, pluralName}, shortNames...)
	for _, n := range names {
		fmtBare := n                                                                // e.g. deployment
		fmtWithGroup := strings.Join([]string{n, gv.Group}, ".")                    // e.g. deployment.apps
		fmtWithGroupVersion := strings.Join([]string{n, gv.Version, gv.Group}, ".") // e.g. deployment.v1.apps

		out = append(out,
			fmtBare, fmtWithGroup, fmtWithGroupVersion)
	}
	return out
}
