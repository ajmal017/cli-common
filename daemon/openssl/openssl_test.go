package openssl
import (
	"testing"
)
func TestHmacSha256Signature(t *testing.T) {

	params := "symbol=LTCBTC&side=BUY&type=LIMIT&timeInForce=GTC&quantity=1&price=0.1&recvWindow=5000&timestamp=1499827319559"
	key := "NhqPtmdSJYdKjVHjA7PZj4Mge3R5YNiP1e3UZjInClVN65XAbvqqM6A7H5fATj0j"

	result, err := HmacSha256Signature( params, key)
	if err != nil {
		t.Fatalf("Failed to sign the message %s", err.Error())
	}
	expected := "c8db56825ae71d6d79447849e617115f4a920fa2acdcab2b053c4b2838bd6b71"
	if result != expected {
		t.Errorf("expected %s, but get %s", expected, result)
	}
}
