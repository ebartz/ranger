package store

import (
	"fmt"
	"strings"

	v1 "github.com/ranger/ranger/pkg/generated/norman/core/v1"
	"github.com/ranger/ranger/pkg/namespace"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetValueForPasswordField(name string, secrets v1.SecretInterface) (string, error) {
	var value string
	if name == "" {
		return value, fmt.Errorf("empty secret name %s", name)
	}
	split := strings.SplitN(name, ":", 2)
	if len(split) != 2 {
		return value, fmt.Errorf("not a secret reference %s", name)
	}
	if split[0] != namespace.GlobalNamespace {
		return value, fmt.Errorf("not ranger referenced secret %s", name)
	}
	secret, err := secrets.Controller().Lister().Get(namespace.GlobalNamespace, split[1])
	if err != nil && errors.IsNotFound(err) {
		secret, err = secrets.GetNamespaced(namespace.GlobalNamespace, name, metav1.GetOptions{})
	}
	if err != nil {
		return value, err
	}
	return string(secret.Data[split[1]]), nil
}
