package utils

import (
    "fmt"
    "paystack-payment/models"
)

func FormatBillNotification(bill models.Bill, paymentLink string) string {
    return fmt.Sprintf(
        "New bill notification:\nAmount: %.2f\nDue Date: %s\nDescription: %s\nPay now: %s",
        bill.Amount,
        bill.DueDate.Format("2006-01-02"),
        bill.Description,
        paymentLink,
    )
}