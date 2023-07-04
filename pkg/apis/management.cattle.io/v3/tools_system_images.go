package v3

var (
	ToolsSystemImages = struct {
		AuthSystemImages AuthSystemImages
	}{
		AuthSystemImages: AuthSystemImages{
			KubeAPIAuth: "ranger/kube-api-auth:v0.1.8",
		},
	}
)
