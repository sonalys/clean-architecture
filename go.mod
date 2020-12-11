module github.com/spoonrocker/cart-go-sonalys

replace github.com/spoonrocker/cart-go-alysson/ => ./

go 1.13

require (
	github.com/jinzhu/copier v0.1.0
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.8.0
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.5.1
	golang.org/x/crypto v0.0.0-20201203163018-be400aefbc4c // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	gopkg.in/yaml.v2 v2.3.0
	gorm.io/driver/postgres v1.0.5
	gorm.io/gorm v1.20.8
)
