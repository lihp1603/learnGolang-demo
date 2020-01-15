package main

import (
	"fmt"
)

//求n个元素的全组合算法记录

func main() {
	s := []int{1, 2, 3}
	s1 := s[0 : len(s)-1]
	fmt.Println(s1)
	//求组合数
	// GetAllCombResultBit("abc")
	var pre []byte
	var result []string
	//求k个元素的组合情况
	// findCombinations("abc", 2, 0, pre, &result)
	//求k个元素的组合情况
	// findCombinationsPruning("abc", 2, 0, pre, &result)
	//求全组合
	// findCombinationsPruning("abc", -1, 0, pre, &result)
	//优化了边界后的，进行求k个元素的组合情况
	findCombinationsPruningBord("abc", 2, 0, pre, &result)
	fmt.Printf("%+v", result)
}

//使用循环求全组合的情况-基于位图
//如果不是求字符的所有排列，而是求字符的所有组合应该怎么办呢？还是输入三个字符 a、b、c，则它们的组合有a b c ab ac bc abc。
//当然我们还是可以借鉴全排列的思路，利用问题分解的思路，最终用递归解决。不过这里介绍一种比较巧妙的思路 —— 基于位图。
//假设原有元素 n 个，则最终组合结果是 2^n−1  个。
//我们可以用位操作方法：假设元素原本有：a,b,c 三个，则 1 表示取该元素，0 表示不取。故取a则是001，取ab则是011。
//所以一共三位，每个位上有两个选择 0 和 1。而000没有意义，所以是2^n−1个结果。
//这些结果的位图值都是 1,2…2^n-1。所以从值 1 到值 2^n−1 依次输出结果：
//001,010,011,100,101,110,111 。对应输出组合结果为：a,b,ab,c,ac,bc,abc。
//因此可以循环 1~2^n-1，然后输出对应代表的组合即可
//代码如下
func GetAllCombResultBit(strIn string) []string {
	orgIn := []byte(strIn)  //强制类型转换
	inLen := len(orgIn)     //求出这个字符串中多少个字符
	n := 1 << uint32(inLen) //其中n-1为全组合的数据量大小
	fmt.Println(n - 1)
	//
	var allComRes []string
	for i := 1; i < n; i++ {
		var comRes []byte
		for j := 0; j < inLen; j++ {
			tmp := i
			if tmp&(1<<uint32(j)) > 0 { //如果对应的位上为1的话，就说明这个字符出现，为0表示这个字符不出现
				comRes = append(comRes, orgIn[j])
			}
		}
		fmt.Printf("%s\n", string(comRes))
		allComRes = append(allComRes, string(comRes))
	}
	return allComRes
}

//采用回溯法来求解
//获取strin中取k个字符串的组合结果
func findCombinations(strIn string, k, begin int, pre []byte, result *[]string) {
	orgIn := []byte(strIn) //强制类型转换
	inLen := len(orgIn)    //求出这个字符串中多少个字符
	if len(pre) == k {
		*result = append(*result, string(pre))
		return
	}
	fmt.Printf("%s\n", string(pre))
	for i := begin; i < inLen; i++ {
		pre = append(pre, orgIn[i]) //取原来字符串中的第i个,下一轮就是剩下的子集
		findCombinations(strIn, k, i+1, pre, result)
		pre = pre[0 : len(pre)-1] //pop最后一个
	}
}

//在上面的方法中再进行进一步的优化算法:剪枝
//回溯+剪枝
//
func findCombinationsPruning(strIn string, k, begin int, pre []byte, result *[]string) {
	orgIn := []byte(strIn) //强制类型转换
	inLen := len(orgIn)    //求出这个字符串中多少个字符
	if k > 0 {             //如果这个地方k有设置值的话，就表示是求k个字符的组合情况
		if len(pre) == k {
			*result = append(*result, string(pre))
			return
		}
		//目前pre中装的字符数据个数+剩下的数据个数<k的话，意思是，就算继续往这个分支走，也是不符合要求的，所以我们提前结束
		//还差几个字符，剩下字符数量
		if k-len(pre) > inLen-begin {
			return
		}
	} else if len(pre) > 0 { //否则这里求全量组合，即求出全组合
		*result = append(*result, string(pre))
	}

	fmt.Printf("%s\n", string(pre))
	for i := begin; i < inLen; i++ {
		pre = append(pre, orgIn[i]) //取原来字符串中的第i个,下一轮就是剩下的子集
		findCombinationsPruning(strIn, k, i+1, pre, result)
		pre = pre[0 : len(pre)-1] //pop最后一个
	}
}

//添加边界条件，对求取k个数的组合情况
func findCombinationsPruningBord(strIn string, k, begin int, pre []byte, result *[]string) {
	orgIn := []byte(strIn) //强制类型转换
	inLen := len(orgIn)    //求出这个字符串中多少个字符
	//计算边界
	max_border := inLen //计算这轮循环的边界条件
	if k > 0 {          //如果这个地方k有设置值的话，就表示是求k个字符的组合情况
		max_border = inLen - (k - len(pre)) + 1 //计算这轮循环的边界条件
		if len(pre) == k {
			*result = append(*result, string(pre))
			return
		}
		//目前pre中装的字符数据个数+剩下的数据个数<k的话，意思是，就算继续往这个分支走，也是不符合要求的，所以我们提前结束
		//还差几个字符>剩下字符数量
		if k-len(pre) > inLen-begin {
			return
		}
	} else if len(pre) > 0 { //否则这里求全量组合，即求出全组合
		*result = append(*result, string(pre))
	}

	fmt.Printf("%s\n", string(pre))
	for i := begin; i < max_border; i++ {
		pre = append(pre, orgIn[i]) //取原来字符串中的第i个,下一轮就是剩下的子集
		findCombinationsPruningBord(strIn, k, i+1, pre, result)
		pre = pre[0 : len(pre)-1] //pop最后一个
	}
}
