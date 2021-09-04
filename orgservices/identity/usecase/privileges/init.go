package privileges

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/core"
	"github.com/timoth-y/chainmetric-network/orgservices/identity/model"
)

var config *viper.Viper

// Init performs initialization of the `privileges` package.
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

// Has determines whether the user has certain privilege for given `path`.
func Has(user *model.User, path string) bool {
	return config.GetBool(fmt.Sprintf("%s.%s", user.Role, path))
}

// RoleHas determines whether the role has certain privilege for given `path`.
func RoleHas(role string, path string) bool {
	return config.GetBool(fmt.Sprintf("%s.%s", role, path))
}
