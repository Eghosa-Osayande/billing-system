package handlers

import (
	"blanq_invoice/database"
	"blanq_invoice/internal/repos"
	"blanq_invoice/middlewares"
	"blanq_invoice/util"
	"log"

	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
	config *repos.ApiRepos
}

func NewDashboardHandler(config *repos.ApiRepos) *DashboardHandler {
	return &DashboardHandler{config: config}
}

func (handler *DashboardHandler) RegisterHandlers(router fiber.Router) {
	router = router.Group("/dashboard").Use(middlewares.AuthenticatedUserMiddleware)

	router.Get("/", handler.HandleDashboard)

}

type dashboardResponse struct {
	Summary  *database.GetInvoiceCountsRow  `json:"summary"`
	Invoices []database.InvoiceWithItemsAny `json:"invoices"`
}

func (h *DashboardHandler) HandleDashboard(ctx *fiber.Ctx) error {
	userUUID, err := util.GetUserIdFromContext(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	business, err := h.config.BusinessRepo.FindBusinessByUserID(*userUUID)

	if err != nil {
		return ctx.JSON(util.NewSuccessResponse("User has not created any business yet"))
	}

	r, err := h.config.InvoiceRepo.GetInvoicesCount(business.ID)

	if err != nil {
		return ctx.JSON(util.NewErrorMessage("Error getting invoices count", err))
	}

	var limit int32 = 10

	r2, err := h.config.InvoiceRepo.FindInvoicesWhere(
		&database.FindInvoicesWhereParams{
			BusinessID: &business.ID,
			Limit:      &limit,
		},
	)

	if err != nil {
		return ctx.JSON(util.NewErrorMessage("Error getting invoices", err))
	}

	return ctx.JSON(util.NewSuccessResponseWithData[dashboardResponse]("Invoices count", dashboardResponse{
		Summary:  r,
		Invoices: r2,
	}))

}
