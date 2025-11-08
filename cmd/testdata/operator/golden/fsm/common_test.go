package fsm_test

import (
	"k8s.io/apimachinery/pkg/types"
)

const kind = "TestCRD"

// resourceName to be used by both unit and integration tests if needed
var resourceName = types.NamespacedName{
	Namespace: "default",
	Name:      "test-resource",
}
