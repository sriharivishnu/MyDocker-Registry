package config

import (
	"fmt"
	"os"
	"reflect"
)

type configType struct {
	ENVIRONMENT           string
	PORT                  string
	DATABASE_NAME         string
	DATABASE_HOST         string
	DATABASE_PORT         string
	DATABASE_USER         string
	DATABASE_PASSWORD     string
	S3_BUCKET_KEY         string
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	JWT_SECRET            string
}

var Config *configType = &configType{}

func PopulateConfig() {
	v := reflect.ValueOf(Config).Elem()
	typeOfConfig := v.Type()

	for i := 0; i < v.NumField(); i++ {
		key := typeOfConfig.Field(i).Name
		value, present := os.LookupEnv(key)

		if !present {
			panic(fmt.Sprintf("Required Environment Variable '%s' is not set.", key))
		}
		v.Field(i).SetString(value)
	}
}
