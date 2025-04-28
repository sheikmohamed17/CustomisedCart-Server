module practice/models

replace practice/routes => ../routes

replace CustomizedCart/Emails => ../emails

go 1.22.8

require (
	github.com/go-sql-driver/mysql v1.8.1
	github.com/joho/godotenv v1.5.1
)

require filippo.io/edwards25519 v1.1.0 // indirect
