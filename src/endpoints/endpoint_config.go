package endpoints

import (
	"database/sql"
	"friendMgmt/data"
	"friendMgmt/services"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func initUserEndpoint(db *sql.DB) UserEndpoint {
	var userRepo = data.UserRepository{DB: db}
	userService := services.UserService{IUserRepository: userRepo}
	return UserEndpoint{IUserService: userService}
}

func initRelationshipEndpoint(db *sql.DB) RelationshipEndpoint {
	var relationshipRepo = data.RelationshipRepository{DB: db}
	relationshipService := services.RelationshipService{IRelationshipRepository: relationshipRepo}
	var userRepo = data.UserRepository{DB: db}
	userService := services.UserService{IUserRepository: userRepo}
	return RelationshipEndpoint{IRelationshipService: relationshipService, IUserService: userService}
}

func ConfigRoutes(db *sql.DB) {

	gin.SetMode(gin.ReleaseMode)

	userApi := initUserEndpoint(db)
	relationshipApi := initRelationshipEndpoint(db)

	router := gin.Default()

	router.POST("/api/friends/add", relationshipApi.CreateRelationship)
	router.POST("/api/friends", relationshipApi.FriendList)
	router.POST("/api/friends/common-friends", relationshipApi.CommonFriendList)
	router.POST("/api/friends/subcribe", relationshipApi.Subscribe)
	router.POST("/api/friends/block", relationshipApi.Block)
	router.POST("/api/friends/receive-updates", relationshipApi.ReceiveUpdates)
	router.GET("/api/users", userApi.Users)
	router.POST("/api/users", userApi.CreateUser)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := router.Run(":8081")
	if err != nil {
		panic(err)
	}
}
