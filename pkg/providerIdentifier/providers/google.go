package providers

import (
	"net/http"
	"os"
	"strings"

	"github.com/devtron-labs/devtron/pkg/providerIdentifier/bean"
	"go.uber.org/zap"
)

type IdentifyGoogle struct {
	Logger *zap.SugaredLogger
}

func (impl *IdentifyGoogle) Identify() (string, error) {
	data, err := os.ReadFile("/sys/class/dmi/id/product_name")
	if err != nil {
		return bean.Unknown, err
	}
	if strings.Contains(string(data), "Google") {
		return bean.Google, nil
	}
	return bean.Unknown, nil
}

func (impl *IdentifyGoogle) IdentifyViaMetadataServer(detected chan<- string) {
	req, err := http.NewRequest("GET", "http://metadata.google.internal/computeMetadata/v1/instance/tags", nil)
	if err != nil {
		detected <- bean.Unknown
		return
	}
	req.Header.Set("Metadata-Flavor", "Google")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		detected <- bean.Unknown
		return
	}
	if resp.StatusCode == http.StatusOK {
		detected <- bean.Google
	} else {
		detected <- bean.Unknown
		return
	}
}
