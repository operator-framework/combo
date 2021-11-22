package testdata

var GenerateInput = `---
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
`

var GenerateOutput = []string{
	`apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
	name: deployment-reader
rules:
- apiGroups: ["apps"]
	resources: ["deployments"]
	verbs: ["get", "watch", "list"]`,
	`kind: Namespace
metadata:
	name: foo`,

	`apiVersion: v1
kind: ServiceAccount
metadata:
	name: baz
	namespace: foo`,

	`apiVersion: rbac.authorization.k8s.io/v1
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
	apiGroup: rbac.authorization.k8s.io`,

	`kind: Namespace
metadata:
	name: bar`,

	`apiVersion: v1
kind: ServiceAccount
metadata:
	name: baz
	namespace: bar`,

	`apiVersion: rbac.authorization.k8s.io/v1
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
	apiGroup: rbac.authorization.k8s.io`,
}
