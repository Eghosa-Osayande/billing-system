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

type FetchInvoiceFilter struct {
	InvoiceID *uuid.UUID `json:"invoice_id" validate:"omitnil"`
	ClientID  *uuid.UUID `json:"client_id" validate:"omitnil"`
	Limit     *int       `json:"limit"`
	Cursor    *string    `json:"cursor"`
}

func (handler *InvoiceHandler) HandleAll(ctx *fiber.Ctx) error {
	input, valerr := util.ValidateRequestBody[*FetchInvoiceFilter](ctx.Body(), &FetchInvoiceFilter{})

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
			return ctx.JSON(util.NewSuccessResponseWithData[any]("Create a business first", nil))
		}

		var cursor_time *time.Time
		var cursor_id *uuid.UUID

		if input.Cursor != nil {
			created_at, id, errCsr := util.DecodeCursor(*input.Cursor)

			if errCsr != nil {
				return fiber.NewError(fiber.ErrBadRequest.Code, "invalid-cursor")
			}

			parsedId, err := uuid.Parse(id)
			if err != nil {
				return fiber.NewError(fiber.ErrBadRequest.Code, "invalid-cursor")
			}

			cursor_time, cursor_id = &created_at, &parsedId

		}

		params := database.FindInvoicesWhereParams{
			BusinessID: &business.ID,
			ClientID:   input.ClientID,
			InvoiceID:  input.InvoiceID,
			CursorTime: cursor_time,
			CursorID:   cursor_id,
			Limit:      input.Limit,
		}
		invoices, err := handler.config.InvoiceRepo.FindInvoicesWhere(&params)

		if err != nil {
			log.Println(err)
			return fiber.NewError(fiber.ErrInternalServerError.Code)
		}

		return ctx.JSON(
			util.NewSuccessResponseWithData[any](
				"Success",
				util.ListToPagedResult(
					invoices,
					func(
						item database.InvoiceWithItemsT[any],
					) (t time.Time, uuid string) {
						return item.CreatedAt, item.ID.String()
					}),
			),
		)
	}
}

type InvoiceItemInput struct {
	Name         string           `json:"name" validate:"required" db:"name"`
	Price        decimal.Decimal  `json:"price" validate:"required,number" db:"price"`
	Quantity     int              `json:"quantity" validate:"required,number" db:"quantity"`
	Discount     *decimal.Decimal `json:"discount" db:"discount,number"`
	DiscountType *string          `json:"discount_type" db:"discount_type"`
}

type CreateInvoiceInput struct {
	Currency        *string             `json:"currency"`
	PaymentDueDate  *string             `json:"payment_due_date" validate:"omitnil,datetime=2006-01-02"`
	DateOfIssue     *string             `json:"date_of_issue" validate:"omitnil,datetime=2006-01-02"`
	Notes           *string             `json:"notes"`
	PaymentMethod   *string             `json:"payment_method"`
	PaymentStatus   *string             `json:"payment_status"`
	Items           *[]InvoiceItemInput `json:"items"`
	ClientID        *uuid.UUID          `json:"client_id"`
	ShippingFeeType *string             `json:"shipping_fee_type"`
	ShippingFee     *decimal.Decimal    `json:"shipping_fee"`
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
			return ctx.JSON(util.NewSuccessResponseWithData[any]("Create a business first", nil))
		}

		itemsParams := make([]database.CreateInvoiceItemParams, 0)

		if input.Items != nil {
			itemsParams = make([]database.CreateInvoiceItemParams, len(*input.Items))
			for index := range *input.Items {
				item := (*input.Items)[index]

				itemsParams[index] = database.CreateInvoiceItemParams{
					Title:        item.Name,
					Price:        item.Price,
					Quantity:     decimal.NewFromInt(int64(item.Quantity)),
					Discount:     item.Discount,
					DiscountType: item.DiscountType,
				}
			}
		}
		var paymentDueDate *time.Time
		if input.PaymentDueDate != nil {
			paymentDueDate = new(time.Time)
			*paymentDueDate, err = time.Parse("2006-01-02", *input.PaymentDueDate)
			if err != nil {
				return fiber.NewError(fiber.ErrBadRequest.Code, "Invalid date format")
			}
		}
		var issueDate *time.Time
		if input.DateOfIssue != nil {
			issueDate = new(time.Time)
			*issueDate, err = time.Parse("2006-01-02", *input.DateOfIssue)
			if err != nil {
				return fiber.NewError(fiber.ErrBadRequest.Code, "Invalid date format")
			}
		}
		invoice, err := handler.config.InvoiceRepo.CreateInvoice(
			&database.CreateInvoiceParams{
				BusinessID:      business.ID,
				Currency:        input.Currency,
				PaymentDueDate:  paymentDueDate,
				DateOfIssue:     issueDate,
				Notes:           input.Notes,
				PaymentMethod:   input.PaymentMethod,
				PaymentStatus:   input.PaymentStatus,
				ClientID:        input.ClientID,
				ShippingFeeType: input.ShippingFeeType,
				ShippingFee:     input.ShippingFee,
			},
			itemsParams)

		if err != nil {
			log.Println(err)
			return fiber.NewError(fiber.ErrInternalServerError.Code)
		}

		return ctx.JSON(util.NewSuccessResponseWithData[any]("Success", invoice))
	}
}