package handler

import (
	"fmt"
	"go-generate-sk/model"
	"go-generate-sk/service"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GenerateSK(c *fiber.Ctx) error {
	var user model.User
	// if err := c.BodyParser(&user); err != nil {
	// 	fmt.Println("Error parsing body:", err.Error())
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	// }

	// Generate SK logic here
	user = model.User{
		ID:           1,
		Name:         "Tajri Mintahtihal Anhaar, S.Kom",
		BobotPertama: 30,
		BobotKedua:   35,
		BobotKetiga:  35,
		LabelPertama: "A",
		LabelKedua:   "B",
		LabelKetiga:  "C",
		Nilai:        85,
	}

	pdfBytes, err := service.GeneratePDF(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to generate PDF: %v", err),
		})
	}

	// --- MULAI Fungsionalitas Menyimpan File ---

	// 4. Tentukan nama file dan direktori output
	outputDir := "generated_pdfs"
	filename := fmt.Sprintf("SK_%s_%s.pdf", user.Name, time.Now().Format("20060102_150405"))
	filePath := filepath.Join(outputDir, filename)

	// 5. Buat direktori jika belum ada
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Printf("Gagal membuat direktori: %v", err)
		// Anda bisa memilih untuk tetap lanjut atau mengembalikan error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to create directory: %v", err),
		})
	}

	// 6. Tulis byte PDF ke file
	if err := os.WriteFile(filePath, pdfBytes, 0644); err != nil {
		log.Printf("Gagal menyimpan file PDF: %v", err)
		// Anda bisa memilih untuk tetap lanjut atau mengembalikan error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to save PDF file: %v", err),
		})
	}

	log.Printf("File PDF berhasil disimpan di: %s", filePath)

	// --- SELESAI Fungsionalitas Menyimpan File ---

	// Set headers for PDF download
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Set("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))

	return c.Send(pdfBytes)
}
