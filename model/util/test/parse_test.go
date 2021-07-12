package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/figment-networks/mina-indexer/model/util"
)

func TestParseMemoText(t *testing.T) {
	examples := map[string]string{
		"E4Ygr3AhYC4HsXTypxjUXudZYuivoiurTLZEkd6zHt7hGijCTSbjP": "1594685079",
		"E4Ygr3AhXttkpiWqvshsMH4Y8Fh5FzW7zE9wdgeMRmZp8GQw37Zqk": "1593702730",
		"E4Ygr3AhYBijB1btFi7cCLtMkXwDbAiWtZFa6k93oZwRoQyFPqpc1": "1594153648",
		"E4YM2vTHhWEg66xpj52JErHUBU4pZ1yageL4TVDDpTTSsv8mK6YaH": "",
		"E4YXJXbj1YBnwfCF3LDb8muzU7u8VWD32yHc7cAnRgduLi2pUsx8U": "Hello",
		"E4YiyKK5dkD5mj4A6FnoZ77929s6vFCteoXnxSHbEMHUPpnzEs9kH": "I am a memo",
		"E4YoBsLAeXCBHfU1puDrjpVsNyGM9sTj5HpZoiZpddaboj453bBMA": "timsmith#2774",
	}
	for given, expected := range examples {
		assert.Equal(t, expected, util.ParseMemoText(given))
	}
}
