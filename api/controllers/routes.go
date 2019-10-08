package controllers

import (
	"github.com/victorsteven/fullstack/api/middlewares"
)

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.GET("/", s.Home)

	v1 := s.Router.Group("/api/v1")
	{
		// Login Route
		v1.POST("/login", s.Login)
		v1.POST("/logout", s.Logout)

		//Users routes
		v1.POST("/users", s.CreateUser)
		v1.GET("/users", s.GetUsers)
		// v1.GET("/users/:id", s.GetUser)
		v1.PUT("/users/:id", middlewares.TokenAuthMiddleware(), s.UpdateUser)
		v1.PUT("/avatar/users/:id", middlewares.TokenAuthMiddleware(), s.UpdateAvatar)

		// v1.DELETE("/users/:id", middlewares.TokenAuthMiddleware(), s.DeleteUser)

		//Posts routes
		v1.POST("/posts", middlewares.TokenAuthMiddleware(), s.CreatePost)
		v1.GET("/posts", s.GetPosts)
		v1.GET("/posts/:id", s.GetPost)
		v1.PUT("/posts/:id", middlewares.TokenAuthMiddleware(), s.UpdatePost)
		v1.DELETE("/posts/:id", middlewares.TokenAuthMiddleware(), s.DeletePost)

		//profile
		//v1.POST("/posts/:id", middlewares.TokenAuthMiddleware(), s.UpdateProfile)


	}
}
