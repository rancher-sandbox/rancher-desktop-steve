module github.com/rancher/steve

go 1.19

replace (
	github.com/crewjam/saml => github.com/rancher/saml v0.2.0
	github.com/knative/pkg => github.com/rancher/pkg v0.0.0-20181214184433-b04c0947ad2f
	github.com/matryer/moq => github.com/rancher/moq v0.0.0-20190404221404-ee5226d43009

	k8s.io/client-go => github.com/rancher/client-go v1.25.4-rancher1
)

require (
	github.com/adrg/xdg v0.4.0
	github.com/golang/mock v1.6.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.5.0
	github.com/pborman/uuid v1.2.1
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.16.0
	github.com/rancher/apiserver v0.0.0-20230809162532-de731cd43b96
	github.com/rancher/dynamiclistener v0.3.5
	github.com/rancher/kubernetes-provider-detector v0.1.5
	github.com/rancher/norman v0.0.0-20230809232250-ccd31a468197
	github.com/rancher/remotedialer v0.3.0
	github.com/rancher/wrangler v1.1.1
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.8.4
	github.com/urfave/cli v1.22.14
	github.com/urfave/cli/v2 v2.25.7
	golang.org/x/sync v0.3.0
	k8s.io/api v0.25.12
	k8s.io/apiextensions-apiserver v0.25.12
	k8s.io/apimachinery v0.25.12
	k8s.io/apiserver v0.25.12
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog v1.0.0
	k8s.io/kube-aggregator v0.25.12
	k8s.io/kube-openapi v0.0.0-20220803162953-67bda5d908f1
)

require (
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emicklei/go-restful/v3 v3.8.0 // indirect
	github.com/evanphx/json-patch v4.12.0+incompatible // indirect
	github.com/felixge/httpsnoop v1.0.1 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.5 // indirect
	github.com/go-openapi/swag v0.19.14 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/gnostic v0.5.7-v3refs // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mxk/go-flowrate v0.0.0-20140419014527-cca7078d478f // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.42.0 // indirect
	github.com/prometheus/procfs v0.10.1 // indirect
	github.com/rancher/lasso v0.0.0-20221227210133-6ea88ca2fbcc // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	go.opentelemetry.io/contrib v0.20.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.20.0 // indirect
	go.opentelemetry.io/otel v0.20.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp v0.20.0 // indirect
	go.opentelemetry.io/otel/metric v0.20.0 // indirect
	go.opentelemetry.io/otel/sdk v0.20.0 // indirect
	go.opentelemetry.io/otel/sdk/export/metric v0.20.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v0.20.0 // indirect
	go.opentelemetry.io/otel/trace v0.20.0 // indirect
	go.opentelemetry.io/proto/otlp v0.7.0 // indirect
	golang.org/x/crypto v0.11.0 // indirect
	golang.org/x/net v0.12.0 // indirect
	golang.org/x/oauth2 v0.10.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/term v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220502173005-c8bf987b8c21 // indirect
	google.golang.org/grpc v1.47.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/component-base v0.25.12 // indirect
	k8s.io/klog/v2 v2.100.1 // indirect
	k8s.io/utils v0.0.0-20221011040102-427025108f67 // indirect
	sigs.k8s.io/apiserver-network-proxy/konnectivity-client v0.0.37 // indirect
	sigs.k8s.io/cli-utils v0.27.0 // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)
