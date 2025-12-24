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
		publicGroup.POST("/login", authhandler.LoginHandler)
		publicGroup.POST("/register", authhandler.RegisterUserHandler)
		publicGroup.GET("/books", bookhandler.GetBooksHandler)
		publicGroup.GET("/books/:id", bookhandler.GetBookByIdHandler)
	}

	authGroup := r.Group("/api")
	authGroup.Use(middleware.AuthMiddleware(sessionservice))
	{
		authGroup.GET("/logout", authhandler.LogoutHandler)
		authGroup.POST("/user/del", authhandler.DelHandler)
		authGroup.GET("/user", userhandler.GetUserInfoHandler)
		authGroup.POST("/user/change_password", userhandler.UserChangePasswordHandler)
		authGroup.POST("/user/change_info", userhandler.ChangeUserInfoHandler)

		// 借阅相关
		authGroup.GET("/borrow_records", borrowhandler.GetUserBorrowRecordsHandler)
		authGroup.POST("/borrow", borrowhandler.BorrowBookHandler)
		authGroup.POST("/return", borrowhandler.ReturnBookHandler)

		adminGroup := authGroup.Group("/admin")
		adminGroup.Use(middleware.AdminMiddleware())
		{
			// 用户管理
			adminGroup.GET("/users", userhandler.AdminGetAllUserInfoHandler)
			adminGroup.POST("/users/detail", userhandler.AdminGetUserInfoHandler)
			adminGroup.POST("/users/password", userhandler.AdminChangePasswordHandler)
			adminGroup.POST("/users/delete", authhandler.DelHandler)

			// 图书管理
			adminGroup.POST("/books", bookhandler.AddBookHandler)
			adminGroup.PUT("/books/:id", bookhandler.UpdateBookHandler)
			adminGroup.DELETE("/books/:id", bookhandler.DeleteBookHandler)

			// 借阅管理
			adminGroup.GET("/borrow_records", borrowhandler.GetAllBorrowRecordsHandler)
		}
	}
}
