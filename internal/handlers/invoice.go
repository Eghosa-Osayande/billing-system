package handlers

import (
	"blanq_invoice/database"
	"blanq_invoice/internal/repos"
	"blanq_invoice/middlewares"
	"blanq_invoice/util"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

type InvoiceHandler struct {
	config *repos.ApiRepos
}

func NewInvoiceHandler(config *repos.ApiRepos) *InvoiceHandler {
	return &InvoiceHandler{config: config}
}

func (handler *InvoiceHandler) RegisterHandlers(router fiber.Router) {
	router = router.Group("/invoices").Use(middlewares.AuthenticatedUserMiddleware).Use(middlewares.UserMustHaveBusinessMiddlewareInstance().Use)

	router.Get("/all", handler.HandleAll)
	router.Post("/new", handler.HandleCreateInvoice)
	router.Post("/update", handler.HandleUpdateInvoice)

}

type FetchInvoiceFilter struct {
	InvoiceID *uuid.UUID `json:"invoice_id" validate:"omitnil"`
	ClientID  *uuid.UUID `json:"client_id" validate:"omitnil"`
	Limit     *int32     `json:"limit"`
	Cursor    *string    `json:"cursor"`
}

func (handler *InvoiceHandler) HandleAll(ctx *fiber.Ctx) error {
	input, valerr := ValidateRequestBody[*FetchInvoiceFilter](ctx.Body(), &FetchInvoiceFilter{})

	if valerr != nil {
		return valerr
	}

	businessId, err := util.GetUserBusinessIdFromContext(ctx)
	if err != nil {
		return err
	}

	var cursor_time pgtype.Timestamptz
	var cursor_id *int64

	if input.Cursor != nil {
		created_at, id, errCsr := util.DecodeCursor(*input.Cursor)

		if errCsr != nil {
			return fiber.NewError(fiber.ErrBadRequest.Code, "invalid-cursor")
		}

		cursor_time, cursor_id = pgtype.Timestamptz{Time: created_at, Valid: true}, &id

	}

	params := database.FindInvoicesWhereParams{
		BusinessID: businessId,
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
					item database.FullInvoice,
				) (t time.Time, uuid int64) {
					return item.CreatedAt.Time, item.CountID
				}),
		),
	)

}

type InvoiceItemInput struct {
	Name         string           `json:"name" validate:"required" db:"name"`
	Price        decimal.Decimal  `json:"price" validate:"required" db:"price"`
	Quantity     int              `json:"quantity" validate:"required" db:"quantity"`
	Discount     *decimal.Decimal `json:"discount" db:"discount" validate:"required_with=DiscountType"`
	DiscountType *string          `json:"discount_type" db:"discount_type" validate:"oneof=fixed percent,required_with=Discount"`
}

type CreateInvoiceInput struct {
	Currency        *string             `json:"currency"`
	CurrencySymbol  *string             `json:"currency_symbol"`
	PaymentDueDate  *string             `json:"payment_due_date" validate:"omitnil,datetime=2006-01-02"`
	DateOfIssue     *string             `json:"date_of_issue" validate:"omitnil,datetime=2006-01-02"`
	Notes           *string             `json:"notes"`
	PaymentMethod   *string             `json:"payment_method"`
	Items           *[]InvoiceItemInput `json:"items"`
	ClientID        *uuid.UUID          `json:"client_id"`
	ShippingFeeType *string             `json:"shipping_fee_type" validate:"omitnil,oneof=fixed percent"`
	ShippingFee     *decimal.Decimal    `json:"shipping_fee" validate:"required_with=ShippingFeeType"`
	Tax             *decimal.Decimal    `json:"tax"`
	PaymentStatus   *string             `json:"payment_status" validate:"omitempty,oneof=paid unpaid partial_paid over_due"`
}

