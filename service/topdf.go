package service

import (
	"bytes"
	"context"
	"fmt"
	"go-generate-sk/model"
	"html/template"
	"log"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type PDFData struct {
	model.User
	ID             uint
	Year           int
	Location       string
	Date           string
	SignatureName  string
	SignatureTitle string
	GeneratedAt    string
}

func GeneratePDF(user model.User) ([]byte, error) {
	// Prepare data for template
	data := PDFData{
		User:           user,
		ID:             user.ID,
		Year:           time.Now().Year(),
		Location:       "Jakarta",
		Date:           time.Now().Format("02 January 2006"),
		SignatureName:  "Manager",
		SignatureTitle: "General Manager",
		GeneratedAt:    time.Now().Format("02/01/2006 15:04:05"),
	}

	// Generate HTML from template
	htmlContent, err := generateHTML(data)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HTML: %v", err)
	}

	// Generate PDF from HTML using chromedp
	pdfBytes, err := htmlToPDFChromedp(htmlContent)
	if err != nil {
		log.Printf("Chromedp failed, trying WKHTML: %v", err)
		// Fallback to WKHTML if chromedp fails
		return GeneratePDFWithWKHTML(user)
	}

	return pdfBytes, nil
}

func generateHTML(data PDFData) ([]byte, error) {
	templateContent := `
<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <title>Surat Keterangan</title>
    <style>
        @page {
            /* Atur margin standar untuk surat resmi */
            margin: 2.5cm; 
        }
        body {
            font-family: 'Times New Roman', Times, serif;
            font-size: 12pt;
            line-height: 1.5;
        }
        .kop-surat {
            text-align: center;
            line-height: 1.2;
            border-bottom: 4px double #000;
            padding-bottom: 10px;
        }
        .kop-surat .logo {
            /* Posisikan logo di kiri atas kop surat */
            position: absolute;
            top: 1.5cm;
            left: 2.5cm;
            width: 70px; /* Sesuaikan ukuran logo */
        }
        .kop-surat .nama-instansi {
            font-size: 16pt;
            font-weight: bold;
        }
        .kop-surat .alamat-instansi {
            font-size: 11pt;
        }
        .judul-surat {
            text-align: center;
            font-weight: bold;
            text-decoration: underline;
            font-size: 14pt;
            margin-top: 30px;
            margin-bottom: 5px;
        }
        .nomor-surat {
            text-align: center;
            margin-bottom: 25px;
        }
        .paragraf {
            text-align: justify;
            margin-bottom: 15px;
        }
        .tabel-data {
            /* Tabel untuk merapikan data tanpa terlihat ada border */
            border-collapse: collapse;
            width: 100%;
            margin-left: 20px; /* Sedikit menjorok ke dalam */
        }
        .tabel-data td {
            padding: 2px 0;
            vertical-align: top;
        }
        .tabel-data .label {
            width: 25%;
        }
        .tabel-data .separator {
            width: 2%;
        }
        .blok-tanda-tangan {
            margin-top: 40px;
            /* Posisi tanda tangan di kanan bawah */
            width: 40%; 
            margin-left: 60%;
            text-align: left;
        }
        .blok-tanda-tangan .nama-penandatangan {
            font-weight: bold;
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <div class="kop-surat">
        <img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII=" class="logo" alt="Logo">
        <div class="nama-instansi">NAMA INSTANSI ANDA</div>
        <div class="alamat-instansi">
            Jalan Alamat Instansi No. 123, Kota, Kode Pos<br>
            Telepon: (021) 1234567, Website: www.websiteinstansi.com
        </div>
    </div>

    <div class="judul-surat">SURAT KETERANGAN</div>
    <div class="nomor-surat">NOMOR: {{.ID}}/SK/{{.Year}}</div>

    <p class="paragraf">Yang bertanda tangan di bawah ini menerangkan bahwa:</p>
    
    <table class="tabel-data">
        <tr>
            <td class="label">Nama</td>
            <td class="separator">:</td>
            <td><strong>{{.Name}}</strong></td>
        </tr>
        <tr>
            <td class="label">Nilai Akhir</td>
            <td class="separator">:</td>
            <td><strong>{{.Nilai}}</strong></td>
        </tr>
    </table>

    <p class="paragraf">Dengan rincian penilaian sebagai berikut:</p>
    <table class="tabel-data">
        <tr>
            <td class="label">{{.LabelPertama}}</td>
            <td class="separator">:</td>
            <td>Bobot {{.BobotPertama}}%</td>
        </tr>
        <tr>
            <td class="label">{{.LabelKedua}}</td>
            <td class="separator">:</td>
            <td>Bobot {{.BobotKedua}}%</td>
        </tr>
        <tr>
            <td class="label">{{.LabelKetiga}}</td>
            <td class="separator">:</td>
            <td>Bobot {{.BobotKetiga}}%</td>
        </tr>
    </table>

    <p class="paragraf">Demikian surat keterangan ini dibuat untuk dapat dipergunakan sebagaimana mestinya.</p>
    
    <div class="blok-tanda-tangan">
        <p>{{.Location}}, {{.Date}}</p>
        <p>{{.SignatureTitle}}</p>
        <br><br><br><br>
        <p class="nama-penandatangan">{{.SignatureName}}</p>
    </div>
</body>
</html>
`

	tmpl, err := template.New("sk").Parse(templateContent)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Implementasi Chromedp yang lebih andal
func htmlToPDFChromedp(htmlContent []byte) ([]byte, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var pdfBuffer []byte
	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Menggunakan SetDocumentContent lebih stabil daripada data URL
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			if err := page.SetDocumentContent(frameTree.Frame.ID, string(htmlContent)).Do(ctx); err != nil {
				return err
			}

			// Beri waktu sejenak untuk rendering
			time.Sleep(500 * time.Millisecond)

			// Print to PDF
			pdfBuffer, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				Do(ctx)
			return err
		}),
	)
	if err != nil {
		return nil, err
	}
	return pdfBuffer, nil
}

// Simplified WKHTML implementation
func GeneratePDFWithWKHTML(user model.User) ([]byte, error) {
	// Generate HTML
	data := PDFData{
		User:           user,
		ID:             user.ID,
		Year:           time.Now().Year(),
		Location:       "Jakarta",
		Date:           time.Now().Format("02 January 2006"),
		SignatureName:  "Manager",
		SignatureTitle: "General Manager",
		GeneratedAt:    time.Now().Format("02/01/2006 15:04:05"),
	}

	htmlContent, err := generateHTML(data)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HTML: %v", err)
	}

	// Generate PDF using WKHTML
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, fmt.Errorf("failed to create PDF generator: %v", err)
	}

	// Create page from HTML string
	page := wkhtmltopdf.NewPageReader(bytes.NewReader([]byte(htmlContent)))

	// Set page options
	page.FooterRight.Set("[page]")
	page.FooterFontSize.Set(10)
	page.Zoom.Set(0.95)

	// Set PDF options
	pdfg.AddPage(page)
	pdfg.MarginTop.Set(10)
	pdfg.MarginBottom.Set(10)
	pdfg.MarginLeft.Set(10)
	pdfg.MarginRight.Set(10)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)

	// Create PDF
	err = pdfg.Create()
	if err != nil {
		return nil, fmt.Errorf("failed to create PDF: %v", err)
	}

	return pdfg.Bytes(), nil
}
