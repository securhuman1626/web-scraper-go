package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

const (
	outputHTMLFile = "site_data.html"
	outputTextFile = "output.txt"
	screenshotFile = "screenshot.png"
	urlsFile       = "urls.txt"
)

func main() {
	// Komut satırı argümanını al
	targetURL := flag.String("url", "", "Hedef web sitesi URL'si (zorunlu)")
	flag.Parse()

	if *targetURL == "" {
		fmt.Println("Hata: URL belirtilmedi!")
		fmt.Println("Kullanım: go run main.go -url <URL>")
		os.Exit(1)
	}

	// URL geçerlilik kontrolü
	parsedURL, err := url.Parse(*targetURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		log.Fatalf("Hata: Geçersiz URL formatı: %s\n", *targetURL)
	}

	// 1. URL'ye özel klasör ismini oluştur ve klasörü aç
	folderName := sanitizeFolderName(*targetURL)
	if err := os.MkdirAll(folderName, 0755); err != nil {
		log.Fatalf("Klasör oluşturma hatası: %v\n", err)
	}

	fmt.Printf("Web scraper başlatılıyor...\n")
	fmt.Printf("Hedef URL: %s\n", *targetURL)
	fmt.Printf("Çıktı Klasörü: %s\n\n", folderName)

	// HTML içeriğini çek
	htmlContent, err := fetchHTML(*targetURL)
	if err != nil {
		log.Fatalf("HTML çekme hatası: %v\n", err)
	}

	// 2. HTML içeriğini klasör içine kaydet
	htmlPath := filepath.Join(folderName, outputHTMLFile)
	if err := saveHTML(htmlContent, htmlPath); err != nil {
		log.Fatalf("HTML kaydetme hatası: %v\n", err)
	}
	fmt.Printf("✓ HTML içeriği '%s' içine kaydedildi\n", htmlPath)

	// 3. Metin içeriğini klasör içine kaydet
	textPath := filepath.Join(folderName, outputTextFile)
	if err := saveText(htmlContent, textPath); err != nil {
		log.Fatalf("Metin kaydetme hatası: %v\n", err)
	}
	fmt.Printf("✓ Metin içeriği '%s' içine kaydedildi\n", textPath)

	// 4. Ekran görüntüsü al (1920x1080) ve klasör içine kaydet
	ssPath := filepath.Join(folderName, screenshotFile)
	if err := takeScreenshot(*targetURL, ssPath); err != nil {
		log.Fatalf("Ekran görüntüsü alma hatası: %v\n", err)
	}
	fmt.Printf("✓ Ekran görüntüsü '%s' içine kaydedildi\n", ssPath)

	// 5. URL'leri çıkar ve klasör içine kaydet (Bonus görev)
	urls, err := extractURLs(htmlContent, *targetURL)
	if err == nil {
		urlsPath := filepath.Join(folderName, urlsFile)
		if err := saveURLs(urls, urlsPath); err == nil {
			fmt.Printf("✓ %d adet URL '%s' içine kaydedildi\n", len(urls), urlsPath)
		}
	}

	fmt.Println("\n✓ Tüm işlemler başarıyla tamamlandı!")
}

func sanitizeFolderName(targetURL string) string {
	u, err := url.Parse(targetURL)
	if err != nil {
		return "downloaded_site"
	}
	name := u.Host + u.Path
	name = strings.ReplaceAll(name, ".", "_")
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, ":", "_")
	return strings.Trim(name, "_")
}

func fetchHTML(targetURL string) (string, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d hatası", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	return string(body), err
}

func saveHTML(content, filename string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

func saveText(htmlContent, filename string) error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return os.WriteFile(filename, []byte(htmlContent), 0644)
	}
	var textContent strings.Builder
	doc.Find("body").Each(func(i int, s *goquery.Selection) {
		textContent.WriteString(s.Text() + "\n")
	})
	return os.WriteFile(filename, []byte(textContent.String()), 0644)
}

func takeScreenshot(targetURL, filename string) error {
	// Düzeltilen kısım: DefaultExecAllocatorOptions kullanıldı
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.WindowSize(1920, 1080),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var buf []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate(targetURL),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.Sleep(20*time.Second),
		chromedp.CaptureScreenshot(&buf),
	)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, buf, 0644)
}

func extractURLs(htmlContent, baseURL string) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}
	base, _ := url.Parse(baseURL)
	var urls []string
	urlMap := make(map[string]bool)
	doc.Find("a[href], img[src]").Each(func(i int, s *goquery.Selection) {
		val, _ := s.Attr("href")
		if val == "" {
			val, _ = s.Attr("src")
		}
		parsed, err := url.Parse(val)
		if err == nil {
			abs := base.ResolveReference(parsed).String()
			if (strings.HasPrefix(abs, "http")) && !urlMap[abs] {
				urls = append(urls, abs)
				urlMap[abs] = true
			}
		}
	})
	return urls, nil
}

func saveURLs(urls []string, filename string) error {
	return os.WriteFile(filename, []byte(strings.Join(urls, "\n")), 0644)
}
