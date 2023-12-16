package main

import (
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
)

func PageViewVideo(c *fiber.Ctx) error {
	// Здесь можно получить данные о видео, например, из базы данных или другого источника
	videoData := fiber.Map{
		"Title":       "Просмотр видео",
		"VideoTitle":  "Название видео",
		"Description": "Описание видео",
	}

	return c.Render("video_view", videoData)
}

func VideoExists(videoID string) bool {
	path := filepath.Join("hls", videoID, "output.m3u8")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false // Файл не существует
	}

	return true // Файл существует
}

func main() {

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		// Рендеринг index.html
		return c.Render("index", fiber.Map{
			"Title":   "Главная страница",
			"Message": "тест",
		})
	})

	app.Get("/video/:id", func(c *fiber.Ctx) error {
		// Получаем ID видео из URL
		videoID := c.Params("id")
		LoadSource := path.Join(videoID)
		if !VideoExists(videoID) {
			// Если видео не найдено, отправляем ответ 404
			return c.SendStatus(fiber.StatusNotFound)
		}
		// Рендеринг video.html с передачей ID видео

		return c.Render("page_video", fiber.Map{
			"Title":      "Видео " + videoID,
			"VideoID":    videoID,
			"LoadSource": LoadSource,
		})
	})

	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.SendFile("./index.html") // Укажите путь к вашему HTML-файлу
	// })

	// app.Get("/page", func(c *fiber.Ctx) error {
	// 	return c.Render("index", fiber.Map{
	// 		"Title":   "Главная страница",
	// 		"Message": "Добро пожаловать!",
	// 	})
	// })

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "*", // Или укажите конкретные домены, например, "http://example.com"
		AllowHeaders:     "Origin, Content-Type, Accept",
	}))

	app.Static("/", ".")

	log.Fatal(app.Listen(":8080"))
}
