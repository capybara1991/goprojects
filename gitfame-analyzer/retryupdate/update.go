package retryupdate

import (
	"errors"

	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

func UpdateValue(
	c kvapi.Client,
	key string,
	updateFn func(oldValue *string) (newValue string, err error),
) error {
	type curState struct {
		val     *string
		version uuid.UUID
	}

	read := func() (curState, error) {
		for {
			resp, err := c.Get(&kvapi.GetRequest{Key: key})
			switch {
			case err == nil:
				v := resp.Value
				return curState{val: &v, version: resp.Version}, nil

			case errors.Is(err, kvapi.ErrKeyNotFound):
				return curState{val: nil, version: uuid.UUID{}}, nil

			default:
				var auth *kvapi.AuthError
				if errors.As(err, &auth) {
					return curState{}, err
				}
				var api *kvapi.APIError
				if errors.As(err, &api) {
					continue
				}
				continue
			}
		}
	}

	for {
		cur, err := read()
		if err != nil {
			return err
		}

		newVal, uerr := updateFn(cur.val)
		if uerr != nil {
			return uerr
		}

		if cur.val != nil && *cur.val == newVal {
			return nil
		}

		req := &kvapi.SetRequest{
			Key:        key,
			Value:      newVal,
			OldVersion: cur.version,
			NewVersion: uuid.Must(uuid.NewV4()),
		}

	setRetry:
		for {
			_, err := c.Set(req)
			if err == nil {
				return nil
			}

			var auth *kvapi.AuthError
			if errors.As(err, &auth) {
				return err
			}

			if errors.Is(err, kvapi.ErrKeyNotFound) {
				valFromNil, uerr := updateFn(nil)
				if uerr != nil {
					return uerr
				}
				req = &kvapi.SetRequest{
					Key:        key,
					Value:      valFromNil,
					OldVersion: uuid.UUID{},
					NewVersion: uuid.Must(uuid.NewV4()),
				}
				continue
			}

			var api *kvapi.APIError
			if errors.As(err, &api) {
				var conflict *kvapi.ConflictError
				if errors.As(api, &conflict) {
					if conflict.ExpectedVersion == req.NewVersion {
						return nil
					}
					break
				}
				continue setRetry
			}

			continue setRetry
		}

	}
}
