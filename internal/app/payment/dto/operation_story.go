package dto

import (
	cardEntities "go-user-microservice/internal/app/card/entities"
	"go-user-microservice/internal/app/payment/entities"
)

type OperationStory struct {
	entities.OperationStory
	cardEntities.Card
}
