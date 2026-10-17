package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/NASA_API/internal/core"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service core.ServiceActions
}

func NewHandler(service core.ServiceActions) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Mount(app *fiber.App) {
	app.Get("/api/v1/apods", h.GetAllApods)

	app.Get("/api/v1/apod/today", h.GetToday)
	app.Get("/api/v1/apod/today/image", h.GetTodayImage)

	app.Get("/api/v1/apod/image/:date", h.GetByDateImage)
	app.Get("/api/v1/apod/:date", h.GetByDate)
}

func (h *Handler) GetByDate(c fiber.Ctx) error {
	date := fiber.Params(c, "date", "")
	apod, err := h.service.GetByDate(date, context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(apod)
}

func (h *Handler) GetToday(c fiber.Ctx) error {
	time := time.Now()
	date := time.Format("2006-01-02")
	fmt.Println(date)
	apod, err := h.service.GetByDate(date, context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(apod)
}

func (h *Handler) GetTodayImage(c fiber.Ctx) error {
	time := time.Now()
	date := time.Format("2006-01-02")
	image, err := h.service.GetByDateImage(date, context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON((fiber.Map{
			"error": err.Error(),
		}))
	}

	c.Set("Content-Type", "image/jpeg")
	return c.Send(image)
}

func (h *Handler) GetByDateImage(c fiber.Ctx) error {
	date := fiber.Params(c, "date", "")
	image, err := h.service.GetByDateImage(date, context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON((fiber.Map{
			"error": err.Error(),
		}))
	}

	c.Set("Content-Type", "image/jpeg")
	return c.Send(image)
}

func (h *Handler) GetAllApods(c fiber.Ctx) error {
	apods, err := h.service.GetAllApods()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON((fiber.Map{
			"error": err.Error(),
		}))
	}

	return c.JSON(apods)
}
