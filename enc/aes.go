package enc

import (
	"bytes"
	"crypto/aes"
	"fmt"
	"github.com/pkg/errors"
)

// New :新建待加密类
func (r *Enc) New(key, data []byte) {
	r.key = key
	r.data = data
}

// New2 :接受变长参数, 类型: string or []byte, 第一参数: 默认text，双参数时为Key
func (r *Enc) New2(args ...interface{}) *Enc {
	switch len(args) {
	case 1:
		if r.key == nil {
			panic("No aes Key")
		} else {
			switch args[0].(type) {
			case []byte:
				b := args[0].([]byte)
				r.New(r.key, b)
			case string:
				r.New(r.key, []byte(args[0].(string)))
			default:
				panic("New2 expect string or []byte")
			}
		}
	case 2:
		r.New(args[0].([]byte), args[1].([]byte))
	default:
		panic(fmt.Sprint("New2 expect two params, but receive ", len(args)))
	}
	return r
}

// EcbEncrypt :AES ECB 模式 采用PKCS#7 填充
func (r *Enc) EcbEncrypt() ([]byte, error) {
	if b, err := aes.NewCipher(r.key); err == nil {
		bs := b.BlockSize()
		paddingData := PKCS7Padding(r.data, bs)
		decrypted := make([]byte, len(paddingData))
		for s, e := 0, bs; s < len(paddingData); s, e = s+bs, e+bs {
			b.Encrypt(decrypted[s:e], paddingData[s:e])
		}
		return decrypted, nil
	} else {
		return nil, err
	}
}

//  EcbDecrypt :AES ECB 模式 采用PKCS#7 解填充
func (r *Enc) EcbDecrypt() ([]byte, error) {
	if b, err := aes.NewCipher(r.key); err == nil {
		bs := b.BlockSize()

		decrypted := make([]byte, len(r.data))
		for s, e := 0, bs; s < len(r.data); s, e = s+bs, e+bs {
			b.Decrypt(decrypted[s:e], r.data[s:e])
		}
		length := len(decrypted)
		unpadding := int(decrypted[length-1])

		if (length - unpadding) < 0 {
			return nil, errors.New("error in unPadding")
		}
		return PKCS7UnPadding(decrypted), nil
	} else {
		return nil, err
	}
}

// PKCS7Padding :PKCS#7 填充
func PKCS7Padding(ciphertext []byte, bs int) []byte {
	padding := bs - len(ciphertext)%bs
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding :PKCS#7 解填充
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
