module github.com/1mpuser/assembly

go 1.26.0

require (
	github.com/1mpuser/platform v0.0.0-00010101000000-000000000000
	github.com/ilyakaznacheev/cleanenv v1.5.0
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

replace github.com/1mpuser/platform => ../platform

replace github.com/1mpuser/shared => ../shared
