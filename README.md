
---

###  Proje Mimari Diyagramı (Mermaid)


```mermaid
graph TD
    subgraph Clients
        User[User / Client]
        Postman[Postman / cURL]
    end

    subgraph "API Gateway (Port: 3000)"
        Gateway[Go Fiber Gateway]
    end

    subgraph "Microservices Cluster"
        Identity[Identity Service :8080]
        Event[Event Service :8082]
        Booking[Booking Service :50051]
        Payment[Payment Service Worker]
    end

    subgraph "Infrastructure & Data"
        PG_ID[(Postgres Auth DB)]
        PG_BOOK[(Postgres Order DB)]
        MONGO[(MongoDB Event DB)]
        REDIS[(Redis Cache & Lock)]
        RABBIT[RabbitMQ Message Broker]
    end

    %% Flow Connections
    User -->|HTTP JSON| Gateway
    Postman -->|HTTP JSON| Gateway

    %% Gateway Routing
    Gateway -->|HTTP Proxy| Identity
    Gateway -->|HTTP Proxy| Event
    Gateway -->|gRPC Protobuf| Booking

    %% Database Connections
    Identity -->|GORM| PG_ID
    Event -->|Mongo Driver| MONGO
    Booking -->|GORM| PG_BOOK
    
    %% Redis Locking Flow
    Booking -->|Acquire Lock| REDIS
    
    %% Async Messaging Flow
    Booking -.->|Publish Event| RABBIT
    RABBIT -.->|Consume Event| Payment

    %% Styling
    style Gateway fill:#f9f,stroke:#333,stroke-width:2px
    style Booking fill:#bbf,stroke:#333,stroke-width:2px
    style REDIS fill:#d65,stroke:#333,stroke-width:2px
    style RABBIT fill:#d65,stroke:#333,stroke-width:2px

```

---




#  GoTicketra - Scalable Microservices Event Ticketing Platform

