package event

import (
	"fmt"
	"log/slog"
	"net/smtp"
	"os"

	"product-api/internal/model"
)

func notifyAdminNewProduct(product model.Product) error {
	from := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASSWORD")
	to := os.Getenv("ADMIN_EMAIL")

	if from == "" || pass == "" || to == "" {
		return fmt.Errorf("konfigurasi SMTP tidak lengkap (cek env var SMTP_USER, SMTP_PASSWORD, ADMIN_EMAIL)")
	}

	subject := "Produk Baru Ditambahkan"
	body := fmt.Sprintf("Produk baru telah ditambahkan ke sistem:\n\nNama: %s\nHarga: Rp%.0f\nStok: %d\nID: %d",
		product.Name, 
		product.Price, 
		product.Stock, 
		product.ID)

	msg := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)

	auth := smtp.PlainAuth("", from, pass, "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, []byte(msg))
	if err != nil {
		return fmt.Errorf("gagal kirim email: %w", err)
	}

	slog.Info("email notifikasi terkirim", "to", to, "product_id", product.ID)
	return nil
}