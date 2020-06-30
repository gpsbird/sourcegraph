package api

import (
	"context"
	"strings"

	"github.com/inconshreveable/log15"
	"github.com/pkg/errors"
	"github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/bundles/client"
	bundles "github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/bundles/client"
)

type ResolvedRangeView struct {
	Range       bundles.Range
	Definitions []ResolvedLocation
	References  []ResolvedLocation
	HoverText   string
}

// TODO - document
func (api *codeIntelAPI) Ranges(ctx context.Context, file string, uploadID int) ([]ResolvedRangeView, error) {
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

	var views []ResolvedRangeView
	for _, r := range ranges {
		views = append(views, ResolvedRangeView{
			Range:       r.Range,
			Definitions: resolveLocationsWithDump(dump, r.Definitions),
			References:  resolveLocationsWithDump(dump, r.References),
			HoverText:   r.HoverText,
		})
	}

	return views, nil
}
