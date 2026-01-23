# Psiphon Inproxy Node / نود پروکسی Psiphon Inproxy

> **Note:** This README was created by CURSOR AI

Docker container for running a Psiphon inproxy proxy node on Linux. This provides the same functionality as the Conduit mobile app, but runs as a headless service on a server.

کانتینر Docker برای اجرای نود پروکسی Psiphon Inproxy روی لینوکس. این همان عملکردی است که اپلیکیشن Conduit در موبایل ارائه می‌دهد، اما به صورت سرویس بدون رابط کاربری روی سرور اجرا می‌شود.

---

## What is an Inproxy Node? / نود Inproxy چیست؟

An inproxy node is a proxy that helps Psiphon clients connect to the Psiphon network. It acts as an intermediary, relaying traffic between clients and Psiphon servers without being able to see the encrypted tunnel contents. This is the same service that Conduit provides on mobile devices.

نود Inproxy یک پروکسی است که به کلاینت‌های Psiphon کمک می‌کند تا به شبکه Psiphon متصل شوند. این نود به عنوان واسطه عمل می‌کند و ترافیک بین کلاینت‌ها و سرورهای Psiphon را رله می‌کند بدون اینکه بتواند محتوای تونل رمزگذاری شده را ببیند. این همان سرویسی است که Conduit در دستگاه‌های موبایل ارائه می‌دهد.

---

## Requirements / نیازمندی‌ها

- Docker
- psiphon-tunnel-core repository (cloned locally)
- Internet connection

---

## Installation / نصب

### Prerequisites / پیش‌نیازها

Clone the psiphon-tunnel-core repository:

```bash
cd ..
git clone https://github.com/Psiphon-Labs/psiphon-tunnel-core.git
```

The directory structure should be:
```
cond/
├── psiphon-tunnel-core/
└── psiphon-inproxy-node/
```

ساختار دایرکتوری باید به این صورت باشد:
```
cond/
├── psiphon-tunnel-core/
└── psiphon-inproxy-node/
```

### Building / ساخت

#### Option 1: Using docker-compose

```bash
docker-compose build
```

#### Option 2: Using docker build

```bash
docker build -t psiphon-inproxy-node .
```

If the psiphon-tunnel-core path is different:

```bash
docker build --build-arg PSIPHON_REPO_PATH=/path/to/psiphon-tunnel-core -t psiphon-inproxy-node .
```

اگر مسیر psiphon-tunnel-core متفاوت است:

```bash
docker build --build-arg PSIPHON_REPO_PATH=/path/to/psiphon-tunnel-core -t psiphon-inproxy-node .
```

---

## Running / اجرا

### Basic Usage / استفاده پایه

```bash
docker run -d \
  --name psiphon-inproxy \
  --restart unless-stopped \
  psiphon-inproxy-node
```

### With docker-compose

```bash
docker-compose up -d
```

### With Custom Limits / با محدودیت‌های سفارشی

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

### With Configuration File / با فایل تنظیمات

```bash
docker run -d \
  --name psiphon-inproxy \
  --restart unless-stopped \
  -v /path/to/config.json:/config.json:ro \
  -v inproxy-data:/data \
  psiphon-inproxy-node \
  -config /config.json
```

---

## Parameters / پارامترها

- `-config`: Path to Psiphon configuration JSON file (optional) / مسیر فایل تنظیمات JSON (اختیاری)
- `-dataRootDirectory`: Directory for persistent data (default: `/data`) / دایرکتوری برای داده‌های پایدار (پیش‌فرض: `/data`)
- `-maxClients`: Maximum number of concurrent clients (default: 10) / حداکثر تعداد کلاینت‌های همزمان (پیش‌فرض: 10)
- `-limitUpstream`: Upstream bandwidth limit in bytes/sec (default: 0 = unlimited) / محدودیت پهنای باند آپلود به بایت بر ثانیه (پیش‌فرض: 0 = نامحدود)
- `-limitDownstream`: Downstream bandwidth limit in bytes/sec (default: 0 = unlimited) / محدودیت پهنای باند دانلود به بایت بر ثانیه (پیش‌فرض: 0 = نامحدود)
- `-version`: Print version information and exit / نمایش اطلاعات نسخه و خروج

---

## Network Requirements / نیازمندی‌های شبکه

The inproxy node uses WebRTC for client connections, which requires:
- UDP ports for WebRTC traffic (dynamically allocated)
- Outbound HTTPS connections to Psiphon brokers
- Outbound connections to Psiphon servers

The container should have network access, but you typically don't need to expose specific ports since WebRTC handles NAT traversal.

