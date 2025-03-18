# Transaction Decoder

This script decodes a binary-encoded XRP Ledger object (blob) and prints it as a human-readable JSON. It supports decoding blobs encoded in hex, base64, or base64(hex(blob)).

## Usage

The script expects a single command-line argument: the encoded blob.

```bash
go run main.go [blob]
```

Where `[blob]` is the encoded data.  The script will attempt to decode the blob using the following logic:

1. **Hex Decoding:** It first checks if the blob is a valid hexadecimal string. If so, it decodes it as hex and prints the JSON representation of the decoded object.
2. **Base64 Decoding:** If the blob is *not* a valid hex string, it attempts to decode it as a Base64 string.
    *   If Base64 decoding is successful, it then checks if *the decoded result* is a valid hex string. If it is, the script decodes the hex string and prints the JSON.
    *   If the Base64 decoded result is *not* a hex string, the script encodes *the result* in hex. It will try to decode this value and print as JSON.
3. **Failure:** If none of the above steps succeed, the script prints an error message and exits.

## Notes

Install dependencies as follows:

```bash
go get ./...
```

