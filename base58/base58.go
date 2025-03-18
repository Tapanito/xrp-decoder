package base58

import (
	"fmt"
)

// Encode encodes the passed bytes into a base58 encoded string.
func Encode(bin []byte) string {
	return FastBase58EncodingAlphabet(bin, XRPLAlphabet)
}

// EncodeAlphabet encodes the passed bytes into a base58 encoded string with the
// passed alphabet.
func EncodeAlphabet(bin []byte, alphabet *Alphabet) string {
	return FastBase58EncodingAlphabet(bin, alphabet)
}

// FastBase58Encoding encodes the passed bytes into a base58 encoded string.
func FastBase58Encoding(bin []byte) string {
	return FastBase58EncodingAlphabet(bin, XRPLAlphabet)
}

// FastBase58EncodingAlphabet encodes the passed bytes into a base58 encoded
// string with the passed alphabet.
func FastBase58EncodingAlphabet(bin []byte, alphabet *Alphabet) string {

	size := len(bin)

	zcount := 0
	for zcount < size && bin[zcount] == 0 {
		zcount++
	}

	// It is crucial to make this as short as possible, especially for
	// the usual case of bitcoin addrs
	size = zcount +
		// This is an integer simplification of
		// ceil(log(256)/log(58))
		(size-zcount)*555/406 + 1

	out := make([]byte, size)

	var i, high int
	var carry uint32

	high = size - 1
	for _, b := range bin {
		i = size - 1
		for carry = uint32(b); i > high || carry != 0; i-- {
			carry = carry + 256*uint32(out[i])
			out[i] = byte(carry % 58)
			carry /= 58
		}
		high = i
	}

	// Determine the additional "zero-gap" in the buffer (aside from zcount)
	for i = zcount; i < size && out[i] == 0; i++ {
	}

	// Now encode the values with actual alphabet in-place
	val := out[i-zcount:]
	size = len(val)
	for i = 0; i < size; i++ {
		out[i] = alphabet.encode[val[i]]
	}

	return string(out[:size])
}

// Decode decodes the base58 encoded bytes.
func Decode(str string) ([]byte, error) {
	return FastBase58DecodingAlphabet(str, XRPLAlphabet)
}

// DecodeAlphabet decodes the base58 encoded bytes using the given b58 alphabet.
func DecodeAlphabet(str string, alphabet *Alphabet) ([]byte, error) {
	return FastBase58DecodingAlphabet(str, alphabet)
}

// FastBase58Decoding decodes the base58 encoded bytes.
func FastBase58Decoding(str string) ([]byte, error) {
	return FastBase58DecodingAlphabet(str, XRPLAlphabet)
}

// FastBase58DecodingAlphabet decodes the base58 encoded bytes using the given
// b58 alphabet.
func FastBase58DecodingAlphabet(str string, alphabet *Alphabet) ([]byte, error) {
	if len(str) == 0 {
		return nil, fmt.Errorf("zero length string")
	}

	zero := alphabet.encode[0]
	b58sz := len(str)

	var zcount int
	for i := 0; i < b58sz && str[i] == zero; i++ {
		zcount++
	}

	var t, c uint64

	// the 32bit algo stretches the result up to 2 times
	binu := make([]byte, 2*((b58sz*406/555)+1))
	outi := make([]uint32, (b58sz+3)/4)

	for _, r := range str {
		if r > 127 {
			return nil, fmt.Errorf("high-bit set on invalid digit")
		}
		if alphabet.decode[r] == -1 {
			return nil, fmt.Errorf("invalid base58 digit (%q)", r)
		}

		c = uint64(alphabet.decode[r])

		for j := len(outi) - 1; j >= 0; j-- {
			t = uint64(outi[j])*58 + c
			c = t >> 32
			outi[j] = uint32(t & 0xffffffff)
		}
	}

	// initial mask depends on b58sz, on further loops it always starts at 24 bits
	mask := (uint(b58sz%4) * 8)
	if mask == 0 {
		mask = 32
	}
	mask -= 8

	outLen := 0
	for j := 0; j < len(outi); j++ {
		for mask < 32 { // loop relies on uint overflow
			binu[outLen] = byte(outi[j] >> mask)
			mask -= 8
			outLen++
		}
		mask = 24
	}

	// find the most significant byte post-decode, if any
	for msb := zcount; msb < len(binu); msb++ {
		if binu[msb] > 0 {
			return binu[msb-zcount : outLen], nil
		}
	}

	// it's all zeroes
	return binu[:outLen], nil
}

// Alphabet is a a b58 alphabet.
type Alphabet struct {
	decode [128]int8
	encode [58]byte
}

// NewAlphabet creates a new alphabet from the passed string.
//
// It panics if the passed string is not 58 bytes long, isn't valid ASCII,
// or does not contain 58 distinct characters.
func NewAlphabet(s string) *Alphabet {
	if len(s) != 58 {
		panic("base58 alphabets must be 58 bytes long")
	}
	ret := new(Alphabet)
	copy(ret.encode[:], s)
	for i := range ret.decode {
		ret.decode[i] = -1
	}

	distinct := 0
	for i, b := range ret.encode {
		if ret.decode[b] == -1 {
			distinct++
		}
		ret.decode[b] = int8(i)
	}

	if distinct != 58 {
		panic("provided alphabet does not consist of 58 distinct characters")
	}

	return ret
}

// XRPLAlphabet is the XRP Ledger base58 alphabet.
var XRPLAlphabet = NewAlphabet("rpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65jkm8oFqi1tuvAxyz")
