package catalogv2

import (
	v1 "github.com/ranger/ranger/pkg/apis/catalog.cattle.io/v1"
	corev1controllers "github.com/ranger/wrangler/pkg/generated/controllers/core/v1"
	corev1 "k8s.io/api/core/v1"
)

func GetSecret(secrets corev1controllers.SecretCache, repoSpec *v1.RepoSpec, repoNamespace string) (*corev1.Secret, error) {
	if repoSpec.ClientSecret == nil {
		return nil, nil
	}
	ns := repoSpec.ClientSecret.Namespace
	if repoNamespace != "" {
		ns = repoNamespace
	}

	return secrets.Get(ns, repoSpec.ClientSecret.Name)
}
