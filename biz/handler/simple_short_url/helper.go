package simple_short_url

import (
	"fmt"
	"os"
)

func CheckToken(token string) error {
	if len(token) == 0 {
		return fmt.Errorf("token is empty")
	}
	sysToken := os.Getenv("ACCESS_TOKEN")
	if token != sysToken {
		return fmt.Errorf("access token: %v is invalid", token)
	}
	return nil
}
