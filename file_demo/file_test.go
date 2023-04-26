// create by chencanhua in 2023/4/25
package file_demo

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestFile(t *testing.T) {
	fmt.Println(os.Getwd())
	f, err := os.Open("testdata/aivier.txt")
	require.NoError(t, err)
	data := make([]byte, 64)
	n, err := f.Read(data)
	fmt.Println(n)
	require.NoError(t, err)
	f.Close()

	//n, err = f.WriteString("yes I am")
	//fmt.Println(n)
	//// Access is denied 无权限
	//require.NoError(t, err)
	//f.Close()

	f, err = os.OpenFile("testdata/aivier.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	require.NoError(t, err)
	n, err = f.WriteString("hello")
	fmt.Println(n)
	require.NoError(t, err)
	f.Close()

	f, err = os.Create("testdata/my_file_copy.txt")
	require.NoError(t, err)
	n, err = f.WriteString("hello, world")
	fmt.Println(n)
	require.NoError(t, err)
}
