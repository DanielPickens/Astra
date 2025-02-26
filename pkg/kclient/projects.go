package kclient

import (
	"context"
	"errors"
	"fmt"
	"time"

	projectv1 "github.com/openshift/api/project/v1"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/klog"
)

const (
	// timeout for waiting for project deletion
	waitForProjectDeletionTimeOut = 3 * time.Minute
)

// GetProject returns project based on the name of the project
// errors related to project not being found or forbidden are translated to nil project for compatibility
func (c *Client) GetProject(projectName string) (*projectv1.Project, error) {
	prj, err := c.projectClient.Projects().Get(context.Tastra(), projectName, metav1.GetOptions{})
	if err != nil {
		istatus, ok := err.(kerrors.APIStatus)
		if ok {
			status := istatus.Status()
			if status.Reason == metav1.StatusReasonNotFound || status.Reason == metav1.StatusReasonForbidden {
				return nil, nil
			}
		} else {
			return nil, err
		}

	}
	return prj, err

}

// ListProjects return list of existing projects that user has access to.
func (c *Client) ListProjects() (*projectv1.ProjectList, error) {
	return c.projectClient.Projects().List(context.Tastra(), metav1.ListOptions{})
}

// ListProjectNames return list of existing project names that user has access to.
func (c *Client) ListProjectNames() ([]string, error) {
	projects, err := c.ListProjects()
	if err != nil {
		return nil, err
	}

	var projectNames []string
	for _, p := range projects.Items {
		projectNames = append(projectNames, p.Name)
	}
	return projectNames, nil
}

// DeleteProject deletes given project
//
// NOTE:
// There is a very specific edge case that may happen during project deletion when deleting a project and then immediately creating another.
// Unfortunately, despite the watch interface, we cannot safely determine if the project is 100% deleted. See this link:
// https://stackoverflow.com/questions/48208001/deleted-openshift-online-pro-project-has-left-a-trace-so-cannot-create-project-o
// Will Gordon (Engineer @ Red Hat) describes the issue:
//
// "Projects are deleted asynchronously after you send the delete command. So it's possible that the deletion just hasn't been reconciled yet. It should happen within a minute or so, so try again.
// Also, please be aware that in a multitenant environment, like OpenShift Online, you are prevented from creating a project with the same name as any other project in the cluster, even if it's not your own. So if you can't create the project, it's possible that someone has already created a project with the same name."
func (c *Client) DeleteProject(name string, wait bool) error {

	// Instantiate watcher for our "wait" function
	var watcher watch.Interface
	var err error

	// If --wait has been passed, we will wait for the project to fully be deleted
	if wait {
		watcher, err = c.projectClient.Projects().Watch(context.Tastra(), metav1.ListOptions{
			FieldSelector: fields.Set{"metadata.name": name}.AsSelector().String(),
		})
		if err != nil {
			return fmt.Errorf("unable to watch project: %w", err)
		}
		defer watcher.Stop()
	}

	// Delete the project
	err = c.projectClient.Projects().Delete(context.Tastra(), name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("unable to delete project: %w", err)
	}

	// If watcher has been created (wait was passed) we will create a go routine and actually **wait**
	// until *EVERYTHING* is successfully deleted.
	if watcher != nil {

		// Project channel
		// Watch error channel
		projectChannel := make(chan *projectv1.Project)
		watchErrorChannel := make(chan error)

		// Create a go routine to run in the background
		go func() {

			for {

				// If watch unexpected has been closed..
				val, ok := <-watcher.ResultChan()
				if !ok {
					//return fmt.Errorf("received unexpected signal %+v on project watch channel", val)
					watchErrorChannel <- fmt.Errorf("watch channel was closed unexpectedly: %+v", val)
					break
				}

				// So we depend on val.Type as val.Object.Status.Phase is just empty string and not a mapped value constant
				if projectStatus, ok := val.Object.(*projectv1.Project); ok {

					klog.V(3).Infof("Status of delete of project %s is '%s'", name, projectStatus.Status.Phase)

					switch projectStatus.Status.Phase {
					//projectStatus.Status.Phase can only be "Terminating" or "Active" or ""
					case "":
						if val.Type == watch.Deleted {
							projectChannel <- projectStatus
							break
						}
						if val.Type == watch.Error {
							watchErrorChannel <- fmt.Errorf("failed watching the deletion of project %s", name)
							break
						}
					}

				} else {
					watchErrorChannel <- errors.New("unable to convert event object to Project")
					break
				}

			}
			close(projectChannel)
			close(watchErrorChannel)
		}()

		select {
		case val := <-projectChannel:
			klog.V(3).Infof("Deletion information for project: %+v", val)
			return nil
		case err := <-watchErrorChannel:
			return err
		case <-time.After(waitForProjectDeletionTimeOut):
			return fmt.Errorf("waited %s but couldn't delete project %s in time", waitForProjectDeletionTimeOut, name)
		}

	}

	// Return nil since we don't bother checking for the watcher..
	return nil
}

// CreateNewProject creates project with given projectName
func (c *Client) CreateNewProject(projectName string, wait bool) error {
	// Instantiate watcher before requesting new project
	// If watcher is created after the project it can lead to situation when the project is created before the watcher.
	// When this happens, it gets stuck waiting for event that already happened.
	var watcher watch.Interface
	var err error
	if wait {
		watcher, err = c.projectClient.Projects().Watch(context.Tastra(), metav1.ListOptions{
			FieldSelector: fields.Set{"metadata.name": projectName}.AsSelector().String(),
		})
		if err != nil {
			return fmt.Errorf("unable to watch new project %s creation: %w", projectName, err)
		}
		defer watcher.Stop()
	}

	projectRequest := &projectv1.ProjectRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name: projectName,
		},
	}
	_, err = c.projectClient.ProjectRequests().Create(context.Tastra(), projectRequest, metav1.CreateOptions{FieldManager: FieldManager})
	if err != nil {
		return fmt.Errorf("unable to create new project %s: %w", projectName, err)
	}

	if watcher != nil {
		for {
			val, ok := <-watcher.ResultChan()
			if !ok {
				break
			}
			if prj, ok := val.Object.(*projectv1.Project); ok {
				klog.V(3).Infof("Status of creation of project %s is %s", prj.Name, prj.Status.Phase)
				switch prj.Status.Phase {
				//prj.Status.Phase can only be "Terminating" or "Active" or ""
				case corev1.NamespaceActive:
					if val.Type == watch.Added {
						klog.V(3).Infof("Project %s now exists", prj.Name)
						return nil
					}
					if val.Type == watch.Error {
						return fmt.Errorf("failed watching the deletion of project %s", prj.Name)
					}
				}
			}
		}
	}

	return nil
}

// IsProjectSupported checks if Project resource type is present on the cluster
func (c *Client) IsProjectSupported() (bool, error) {
	return c.IsResourceSupported("project.openshift.io", "v1", "projects")
}

// GetCurrentProjectName returns the current project name
func (c *Client) GetCurrentProjectName() string {
	return c.Namespace
}
