operator-sdk init \
    --domain fabiomnk.co \
    --project-version 3 \
    --plugins go/v4 \
    --repo github.com/fabiomnk/github-operator 

operator-sdk create api \
    --version v1alpha1 \
    --group gitlab \
    --kind Project \
    --resource \
    --controller 


> make deploy
test -s /mnt/c/Users/fabio.mancini/git/fabiomnk/gitlab-operator/bin/controller-gen && /mnt/c/Users/fabio.mancini/git/fabiomnk/gitlab-operator/bin/controller-gen --version | grep -q v0.13.0 || \
GOBIN=/mnt/c/Users/fabio.mancini/git/fabiomnk/gitlab-operator/bin go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.13.0
/mnt/c/Users/fabio.mancini/git/fabiomnk/gitlab-operator/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
cd config/manager && /mnt/c/Users/fabio.mancini/git/fabiomnk/gitlab-operator/bin/kustomize edit set image controller=controller:latest
/mnt/c/Users/fabio.mancini/git/fabiomnk/gitlab-operator/bin/kustomize build config/default | kubectl apply -f -
namespace/gitlab-operator-system created
customresourcedefinition.apiextensions.k8s.io/projects.gitlab.fabiomnk.co unchanged
serviceaccount/gitlab-operator-controller-manager created
role.rbac.authorization.k8s.io/gitlab-operator-leader-election-role created
clusterrole.rbac.authorization.k8s.io/gitlab-operator-manager-role created
clusterrole.rbac.authorization.k8s.io/gitlab-operator-metrics-reader created
clusterrole.rbac.authorization.k8s.io/gitlab-operator-proxy-role created
rolebinding.rbac.authorization.k8s.io/gitlab-operator-leader-election-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/gitlab-operator-manager-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/gitlab-operator-proxy-rolebinding created
service/gitlab-operator-controller-manager-metrics-service created
deployment.apps/gitlab-operator-controller-manager created


