package config

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func isBcryptHash(input string) bool {
	// bcrypt加密后的文本以$2a$开头，后面跟着一串字符
	pattern := `^\$2a\$[0-9]{1,2}\$[A-Za-z0-9./]{53}$`
	matched, err := regexp.MatchString(pattern, input)
	if err != nil {
		fmt.Println("Error matching regex:", err)
		return false
	}
	return matched
}

func encodePassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //加密处理
	encodePWD := string(hash)                                                    // 保存在数据库的密码，虽然每次生成都不同，只需保存一份即可
	return encodePWD
}

func AuthUser(u string, p string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(Cfg.Auth.Password), []byte(p))
	if u == Cfg.Auth.Username && err == nil {
		return true
	}
	return false
}

func RequireAuth() bool {
	return Cfg.Auth.Authentication && Cfg.Auth.Username != "" && Cfg.Auth.Password != ""
}
