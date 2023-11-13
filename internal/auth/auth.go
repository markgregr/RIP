package auth

type AuthManager struct {
	UserID uint
	Role   string
}

var authInstance *AuthManager

func GetAuthInstance() *AuthManager {
	if authInstance == nil {
		authInstance = &AuthManager{
			UserID: 1,      // Пример значения для пользователя
			Role:   "user", // Пример роли пользователя
		}
	}
	return authInstance
}
