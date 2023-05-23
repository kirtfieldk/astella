package conf

type AwsConfig struct {
	AwsAccessKeyId     string `yaml:"aws_access_key"`
	AwsAccessKeySecret string `yaml:"aws_secret_access_key"`
	Region             string `yaml:"region"`
	BucketName         string `yaml:"bucketname"`
	UploadTimeout      int    `yaml:"upload_timeout"`
	BaseUrl            string `yaml:"base_url"`
}
