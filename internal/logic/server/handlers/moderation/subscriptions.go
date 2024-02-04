package moderation

import (
	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/server"
	"github.com/cufee/aftermath-core/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserSubscriptionsHandler(c *fiber.Ctx) error {
	userId := c.Params("userId")
	if userId == "" {
		return c.Status(400).JSON(server.NewErrorResponse("user id required", ""))
	}

	subs, err := database.FindSubscriptionsByUserID(userId)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "database.FindSubscriptionsByUserID"))
	}

	return c.JSON(server.NewResponse(subs))
}

func GetSubscriptionHandler(c *fiber.Ctx) error {
	subId := c.Params("id")
	if subId == "" {
		return c.Status(400).JSON(server.NewErrorResponse("user id required", ""))
	}

	id, err := primitive.ObjectIDFromHex(subId)
	if err != nil {
		return c.Status(400).JSON(server.NewErrorResponseFromError(err, "primitive.ObjectIDFromHex"))
	}

	subs, err := database.GetSubscriptionByID(id)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "database.FindSubscriptionsByUserID"))
	}

	return c.JSON(server.NewResponse(subs))
}

func UpdateSubscriptionHandler(c *fiber.Ctx) error {
	subId := c.Params("id")
	if subId == "" {
		return c.Status(400).JSON(server.NewErrorResponse("user id required", ""))
	}
	id, err := primitive.ObjectIDFromHex(subId)
	if err != nil {
		return c.Status(400).JSON(server.NewErrorResponseFromError(err, "primitive.ObjectIDFromHex"))
	}

	var body types.SubscriptionPayload
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(server.NewErrorResponseFromError(err, "c.BodyParser"))
	}

	userId := body.UserID
	if userId == "" {
		return c.Status(400).JSON(server.NewErrorResponse("user id required", ""))
	}
	if !body.Valid() {
		return c.Status(400).JSON(server.NewErrorResponse("invalid subscription payload", ""))
	}

	subscription, err := database.UpdateUserSubscription(id, body.ToSubscriptionUpdate())
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "database.FindSubscriptionsByUserID"))
	}

	return c.JSON(server.NewResponse(subscription))
}

func CreateUserSubscriptionsHandler(c *fiber.Ctx) error {
	var body types.SubscriptionPayload
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(server.NewErrorResponseFromError(err, "c.BodyParser"))
	}

	userId := body.UserID
	if userId == "" {
		return c.Status(400).JSON(server.NewErrorResponse("user id required", ""))
	}
	if !body.Valid() {
		return c.Status(400).JSON(server.NewErrorResponse("invalid subscription payload", ""))
	}

	subscription, err := database.AddNewUserSubscription(userId, body.ToUserSubscription())
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "database.FindSubscriptionsByUserID"))
	}

	return c.JSON(server.NewResponse(subscription))
}
