package database

import (
	"log"
	"strconv"

	"github.com/vsynclabs/billsoft/internals/models"
)

func (q *Query) CreateInvoice(invoice *models.Invoice) error {
	query := `INSERT INTO invoice (
				invoice_id,
				invoice_name,
				invoice_payment_status,
				invoice_reverse_charge,
				invoice_date,
				invoice_state,
				invoice_state_code,
				invoice_challan_number,
				invoice_vehicle_number,
				invoice_date_of_supply,
				invoice_place_of_supply,
				invoice_gst,
				user_id,
				billed_id,
				shipped_id
				) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`
	_, err := q.db.Exec(
		query,
		invoice.InvoiceId,
		invoice.InvoiceName,
		invoice.InvoicePaymentStatus,
		invoice.InvoiceReverseRecharge,
		invoice.InvoiceDate,
		invoice.InvoiceState,
		invoice.InvoiceStateCode,
		invoice.InvoiceChallanNumber,
		invoice.InvoiceVehicleNumber,
		invoice.InvoiceDateOfSupply,
		invoice.InvoicePlaceOfSupply,
		invoice.InvoiceGst,
		invoice.UserId,
		invoice.BilledId,
		invoice.ShippedId,
	)

	return err
}

func (q *Query) DeleteInvoice(invoiceId string) error {
	query := `DELETE FROM invoice WHERE invoice_id=$1`

	_, err := q.db.Exec(query, invoiceId)

	return err
}

func (q *Query) GetInvoices(userId string) ([]*models.InvoiceResponse, error) {
	query := `SELECT invoice_id,invoice_name,invoice_payment_status FROM invoice WHERE user_id = $1`

	rows, err := q.db.Query(query, userId)

	if err != nil {
		return nil, err
	}

	var invoices []*models.InvoiceResponse

	for rows.Next() {
		var invoice models.InvoiceResponse

		err := rows.Scan(&invoice.InvoiceId, &invoice.Name, &invoice.PaymentStatus)

		if err != nil {
			return nil, err
		}

		invoices = append(invoices, &invoice)
	}

	return invoices, nil
}

func (q *Query) DownloadInvoice(invoiceId string) (*models.InvoicePdf, error) {

	query1 := `SELECT
				product_name,
				product_hsn,
				product_quantity,
				product_unit,
				product_rate,
				product_total
			FROM product WHERE invoice_id=$1`

	rows, err := q.db.Query(query1, invoiceId)

	if err != nil {
		return nil, err
	}

	var productPdfs []*models.ProductPdf

	var grandTotal int32 = 0
	var totalQty int32 = 0

	for rows.Next() {
		var productPdf models.ProductPdf

		if err := rows.Scan(
			&productPdf.ProductName,
			&productPdf.ProductHsn,
			&productPdf.ProductQty,
			&productPdf.ProductUnit,
			&productPdf.ProductRate,
			&productPdf.Total,
		); err != nil {
			return nil, err
		}

		totalPrice, err := strconv.Atoi(productPdf.Total)

		if err != nil {
			return nil, err
		}

		grandTotal += int32(totalPrice)

		productQty, err := strconv.Atoi(productPdf.ProductQty)

		if err != nil {
			log.Println(err)
		}

		totalQty += int32(productQty)

		productPdfs = append(productPdfs, &productPdf)
	}

	query2 := `SELECT
				u.user_name,
				u.user_phone,
				u.user_email,
			

				i.invoice_reverse_charge,
				i.invoice_number,
				i.invoice_date,
				i.invoice_state,
				i.invoice_state_code,
				i.invoice_challan_number,
				i.invoice_vehicle_number,
				i.invoice_date_of_supply,
				i.invoice_place_of_supply,
				i.invoice_gst,

				r.billed_name,
				r.billed_address,
				r.billed_gstin,
				r.billed_state,
				r.billed_state_code,

				c.shipped_name,
				c.shipped_address,
				c.shipped_gstin,
				c.shipped_mobile,
				c.shipped_state,
				c.shipped_state_code

			FROM invoice i
			JOIN users u ON i.user_id=u.user_id
		 	JOIN billed r ON i.billed_id=r.billed_id
			JOIN shipped c ON i.shipped_id=c.shipped_id

			WHERE i.invoice_id=$1
			`
	var invoicePdf models.InvoicePdf

	invoicePdf.TotalQty = strconv.Itoa(int(totalQty))
	invoicePdf.GrandTotal = strconv.Itoa(int(grandTotal))

	invoicePdf.Products = productPdfs
	var invoiceNumber int32

	err = q.db.QueryRow(query2, invoiceId).Scan(
		&invoicePdf.UserName,
		&invoicePdf.UserPhone,
		&invoicePdf.UserEmail,

		&invoicePdf.InvoiceReverseCharge,
		&invoiceNumber,
		&invoicePdf.InvoiceDate,
		&invoicePdf.InvoiceState,
		&invoicePdf.InvoiceStateCode,
		&invoicePdf.InvoiceChallanNumber,
		&invoicePdf.InvoiceVehicleNumber,
		&invoicePdf.InvoiceDateOfSupply,
		&invoicePdf.InvoicePlaceOfSupply,
		&invoicePdf.InvoiceGst,
		&invoicePdf.ReceiverName,
		&invoicePdf.ReceiverAdddress,
		&invoicePdf.ReceiverGstin,
		&invoicePdf.ReceiverState,
		&invoicePdf.ReceiverStateCode,
		&invoicePdf.ConsigneeName,
		&invoicePdf.ConsigneeAddress,
		&invoicePdf.ConsigneeGstin,
		&invoicePdf.ConsigneeMobile,
		&invoicePdf.ConsigneeState,
		&invoicePdf.ConsigneeStateCode,
	)

	if err != nil {
		return nil, err
	}

	invoicePdf.InvoiceNumber = strconv.Itoa(int(invoiceNumber))

	return &invoicePdf, nil

}

func (q *Query) UpdatePaymentStatus(invoiceId string) error {
	query := `UPDATE invoice SET invoice_payment_status = TRUE WHERE invoice_id=$1`
	_, err := q.db.Exec(query, invoiceId)
	return err
}
