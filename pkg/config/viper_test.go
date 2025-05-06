package config

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试配置结构体
type testConfig struct {
	Name  string `mapstructure:"name"`
	Value int    `mapstructure:"value"`
}

// 测试数据库配置结构体
type testDBConfig struct {
	Master struct {
		Type        string        `mapstructure:"type"`
		DSN         string        `mapstructure:"dsn"`
		MaxOpenConn int           `mapstructure:"maxopen"`
		MaxIdleConn int           `mapstructure:"maxidle"`
		MaxLifeTime time.Duration `mapstructure:"maxlifetime"`
	} `mapstructure:"master"`
}

func TestInitConfig(t *testing.T) {
	// 创建临时配置文件
	dir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// 创建配置目录
	configDir := filepath.Join(dir, "configs")
	err = os.Mkdir(configDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// 创建测试配置文件
	apiConfigContent := `
app:
  name: "test-api"
  value: 123
database:
  master:
    type: "mysql"
    dsn: "user:pass@tcp(localhost:3306)/testdb"
    maxopen: 100
    maxidle: 10
    maxlifetime: 3600
`
	apiConfigPath := filepath.Join(configDir, "api-server.yaml")
	err = os.WriteFile(apiConfigPath, []byte(apiConfigContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 保存当前工作目录
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// 切换到测试目录
	err = os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd) // 测试结束后恢复工作目录

	// 准备测试配置结构体
	testAppConfig := &testConfig{}
	testDbConfig := &testDBConfig{}

	configs := map[string]interface{}{
		"app":      testAppConfig,
		"database": testDbConfig,
	}

	// 执行测试
	InitConfig("api", configs)

	// 验证配置是否正确解析
	assert.Equal(t, "test-api", testAppConfig.Name)
	assert.Equal(t, 123, testAppConfig.Value)
	assert.Equal(t, "mysql", testDbConfig.Master.Type)
	assert.Equal(t, "user:pass@tcp(localhost:3306)/testdb", testDbConfig.Master.DSN)
	assert.Equal(t, 100, testDbConfig.Master.MaxOpenConn)
	assert.Equal(t, 10, testDbConfig.Master.MaxIdleConn)
	assert.Equal(t, time.Duration(3600), testDbConfig.Master.MaxLifeTime)
}

func TestParseConfig(t *testing.T) {
	// 创建临时配置文件
	dir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// 创建配置目录
	configDir := filepath.Join(dir, "configs")
	err = os.Mkdir(configDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// 创建测试配置文件
	configContent := `
test:
  name: "test-config"
  value: 456
`
	configPath := filepath.Join(configDir, "test-config.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 保存当前工作目录
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// 切换到测试目录
	err = os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd) // 测试结束后恢复工作目录

	// 创建viper实例
	v := createTestViper("test-config")
	if v == nil {
		t.Fatal("无法创建Viper实例")
	}

	// 测试正确的配置解析
	testConfig := &testConfig{}
	configs := map[string]interface{}{
		"test": testConfig,
	}

	err = parseConfig(v, configs)
	assert.NoError(t, err)
	assert.Equal(t, "test-config", testConfig.Name)
	assert.Equal(t, 456, testConfig.Value)

	// 测试非指针类型的错误情况
	invalidConfig := *testConfig // 注意这里使用值而不是指针
	invalidConfigs := map[string]interface{}{
		"test": invalidConfig,
	}

	err = parseConfig(v, invalidConfigs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "配置必须是指针类型")

	// 测试空配置映射
	err = parseConfig(v, nil)
	assert.NoError(t, err)
}

// 辅助函数：创建测试用的viper实例
func createTestViper(configName string) *viper.Viper {
	v := viper.New()
	v.AddConfigPath("./configs")
	v.SetConfigName(configName)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		return nil // 返回nil而不是panic，让测试函数处理错误
	}
	return v
}
