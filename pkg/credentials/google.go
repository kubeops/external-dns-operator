package credentials

import (
	"errors"
	"fmt"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"os"
)

func validGoogleSecret(secret *core.Secret) bool {
	_, found := secret.Data["credentials.json"]
	return found
}

func setGoogleCredential(secret *core.Secret, ednsKey types.NamespacedName) error {
	if !validGoogleSecret(secret) {
		return errors.New("invalid Google provider secret")
	}
	fileName := fmt.Sprintf("%s-%s-credential", ednsKey.Namespace, ednsKey.Name)
	filepath := fmt.Sprintf("/tmp/%s", fileName)

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	b := secret.Data["credentials.json"]
	_, err = file.Write(b)
	if err != nil {
		return err
	}

	return nil
}
