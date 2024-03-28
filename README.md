mkdir gitlab-operator
cd gitlab-operator
operator-sdk init --domain github.com/fabiomnk/gitlab-operator --plugins helm

operator-sdk create api --group apps --version v1alpha1 --kind Gitlab
