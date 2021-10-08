# go-tink Demo

## How to use
- get tink
```
go get github.com/google/tink/go/...
```
- Set up KMS and create service account
- Edit `file.tar` on main.go
- Edit  main.go and chose Encode/Decode
- Run
```
go run ./main.go
```

### Encode
First time or without `key.json`
```
func main() {
	aad := []byte("Hello world")
	kh, _ := keyset.NewHandle(aead.AES128GCMKeyTemplate())
	createKeyJson("./key.json", kh)

	encode("file.tar", aad, kh)
}
```
with `key.json`
```
func main() {
	aad := []byte("Hello world")
  kh := getFromKMS("./key.json")

	encode("file.tar", aad, kh)
}
```

### Decode
with `key.json`
```
func main() {
	aad := []byte("Hello world")
  kh := getFromKMS("./key.json")

	decode("file.tar", aad, kh)
}
```
