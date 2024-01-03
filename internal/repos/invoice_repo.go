package repos

import (
	"blanq_invoice/database"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type InvoiceRepo struct {
	db *database.Queries
}

func NewInvoiceRepo(db *database.Queries) *InvoiceRepo {

	return &InvoiceRepo{
		db: db,
	}

}

func (repo *InvoiceRepo) FindAllInvoicesByBusinessID(businessID uuid.UUID) ([]database.InvoiceWithItems, error) {
	ctx := context.Background()
	invoices, err := repo.db.FindAllBusinessInvoices(ctx, businessID)
	if err != nil {
		return nil, err
	}

	invoiceWithItems := make([]database.InvoiceWithItems, len(invoices))
	for index := range invoices {
		invoice := invoices[index]
		invoiceitems, err := repo.db.FindInvoiceItemsByInvoiceId(ctx, invoice.ID)
		if err != nil {
			return nil, err
		}
		invoiceWithItems[index]= database.InvoiceWithItems{
			Invoice: invoice,
			Items: invoiceitems,
		}
	}
	return invoiceWithItems, nil
}

func (repo *InvoiceRepo) CreateInvoice(input *database.CreateInvoiceParams, items []database.CreateInvoiceItemParams) (*database.InvoiceWithItems, error) {
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Println(pgErr.Message) // => syntax error at end of input
			fmt.Println(pgErr.Code)    // => 42601
		}
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

	fmt.Println(newItems)

	tx.Commit(ctx)

	return &database.InvoiceWithItems{
		Invoice: newinvoice,
		Items:   newItems,
	}, nil

}
