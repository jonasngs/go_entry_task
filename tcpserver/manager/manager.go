package manager

import (
	"github.com/jonasngs/go_entry_task/tcpserver/manager/services"
	"github.com/jonasngs/go_entry_task/tcpserver/storage"
)

type Manager struct {
	Authentication services.AuthenticationInterface
	Session        services.SessionInterface
	User           services.UserInterface
}

// Manager will be used to initialize and serves as an interface to access all services
func InitializeManager() *Manager {
	db := storage.InitializeDatabase()
	dao := services.InitializeDAOservice(db)
	redis := storage.InitializeCache()
	cacheService := services.InitializeCacheService(redis)
	imageService := services.InitializeImageService()

	return &Manager{
		Authentication: services.InitializeAuthService(dao),
		Session:        services.InitializeSession(cacheService, dao),
		User:           services.InitializeUserService(dao, cacheService, imageService),
	}
}
