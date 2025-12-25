package router

import (
	"github.com/Cheerdoge/library-manage-system/internal/handler"
	"github.com/Cheerdoge/library-manage-system/internal/middleware"
	"github.com/Cheerdoge/library-manage-system/internal/service"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.Engine,
	bookhandler *handler.BookHandler,
	borrowhandler *handler.BorrowHandler,
	userhandler *handler.UserHandler,
	authhandler *handler.AuthHandler,
	sessionservice *service.SessionService,
) {
	publicGroup := r.Group("/api")
	{
		// 用户相关
		publicGroup.POST("/login", authhandler.LoginHandler)
		publicGroup.POST("/register", authhandler.RegisterUserHandler)
		// 图书相关
		publicGroup.GET("/books", bookhandler.GetBooksHandler)
		publicGroup.GET("/books/:id", bookhandler.GetBookByIdHandler)
	}

	authGroup := r.Group("/api")
	authGroup.Use(middleware.AuthMiddleware(sessionservice))
	{
		// 用户相关
		authGroup.GET("/logout", authhandler.LogoutHandler)
		authGroup.DELETE("/user", authhandler.DelHandler)
		authGroup.GET("/user", userhandler.GetUserInfoHandler)
		authGroup.PUT("/user/change_password", userhandler.UserChangePasswordHandler)
		authGroup.PUT("/user/change_info", userhandler.ChangeUserInfoHandler)

		// 借阅相关
		authGroup.GET("/borrow_records", borrowhandler.GetUserBorrowRecordsHandler)
		authGroup.POST("/borrow_records/borrow", borrowhandler.BorrowBookHandler)
		authGroup.POST("/borrow_records/return", borrowhandler.ReturnBookHandler)

		adminGroup := authGroup.Group("/admin")
		adminGroup.Use(middleware.AdminMiddleware())
		{
			// 用户管理
			adminGroup.GET("/users", userhandler.AdminGetAllUserInfoHandler)
			adminGroup.GET("/users/:userid", userhandler.AdminGetUserInfoHandler)
			adminGroup.PUT("/users/:userid/password", userhandler.AdminChangePasswordHandler)
			adminGroup.DELETE("/users/:userid", authhandler.DelHandler)

			// 图书管理
			adminGroup.POST("/books", bookhandler.AddBookHandler)
			adminGroup.PUT("/books/:id", bookhandler.UpdateBookHandler)
			adminGroup.DELETE("/books/:id", bookhandler.DeleteBookHandler)

			// 借阅管理
			adminGroup.GET("/borrow_records", borrowhandler.GetAllBorrowRecordsHandler)
		}
	}
}
