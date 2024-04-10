/*
Copyright 2024.

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

package controller

import (
	"context"
	"fmt"
	"os"

	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	gitlabv1alpha1 "github.com/fabiomnk/github-operator/api/v1alpha1"
	"github.com/pkg/errors"

	gitlab "github.com/xanzy/go-gitlab"
)

// ProjectReconciler reconciles a Project object
type ProjectReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=gitlab.fabiomnk.co,resources=projects,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=gitlab.fabiomnk.co,resources=projects/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=gitlab.fabiomnk.co,resources=projects/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Project object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *ProjectReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	// init gitlab client
	gitlab_pat := os.Getenv("GITLAB_TOKEN")
	git, err := gitlab.NewClient(string(gitlab_pat))
	if err != nil {
		l.Error(err, "error creating gitlab client")
	}

	l.Info("Start Reconciliation")
	project := &gitlabv1alpha1.Project{}
	err = r.Client.Get(ctx, req.NamespacedName, project)
	l.Info(fmt.Sprintf("reconciliate project: %s", project.Spec.Name))
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			l.Info("Deleting Project")

			return ctrl.Result{}, nil
		}
		l.Error(err, "Could not fetch request")
		return ctrl.Result{Requeue: true}, err
	}

	// check if already exists
	proj, err := ProjectsExists(ctx, git, project.Spec.Name)
	if err != nil {
		if errors.Is(err, ErrProjectNotFound) {

			l.Info("Project Not found, Creating...")

			// Create new project
			cpo := &gitlab.CreateProjectOptions{
				Name:                 &project.Spec.Name,
				Description:          &project.Spec.Description,
				MergeRequestsEnabled: gitlab.Ptr(true),
				SnippetsEnabled:      gitlab.Ptr(true),
				Visibility:           gitlab.Ptr(gitlab.PrivateVisibility),
				InitializeWithReadme: gitlab.Ptr(true),
			}

			proj, _, err = git.Projects.CreateProject(cpo)
			if err != nil {
				l.Error(err, "error creating project")
			}

			l.Info("Project Created")
		} else {
			l.Error(err, "error check project existence")
		}
	}

	l.Info(fmt.Sprintf("configuring project: %s with id: %b", proj.Name, proj.ID))

	// edit existing project
	epo := &gitlab.EditProjectOptions{
		Description: &project.Spec.Description,
	}

	_, res, err := git.Projects.EditProject(proj.ID, epo)
	if err != nil {
		l.Error(err, fmt.Sprintf("error reconfiguring project %+v", res))

	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ProjectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gitlabv1alpha1.Project{}).
		Complete(r)
}
