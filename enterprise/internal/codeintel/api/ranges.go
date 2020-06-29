package api

import (
	"context"
	"strings"

	"github.com/inconshreveable/log15"
	"github.com/pkg/errors"
	"github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/bundles/client"
	bundles "github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/bundles/client"
)

// TODO - document
func (api *codeIntelAPI) Ranges(ctx context.Context, file string, uploadID int) ([]bundles.Range, error) {
	dump, exists, err := api.store.GetDumpByID(ctx, uploadID)
	if err != nil {
		return nil, errors.Wrap(err, "store.GetDumpByID")
	}
	if !exists {
		return nil, ErrMissingDump
	}

	pathInBundle := strings.TrimPrefix(file, dump.Root)
	bundleClient := api.bundleManagerClient.BundleClient(dump.ID)

	ranges, err := bundleClient.Ranges(ctx, pathInBundle)
	if err != nil {
		if err == client.ErrNotFound {
			log15.Warn("Bundle does not exist")
			return nil, nil
		}
		return nil, errors.Wrap(err, "bundleClient.Ranges")
	}

	// TODO - need to attach dump?
	return ranges, nil
}
