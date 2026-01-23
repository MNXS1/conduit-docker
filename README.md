# نود پروکسی Psiphon Inproxy

کانتینر Docker برای اجرای نود پروکسی Psiphon Inproxy روی لینوکس. این همان عملکردی است که اپلیکیشن Conduit در موبایل ارائه می‌دهد، اما به صورت سرویس بدون رابط کاربری روی سرور اجرا می‌شود.

## نود Inproxy چیست؟

نود Inproxy یک پروکسی است که به کلاینت‌های Psiphon کمک می‌کند تا به شبکه Psiphon متصل شوند. این نود به عنوان واسطه عمل می‌کند و ترافیک بین کلاینت‌ها و سرورهای Psiphon را رله می‌کند بدون اینکه بتواند محتوای تونل رمزگذاری شده را ببیند. این همان سرویسی است که Conduit در دستگاه‌های موبایل ارائه می‌دهد.

## نیازمندی‌ها

- Docker
- مخزن psiphon-tunnel-core (کلون شده به صورت محلی)
- اتصال به اینترنت

## نصب

### پیش‌نیازها

مخزن psiphon-tunnel-core را کلون کنید:

```bash
cd ..
git clone https://github.com/Psiphon-Labs/psiphon-tunnel-core.git
```

ساختار دایرکتوری باید به این صورت باشد:
```
cond/
├── psiphon-tunnel-core/
└── psiphon-inproxy-node/
```

### ساخت با docker-compose

```bash
docker-compose build
```

### ساخت با docker build

```bash
docker build -t psiphon-inproxy-node .
```

اگر مسیر psiphon-tunnel-core متفاوت است:

```bash
docker build --build-arg PSIPHON_REPO_PATH=/path/to/psiphon-tunnel-core -t psiphon-inproxy-node .
```

## اجرا

### استفاده پایه

```bash
docker run -d \
  --name psiphon-inproxy \
  --restart unless-stopped \
  psiphon-inproxy-node
```

### با docker-compose

```bash
docker-compose up -d
```

### با محدودیت‌های سفارشی

```bash
docker run -d \
  --name psiphon-inproxy \
  --restart unless-stopped \
  -v inproxy-data:/data \
  psiphon-inproxy-node \
  -maxClients 20 \
  -limitUpstream 10485760 \
  -limitDownstream 10485760
```

### با فایل تنظیمات

```bash
docker run -d \
  --name psiphon-inproxy \
  --restart unless-stopped \
  -v /path/to/config.json:/config.json:ro \
  -v inproxy-data:/data \
  psiphon-inproxy-node \
  -config /config.json
```

## پارامترها

- `-config`: مسیر فایل تنظیمات JSON (اختیاری)
- `-dataRootDirectory`: دایرکتوری برای داده‌های پایدار (پیش‌فرض: `/data`)
- `-maxClients`: حداکثر تعداد کلاینت‌های همزمان (پیش‌فرض: 10)
- `-limitUpstream`: محدودیت پهنای باند آپلود به بایت بر ثانیه (پیش‌فرض: 0 = نامحدود)
- `-limitDownstream`: محدودیت پهنای باند دانلود به بایت بر ثانیه (پیش‌فرض: 0 = نامحدود)
- `-version`: نمایش اطلاعات نسخه و خروج

## نیازمندی‌های شبکه

نود Inproxy از WebRTC برای اتصالات کلاینت استفاده می‌کند که نیاز به:
- پورت‌های UDP برای ترافیک WebRTC (تخصیص پویا)
- اتصالات HTTPS خروجی به بروکرهای Psiphon
- اتصالات خروجی به سرورهای Psiphon

کانتینر باید دسترسی به شبکه داشته باشد، اما معمولاً نیازی به باز کردن پورت‌های خاص نیست چون WebRTC از NAT traversal استفاده می‌کند.

## ذخیره‌سازی داده

کانتینر داده‌های پایدار را در `/data` ذخیره می‌کند. برای حفظ این داده‌ها در راه‌اندازی مجدد کانتینر، از volume استفاده کنید:

```bash
docker run -d \
  --name psiphon-inproxy \
  --restart unless-stopped \
  -v psiphon-data:/data \
  psiphon-inproxy-node
```

## مانیتورینگ

### مشاهده لاگ‌ها

```bash
docker logs psiphon-inproxy
```

### دنبال کردن لاگ‌ها

```bash
docker logs -f psiphon-inproxy
```

### بررسی وضعیت

```bash
docker ps --filter "name=psiphon-inproxy"
```

## توقف

```bash
docker stop psiphon-inproxy
docker rm psiphon-inproxy
```

یا با docker-compose:

```bash
docker-compose down
```

## نحوه کار

1. نود Inproxy شروع می‌شود و به بروکرهای Psiphon متصل می‌شود
2. کلاینت‌های Psiphon (کاربران در مناطق سانسور شده) با بروکر تماس می‌گیرند
3. بروکر کلاینت‌ها را با نودهای Inproxy موجود تطبیق می‌دهد
4. کلاینت‌ها از طریق WebRTC به نود شما متصل می‌شوند
5. نود شما ترافیک آن‌ها را به سرورهای Psiphon رله می‌کند
6. کلاینت‌ها می‌توانند از طریق شبکه Psiphon به اینترنت دسترسی پیدا کنند

## تنظیمات

نود به صورت خودکار تنظیمات بروکر را از تاکتیک‌های Psiphon دریافت می‌کند. همچنین می‌توانید یک فایل تنظیمات سفارشی با مشخصات بروکر ارائه دهید.

## نکات

- نود Inproxy نیاز به build tag `PSIPHON_ENABLE_INPROXY` دارد (در Dockerfile گنجانده شده)
- تنظیمات معمولاً از طریق تاکتیک‌های شبکه Psiphon ارائه می‌شود
- نود به صورت خودکار به بروکرهای Psiphon متصل می‌شود و شروع به پذیرش اتصالات کلاینت می‌کند
- محدودیت‌های پهنای باند برای هر کلاینت است، نه مجموع
- کانتینر به عنوان کاربر `inproxy` (UID 1000) برای امنیت اجرا می‌شود

## عیب‌یابی

اگر کانتینر بلافاصله متوقف می‌شود، لاگ‌ها را بررسی کنید:

```bash
docker logs psiphon-inproxy
```

مشکلات رایج:
- تنظیمات ناقص: نود نیاز به تنظیمات بروکر از تاکتیک‌ها دارد
- مشکلات شبکه: مطمئن شوید کانتینر دسترسی به اینترنت دارد
- مشکلات دسترسی: کانتینر به عنوان کاربر `inproxy` (UID 1000) اجرا می‌شود

## لایسنس

این پروژه از کتابخانه psiphon-tunnel-core استفاده می‌کند که تحت لایسنس GPL-3.0 است.
