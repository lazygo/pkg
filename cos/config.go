package cos

type Config struct {
	BucketURL string `json:"bucket_url" toml:"bucket_url"`
	SecretID  string `json:"secret_id" toml:"secret_id"`
	SecretKey string `json:"secret_key" toml:"secret_key"`
}

/**
bucket_url := "https://img-domain-1303896251.cos.ap-nanjing.myqcloud.com"
secret_id := "NKIDXYIn113YFgtqXypfsqX02rRkwud5AJ8j"
secret_key := "Yw50e7MQIuQkpgRyKuDUYOxFLOCnFTWJ"
*/
