package quota

import (
	"testing"

	configv1alpha1 "github.com/kiosk-sh/kiosk/pkg/apis/config/v1alpha1"
	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type getResourceQuotasStatusByNamespaceTestCase struct {
	name string

	namespaceStatuses configv1alpha1.AccountQuotasStatusByNamespace
	namespace         string

	expectedStatus corev1.ResourceQuotaStatus
	expectedFound  bool
}

func TestGetResourceQuotasStatusByNamespace(t *testing.T) {
	testCases := []getResourceQuotasStatusByNamespaceTestCase{
		{
			name:           "Empty arr",
			expectedStatus: corev1.ResourceQuotaStatus{},
		},
		{
			name: "Namespace found",
			namespaceStatuses: []configv1alpha1.AccountQuotaStatusByNamespace{
				{
					Namespace: "needle",
					Status: corev1.ResourceQuotaStatus{
						Hard: map[corev1.ResourceName]resource.Quantity{
							"test": resource.Quantity{
								Format: "test",
							},
						},
					},
				},
			},
			namespace: "needle",
			expectedStatus: corev1.ResourceQuotaStatus{
				Hard: map[corev1.ResourceName]resource.Quantity{
					"test": resource.Quantity{
						Format: "test",
					},
				},
			},
			expectedFound: true,
		},
	}

	for _, testCase := range testCases {
		status, found := GetResourceQuotasStatusByNamespace(testCase.namespaceStatuses, testCase.namespace)

		statusAsYaml, err := yaml.Marshal(status)
		assert.NilError(t, err, "Error parsing status in testCase %s", testCase.name)
		expectedAsYaml, err := yaml.Marshal(testCase.expectedStatus)
		assert.NilError(t, err, "Error parsing expectation in testCase %s", testCase.name)
		assert.Equal(t, string(statusAsYaml), string(expectedAsYaml), "Unexpected status in testCase %s", testCase.name)

		assert.Equal(t, found, testCase.expectedFound, "Unexpected found bool in testCase %s", testCase.name)
	}
}
