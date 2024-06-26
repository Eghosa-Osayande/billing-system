package repos

import (
	"blanq_invoice/database"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type InvoiceRepo struct {
	db *database.Queries
}

func NewInvoiceRepo(db *database.Queries) *InvoiceRepo {

	return &InvoiceRepo{
		db: db,
	}

}

func (repo *InvoiceRepo) CreateInvoice(input *database.CreateInvoiceParams, items []database.CreateInvoiceItemParams) (*database.FullInvoice, error) {

	ctx := context.Background()
	sqDB := repo.db.GetSqlDB()
	tx, err := sqDB.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)
	db := repo.db.WithTx(tx)

	newinvoice, err := db.CreateInvoice(ctx, *input)

	if err != nil {

		fmt.Println(err)

		return nil, err
	}

	newItems := make([]database.Invoiceitem, len(items))
	for index := range items {
		item := items[index]
		item.InvoiceID = newinvoice.ID

		newItem, err := db.CreateInvoiceItem(ctx, item)
		if err != nil {
			return nil, err
		}
		newItems[index] = newItem

	}

	invoiceWithTotal, err := calculateInvoiceTotal(db, newinvoice.ID)

	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)

	if err != nil {
		return nil, err

	}

	return invoiceWithTotal.ToFullInvoice()

}

func (repo *InvoiceRepo) FindInvoicesWhere(input *database.FindInvoicesWhereParams) ([]database.FullInvoice, error) {
	ctx := context.Background()

	result, err := repo.db.FindInvoicesWhere(ctx, *input)

	if err != nil {
		return nil, err
	}

	invoiceList := make([]database.FullInvoice, len(result))

	for index := range result {
		i, err := result[index].ToFullInvoice()
		if err != nil {
			return nil, err
		}
		invoiceList[index] = *i
	}

	return invoiceList, nil

}

func (repo *InvoiceRepo) UpdateInvoice(input *database.UpdateInvoiceParams, items []database.CreateInvoiceItemParams) (*database.FullInvoice, error) {

	ctx := context.Background()
	sqDB := repo.db.GetSqlDB()
	tx, err := sqDB.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)
	db := repo.db.WithTx(tx)
	newInvoice, err := db.UpdateInvoice(ctx, *input)
	if err != nil {
		return nil, err
	}

	err = db.DeleteInvoiceItemByInvoiceID(ctx, newInvoice.ID)
	if err != nil {
		return nil, err
	}

	newItems := make([]database.Invoiceitem, len(items))
	for index := range items {
		item := items[index]
		item.InvoiceID = newInvoice.ID

		newItem, err := db.CreateInvoiceItem(ctx, item)
		if err != nil {
			return nil, err
		}
		newItems[index] = newItem

	}

	invoiceWithTotal, err := calculateInvoiceTotal(db, newInvoice.ID)

	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)

	if err != nil {
		return nil, err

	}

	return invoiceWithTotal.ToFullInvoice()

}

func (repo *InvoiceRepo) GetInvoicesCount(businessId uuid.UUID) (*database.GetInvoiceCountsRow, error) {
	ctx := context.Background()

	count, err := repo.db.GetInvoiceCounts(ctx, businessId)

	if err != nil {
		return nil, err
	}

	return &count, nil
}

func calculateInvoiceTotal(db *database.Queries, invoiceId uuid.UUID) ( *database.FindInvoicesWhereRow,  error) {

	ctx := context.Background()
	
	invoice, err := db.FindInvoiceById(ctx, invoiceId)
	if err != nil {
		return nil, err
	}

	items, err := db.FindInvoiceItemsByInvoiceId(ctx, invoice.ID)

	if err != nil {
		return nil, err
	}

	itemTotals := make([]decimal.Decimal, len(items))
	for i := range items {
		item := items[i]
		amt := item.Price.Mul(item.Quantity)
		if item.Discount != nil {

			if item.DiscountType == nil {
				return nil, errors.New("discount type must be specified")
			}

			switch *item.DiscountType {
			case "percent":
				di := amt.Mul(*item.Discount).Div(decimal.NewFromInt(100))

				amt = amt.Sub(di)
			case "fixed":
				amt = amt.Sub(*item.Discount)

			default:
				return nil, errors.New("invalid discount type for an item, must be either percent or fixed")
			}

		}
		itemTotals[i] = amt

	}
	zero := decimal.NewFromInt(0)
	subtotal := decimal.Sum(zero, itemTotals...)

	shippingfee := zero
	tax := zero

	if invoice.ShippingFee != nil {
		shippingfee = *invoice.ShippingFee
		if invoice.ShippingFeeType == nil {
			return nil, errors.New("shipping fee type must be specified")
		}

		switch *invoice.ShippingFeeType {
		case "percent":
			shippingfee = subtotal.Mul(shippingfee).Div(decimal.NewFromInt(100))

		case "fixed":

		default:
			return nil, errors.New("invalid shipping fee type, must be either percent or fixed")

		}
	}

	if invoice.Tax != nil {
		tax = subtotal.Mul(*invoice.Tax).Div(decimal.NewFromInt(100))
	}

	total := subtotal.Add(shippingfee).Add(tax)

	invoiceUpdated, err := db.UpdateInvoice(ctx, database.UpdateInvoiceParams{ID: invoice.ID, Total: &total})

	if err != nil {
		return nil, err
	}

	finalInvoice, err := db.FindInvoicesWhere(ctx, database.FindInvoicesWhereParams{InvoiceID: &invoiceUpdated.ID})
	if err != nil {
		return nil, err
	}
	if len(finalInvoice) != 1 {
		return nil, errors.New("invoice not found")
	}
	updatedInvoice := &finalInvoice[0]

	return updatedInvoice, nil
}
