package controllers

import "github.com/victorsteven/fullstack/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.GET("/", s.Home)

	// Login Route
	s.Router.POST("/login", s.Login)

	//Users routes
	s.Router.POST("/users", s.CreateUser)
	s.Router.GET("/users", s.GetUsers)
	s.Router.GET("/users/:id", s.GetUser)
	s.Router.PUT("/users/:id", middlewares.TokenAuthMiddleware(), s.UpdateUser)
	s.Router.DELETE("/users/:id", middlewares.TokenAuthMiddleware(), s.DeleteUser)

	// //Posts routes
	s.Router.POST("/posts", middlewares.TokenAuthMiddleware(), s.CreatePost)
	s.Router.GET("/posts", s.GetPosts)
	s.Router.GET("/posts/:id", s.GetPost)
	s.Router.PUT("/posts/:id", middlewares.TokenAuthMiddleware(), s.UpdatePost)
	s.Router.DELETE("/posts/:id", middlewares.TokenAuthMiddleware(), s.DeletePost)
}
