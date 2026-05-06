package api

import (
	"context"
	"discord-server-go/config"
	"discord-server-go/handler"
	"discord-server-go/model"
	"discord-server-go/repository"
	"discord-server-go/service"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	imagekit "github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/option"
	cors "github.com/rs/cors/wrapper/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	router *gin.Engine
	once   sync.Once
)

type dataSources struct {
	DB       *gorm.DB
	ImageKit *imagekit.Client
}

func initDS(ctx context.Context, cfg config.Config) (*dataSources, error) {
	log.Printf("Initializing data sources\n")

	log.Printf("Connecting to Postgresql\n")

	db, err := gorm.Open(postgres.Open(cfg.DatabaseUrl), &gorm.Config{
		PrepareStmt: false,
	})
	if err != nil {
		return nil, fmt.Errorf("Error opening database: %w", err)
	}

	if gin.Mode() != gin.ReleaseMode {
		if err = db.AutoMigrate(&model.User{}); err != nil {
			return nil, fmt.Errorf("Error migrating models: %w", err)
		}
	}
	ik := imagekit.NewClient(
		option.WithPrivateKey(cfg.ImageKitPrivateKey),
		option.WithHeader("X-ImageKit-PublicKey", cfg.ImageKitPublicKey),
	)

	return &dataSources{
		DB:       db,
		ImageKit: &ik,
	}, nil
}

func (d *dataSources) close() error {
	return nil
}
func inject(d *dataSources, cfg config.Config) (*gin.Engine, error) {
	log.Printf("Injecting data sources")

	userRepository := repository.NewUserRepository(d.DB)

	fileRepository := repository.NewFileRepository(d.ImageKit)

	userService := service.NewUserService(&service.USConfig{
		UserRepository: userRepository,
		FileRepository: fileRepository,
	})

	router := gin.Default()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{cfg.CorsOrigin},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Authorization"},
	})

	router.Use(c)

	store := cookie.NewStore([]byte(cfg.SessionSecret))
	store.Options(sessions.Options{
		Domain:   "",
		MaxAge:   60 * 60 * 24 * 7,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	})
	router.Use(sessions.Sessions(model.CookieName, store))
	rate := limiter.Rate{
		Period: 1 * time.Hour,
		Limit:  1500,
	}
	memoryStore := memory.NewStore()
	rateLimiter := mgin.NewMiddleware(limiter.New(memoryStore, rate))

	router.Use(rateLimiter)
	handler.NewHandler(&handler.Config{
		R:           router,
		UserService: userService,
	})

	return router, nil
}

// Hàm này thay thế cho main() để khởi tạo hệ thống 1 lần duy nhất
func setup() {
	log.Println("Initializing Serverless Setup...")
	ctx := context.Background()

	// 1. Load config (Vercel sẽ lấy từ Environment Variables)
	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// 2. Init Data Sources (Gorm, ImageKit...)
	ds, err := initDS(ctx, cfg) // Tận dụng hàm initDS của bạn
	if err != nil {
		log.Fatalf("Unable to initialize data sources: %v", err)
	}

	// 3. Inject (Tạo router Gin)
	// Tận dụng hàm inject của bạn, nó sẽ trả về *gin.Engine
	router, err = inject(ds, cfg)
	if err != nil {
		log.Fatalf("Failure to inject: %v", err)
	}
}

// Vercel sẽ gọi hàm này cho mỗi request
func Handler(w http.ResponseWriter, r *http.Request) {
	// Đảm bảo chỉ khởi tạo router 1 lần để tiết kiệm tài nguyên
	once.Do(setup)

	// Để Gin xử lý request
	router.ServeHTTP(w, r)
}
