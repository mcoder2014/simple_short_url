package model

type ShortURLConfig struct {
	Short      string `json:"short"`
	Long       string `json:"long"`
	Enable     bool   `json:"enable"`
	Desp       string `json:"desp"`
	Creator    string `json:"creator"`
	CreateTime int64  `json:"create_time"`
}
