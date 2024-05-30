package stringutils

import (
	"testing"
)

func FuzzIndexAllIgnoreCase(f *testing.F) {
	type caseType struct {
		Haystack string
		Needle   string
		Count    int
	}

	testcases := []caseType{{
		Haystack: testUnicodeMatchEndCaseLarge,
		Needle:   testUnicodeMatchEndCaseLarge[:5],
		Count:    9,
	},
		{
			Haystack: testMatchEndCaseLarge,
			Needle:   testUnicodeMatchEndCaseLarge[:9],
			Count:    10,
		},
		{
			Haystack: `list1=[0,1,2,3,4,5,6,7,8,9]#制作一个0-9的列表
list1.reverse()#reverse()函数直接对列表中的元素践行反向
print(list1)

# the following line is where it is breaking
list2=[str(i) for i in list1]#将列表中的每一个数字转换成字符串
print(list2)

str1="".join(list2)#通过join()函数，将列表中的单个字符串拼接成一整个字符串
print(str1)

str2=str1[2:8]#对字符串中的第三到第八字符进行切片
print(str2)

str3=str2[::-1]#通过右边第一个开始对整个字符串开始切片，以实现其翻转
print(str3)

i=int(str3)#int()函数试讲字符串转换为整数
print(i)#这里输出的结果虽然与print(str3)相同，但是性质是不同的

#转换成二进制、八进制、十六进制
print('转换成二进制:',bin(i),'转换成八进制:',oct(i), '转换成十六进制:',hex(i))
#二进制、八进制、十六进制这几个进制相互转换的时候，都要先转换为十进制int()`,
			Count: -1,
		},
		{
			Haystack: `func MultipleInputs(a, b int, name string) {
    // ... fancy code goes here 
}

func FuzzMultipleInputs(f *testing.F) {
  // We can add Multiple Seeds, but it has to be the same order as the input parameters for MultipleInputs
  f.Add(10,20,"John the Ripper")
  f.Fuzz(func(t *testing.T,a int,b int,name string){
      MultipleInputs(a,b,name)      
  })
}`,
			Needle: `jp morgan प्रीपेड`,
			Count:  10,
		},
		{
			Haystack: `func MultipleInputs(a, b int, name string) {
    // ... fancy code goes here 
}

func FuzzMultipleInputs(f *testing.F) {
  // We can add Multiple Seeds, but it has to be the same order as the input parameters for MultipleInputs
  f.Add(10,20,"John the Ripper")
  f.Fuzz(func(t *testing.T,a int,b int,name string){
      MultipleInputs(a,b,name)      
  })
}`,
			Needle: `tes`,
			Count:  -1,
		},
	}
	for _, tc := range testcases {
		f.Add(tc.Haystack, tc.Needle, tc.Count) // Use f.Add to provide a seed corpus
	}

	f.Fuzz(func(t *testing.T, h string, n string, c int) {
		_ = IndexAllIgnoreCase(h, n, c)
	})
}
