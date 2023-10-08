package cdp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/util"
)

// OutputFile auto creates file if not exists, it will try to detect the data type and
// auto output binary, string or json
func OutputFile(p string, data interface{}) error {
	var (
		bin []byte
		err error
	)

	if err = util.EnsureDir(filepath.Dir(p)); err != nil {
		return fmt.Errorf("ensure dir: %v", err)
	}

	switch t := data.(type) {
	case []byte:
		bin = t

	case string:
		bin = []byte(t)

	case io.Reader:
		f, _ := util.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o664)
		_, err = util.CopyReader(f, t)

		return err

	default:
		bin, err = MustToJSONBytes(data)

		if err != nil {
			return err
		}
	}

	return util.WriteFile(p, bin, 0o664)
}

// MustToJSONBytes encode data to json bytes
func MustToJSONBytes(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(data); err != nil {
		return nil, fmt.Errorf("encode: %v", err)
	}

	b := buf.Bytes()

	return b[:len(b)-1], nil
}
