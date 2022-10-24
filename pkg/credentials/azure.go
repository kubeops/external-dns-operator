package credentials

import (
	"errors"
	"fmt"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"os"
)

func validAzureSecret(secret *core.Secret) bool {
	_, found := secret.Data["azure.json"]
	return found
}

func setAzureCredential(secret *core.Secret, ednsKey types.NamespacedName) error {
	if !validAzureSecret(secret) {
		return errors.New("invalid Azure provider secret")
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
