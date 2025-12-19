# 🛡️ Gelişmiş Web Scraper & CTI Tool (Go)

Bu proje, Siber Tehdit İstihbaratı (CTI) toplama süreçlerinin temel bir adımı olarak, hedef web sitelerinden veri çekmek ve görsel kanıt (ekran görüntüsü) oluşturmak amacıyla Go (Golang) diliyle geliştirilmiştir. Program, girilen her URL için izole bir çalışma alanı oluşturarak verilerin düzenli bir şekilde saklanmasını sağlar.

## 🚀 Öne Çıkan Özellikler

* **Otomatik İzole Klasörleme:** Her tarama işlemi, hedef URL'den türetilen benzersiz bir klasör içinde saklanır (Örn: `www_haberturk_com`).
* **Tam HTML Çekimi:** Sayfanın ham HTML içeriği `site_data.html` olarak kaydedilir.
* **Temiz Metin Analizi:** HTML etiketlerinden arındırılmış saf metin içeriği `output.txt` dosyasına aktarılır.
* **Headless Ekran Görüntüsü:** Chrome/Chromium altyapısı kullanılarak sitenin anlık görüntüsü `screenshot.png` olarak alınır.
* **Bonus - Akıllı URL Ayıklama:** Sayfa içindeki tüm `<a>` linkleri ve `<img>` kaynakları otomatik olarak tespit edilip `urls.txt` dosyasına listelenir.
* **Gelişmiş Hata Yönetimi:** HTTP 404/500 hataları, bağlantı zaman aşımları ve geçersiz URL formatları kullanıcıya detaylıca raporlanır.

## 🛠️ Gereksinimler

* **Go 1.21** veya üzeri sürüm.
* **Chrome/Chromium** tarayıcısı (Arka planda ekran görüntüsü motoru olarak kullanılır).

## 📦 Kurulum ve Hazırlık

1. Proje klasörüne gidin:
```bash
cd webscraper

```


2. Gerekli kütüphaneleri ve bağımlılıkları yükleyin:
```bash
go mod tidy

```



## 📖 Kullanım ve Örnekler

Programı çalıştırmak için `-url` parametresi ile hedef adresi belirtmeniz yeterlidir:

```bash
go run main.go -url https://www.haberturk.com

```

### Örnek Senaryolar

```bash
# Teknoloji haberleri analizi
go run main.go -url https://news.ycombinator.com

# Kurumsal site incelemesi
go run main.go -url https://example.com

```

## 📂 Çıktı Yapısı ve Organizasyon

Program her tarama için şu hiyerarşiyi oluşturur:

```text
/webscraper (Proje Ana Dizini)
  ├── www_haberturk_com/         <-- Otomatik oluşturulan klasör
  │   ├── screenshot.png         <-- Sitenin ekran görüntüsü
  │   ├── site_data.html         <-- Ham HTML kodu
  │   ├── output.txt             <-- Ayıklanmış metinler
  │   └── urls.txt               <-- Sayfadaki tüm linkler
  └── google_com/                <-- Başka bir tarama sonucu

```

## 🔬 Teknik Detaylar

* **chromedp:** Tarayıcıyı "headless" modda kontrol ederek JavaScript yüklü dinamik sayfaların tam görüntüsünü alır.
* **goquery:** HTML doküman ağacını (DOM) analiz ederek metinleri ve linkleri hızlıca ayıklar.
* **net/http:** Düşük seviyeli HTTP istekleri oluşturarak ham veriyi en güvenilir yoldan çeker.
* **path/filepath:** Dosya ve klasör yollarını işletim sisteminden bağımsız (Windows/Linux) yönetir.

## ⚠️ Önemli Notlar

* **Bot Koruması:** Bazı siteler yoğun istekleri engelleyebilir; program bu durumu aşmak için gerçekçi bir **User-Agent** başlığı kullanır.
* **Dinamik İçerik:** JavaScript ile sonradan yüklenen içeriklerin görüntülenebilmesi için tarayıcı motoruna 3 saniyelik bekleme süresi eklenmiştir.
* **Üst Üste Yazma:** Her URL kendi klasöründe saklandığı için farklı sitelerin verileri asla birbirine karışmaz.

---