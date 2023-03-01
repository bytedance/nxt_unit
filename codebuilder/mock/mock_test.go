package mock

import (
	"testing"
	"time"
)

func Mock(name string) (int, error) {
	switch time.Now().UnixNano() {
	case 0:
		return 89, nil
	case 78:
		return 77, nil
	}
	return 0, nil
}

type s struct {
}

func (s *s) Mock(name string) (int, error) {
	switch time.Now().UnixNano() {
	case 0:
		return 89, nil
	case 78:
		return 77, nil
	}
	return 0, nil
}

// TODO: @guancheng.liu write a wrong unit test
// func TestMockCall(t *testing.T) {
// 	mockRender := &StatementRender{
// 		MockStatement: []string{},
// 	}
// 	mockito.Mock(Mock).To(MakeCall(context.Background(),"Mock", mockRender, Mock)).Build()
// 	i, err := Mock("name")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Logf("number is  %v", i)
// 	t.Log(mockRender.MockStatement)
// }

// TODO: @guancheng.liu write a wrong unit test
//func TestReceiverMockCall(t *testing.T) {
//	mockRender := &StatementRender{
//		MockStatement: []string{},
//	}
//	mockito.Mock((*s).Mock).To(MakeCall("(*s).Mock", mockRender, (*s).Mock)).Build()
//	var sf *s
//	i, err := sf.Mock("name")
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Logf("number is  %v", i)
//}

func TestMakeCall(t *testing.T) {

}
