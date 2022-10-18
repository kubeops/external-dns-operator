package credentials

import (
	"fmt"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"os"
)

func setAWSCredential(secret *core.Secret, endsKey types.NamespacedName) error {
	fileName := fmt.Sprintf("%s-%s-credential", endsKey.Namespace, endsKey.Name)

	filePath := fmt.Sprintf("/tmp/%s", fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	b := secret.Data["credentials"]
	_, err = file.Write(b)
	if err != nil {
		return err
	}

	err = os.Setenv("AWS_SHARED_CREDENTIALS_FILE", filePath)
	if err != nil {
		klog.Error("failed to set the environment variables")
		return err
	}
	return nil
}
