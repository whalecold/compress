package lz77

import (
	"github.com/whalecold/zlip/pkg/huffman"
	"github.com/whalecold/zlip/pkg/stack"
)

func dealWithBytesAndStack(stackNode *stack.Stack, r *[]byte) {
	temp := make([]byte, 0, RLCMaxLength)
	for node := stackNode.Pop(); node != nil; {
		temp = append(temp, node.(byte))
		node = stackNode.Pop()
	}
	if temp[0] == RLCZero && len(temp) >= RLCLength {
		tempLen := len(temp)
		for tempLen > huffman.CCLLen {
			*r = append(*r, RLCSpecial)
			*r = append(*r, huffman.CCLLen)
			tempLen -= huffman.CCLLen
		}

		if tempLen != 0 {
			*r = append(*r, RLCSpecial)
			*r = append(*r, byte(tempLen))
		}

	} else {
		*r = append(*r, temp...)
	}
}

//RLC rlc
//游程编码
//run length coding 这里感觉非0的重复不会很多 只对0进行编码
//简单处理下 17表示0 后面的数字表示0重复的个数 多于3个重复才开始编码
func RLC(bytes []byte) []byte {
	result := make([]byte, 0, len(bytes))
	stackNode := stack.NewStack()

	stackNode.Push(bytes[0])

	for i := 1; i < len(bytes); i++ {
		lastNode := stackNode.Back()
		if lastNode != nil {
			lastData := lastNode.(byte)
			if lastData != bytes[i] {
				dealWithBytesAndStack(stackNode, &result)
			}
		}
		stackNode.Push(bytes[i])
	}
	dealWithBytesAndStack(stackNode, &result)
	return bytes
}

// UnRLC ..
func UnRLC(bytes []byte) []byte {
	result := make([]byte, 0, RLCMaxLength)
	for i := 0; i < len(bytes); i++ {
		if bytes[i] == RLCSpecial {
			tempLen := bytes[i+1]
			for k := byte(0); k < tempLen; k++ {
				result = append(result, RLCZero) //nolint
			}
			i++
		} else {
			result = append(result, bytes[i]) //nolint
		}
	}
	return bytes
}
