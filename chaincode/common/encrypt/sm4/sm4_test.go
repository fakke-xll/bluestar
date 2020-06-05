package sm4

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestSM4(t *testing.T) {
	//SMS4算法密钥长度也是128bit
	bykey := []byte("1234567890abcdef")
	fmt.Printf("原始密码key = %s %v\n", string(bykey), bykey)
	data := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10}
	WriteKeyToPem("key.pem", bykey, nil)

	key, err := ReadKeyFromPem("key.pem", nil)
	fmt.Printf("读出密码key = %s %v\n", string(key), key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("待加密数据data\t= %x\n", data)
	c, err := NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}
	d0 := make([]byte, 16)
	c.Encrypt(d0, data)
	fmt.Printf("加密后数据d0\t= %x\n", d0)
	d1 := make([]byte, 16)
	c.Decrypt(d1, d0)
	fmt.Printf("解密后数据d1\t= %x\n", d1)
	if sa := testCompare(data, d1); sa != true {
		fmt.Printf("Error data!")
	}
}

func TestSM4ex(t *testing.T) {
	//SMS4算法密钥长度也是128bit
	bykey := []byte("1234567890abcdef")
	fmt.Printf("原始密码key = %s %v\n", string(bykey), bykey)
	data := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10}
	fmt.Printf("待加密数据data\t= %x\n", data)

	d0 := make([]byte, 16)
	EncryptBlock(bykey, d0, data)
	fmt.Printf("加密后数据d0\t= %x\n", d0)

	d1 := make([]byte, 16)
	DecryptBlock(bykey, d1, d0)
	fmt.Printf("解密后数据d1\t= %x\n", d1)
	if sa := testCompare(data, d1); sa != true {
		log.Fatal("Error data!")
	}
}

func BenchmarkSM4(t *testing.B) {
	t.ReportAllocs()
	key := []byte("1234567890abcdef")
	data := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10}
	WriteKeyToPem("key.pem", key, nil)
	key, err := ReadKeyFromPem("key.pem", nil)
	if err != nil {
		log.Fatal(err)
	}
	c, err := NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < t.N; i++ {
		d0 := make([]byte, 16)
		c.Encrypt(d0, data)
		d1 := make([]byte, 16)
		c.Decrypt(d1, d0)
	}
}

func TestErrKeyLen(t *testing.T) {
	fmt.Printf("\n--------------test key len------------------")
	key := []byte("1234567890abcdefg")
	_, err := NewCipher(key)
	if err != nil {
		fmt.Println("\nError key len !")
	}
	key = []byte("1234")
	_, err = NewCipher(key)
	if err != nil {
		fmt.Println("Error key len !")
	}
	fmt.Println("------------------end----------------------")
}

func testCompare(key1, key2 []byte) bool {
	if len(key1) != len(key2) {
		return false
	}
	for i, v := range key1 {
		if i == 1 {
			fmt.Println("type of v", reflect.TypeOf(v))
		}
		a := key2[i]
		if a != v {
			return false
		}
	}
	return true
}
