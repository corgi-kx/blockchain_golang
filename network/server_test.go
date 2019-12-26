package network

import "testing"

func TestBytes(t *testing.T) {
	t.Log("测试命令拼接、拆分功能")
	{
		t.Log("\t测试拼接功能：")
		v := version{versionInfo, 10, nil}
		b := jointMessage(cVersion, v.serialize())
		t.Log("\t拼接后的字节数组为:", b)
		t.Log("\t测试拆分功能：")
		cmd, content := splitMessage(b)
		newV := version{}
		newV.deserialize(content)
		t.Logf("\t命令为：%s,长度：%d ", cmd, len(cmd))
		t.Logf("\t版本信息：%v ", newV)
	}
}
