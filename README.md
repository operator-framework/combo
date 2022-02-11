# combo

`combo` is a [Kubernetes controller](https://kubernetes.io/docs/concepts/architecture/controller/) that generates and applies resources for __all combinations__ of a manifest template and its arguments.

## What on earth does "all combinations" mean!?

For example, say `combo` is given a template containing 1 _parameter_  -- i.e. distinct token to replaced -- and 2 _arguments_  for that parameter (value to replace a parameter with), then `combo` will output 2 _evaluations_ (one for each argument).

e.g.

_template:_

```yaml
PARAM_1
```

_arguments:_

```
PARAM_1:
- a
- b
```

_evaluations:_

```yaml
---
a
---
b
```

Simple enough, right? What about something more advanced...

Suppose we give `combo` a template with 2 parameters and 3 arguments for each parameter.

e.g.

_template:_

```yaml
PARAM_1: PARAM_2
```
_arguments:_

```yaml
PARAM_1:
- a
- b
PARAM_2:
- c
- d
- e
```

_evaluations:_

```yaml
---
a: c
---
a: d
---
a: e
---
b: c
---
b: d
---
b: e
```

## Wait, can't Helm do this?

It could, __but__ getting this experience from Helm would require:

- the use of nested Go Template loops
- carrying along the rest of the ecosystem; e.g. I don't want to think about Helm charts for something this simple

## Can I generate combinations locally?

Yes! There is a built in CLI interaction via the binary. Let's look at the same example we used for the controller in this context.

First, create a simple YAML file that defines two arguments.

```yaml
# ./sample_input.yaml
PARAM_1: PARAM_2
```

Next, go ahead and run the `eval` subcommand.

```shell
make build-cli
./combo eval -r PARAM_1=a,b -r PARAM_2=c,d,e sample_input.yaml
```

This will run the same logic that the controller utilizes to generate combinations and output them to stdout. The above command will produce the following:

```yaml
---
a: c
---
a: d
---
a: e
---
b: c
---
b: d
---
b: e
```

## Primary use cases

To parameterize RBAC and other namespace-scoped resources so they can be stamped out as necessary later on.

```shell
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
  - TARGET_GROUP
  - TARGET_NAMESPACE
EOF
```

Assuming the existence of the `feature-controller` and `feature-user` `ClusterRoles` as well as the `feature`, `staging`, and `prod` `Namespaces`, instantiate all resource/argument combinations with a `Combination`:

```shell
$ cat <<EOF | kubectl create -f -
apiVersion: combo.io/v1alpha1
kind: Combination
metadata:
  name: enable-feature
spec:
  template: feature
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

```shell
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

## Ulterior motives

Our "hidden" agenda with `combo` is for it to:

- serve as an example of best-practices when building operators
- be used to dog food operator-framework packaging formats (e.g. OLM, rukpak, etc)
- become a dependency of operators that require post-install configuration (e.g. RBAC scoping, secret creation, etc)
- serve as a cruft-free way to up-skill new maintainers through easy contributions

## Planned Features

The combo project is currently in active development and you can find up-to-date, tracked work in the [Combo Project Board][roadmap].

Introducing or discussing new roadmap items can be done by opening an issue.

Before opening an issue, it's recommended to check whether a similar issue is already open to avoid duplicating roadmap efforts.

[roadmap]: <https://github.com/operator-framework/combo/projects/1>