![Go Version](https://img.shields.io/badge/Go-1.24-blue)
![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED)
![Architecture](https://img.shields.io/badge/Architecture-Microservices-orange)

**GoTicketra**, yüksek trafikli bilet satış senaryolarını simüle etmek için tasarlanmış, modern ve dağıtık bir **Biletleme Platformudur**. Bu proje, gerçek dünya senaryolarında karşılaşılan *Race Condition (Yarış Durumu)*, *Veri Tutarlılığı* ve *Asenkron İletişim* problemlerine endüstri standartlarında çözümler sunar.

---

##  Mimari & Teknoloji Yığını

Proje, **Clean Architecture** prensiplerine sadık kalınarak geliştirilmiş olup, **Monorepo** yapısında yönetilen 4 farklı mikroservis ve 1 API Gateway'den oluşmaktadır.

###  Kullanılan Teknolojiler

* **Dil:** Golang (Go 1.24)
* **Web Framework:** Fiber (v2)
* **İletişim Protokolleri:**
    * **REST (HTTP):** Frontend <-> Gateway iletişimi.
    * **gRPC (Protobuf):** Gateway <-> Booking Servisi (Yüksek performanslı dahili iletişim).
    * **AMQP (RabbitMQ):** Servisler arası asenkron olay yönetimi (Event-Driven).
* **Veritabanları & Depolama:**
    * **PostgreSQL:** Kullanıcı ve Sipariş verileri (İlişkisel Veri).
    * **MongoDB:** Etkinlik detayları (Esnek/Doküman bazlı Veri).
    * **Redis:** Distributed Locking (Stok yönetimi ve Race Condition önleme) ve Caching.
* **DevOps:** Docker, Docker Compose, Multi-Stage Builds (Alpine Linux).

---

##  Servisler ve Görevleri

| Servis Adı | Port | Teknoloji | Açıklama |
| :--- | :--- | :--- | :--- |
| **API Gateway** | `:3000` | Fiber | Tek giriş noktası. İstekleri yönlendirir (Reverse Proxy) ve JSON'u gRPC'ye çevirir. |
| **Identity Service** | `:8080` | Postgres, JWT | Kullanıcı kaydı, giriş işlemleri ve JWT Token üretimi. |
| **Event Service** | `:8082` | MongoDB | Etkinlik oluşturma ve listeleme işlemleri (NoSQL). |
| **Booking Service** | `:50051` | Postgres, Redis, gRPC | **Projenin kalbi.** Bilet satın alma, Redis ile stok kilitleme (Concurrency Control). |
| **Payment Service** | `Worker` | RabbitMQ | Arka planda çalışan tüketici (Consumer). Ödemeleri asenkron işler. |

---

##  Kurulum ve Çalıştırma

Proje tamamen **Dockerize** edilmiştir. Tek bir komutla tüm altyapıyı (DB'ler dahil) ve servisleri ayağa kaldırabilirsiniz.

### Gereksinimler
* Docker & Docker Compose

### Adım Adım Çalıştırma

1.  Repoyu klonlayın:
    ```bash
    git clone [https://github.com/keremkartal/goticketra.git](https://github.com/keremkartal/goticketra.git)
    cd goticketra
    ```

2.  Sistemi başlatın:
    ```bash
    docker-compose up --build -d
    ```
    *(İlk çalıştırmada Go modüllerinin indirilmesi ve imajların oluşturulması birkaç dakika sürebilir.)*

3.  Logları izleyin (Opsiyonel):
    ```bash
    docker-compose logs -f
    ```

---

##  API Kullanım Örnekleri (Testing)

Sistem ayağa kalktığında **http://localhost:3000** üzerinden tüm servislere erişebilirsiniz.

### 1. Kullanıcı Kaydı (Identity Service)
```bash
curl -X POST http://localhost:3000/api/auth/register \
-H "Content-Type: application/json" \
-d '{"email": "kerem@test.com", "password": "123"}'

```

### 2. Etkinlik Oluşturma (Event Service)

```bash
curl -X POST http://localhost:3000/api/events \
-H "Content-Type: application/json" \
-d '{
    "title": "Tarkan Konseri",
    "location": "Harbiye",
    "date": "2026-07-20T21:00:00Z",
    "total_tickets": 100
}'

```

*(Dönen cevaptaki `id` değerini kopyalayın!)*

### 3. Bilet Satın Alma (Booking Service -> gRPC -> RabbitMQ)

Burada **Redis Distributed Lock** devreye girer. Aynı anda aynı etkinliğe istek gelse bile tutarlılık sağlanır.

```bash
curl -X POST http://localhost:3000/api/bookings \
-H "Content-Type: application/json" \
-d '{
    "user_id": "kerem_user",
    "event_id": "BURAYA_EVENT_ID_YAPISTIR",
    "ticket_count": 2
}'

```

---
---

## Detaylı Servis İş Akışları

Aşağıda her bir mikroservisin kendi içindeki çalışma mantığı ve veri akışı detaylandırılmıştır.

### 1. API Gateway (The Entry Point)
Gateway, gelen HTTP isteklerini yönlendirir. Booking servisi için özel bir çevirici (JSON -> Protobuf) görevi görür.

```mermaid
graph LR
    Client[Client / Frontend]
    
    subgraph "API Gateway Logic"
        Router{Route Matcher}
        Proxy[Reverse Proxy]
        Adapter[gRPC Adapter]
    end
    
    Client -->|HTTP Request| Router
    
    Router -- "/api/auth/*" --> Proxy
    Proxy -->|Forward HTTP| Identity[Identity Service]
    
    Router -- "/api/events/*" --> Proxy
    Proxy -->|Forward HTTP| Event[Event Service]
    
    Router -- "/api/bookings" --> Adapter
    Adapter -->|Convert to Protobuf| gRPC[Booking Service]

```

### 2. Identity Service (Authentication)

Kullanıcı doğrulama ve JWT yönetimi akışı.

```mermaid
sequenceDiagram
    participant Client
    participant Handler
    participant Service
    participant Repo
    participant DB as PostgreSQL

    Client->>Handler: Login Request (Email, Pass)
    Handler->>Service: Login(email, pass)
    Service->>Repo: FindByEmail(email)
    Repo->>DB: SQL Query
    DB-->>Repo: User Record (Hash)
    Repo-->>Service: User Struct
    
    Service->>Service: bcrypt.Compare(Hash, Pass)
    
    alt Password Match
        Service->>Service: Generate JWT (Claims: ID, Role)
        Service-->>Handler: Return Token
        Handler-->>Client: 200 OK {token: ...}
    else Mismatch
        Service-->>Handler: Error "Invalid Credentials"
        Handler-->>Client: 401 Unauthorized
    end

```

### 3. Booking Service (Concurrency Control)

**Projenin en kritik akışı.** Redis kilitleme mekanizması sayesinde Race Condition engellenir.

```mermaid
sequenceDiagram
    participant Gateway
    participant Service
    participant Redis as Redis (Lock)
    participant PG as PostgreSQL
    participant MQ as RabbitMQ

    Gateway->>Service: gRPC CreateBooking(User, Event, Count)
    
    note right of Service: 1. Distributed Lock
    Service->>Redis: SETNX lock:event:{id} (TTL=5s)
    
    alt Lock Acquired (True)
        Service->>PG: INSERT INTO bookings (Status: PENDING)
        PG-->>Service: Success
        
        note right of Service: 2. Async Event
        Service->>MQ: Publish "BookingCreatedEvent"
        
        note right of Service: 3. Unlock
        Service->>Redis: DEL lock:event:{id}
        Service-->>Gateway: gRPC Response (Success)
    else Lock Failed (False)
        Service-->>Gateway: Error "Event is busy, try again"
    end

```

### 4. Payment Service (Event Driven Consumer)

Arka planda çalışan, mesaj kuyruğunu dinleyen ve telafi işlemini (Compensating Transaction) yöneten işçi servis.

```mermaid
stateDiagram-v2
    [*] --> Listening
    
    state "RabbitMQ Queue" as Q
    Listening --> Q: Consume Message
    
    state "Process Payment" as P {
        [*] --> MockBankCheck
        MockBankCheck --> Decision
        
        state Decision <<choice>>
        Decision --> Success: TicketCount <= 5
        Decision --> Failure: TicketCount > 5
    }
    
    Q --> P: Event Received
    
    Success --> LogSuccess: Payment Confirmed
    Failure --> Compensate: Insufficient Funds / Fraud
    
    state "Compensating Transaction" as Compensate {
        [*] --> LogCancel
        LogCancel --> RequestCancel
        RequestCancel --> [*]: TODO - Call Booking Service to Cancel
    }
    
    LogSuccess --> Listening
    Compensate --> Listening
```

### 5. Event Service (NoSQL Data)

MongoDB üzerinde esnek veri yönetimi.

```mermaid
graph TD
    Request[GET /api/events] --> Handler
    Handler --> Service
    
    subgraph "Event Logic"
        Service --> |Calc Skip/Limit| Repository
        Repository --> |BSON Filter| Mongo[(MongoDB)]
        Mongo --> |Documents| Repository
        Repository --> |Event Structs| Service
    end
    
    Service --> Response[JSON List]

```

```

