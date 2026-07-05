// TODO: Поменяй имя модуля github.com/1mpuser на своё и обнови все импорты
module github.com/1mpuser/inventory

go 1.26.0

require (
	github.com/1mpuser/platform v0.0.0-00010101000000-000000000000
	github.com/1mpuser/shared v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.11.1
	google.golang.org/grpc v1.79.2
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/1mpuser/platform v0.0.0-00010101000000-000000000000
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/ilyakaznacheev/cleanenv v1.5.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.42.0 // indirect
	golang.org/x/net v0.51.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260226221140-a57be14db171 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

replace github.com/1mpuser/shared => ../shared

replace github.com/1mpuser/platform => ../platform