func (handler *InvoiceHandler) HandleCreateInvoice(ctx *fiber.Ctx) error {

	input, valerr := ValidateRequestBody[*CreateInvoiceInput](ctx.Body(), &CreateInvoiceInput{})

	if valerr != nil {
		return valerr
	}

	businessIdPtr, err := util.GetUserBusinessIdFromContext(ctx)
	if err != nil {
		return err
	}
	businessId := *businessIdPtr

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
	var paymentDueDate pgtype.Timestamptz
	if input.PaymentDueDate != nil {

		d, err := time.Parse("2006-01-02", *input.PaymentDueDate)
		if err != nil {
			return fiber.NewError(fiber.ErrBadRequest.Code, "Invalid date format")
		}
		paymentDueDate = pgtype.Timestamptz{Time: d, Valid: true}

	}
	var issueDate pgtype.Timestamptz
	if input.DateOfIssue != nil {
		d, err := time.Parse("2006-01-02", *input.DateOfIssue)
		if err != nil {
			return fiber.NewError(fiber.ErrBadRequest.Code, "Invalid date format")
		}
		issueDate = pgtype.Timestamptz{Time: d, Valid: true}
	}
	if input.ClientID != nil {
		cl, err := handler.config.ClientRepo.FindBusinessClientById(*input.ClientID, businessId)
		log.Println(cl, err)
		if err != nil || cl == nil {
			return fiber.NewError(fiber.ErrBadRequest.Code, "Client not found")
		}
	}

	invoice, err := handler.config.InvoiceRepo.CreateInvoice(
		&database.CreateInvoiceParams{
			BusinessID:      businessId,
			Currency:        input.Currency,
			CurrencySymbol:  input.CurrencySymbol,
			PaymentDueDate:  paymentDueDate,
			DateOfIssue:     issueDate,
			Notes:           input.Notes,
			PaymentMethod:   input.PaymentMethod,
			ClientID:        input.ClientID,
			ShippingFeeType: input.ShippingFeeType,
			ShippingFee:     input.ShippingFee,
			Total:           &decimal.Decimal{},
			Tax:             input.Tax,
			PaymentStatus:   input.PaymentStatus,
		},
		itemsParams)

	if err != nil {
		log.Println(err)
		return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	}

	return ctx.JSON(util.NewSuccessResponseWithData[any]("Success", invoice))
}

type UpdateInvoiceInput struct {
	InvoiceID uuid.UUID `json:"invoice_id" validate:"required"`
	CreateInvoiceInput
}

func (handler *InvoiceHandler) HandleUpdateInvoice(ctx *fiber.Ctx) error {

	input, valerr := ValidateRequestBody[*UpdateInvoiceInput](ctx.Body(), &UpdateInvoiceInput{})

	if valerr != nil {
		return valerr
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
	var paymentDueDate pgtype.Timestamptz
	if input.PaymentDueDate != nil {

		d, err := time.Parse("2006-01-02", *input.PaymentDueDate)
		if err != nil {
			return fiber.NewError(fiber.ErrBadRequest.Code, "Invalid date format")
		}
		paymentDueDate = pgtype.Timestamptz{Time: d, Valid: true}

	}
	var issueDate pgtype.Timestamptz
	if input.DateOfIssue != nil {
		d, err := time.Parse("2006-01-02", *input.DateOfIssue)
		if err != nil {
			return fiber.NewError(fiber.ErrBadRequest.Code, "Invalid date format")
		}
		issueDate = pgtype.Timestamptz{Time: d, Valid: true}
	}

	result, err := handler.config.InvoiceRepo.UpdateInvoice(
		&database.UpdateInvoiceParams{
			ID:              input.InvoiceID,
			Currency:        input.Currency,
			PaymentDueDate:  paymentDueDate,
			DateOfIssue:     issueDate,
			Notes:           input.Notes,
			PaymentMethod:   input.PaymentMethod,
			ClientID:        input.ClientID,
			ShippingFeeType: input.ShippingFeeType,
			ShippingFee:     input.ShippingFee,
			Total:           &decimal.Decimal{},
			PaymentStatus:   input.PaymentStatus,
		}, itemsParams)

	if err != nil {
		log.Println(err)
		return fiber.NewError(fiber.ErrInternalServerError.Code, "Internal Server Error")
	}

	return ctx.JSON(util.NewSuccessResponseWithData[*database.FullInvoice]("Invoice updated successfully", result))
}
