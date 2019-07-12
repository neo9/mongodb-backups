package config


type Plan struct {
	Backups []Backup `json:"backups,omitempty"`
	S3 S3 `json:"s3"`
}

type Backup struct {
	Name string `json:"name"`
	Schedule string `json:"schedule"`
	Retention string `json:"retention"`
	Timeout string `json:"timeout"`
	MongoDB MongoDB `json:"mongodb"`
}

type S3 struct {
	Name string `json:"name"`
	Region string `json:"region"`
}

type MongoDB struct {
	Host string `json:"host"`
	Port string `json:"port"`
}


