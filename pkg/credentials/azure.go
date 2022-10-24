package credentials

import (
	"errors"
	"fmt"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"os"
)

func validateAzureSecret(secret *core.Secret) error {
	if secret.Data["azure.json"] != nil {
		return nil
	} else {
		return errors.New("invalid secret format(s)")
	}
}

func setAzureCredential(secret *core.Secret, ednsKey types.NamespacedName) error {
	if err := validateAzureSecret(secret); err != nil {
		return err
	}
	fileName := fmt.Sprintf("%s-%s-credential", ednsKey.Namespace, ednsKey.Name)
	filepath := fmt.Sprintf("/tmp/%s", fileName)

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	b := secret.Data["azure.json"]
	_, err = file.Write(b)
	if err != nil {
		return err
	}

	return nil
}
