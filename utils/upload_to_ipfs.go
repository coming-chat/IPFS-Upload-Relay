package utils

import (
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/global"
	shell "github.com/ipfs/go-ipfs-api"
	files "github.com/ipfs/go-ipfs-files"
	"io"
	"net/http"
)

func NewClient(projectId, projectSecret string) *http.Client {
	return &http.Client{
		Transport: authTransport{
			RoundTripper:  http.DefaultTransport,
			ProjectId:     projectId,
			ProjectSecret: projectSecret,
		},
	}
}

// authTransport decorates each request with a basic auth header.
type authTransport struct {
	http.RoundTripper
	ProjectId     string
	ProjectSecret string
}

func (t authTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.SetBasicAuth(t.ProjectId, t.ProjectSecret)
	return t.RoundTripper.RoundTrip(r)
}

func UploadToIpfs(r io.ReadSeeker) (string, int64, error) {

	ipfs := shell.NewShellWithClient(global.IPFS_UPLOAD_URL, NewClient(global.PROJECT_ID, global.PROJECT_SECRET))

	fr := files.NewReaderFile(r)
	resp, err := ipfs.Add(fr, shell.CidVersion(1), shell.Pin(true))
	if err != nil {
		return "", 0, err
	}

	return resp, 0, nil
}
