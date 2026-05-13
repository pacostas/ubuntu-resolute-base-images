package acceptance_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/onsi/gomega/format"
	"github.com/paketo-buildpacks/packit/v2/pexec"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
)

var tinyImages struct {
	RunArchive string
	RunImageID string
}

var baseImages struct {
	BuildArchive string
	RunArchive   string
	BuildImageID string
	RunImageID   string
}

var staticImages struct {
	RunArchive string
	RunImageID string
}

var RegistryUrl string

func by(_ string, f func()) { f() }

func TestAcceptance(t *testing.T) {

	format.MaxLength = 0
	SetDefaultEventuallyTimeout(30 * time.Second)

	Expect := NewWithT(t).Expect

	RegistryUrl = os.Getenv("REGISTRY_URL")
	Expect(RegistryUrl).NotTo(Equal(""))

	root, err := filepath.Abs(".")
	Expect(err).ToNot(HaveOccurred())

	baseImages.BuildArchive = filepath.Join(root, "builds", "resolute-base-images", "build.oci")
	baseImages.BuildImageID = fmt.Sprintf("%s/resolute-base-build-image-%s", RegistryUrl, uuid.NewString())

	baseImages.RunArchive = filepath.Join(root, "builds", "resolute-base-images", "run.oci")
	baseImages.RunImageID = fmt.Sprintf("%s/resolute-base-run-image-%s", RegistryUrl, uuid.NewString())

	tinyImages.RunArchive = filepath.Join(root, "builds", "resolute-tiny-images", "run.oci")
	tinyImages.RunImageID = fmt.Sprintf("%s/resolute-tiny-run-image-%s", RegistryUrl, uuid.NewString())

	staticImages.RunArchive = filepath.Join(root, "builds", "resolute-static-images", "run.oci")
	staticImages.RunImageID = fmt.Sprintf("%s/resolute-static-run-image-%s", RegistryUrl, uuid.NewString())

	suite := spec.New("Acceptance", spec.Report(report.Terminal{}), spec.Parallel())
	suite("MetadataBaseImages", testMetadataBaseImages)
	suite("MetadataTinyImages", testMetadataTinyImages)
	suite("MetadataStaticImages", testMetadataStaticImages)
	suite("BuildpackIntegrationTinyStack", testBuildpackIntegrationTinyStack)
	suite("BuildpackIntegrationBaseStack", testBuildpackIntegrationBaseStack)
	suite.Run(t)
}

func createBuilder(config string, name string) (string, error) {
	buf := bytes.NewBuffer(nil)

	pack := pexec.NewExecutable("pack")
	err := pack.Execute(pexec.Execution{
		Stdout: buf,
		Stderr: buf,
		Args: []string{
			"builder",
			"create",
			name,
			fmt.Sprintf("--config=%s", config),
		},
	})
	return buf.String(), err
}
