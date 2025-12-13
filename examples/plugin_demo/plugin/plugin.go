package plugin

import (
	"github.com/luxingwen/sgin/pkg/app"
	"github.com/luxingwen/sgin/pkg/config"
	"github.com/spf13/viper"
)

// MyExtra is the plugin configuration structure
type MyExtra struct {
	FeatureFlag bool `mapstructure:"feature_flag"`
	Limit       int  `mapstructure:"limit"`
}

// Conf will hold the decoded configuration after InitConfig runs
var Conf *MyExtra

func init() {
	// register extension to decode `myplugin` key from the central config
	config.RegisterExtension("myplugin", func(v *viper.Viper, cfg *config.Config) error {
		var c MyExtra
		if err := v.UnmarshalKey("myplugin", &c); err != nil {
			return err
		}
		Conf = &c
		return nil
	}, false)
}

// Setup registers plugin routes / middleware onto the provided App.
// It uses App.RegisterPlugin so registration can be replayed by a host.
func Setup(a *app.App) {
	if a == nil {
		return
	}
	// register a route only if feature enabled
	if Conf != nil && Conf.FeatureFlag {
		a.RegisterPlugin(func(a *app.App) {
			a.GET("/plugin/hello", func(c *app.Context) {
				c.JSONSuccess(map[string]interface{}{"plugin": "hello", "limit": Conf.Limit})
			})
		})
	}
}
