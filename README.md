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



# Create Namespace (to create secret beforehand)
kubectl create ns gitlab-operator-system
# Create secret with gitlab-token
kubectl create secret generic gitlab-credentials --from-literal=personalAccessToken=$gitlab_token -n gitlab-operator-system

# Deploy Operator
make deploy

# Watch controller logs 
kubectl logs -f -l control-plane=controller-manager

# Apply Sample Project
kubectl apply -k config/samples/



