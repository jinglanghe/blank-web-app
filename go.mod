module github.com/apulis/bmod/aistudio-aom

go 1.16

require (
	github.com/apulis/go-business v0.0.0-00010101000000-000000000000
	github.com/apulis/sdk/go-utils v0.1.0
	github.com/gin-gonic/gin v1.7.4
	github.com/go-playground/validator/v10 v10.4.1
	github.com/go-redis/redis/v8 v8.11.1
	github.com/go-redsync/redsync/v4 v4.3.0
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/prometheus/alertmanager v0.21.0
	github.com/prometheus/client_golang v1.9.0
	github.com/prometheus/common v0.17.0
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/prometheus/prometheus v1.8.2-0.20200727090838-6f296594a852
	github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.0
	golang.org/x/oauth2 v0.0.0-20210628180205-a41e5a781914 // indirect
	google.golang.org/protobuf v1.26.0-rc.1 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/datatypes v1.0.1
	gorm.io/driver/postgres v1.2.2
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.22.2
	k8s.io/api v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/kubectl v0.18.0
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/apulis/go-business => ./deps/go-business
	github.com/apulis/sdk/go-utils => ./deps/go-utils
	k8s.io/client-go => k8s.io/client-go v0.18.0
)
