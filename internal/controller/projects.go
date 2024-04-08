package controller

import (
	"context"

	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var ErrProjectNotFound = errors.New("Project not found")

func ProjectsExists(ctx context.Context, git *gitlab.Client, name string) (*gitlab.Project, error) {
	l := log.FromContext(ctx)

	projects, res, err := git.Projects.ListProjects(&gitlab.ListProjectsOptions{Search: gitlab.Ptr(name)})
	if err != nil {
		l.Error(err, "error listing projects")
	}

	//l.Info(fmt.Sprintf("project: %+v", projects))
	//l.Info(fmt.Sprintf("res: %+v", res))

	var project = gitlab.Project{}

	err = nil
	if res.TotalItems == 1 {
		project = *projects[0]
	} else {
		err = ErrProjectNotFound
	}
	return &project, err
}
