docker-up:
	@echo "Docker konteynerları başlatılıyor..."
	docker-compose up -d
	@echo "Konteynerlar başarıyla ayağa kalktı."

docker-down:
	@echo "Docker konteynerları durduruluyor..."
	docker-compose down
	@echo "Konteynerlar durduruldu."

run-all:
	@echo "Tüm servisler başlatılıyor..."
	@go run cmd/identity/main.go & 
	@go run cmd/event/main.go & 
	@go run cmd/api-gateway/main.go &
	@echo "Servisler arka planda çalışıyor. Durdurmak için: 'make stop-all'"

stop-all:
	@echo "Servisler durduruluyor..."
	@pkill -f "go run cmd/identity/main.go" || true
	@pkill -f "go run cmd/event/main.go" || true
	@pkill -f "go run cmd/api-gateway/main.go" || true
	@echo "Tüm servisler durduruldu."

run-identity:
	@echo "Identity servisi başlatılıyor..."
	go run cmd/identity/main.go

run-event:
	@echo "Event servisi başlatılıyor..."
	go run cmd/event/main.go

run-gateway:
	@echo "API Gateway başlatılıyor..."
	go run cmd/api-gateway/main.go