package generator

import (
	"context"
	"strings"
	"testing"

	"github.com/operator-framework/combo/pkg/combination"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type expected struct {
	err        error
	evaluation string
}

func TestEvaluate(t *testing.T) {
	for _, tt := range []struct {
		name         string
		template     string
		combinations combination.Stream
		expected     expected
	}{
		{
			name: "can process a template",
			template: `---
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
`,
			expected: expected{
				err: nil,
				evaluation: `---
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
`,
			},
			combinations: combination.NewStream(
				combination.WithArgs(map[string][]string{
					"NAMESPACE": {"foo", "bar"},
					"NAME":      {"baz"},
				}),
				combination.WithSolveAhead(),
			),
		},
		{
			name:     "processes an empty template",
			template: ``,
			expected: expected{
				err:        nil,
				evaluation: ``,
			},
			combinations: combination.NewStream(
				combination.WithArgs(map[string][]string{
					"NAMESPACE": {"foo", "bar"},
					"NAME":      {"baz"},
				}),
				combination.WithSolveAhead(),
			),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			evaluation, err := Evaluate(ctx, tt.template, tt.combinations)
			if err != nil {
				t.Fatal("received an error during evaluation:", err)
			}

			assert.Equal(t, tt.expected.err, err)
			require.ElementsMatch(
				t,
				strings.Split(string(tt.expected.evaluation), "---"),
				strings.Split(string(evaluation), "---"),
				"Document evaluations generated incorrectly")
		})
	}
}
