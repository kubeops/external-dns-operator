package credentials

import (
	"fmt"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"os"
)

var AzureConfigPath string

func setAzureCredential(secret *core.Secret, ednsKey types.NamespacedName) error {
	fileName := fmt.Sprintf("%s-%s-credential", ednsKey.Namespace, ednsKey.Name)
	filepath := fmt.Sprintf("/tmp/%s", fileName)

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	b := secret.Data["credentials"]
	_, err = file.Write(b)
	if err != nil {
		return err
	}

	AzureConfigPath = filepath

	return nil
}
