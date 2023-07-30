package service

import (
	"math/rand"
	"testing"
)

func TestSnowflake_GenerateID(t *testing.T) {
	for r := 0; r < 10; r++ {
		// 使用工作节点 ID 为 1 创建一个新的 Snowflake 实例
		workID := rand.Int63n(1023) + 1
		sf := NewSnowflake(workID)

		// 用于检查标识符是否唯一的映射
		generatedIDs := make(map[int64]bool)

		// 生成 1000 个标识符并验证其唯一性
		for i := 0; i < 1000; i++ {
			id := sf.GenerateID()

			// 检查生成的标识符是否唯一
			if _, ok := generatedIDs[id]; ok {
				t.Errorf("生成了非唯一标识符：%d", id)
			}

			// 将标识符添加到映射中，以确保下次生成的标识符是唯一的
			generatedIDs[id] = true
		}

		// 再生成一个标识符，验证其正确性（格式）
		id := sf.GenerateID()
		// 检查工作节点 ID 是否正确地编码在标识符中
		expectedWorkerID := workID << 12
		if (id & (1023 << 12)) != expectedWorkerID {
			t.Errorf("生成的标识符中的工作节点 ID 不正确：%d", id)
		}
	}
}
