package broker

import (
	"os"
	"time"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	ccWrapper "code.cloudfoundry.org/cli/api/cloudcontroller/wrapper"
	"code.cloudfoundry.org/cli/api/uaa"
	"code.cloudfoundry.org/cli/api/uaa/constant"
	uaaWrapper "code.cloudfoundry.org/cli/api/uaa/wrapper"
	wrapperUtil "code.cloudfoundry.org/cli/api/uaa/wrapper/util"
	"code.cloudfoundry.org/cli/util/configv3"
	// "code.cloudfoundry.org/cli/command/v7/shared"
	"github.com/cloudfoundry-community/go-cfclient"
	"gopkg.in/yaml.v3"
)

const ConfigEnvVarName string = "CONFIG_SERVER_BROKER_CONFIG"

type Config struct {
	Auth                 Auth     `yaml:"broker_auth"`
	ServiceName          string   `yaml:"service_name"`
	ServiceID            string   `yaml:"service_id"`
	BasicPlanId          string   `yaml:"basic_plan_id"`
	BasicPlanName        string   `yaml:"basic_plan_name"`
	Host                 string   `yaml:"host"`
	ServiceInstanceLimit int      `yaml:"service_instance_limit"`
	CfConfig             CfConfig `yaml:"cloud_foundry_config"`
	Description          string   `yaml:"description"`
	LongDescription      string   `yaml:"long_description"`
	ProviderDisplayName  string   `yaml:"provider_display_name"`
	DocumentationURL     string   `yaml:"documentation_url"`
	SupportURL           string   `yaml:"support_url"`
	DisplayName          string   `yaml:"display_name"`
	IconImage            string   `yaml:"icon_image"`
}

type Auth struct {
	Username string `yaml:"user"`
	Password string `yaml:"password"`
}

type CfConfig struct {
	ApiUrl   string `yaml:"api_url"`
	Username string `yaml:"user"`
	Password string `yaml:"password"`
}

func (config *Config) GetCfClient() (*cfclient.Client, error) {
	return cfclient.NewClient(&cfclient.Config{
		ApiAddress: config.CfConfig.ApiUrl,
		Username:   config.CfConfig.Username,
		Password:   config.CfConfig.Password,
	})
}

type UaaConfig struct {
}

func (config UaaConfig) BinaryName() string {
	return "Config Server Broker"
}
func (config UaaConfig) BinaryVersion() string {
	return "0.1.0"
}
func (config UaaConfig) DialTimeout() time.Duration {
	return configv3.DefaultDialTimeout
}
func (config UaaConfig) SetUAAEndpoint(uaaEndpoint string) {}
func (config UaaConfig) SkipSSLValidation() bool {
	return false
}
func (config UaaConfig) UAADisableKeepAlives() bool {
	return true
}
func (config UaaConfig) UAAGrantType() string {
	return string(constant.GrantTypePassword)

}
func (config UaaConfig) UAAOAuthClient() string {
	return ""
}
func (config UaaConfig) UAAOAuthClientSecret() string {
	return ""
}

func (config *Config) GetV3CfClient(ui BrokerLogger) (*ccv3.Client, error) {
	tokenCache := wrapperUtil.NewInMemoryTokenCache()
	authWrapper := ccWrapper.NewUAAAuthentication(nil, tokenCache)
	ccWrappers := []ccv3.ConnectionWrapper{}
	ccWrappers = append(ccWrappers, authWrapper)
	ccWrappers = append(ccWrappers, ccWrapper.NewRequestLogger(ui))
	ccClient := ccv3.NewClient(ccv3.Config{
		AppName:            "Config Server Broker",
		AppVersion:         "0.1.0",
		JobPollingTimeout:  configv3.DefaultOverallPollingTimeout,
		JobPollingInterval: configv3.DefaultPollingInterval,
		Wrappers:           ccWrappers,
	})
	_, _, err := ccClient.TargetCF(ccv3.TargetSettings{
		URL:               config.CfConfig.ApiUrl,
		SkipSSLValidation: false,
		DialTimeout:       configv3.DefaultDialTimeout,
	})
	if err != nil {
		return ccClient, err
	}
	uaaClient := uaa.NewClient(UaaConfig{})
	uaaClient.WrapConnection(uaaWrapper.NewRequestLogger(ui))
	uaaAuthWrapper := uaaWrapper.NewUAAAuthentication(uaaClient, tokenCache)
	uaaClient.WrapConnection(uaaAuthWrapper)
	err = uaaClient.SetupResources(ccClient.UAA())
	credentials := map[string]string{
		"username": config.CfConfig.Username,
		"password": config.CfConfig.Password,
	}
	accessToken, refreshToken, err := uaaClient.Authenticate(credentials, "", constant.GrantTypePassword)
	tokenCache.SetAccessToken(accessToken)
	tokenCache.SetRefreshToken(refreshToken)

	if err != nil {
		return ccClient, err
	}

	return ccClient, nil
}

func ParseConfig() (Config, error) {
	configJson := os.Getenv(ConfigEnvVarName)
	if configJson == "" {
		panic(ConfigEnvVarName + " not set")
	}
	var config Config

	return config, yaml.Unmarshal([]byte(configJson), &config)
}
