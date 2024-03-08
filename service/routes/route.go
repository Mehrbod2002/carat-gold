package routes

import (
	adminAuth "carat-gold/app/controllers/admin"
	adminSetter "carat-gold/app/controllers/admin/setter"
	adminView "carat-gold/app/controllers/admin/views"
	user "carat-gold/app/controllers/user"
	"carat-gold/app/metatrader"

	"carat-gold/models"
	"carat-gold/utils"
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	OnlineClients = make(map[*websocket.Conn]*models.Client)
)

func SetupRouter(dataChannel chan interface{}) *gin.Engine {
	r := gin.Default()

	store := cookie.NewStore([]byte("session-secret"))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   int(time.Hour) * 365,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})

	r.Use(sessions.Sessions("token", store))
	r.Static("/static", "./CDN")
	r.POST(os.Getenv("CALLBACK"), models.HandleIPN)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://0.0.0.0:3001", "http://localhost:3001", "http://127.0.0.1:5173", "https://admin.goldshop24.co", "https://goldshop24.co"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "HEAD", "DELETE"},
		AllowHeaders:     []string{"Origin", "Set-Cookie", "Cookie", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Set-Cookie", "Cookie"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://0.0.0.0:3001" || origin == "http://localhost:3001" || origin == "http://127.0.0.1:5173" || origin == "https://goldshop24.co" || origin == "https://admin.goldshop24.co"
		},
	}))

	r.GET("/feed", func(c *gin.Context) {
		metatrader.HandleWebSocket(c.Writer, c.Request, dataChannel)
	})

	apis := r.Group("/api")
	authRoutes := apis.Group("/auth")
	{
		authRoutes.POST("/user/send_otp", user.SendOTP)
		authRoutes.POST("/admin/login", adminAuth.AdminLogin)
		authRoutes.POST("/user/login_step_1", user.LoginOneTimeLoginStep1)
		authRoutes.POST("/user/login_step_2", user.LoginOneTimeLoginStep2)
		authRoutes.POST("/user/register", user.Register)
	}

	public := apis.Group("/public")
	{
		public.GET("/get_products", user.ViewProducts)
	}

	userRoutes := apis.Group("/user")
	userRoutes.Use(AuthMiddleware())
	{
		userRoutes.POST("/edit_user", user.EditUser)
		userRoutes.GET("/me", user.GetUser)
		userRoutes.POST("/update_currency", user.SetCurrency)
		userRoutes.GET("/general_data", user.GeneralData)
		userRoutes.GET("/products", user.GetProducts)
		userRoutes.POST("/upload_documents", user.SendDocuments)
		userRoutes.POST("/update_fcm", user.UpdateFcm)
	}

	supportRoutes := apis.Group("/support")
	supportRoutes.Use(AdminAuthMiddleware())
	{
	}

	adminRoutes := apis.Group("/admin")
	adminRoutes.Use(AdminAuthMiddleware())
	{
		adminRoutes.GET("/get_users", adminView.ViewAllUsers)
		adminRoutes.POST("/delete_user", adminSetter.SetDeleteUser)
		adminRoutes.POST("/delete_product", adminSetter.SetDeleteProduct)
		adminRoutes.POST("/logout", adminAuth.AdminLogout)
		adminRoutes.POST("/freeze_user", adminSetter.SetFreezeUser)
		adminRoutes.POST("/unfreeze_user", adminSetter.SetUnFreezeUser)
		adminRoutes.POST("/set_user_permissions", adminSetter.SetUserPermissions)
		adminRoutes.POST("/edit_user", adminSetter.SetUser)
		adminRoutes.POST("/get_users", adminView.ViewAllUsers)
		adminRoutes.POST("/define_user", adminSetter.SetDefineUser)
		adminRoutes.GET("/get_symbols", adminView.ViewSymbols)
		adminRoutes.POST("/delete_symbol", adminSetter.SetDeleteSymbol)
		adminRoutes.POST("/set_symbol", adminSetter.SetSymbols)
		adminRoutes.POST("/cancel_order", adminSetter.SetCancelOrder)
		adminRoutes.POST("/set_order", adminSetter.SetOrders)
		adminRoutes.GET("/current_orders", adminView.ViewCurrentOrders)
		adminRoutes.GET("/history_orders", adminView.ViewHistoryOrders)
		adminRoutes.POST("/edit_product", adminSetter.SetEditProduct)
		adminRoutes.POST("/edit_currency", adminSetter.SetEditCurrency)
		adminRoutes.POST("/set_product", adminSetter.SetProduct)
		adminRoutes.GET("/get_currencies", adminView.ViewCurrencies)
		adminRoutes.GET("/get_delivery_methods", adminView.ViewDeliveryMethods)
		adminRoutes.GET("/get_payment_methods", adminView.ViewPaymentMethods)
		adminRoutes.POST("/set_delivery_methods", adminSetter.SetDeliveryMethods)
		adminRoutes.POST("/edit_delivery_methods", adminSetter.SetEditDeliveryMethods)
		adminRoutes.POST("/delete_delivery_method", adminSetter.SetDeleteDeliveryMethodl)
		adminRoutes.POST("/set_payment_method", adminSetter.SetPayment)
		adminRoutes.POST("/edit_payment_method", adminSetter.SetEditPayment)
		adminRoutes.POST("/delete_paymnet_method", adminSetter.SetDeletePayment)
		adminRoutes.POST("/set_call_center", adminSetter.SetCallCenterDatas)
		adminRoutes.GET("/get_call_center", adminView.ViewCallCenter)
		adminRoutes.POST("/set_fandq", adminSetter.SetFANDQ)
		adminRoutes.POST("/edit_fandq", adminSetter.SetEditFANDQ)
		adminRoutes.POST("/delete_fandq", adminSetter.SetDeleteFANDQ)
		adminRoutes.GET("/get_fandq", adminView.ViewFANDQ)
		adminRoutes.GET("/get_metatrader_account", adminView.ViewMetaTrader)
		adminRoutes.POST("/set_metatrader_account", adminSetter.SetMetaData)
		adminRoutes.GET("/get_user_purchases", adminView.ViewPurchase)
	}
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	r.GET("/notification", AdminAuthMiddleware(), func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(err)
			return
		}
		user, _ := models.ValidateSession(c)
		var currentUser models.User
		db, DBerr := utils.GetDBWSS()
		if DBerr != nil {
			log.Println(DBerr)
			models.ErrInSocket(ws, user, "internal_error")
			return
		}
		allowed := false
		err = db.Collection("users").
			FindOne(context.Background(), bson.M{
				"_id": user.ID,
			}).Decode(&currentUser)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				err = db.Collection("admin").
					FindOne(context.Background(), bson.M{
						"phone": user.PhoneNumber,
					}).Decode(&currentUser)
				if err != nil {
					log.Println(err)
					utils.InternalError(c)
					return
				}
				allowed = true
			}
			log.Println(err)
			utils.InternalError(c)
			return
		}
		if !allowed {
			for _, action := range currentUser.Permissions.Actions {
				if action == models.ActionSendNotification {
					allowed = true
				}
			}
		}
		if !allowed {
			models.ErrInSocket(ws, user, "permission_not_allowed")
			return
		}
		for {
			var data models.NotificationAdmin
			err := ws.ReadJSON(&data)
			if err != nil {
				log.Println(err)
				return
			}
			if data.Type == "send_notification" {
				for _, onlineUser := range OnlineClients {
					for _, user := range data.Names {
						if onlineUser.User.Name == user {
							err := onlineUser.Conn.WriteJSON(models.NotificationAdmin{
								Type:    "notification",
								Subject: data.Subject,
								Message: data.Message,
							})
							if err != nil {
								log.Println(err)
								return
							}
						}
					}
				}
			}
		}
	})
	r.GET("/message", AuthMiddleware(), func(c *gin.Context) {
		user, _ := models.ValidateSession(c)
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(err)
			models.ErrInSocket(ws, user, "internal_error")
			return
		}
		if user == nil {
			log.Println(err)
			models.ErrInSocket(ws, user, "internal_error")
			return
		}
		exist := false
		for _, i := range OnlineClients {
			if i.User == user {
				exist = true
			}
		}
		if !exist {
			OnlineClients[ws] = &models.Client{
				Conn: ws,
				User: user,
			}
		}
		if user.Freeze {
			models.ErrInSocket(ws, user, "freezed")
			return
		}
		for {
			var data models.UserMessage
			err := ws.ReadJSON(&data)
			if err != nil {
				log.Println(err)
				models.ErrInSocket(ws, user, "internal_error")
				return
			}
			if data.Type == "initial_messages" {
				if user.Freeze {
					models.ErrInSocket(ws, user, "freezed")
					return
				}
				var messages []models.Message
				db, DBerr := utils.GetDBWSS()
				if DBerr != nil {
					log.Println(DBerr)
					models.ErrInSocket(ws, user, "internal_error")
					return
				}
				id, err := primitive.ObjectIDFromHex(data.ID)
				if err != nil {
					models.ErrInSocket(ws, user, "internal_error")
					return
				}
				var frontUser models.User
				if err := db.Collection("users").FindOne(context.Background(), bson.M{
					"_id": id,
				}).Decode(&frontUser); err != nil {
					log.Println(err)
					models.ErrInSocket(ws, user, "internal_error")
					return
				}
				cursor, err := db.Collection("messages").Find(context.Background(), bson.M{
					"$or": []bson.M{
						{
							"$and": bson.M{
								"receiver": user.ID,
								"sender":   frontUser.ID,
							},
						}, {
							"$and": bson.M{
								"receiver": frontUser.ID,
								"sender":   user.ID,
							},
						}},
				})
				if err != nil {
					if err == mongo.ErrNoDocuments {
						err = ws.WriteJSON(models.Socket{
							ResponseTo: *user,
							Trigger:    "initial_messages",
							Validate:   true,
							Messages:   messages,
						})
						if err != nil {
							log.Println(err)
							models.ErrInSocket(ws, user, "internal_error")
							return
						}
					} else {
						log.Println(err)
						models.ErrInSocket(ws, user, "internal_error")
						return
					}
				}
				if err := cursor.All(context.Background(), &messages); err != nil {
					log.Println(err)
					cursor.Close(context.Background())
					utils.InternalError(c)
					return
				}
				if _, err := db.Collection("messages").UpdateMany(context.Background(), bson.M{
					"$and": bson.M{
						"receiver": frontUser.ID,
						"sender":   user.ID,
					},
				}, bson.M{
					"$set": bson.M{
						"seen": true,
					},
				}); err != nil {
					log.Println(err)
					cursor.Close(context.Background())
					models.ErrInSocket(ws, user, "internal_error")
					return
				}
				var lastMessage []models.Message
				for _, message := range messages {
					if message.Sender == user.ID {
						message.SenderUsername = "my_message"
						message.ReceiverUsername = "front_user"
					}
					lastMessage = append(lastMessage, message)
				}
				err = ws.WriteJSON(models.Socket{
					ResponseTo: *user,
					Trigger:    "initial_messages",
					Validate:   true,
					Messages:   lastMessage,
				})
				if err != nil {
					log.Println(err)
					cursor.Close(context.Background())
					models.ErrInSocket(ws, user, "internal_error")
					return
				}
			}
			if data.Type == "get_message" {
				if user.Freeze {
					models.ErrInSocket(ws, user, "freezed")
					return
				}
				var messages []models.Message
				db, DBerr := utils.GetDBWSS()
				if DBerr != nil {
					log.Println(DBerr)
					models.ErrInSocket(ws, user, "internal_error")
					return
				}
				id, err := primitive.ObjectIDFromHex(data.ID)
				if err != nil {
					log.Println(err)
					models.ErrInSocket(ws, user, "internal_error")
					return
				}
				var frontUser models.User
				if err := db.Collection("users").FindOne(context.Background(), bson.M{
					"_id": id,
				}).Decode(&frontUser); err != nil {
					log.Println(err)
					models.ErrInSocket(ws, user, "internal_error")
					return
				}
				cursor, err := db.Collection("messages").Find(context.Background(), bson.M{
					"$or": []bson.M{
						{
							"$and": bson.M{
								"receiver": user.ID,
								"sender":   frontUser.ID,
							},
						}, {
							"$and": bson.M{
								"receiver": frontUser.ID,
								"sender":   user.ID,
							},
						}},
				})
				if err != nil {
					if err == mongo.ErrNoDocuments {
						log.Println(err)
						err = ws.WriteJSON(models.Socket{
							ResponseTo: *user,
							Trigger:    "initial_messages",
							Validate:   true,
							Messages:   messages,
						})
						if err != nil {
							log.Println(err)
							models.ErrInSocket(ws, user, "internal_error")
							return
						}
					} else {
						models.ErrInSocket(ws, user, "internal_error")
						return
					}
				}
				if err := cursor.All(context.Background(), &messages); err != nil {
					log.Println(err)
					cursor.Close(context.Background())
					utils.InternalError(c)
					return
				}
				if _, err := db.Collection("messages").UpdateMany(context.Background(), bson.M{
					"$and": bson.M{
						"receiver": frontUser.ID,
						"sender":   user.ID,
					},
				}, bson.M{
					"$set": bson.M{
						"seen": true,
					},
				}); err != nil {
					log.Println(err)
					cursor.Close(context.Background())
					models.ErrInSocket(ws, user, "internal_error")
					return
				}
				var lastMessage []models.Message
				for _, message := range messages {
					if message.Sender == user.ID {
						message.SenderUsername = "my_message"
						message.ReceiverUsername = "front_user"
					}
					lastMessage = append(lastMessage, message)
				}
				err = ws.WriteJSON(models.Socket{
					ResponseTo: *user,
					Trigger:    "messages",
					Validate:   true,
					Messages:   lastMessage,
				})
				if err != nil {
					log.Println(err)
					cursor.Close(context.Background())
					models.ErrInSocket(ws, user, "internal_error")
					return
				}
			} else if data.Type == "send_message" {
				if user.Freeze {
					models.ErrInSocket(ws, user, "user_freezed")
					return
				}
				if data.TypeMessage == models.TextType {
					if len(data.Content) > 3000 || len(data.Content) == 0 {
						models.ErrInSocket(ws, user, "invalid_text_length")
						return
					}
				}
				if data.TypeMessage == models.FileType {
					_, err := base64.StdEncoding.DecodeString(data.Content)
					if err != nil {
						log.Println(err)
						models.ErrInSocket(ws, user, "invalid_file_content")
						return
					}
					fileSizeMB := float64(len(data.Content)) / (1024 * 1024)
					if fileSizeMB > 6 {
						models.ErrInSocket(ws, user, "invalid_file_size")
						return
					}
				}
				var sendTo models.User
				db, DBerr := utils.GetDBWSS()
				if DBerr != nil {
					log.Println(DBerr)
					models.ErrInSocket(ws, user, "internal_error")
					return
				}
				if err := db.Collection("users").FindOne(context.Background(), bson.M{}).Decode(&sendTo); err != nil {
					log.Println(err)
					models.ErrInSocket(ws, user, "internal_error")
					return
				}
				_, err := primitive.ObjectIDFromHex(sendTo.ID.String())
				if err != nil || user.ID == sendTo.ID {
					log.Println(err)
					models.ErrInSocket(ws, user, "internal_error")
					return
				}
				var sendMessage = models.Message{
					Sender:      user.ID,
					Receiver:    sendTo.ID,
					Content:     data.Content,
					MessageType: data.TypeMessage,
				}
				result, err := db.Collection("messages").InsertOne(context.Background(), &sendMessage)
				if err != nil {
					log.Println(err)
					models.ErrInSocket(ws, user, "internal_error")
					return
				}
				for _, online_user := range OnlineClients {
					if online_user.User.ID == sendTo.ID {
						sendMessage.Seen = true
						err = online_user.Conn.WriteJSON(models.Socket{
							ResponseTo:    *online_user.User,
							Trigger:       "message_sent",
							Validate:      true,
							SingleMessage: sendMessage,
							Message:       "message sent",
						})
						if err != nil {
							log.Println(err)
							models.ErrInSocket(ws, user, "internal_error")
							return
						}
					}
				}
				if sendMessage.Seen {
					if _, err := db.Collection("messages").UpdateMany(context.Background(), bson.M{
						"_id": result.InsertedID,
					}, bson.M{
						"$set": bson.M{
							"seen": true,
						},
					}); err != nil {
						log.Println(err)
						models.ErrInSocket(ws, user, "internal_error")
						return
					}
				}
				sendMessage.SenderUsername = "my_message"
				sendMessage.ReceiverUsername = "front_user"
				err = ws.WriteJSON(models.Socket{
					ResponseTo:    *user,
					Trigger:       "message_sent",
					Validate:      true,
					SingleMessage: sendMessage,
					Message:       "message sent",
				})
				if err != nil {
					log.Println(err)
					models.ErrInSocket(ws, user, "internal_error")
					return
				}
			}
		}
	})

	return r
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		token := session.Get("token")
		tokenString := c.GetHeader("Authorization")
		if token == nil && tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized",
				"data":    "unauthorized",
			})
			c.Abort()
			return
		}

		if token == nil {
			token = tokenString
		}

		jwtSecret := os.Getenv("SESSION_SECRET")
		if jwtSecret == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized",
				"data":    "unauthorized",
			})
			c.Abort()
			return
		}

		parsedToken, err := jwt.Parse(token.(string), func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !parsedToken.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized",
				"data":    "unauthorized",
			})
			c.Abort()
			return
		}

		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && claims["_id"] != nil {
			userID, ok := claims["_id"].(string)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"message": "Unauthorized",
					"data":    "unauthorized",
				})
				c.Abort()
				return
			}
			if _, err := primitive.ObjectIDFromHex(userID); err == nil {
				c.Next()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized",
				"data":    "unauthorized",
			})
			c.Abort()
			return
		}

		c.Abort()
	}
}

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		token := session.Get("token_admins")
		tokenString := c.GetHeader("Authorization")
		if token == nil && tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": http.StatusUnauthorized,
				"message": "Unauthorized",
				"data":    "unauthorized",
			})
			c.Abort()
			return
		}

		if token == nil {
			token = session.Get("token_supports")
			if token == nil {
				token = tokenString
			}
		}

		jwtSecret := os.Getenv("SESSION_SECRET")
		if jwtSecret == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "JWT secret not configured",
				"data":    "unauthorized",
			})
			c.Abort()
			return
		}

		parsedToken, err := jwt.Parse(token.(string), func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !parsedToken.Valid {
			log.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid token",
				"data":    "unauthorized",
			})
			c.Abort()
			return
		}

		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
			if claims["email"] == os.Getenv("ADMIN_USERNAME") {
				c.Next()
				return
			}
		}

		c.Abort()
	}
}
