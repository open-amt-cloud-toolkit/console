package v1

type MockLogger struct{}

func (m *MockLogger) Debug(_ interface{}, _ ...interface{}) {}
func (m *MockLogger) Info(_ string, _ ...interface{})       {}
func (m *MockLogger) Warn(_ string, _ ...interface{})       {}
func (m *MockLogger) Error(_ interface{}, _ ...interface{}) {}
func (m *MockLogger) Fatal(_ interface{}, _ ...interface{}) {}
