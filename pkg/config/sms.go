package config

type Alibaba struct {
	AccessKeyId     string `yaml:"accessKeyId"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	SignName        string `yaml:"signName"`
	TemplateCode    string `yaml:"templateCode"`
}
