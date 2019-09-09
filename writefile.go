package typegen

import (
	"os"

	"golang.org/x/xerrors"
)

func WriteTupleEncodersToFile(fname, pkg string, types ...interface{}) error {
	fi, err := os.Create(fname)
	if err != nil {
		return xerrors.Errorf("failed to open file: %w", err)
	}
	defer fi.Close()

	if err := PrintHeaderAndUtilityMethods(fi, "types"); err != nil {
		return xerrors.Errorf("failed to write header: %w", err)
	}

	for _, t := range types {
		if err := GenTupleEncodersForType(t, fi); err != nil {
			return xerrors.Errorf("failed to generate encoders: %w", err)
		}
	}

	return nil
}
