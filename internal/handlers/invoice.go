package handlers

import (
	"blanq_invoice/database"
	"blanq_invoice/internal/repos"
	"blanq_invoice/util"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type InvoiceHandler struct {
	config *repos.ApiRepos
}

func NewInvoiceHandler(config *repos.ApiRepos) *InvoiceHandler {
	return &InvoiceHandler{config: config}
}

func (handler *InvoiceHandler) RegisterHandlers(router fiber.Router) {
	router.Get("/all", handler.HandleAll)
	router.Post("/new", handler.HandleCreateInvoice)

}

func (handler *InvoiceHandler) HandleAll(ctx *fiber.Ctx) error {
	if userId, ok := ctx.Context().UserValue("user_id").(string); !ok {
		return fiber.NewError(fiber.ErrUnauthorized.Code, "Unauthorized")
	} else {

		userId, err := uuid.Parse(userId)
		if err != nil {
			log.Println(err)
			return fiber.NewError(fiber.ErrInternalServerError.Code)
		}
		business, err := handler.config.BusinessRepo.FindBusinessByUserID(userId)
		if err != nil {
			log.Println(err)
			return fiber.NewError(fiber.ErrInternalServerError.Code)
		}
		if business == nil {
			return ctx.JSON(util.NewSuccessResponseWithData[*util.PagedResult[database.Invoice]]("Create a business first", nil))
		}
		limit, offset := util.GetPaginationFromQueries(ctx.Queries())

		invoices, err := handler.config.InvoiceRepo.GetInvoices(&database.GetInvoiceWhereParams{
			BusinessID: business.ID,
			ClientID:   nil,
			Limit:      int32(limit),
			Offset:     int32(offset),
		})

		if err != nil {
			log.Println(err)
			return fiber.NewError(fiber.ErrInternalServerError.Code)
		}
		return ctx.JSON(util.NewSuccessResponseWithData[*util.PagedResult[database.Invoice]]("Success", invoices))
	}
}

type InvoiceItem struct {
	Name     string   `json:"name" validate:"required" db:"name"`
	Price    float64  `json:"price" validate:"required" db:"price"`
	Quantity int      `json:"quantity" validate:"required" db:"quantity"`
	Discount *float64 `json:"discount" db:"discount"`
}

type CreateInvoiceInput struct {
	Currency        *string          `json:"currency"`
	PaymentDueDate  *time.Time       `json:"payment_due_date" validate:"omitnil,datetime=2006-01-02"`
	DateOfIssue     *time.Time       `json:"date_of_issue"`
	Notes           *string          `json:"notes"`
	PaymentMethod   *string          `json:"payment_method"`
	PaymentStatus   *string          `json:"payment_status"`
	Items           *[]InvoiceItem   `json:"items"`
	ClientID        *uuid.UUID       `json:"client_id"`
	ShippingFeeType *string          `json:"shipping_fee_type"`
	ShippingFee     *decimal.Decimal `json:"shipping_fee"`
}

func (handler *InvoiceHandler) HandleCreateInvoice(ctx *fiber.Ctx) error {

	input, valerr := util.ValidateRequestBody[*CreateInvoiceInput](ctx.Body(), &CreateInvoiceInput{})

	if valerr != nil {
		return valerr
	}

	if userId, ok := ctx.Context().UserValue("user_id").(string); !ok {
		return fiber.NewError(fiber.ErrUnauthorized.Code, "Unauthorized")
	} else {

		userId, err := uuid.Parse(userId)
		if err != nil {
			log.Println(err)
			return fiber.NewError(fiber.ErrInternalServerError.Code)
		}

		business, err := handler.config.BusinessRepo.FindBusinessByUserID(userId)
		if err != nil {
			log.Println(err)
			return fiber.NewError(fiber.ErrInternalServerError.Code)
		}

		if business == nil {
			return ctx.JSON(util.NewSuccessResponseWithData[*util.PagedResult[database.Invoice]]("Create a business first", nil))
		}

		itemsParams := new([]database.CreateInvoiceItemParams)

		if input.Items != nil {
			*itemsParams = make([]database.CreateInvoiceItemParams, len(*input.Items))
			for index := range *input.Items {
				item := (*input.Items)[index]
				var discount *decimal.Decimal

				if item.Discount == nil {
					discount = nil
				} else {
					discount = new(decimal.Decimal)
					*discount = decimal.NewFromFloat(*item.Discount)
				}
				
				(*itemsParams)[index] = database.CreateInvoiceItemParams{
					Title:        item.Name,
					Price:        decimal.NewFromFloat(item.Price),
					Quantity:     decimal.NewFromInt(int64(item.Quantity)),
					Discount:     discount,
					DiscountType: nil,
				}
			}
		}

		invoice, err := handler.config.InvoiceRepo.CreateInvoice(
			&database.CreateInvoiceParams{
				BusinessID:      business.ID,
				Currency:        input.Currency,
				PaymentDueDate:  input.PaymentDueDate,
				DateOfIssue:     input.DateOfIssue,
				Notes:           input.Notes,
				PaymentMethod:   input.PaymentMethod,
				PaymentStatus:   input.PaymentStatus,
				Items:           nil,
				ClientID:        input.ClientID,
				ShippingFeeType: input.ShippingFeeType,
				ShippingFee:     input.ShippingFee,
			},
			itemsParams)

		if err != nil {
			log.Println(err)
			return fiber.NewError(fiber.ErrInternalServerError.Code)
		}

		return ctx.JSON(util.NewSuccessResponseWithData[*database.InvoiceWithItems]("Success", invoice))
	}
}
