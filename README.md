# combo
`combo` is a [Kubernetes controller](https://kubernetes.io/docs/concepts/architecture/controller/) that generates and applies resources for all combinations of a manifest template and its arguments.

## When to use 
`Combo` is designed to be a sample operator that can be utilized in our CI/CD, packaged in [RukPak](https://github.com/operator-framework/rukpak), and ultimately provide the team a window into the developer experience of being on [OLM](https://github.com/operator-framework/operator-lifecycle-manager). This is the primary purpose of Combo's existance. In addition to this, there were not easy ways to generate combinations of manifests with existing tools - thus the idea of `Combo` emerged. This functionality is not currently found in things like [Helm](https://helm.sh/) and [OpenShift Templates](https://docs.openshift.com/container-platform/4.9/openshift_images/using-templates.html) easily.

If you do not need to generate the combinations of Kubernetes manifests, this tool will not provide different functionality than other templating engines. The most valuable use case of `Combo` is when the combination of manifests with arguments is needed.

## On the command line
Directly evaluate a template from stdin:

```sh
$ cat <<EOF | combo eval -a 'NAMESPACE=foo,bar' -a 'NAME=baz' -
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: deployment-reader
rules:
- apiGroups: ["apps"]
  resources: ["deployments"]
  verbs: ["get", "watch", "list"]
---
kind: Namespace
metadata:
  name: NAMESPACE
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: NAME
  namespace: NAMESPACE
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: deployment-reader
  namespace: NAMESPACE
subjects:
- kind: ServiceAccount
  name: NAME
  namespace: NAMESPACE
roleRef:
  kind: ClusterRole
  name: deployment-reader
  apiGroup: rbac.authorization.k8s.io
EOF
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: deployment-reader
rules:
- apiGroups: ["apps"]
  resources: ["deployments"]
  verbs: ["get", "watch", "list"]
---
kind: Namespace
metadata:
  name: foo
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: baz
  namespace: foo
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: deployment-reader
  namespace: foo
subjects:
- kind: ServiceAccount
  name: baz
  namespace: foo
roleRef:
  kind: ClusterRole
  name: deployment-reader
  apiGroup: rbac.authorization.k8s.io
---
kind: Namespace
metadata:
  name: bar
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: baz
  namespace: bar
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: deployment-reader
  namespace: bar
subjects:
- kind: ServiceAccount
  name: baz
  namespace: bar
roleRef:
  kind: ClusterRole
  name: deployment-reader
  apiGroup: rbac.authorization.k8s.io
```

### As a controller
If `combo` is running as a controller in the current kubectl context's cluster, create a `Template`:

```sh
$ cat <<EOF | kubectl create -f -
apiVersion: combo.io/v1alpha1
kind: Template
metadata:
  name: feature
spec:
  body: |
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: RoleBinding
    metadata:
      name: feature-controller
      namespace: TARGET_NAMESPACE
    subjects:
    - kind: ServiceAccount
      name: controller
      namespace: feature
    roleRef:
      kind: ClusterRole
      name: feature-controller
      apiGroup: rbac.authorization.k8s.io
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: RoleBinding
    metadata:
      generateName: feature-user-
      namespace: TARGET_NAMESPACE
    subjects:
    - kind: Group
      name: TARGET_GROUP
      namespace: TARGET_NAMESPACE
      apiGroup: rbac.authorization.k8s.io
    roleRef:
      kind: ClusterRole
      name: feature-user
      apiGroup: rbac.authorization.k8s.io
  parameters:
  - key: TARGET_GROUP
  - key: TARGET_NAMESPACE
EOF
```

Assuming the existance of the `feature-controller` and `feature-user` `ClusterRoles` as well as the `feature`, `staging`, and `prod` `Namespaces`, instantiate all resource/argument combinations with a `Combination`:

```sh
$ cat <<EOF | kubectl create -f -
apiVersion: combo.io/v1alpha1
kind: Combination
metadata:
  name: enable-feature
spec:
  template:
    name: feature
  arguments:
  - key: TARGET_GROUP
    values:
    - "sre"
    - "system:serviceaccounts:ci"
  - key: TARGET_NAMESPACE
    values:
    - staging
    - prod
EOF
```

combo then surfaces the evaluated template in the status:

```sh
$ kubectl get combination -o yaml
apiVersion: combo.io/v1alpha1
kind: Combination
metadata:
  name: enable-feature
spec:
  template:
    name: feature
  arguments:
  - key: TARGET_GROUP
    values:
    - "sre"
    - "system:serviceaccounts:ci"
  - key: TARGET_NAMESPACE
    values:
    - staging
    - prod
status:
  evaluated:
  - apiVersion: rbac.authorization.k8s.io/v1
    kind: RoleBinding
    metadata:
      name: feature-controller
      namespace: staging
    subjects:
    - kind: ServiceAccount
      name: controller
      namespace: feature
    roleRef:
      kind: ClusterRole
      name: feature-controller
      apiGroup: rbac.authorization.k8s.io
  - apiVersion: rbac.authorization.k8s.io/v1
    kind: RoleBinding
    metadata:
      generateName: feature-user-
      namespace: staging
    subjects:
    - kind: Group
      name: sre
      namespace: staging
      apiGroup: rbac.authorization.k8s.io
    roleRef:
      kind: ClusterRole
      name: feature-user
      apiGroup: rbac.authorization.k8s.io
      ...
```
