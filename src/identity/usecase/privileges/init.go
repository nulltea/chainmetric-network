package privileges

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-contracts/shared/core"
	model "github.com/timoth-y/chainmetric-contracts/shared/model/user"
)

var (
	config *viper.Viper
)

func Init() {
	config = viper.New()

	config.SetConfigType("yaml")
	config.SetConfigName("privileges")
	config.AddConfigPath("data")
	config.AddConfigPath("src/identity/data")

	if err := config.ReadInConfig(); err != nil {
		core.Logrus.Fatal("failed to read privileges config")
	}
}

func Has(user *model.User, path string) bool {
	return config.GetBool(fmt.Sprintf("%s.%s", user.Role, path))
}
