module sigs.k8s.io/external-dns

go 1.22.4

require (
	cloud.google.com/go/compute/metadata v0.3.0
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.12.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.6.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dns/armdns v1.2.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/privatedns/armprivatedns v1.2.0
	github.com/F5Networks/k8s-bigip-ctlr/v2 v2.17.0
	github.com/IBM-Cloud/ibm-cloud-cli-sdk v1.4.0
	github.com/IBM/go-sdk-core/v5 v5.17.3
	github.com/IBM/networking-go-sdk v0.47.1
	github.com/akamai/AkamaiOPEN-edgegrid-golang v1.2.2
	github.com/alecthomas/kingpin/v2 v2.4.0
	github.com/aliyun/alibaba-cloud-sdk-go v1.62.771
	github.com/ans-group/sdk-go v1.17.0
	github.com/aws/aws-sdk-go v1.54.4
	github.com/bodgit/tsig v1.2.2
	github.com/cenkalti/backoff/v4 v4.3.0
	github.com/civo/civogo v0.3.70
	github.com/cloudflare/cloudflare-go v0.98.0
	github.com/cloudfoundry-community/go-cfclient v0.0.0-20190201205600-f136f9222381
	github.com/coredns/caddy v1.1.1
	github.com/datawire/ambassador v1.12.4
	github.com/denverdino/aliyungo v0.0.0-20230411124812-ab98a9173ace
	github.com/digitalocean/godo v1.118.0
	github.com/dnsimple/dnsimple-go v1.7.0
	github.com/exoscale/egoscale v0.102.3
	github.com/ffledgling/pdns-go v0.0.0-20180219074714-524e7daccd99
	github.com/go-gandi/go-gandi v0.7.0
	github.com/go-logr/logr v1.4.2
	github.com/google/go-cmp v0.6.0
	github.com/google/uuid v1.6.0
	github.com/gophercloud/gophercloud v1.12.0
	github.com/hooklift/gowsdl v0.5.0
	github.com/linki/instrumented_http v0.3.0
	github.com/linode/linodego v1.35.0
	github.com/maxatome/go-testdeep v1.14.0
	github.com/miekg/dns v1.1.61
	github.com/nesv/go-dynect v0.6.0
	github.com/nic-at/rc0go v1.1.1
	github.com/onsi/ginkgo v1.16.5
	github.com/openshift/api v0.0.0-20230607130528-611114dca681
	github.com/openshift/client-go v0.0.0-20230607134213-3cd0021bbee3
	github.com/oracle/oci-go-sdk/v65 v65.67.2
	github.com/ovh/go-ovh v1.6.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/pluralsh/gqlclient v1.11.0
	github.com/projectcontour/contour v1.29.1
	github.com/prometheus/client_golang v1.19.1
	github.com/scaleway/scaleway-sdk-go v1.0.0-beta.28
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.9.0
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.0.945
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod v1.0.945
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns v1.0.945
	github.com/transip/gotransip/v6 v6.24.0
	github.com/ultradns/ultradns-sdk-go v1.3.7
	github.com/vinyldns/go-vinyldns v0.9.16
	github.com/vultr/govultr/v2 v2.17.2
	go.etcd.io/etcd/api/v3 v3.5.14
	go.etcd.io/etcd/client/v3 v3.5.14
	go.uber.org/ratelimit v0.3.1
	golang.org/x/net v0.26.0
	golang.org/x/oauth2 v0.21.0
	golang.org/x/sync v0.7.0
	golang.org/x/time v0.5.0
	google.golang.org/api v0.185.0
	gopkg.in/ns1/ns1-go.v2 v2.11.0
	gopkg.in/yaml.v2 v2.4.0
	istio.io/api v1.22.1
	istio.io/client-go v1.22.1
	k8s.io/api v0.30.2
	k8s.io/apimachinery v0.30.2
	k8s.io/client-go v0.30.2
	k8s.io/klog/v2 v2.130.0
	sigs.k8s.io/gateway-api v1.1.0
)

require (
	cloud.google.com/go/auth v0.5.1 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.2 // indirect
	code.cloudfoundry.org/gofileutils v0.0.0-20170111115228-4d0c80011a0f // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.9.0 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.2.2 // indirect
	github.com/Masterminds/semver v1.4.2 // indirect
	github.com/Yamashou/gqlgenc v0.14.0 // indirect
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 // indirect
	github.com/alexbrainman/sspi v0.0.0-20180613141037-e580b900e9f5 // indirect
	github.com/ans-group/go-durationstring v1.2.0 // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/deepmap/oapi-codegen v1.9.1 // indirect
	github.com/emicklei/go-restful/v3 v3.12.0 // indirect
	github.com/evanphx/json-patch v5.7.0+incompatible // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/errors v0.21.0 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/strfmt v0.22.1 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.19.0 // indirect
	github.com/go-resty/resty/v2 v2.13.1 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.4 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.7 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/hcl v1.0.1-vault-5 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/goidentity/v6 v6.0.1 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.3 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/jinzhu/copier v0.3.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/openshift/gssapi v0.0.0-20161010215902-5fb4217df13b // indirect
	github.com/opentracing/opentracing-go v1.2.1-0.20220228012449-10b1cf09e00b // indirect
	github.com/pelletier/go-toml/v2 v2.1.0 // indirect
	github.com/peterhellberg/link v1.1.0 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.53.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/sagikazarmark/locafero v0.3.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/schollz/progressbar/v3 v3.8.6 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/smartystreets/go-aws-auth v0.0.0-20180515143844-0c1422d1fdb9 // indirect
	github.com/smartystreets/gunit v1.3.4 // indirect
	github.com/sony/gobreaker v0.5.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.10.0 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.17.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/terra-farm/udnssdk v1.3.5 // indirect
	github.com/vektah/gqlparser/v2 v2.5.14 // indirect
	github.com/xhit/go-str2duration/v2 v2.1.0 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.14 // indirect
	go.mongodb.org/mongo-driver v1.14.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.49.0 // indirect
	go.opentelemetry.io/otel v1.24.0 // indirect
	go.opentelemetry.io/otel/metric v1.24.0 // indirect
	go.opentelemetry.io/otel/trace v1.24.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/exp v0.0.0-20240416160154-fe59bbe5cc7f // indirect
	golang.org/x/mod v0.18.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/term v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	golang.org/x/tools v0.22.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240610135401-a8a62080eff3 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240617180043-68d350f18fd4 // indirect
	google.golang.org/grpc v1.64.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/resty.v1 v1.12.0 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/kube-openapi v0.0.0-20240423202451-8948a665c108 // indirect
	k8s.io/utils v0.0.0-20240423183400-0849a56e8f22 // indirect
	moul.io/http2curl v1.0.0 // indirect
	sigs.k8s.io/controller-runtime v0.18.4 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)
