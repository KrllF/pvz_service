//go:build suite

package httptestsuite

import "encoding/base64"

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
