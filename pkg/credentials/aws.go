package credentials

import (
	"errors"
	"fmt"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"os"
)

func validateAWSSecret(secret *core.Secret) error {
	if secret.Data["credentials"] != nil {
		return nil
	} else {
		return errors.New("invalid secret format(s)")
	}
}

func setAWSCredential(secret *core.Secret, endsKey types.NamespacedName) error {

	if err := validateAWSSecret(secret); err != nil {
		return err
	}

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
