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

    