نود Inproxy از WebRTC برای اتصالات کلاینت استفاده می‌کند که نیاز به:
- پورت‌های UDP برای ترافیک WebRTC (تخصیص پویا)
- اتصالات HTTPS خروجی به بروکرهای Psiphon
- اتصالات خروجی به سرورهای Psiphon

کانتینر باید دسترسی به شبکه داشته باشد، اما معمولاً نیازی به باز کردن پورت‌های خاص نیست چون WebRTC از NAT traversal استفاده می‌کند.

---

## Data Persistence / ذخیره‌سازی داده

The container stores persistent data in `/data`. To persist this across container restarts, use a volume:

```bash
docker run -d \
  --name psiphon-inproxy \
  --restart unless-stopped \
  -v psiphon-data:/data \
  psiphon-inproxy-node
```

کانتینر داده‌های پایدار را در `/data` ذخیره می‌کند. برای حفظ این داده‌ها در راه‌اندازی مجدد کانتینر، از volume استفاده کنید.

---

## Monitoring / مانیتورینگ

### View Logs / مشاهده لاگ‌ها

```bash
docker logs psiphon-inproxy
```

### Follow Logs / دنبال کردن لاگ‌ها

```bash
docker logs -f psiphon-inproxy
```

### Check Status / بررسی وضعیت

```bash
docker ps --filter "name=psiphon-inproxy"
```

---

## Stopping / توقف

```bash
docker stop psiphon-inproxy
docker rm psiphon-inproxy
```

Or with docker-compose:

```bash
docker-compose down
```

---

## How It Works / نحوه کار

1. The inproxy node starts and connects to Psiphon brokers / نود Inproxy شروع می‌شود و به بروکرهای Psiphon متصل می‌شود
2. Psiphon clients (users in censored regions) contact the broker / کلاینت‌های Psiphon (کاربران در مناطق سانسور شده) با بروکر تماس می‌گیرند
3. The broker matches clients with available inproxy nodes / بروکر کلاینت‌ها را با نودهای Inproxy موجود تطبیق می‌دهد
4. Clients connect to your node via WebRTC / کلاینت‌ها از طریق WebRTC به نود شما متصل می‌شوند
5. Your node relays their traffic to Psiphon servers / نود شما ترافیک آن‌ها را به سرورهای Psiphon رله می‌کند
6. Clients can then access the internet through the Psiphon network / کلاینت‌ها می‌توانند از طریق شبکه Psiphon به اینترنت دسترسی پیدا کنند

---

## Configuration / تنظیمات

The node automatically receives broker configuration from Psiphon tactics. You can also provide a custom configuration file with broker specs if needed.

نود به صورت خودکار تنظیمات بروکر را از تاکتیک‌های Psiphon دریافت می‌کند. همچنین می‌توانید یک فایل تنظیمات سفارشی با مشخصات بروکر ارائه دهید.

---

## Notes / نکات

- The inproxy node requires the `PSIPHON_ENABLE_INPROXY` build tag (included in Dockerfile) / نود Inproxy نیاز به build tag `PSIPHON_ENABLE_INPROXY` دارد (در Dockerfile گنجانده شده)
- Configuration is typically provided via tactics from the Psiphon network / تنظیمات معمولاً از طریق تاکتیک‌های شبکه Psiphon ارائه می‌شود
- The node will automatically connect to Psiphon brokers and start accepting client connections / نود به صورت خودکار به بروکرهای Psiphon متصل می‌شود و شروع به پذیرش اتصالات کلاینت می‌کند
- Bandwidth limits are per-client, not total / محدودیت‌های پهنای باند برای هر کلاینت است، نه مجموع
- The container runs as user `inproxy` (UID 1000) for security / کانتینر به عنوان کاربر `inproxy` (UID 1000) برای امنیت اجرا می‌شود

---

## Troubleshooting / عیب‌یابی

If the container exits immediately, check the logs:

```bash
docker logs psiphon-inproxy
```

Common issues:
- Missing configuration: The node needs broker configuration from tactics / تنظیمات ناقص: نود نیاز به تنظیمات بروکر از تاکتیک‌ها دارد
- Network issues: Ensure the container has internet access / مشکلات شبکه: مطمئن شوید کانتینر دسترسی به اینترنت دارد
- Permission issues: The container runs as user `inproxy` (UID 1000) / مشکلات دسترسی: کانتینر به عنوان کاربر `inproxy` (UID 1000) اجرا می‌شود

اگر کانتینر بلافاصله متوقف می‌شود، لاگ‌ها را بررسی کنید.

---

## License / لایسنس

This project uses the psiphon-tunnel-core library, which is licensed under GPL-3.0.

این پروژه از کتابخانه psiphon-tunnel-core استفاده می‌کند که تحت لایسنس GPL-3.0 است.
