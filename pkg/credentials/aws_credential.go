package credentials

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"os"
)

func setEnvVar(envVarName string, value string) error {
	return os.Setenv(envVarName, value)
}

func setAWSCredential(secret *v1.Secret, key types.NamespacedName) error {
	fileName := fmt.Sprintf(key.Namespace + "-" + key.Name + "-credential")

	//////-------------------------------------------------------------------- Remove before deploy
	filePath := fmt.Sprintf("/home/rasel/Desktop/" + fileName)
	//filePath := fmt.Sprintf("/tmp/" + fileName)
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

	return setEnvVar("AWS_SHARED_CREDENTIALS_FILE", filePath)
}
