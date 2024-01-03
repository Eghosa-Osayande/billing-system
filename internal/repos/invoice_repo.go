package repos

import (
	"blanq_invoice/database"
	"blanq_invoice/util"
	"context"
	"errors"
	"fmt"

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

func (repo *InvoiceRepo) GetInvoices(input *database.GetInvoiceWhereParams) (*util.PagedResult[database.Invoice], error) {
	ctx := context.Background()

	invoices, err := repo.db.GetInvoiceWhere(ctx, *input)

	if database.IsErrNoRows(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	invoiceList := []database.Invoice{}
	total := 0

	for index := range invoices {
		invoiceList = append(invoiceList, invoices[index].Invoice)
		total = int(invoices[index].TotalCount)

	}

	return util.NewPagedResult[database.Invoice](invoiceList, total), nil

}

func (repo *InvoiceRepo) CreateInvoice(input *database.CreateInvoiceParams, items *[]database.CreateInvoiceItemParams) (*database.InvoiceWithItems, error) {
	ctx := context.Background()
	sqDB := repo.db.GetSqlDB()
	tx, err := sqDB.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)
	db := repo.db.WithTx(tx)
	invoice, err := db.CreateInvoice(ctx, *input)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Println(pgErr.Message) // => syntax error at end of input
			fmt.Println(pgErr.Code)    // => 42601
		}
		return nil, err
	}
	
	if items != nil {
		for index := range *items {
			item := (*items)[index]
			item.InvoiceID = invoice.ID

			err := db.CreateInvoiceItem(ctx, item)
			if err != nil {
				return nil, err
			}
		}
	}

	newitems,err:=db.FindInvoiceItemsByInvoiceID(ctx, invoice.ID)
	if err != nil {
		return nil, err
	}
	tx.Commit(ctx)
	return &database.InvoiceWithItems{
		Invoice: invoice,
		Items:   newitems,
	}, nil
	
	
}
