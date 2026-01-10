package main

import (
  "encoding/json"
  "strings"
  "fmt"
  "os"
  dysmsapi20170525  "github.com/alibabacloud-go/dysmsapi-20170525/v5/client"
  openapi  "github.com/alibabacloud-go/darabonba-openapi/v2/client"
  util  "github.com/alibabacloud-go/tea-utils/v2/service"
  credential  "github.com/aliyun/credentials-go/credentials"
  "github.com/alibabacloud-go/tea/tea"
)


// Description:
// 
// 使用凭据初始化账号Client
// 
// @return Client
// 
// @throws Exception
func CreateClient () (_result *dysmsapi20170525.Client, _err error) {
  // 工程代码建议使用更安全的无AK方式，凭据配置方式请参见：https://help.aliyun.com/document_detail/378661.html。
  credential, _err := credential.NewCredential(nil)
  if _err != nil {
    return _result, _err
  }

  config := &openapi.Config{
    Credential: credential,
  }
  // Endpoint 请参考 https://api.aliyun.com/product/Dysmsapi
  config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
  _result = &dysmsapi20170525.Client{}
  _result, _err = dysmsapi20170525.NewClient(config)
  return _result, _err
}

func _main (args []*string) (_err error) {
  client, _err := CreateClient()
  if _err != nil {
    return _err
  }

  sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
    SignName: tea.String("惠州顺视智能科技"),
    TemplateCode: tea.String("SMS_330275014"),
    PhoneNumbers: tea.String("17398467065"),
    TemplateParam: tea.String("{\"code\":\"1234\"}"),
  }
  runtime := &util.RuntimeOptions{}
  tryErr := func()(_e error) {
    defer func() {
      if r := tea.Recover(recover()); r != nil {
        _e = r
      }
    }()
    resp, _err := client.SendSmsWithOptions(sendSmsRequest, runtime)
    if _err != nil {
      return _err
    }

    fmt.Printf("[LOG] %v\n", resp)

    return nil
  }()

  if tryErr != nil {
    var error = &tea.SDKError{}
    if _t, ok := tryErr.(*tea.SDKError); ok {
      error = _t
    } else {
      error.Message = tea.String(tryErr.Error())
    }
    // 此处仅做打印展示，请谨慎对待异常处理，在工程项目中切勿直接忽略异常。
    // 错误 message
    fmt.Println(tea.StringValue(error.Message))
    // 诊断地址
    var data interface{}
    d := json.NewDecoder(strings.NewReader(tea.StringValue(error.Data)))
    d.Decode(&data)
    if m, ok := data.(map[string]interface{}); ok {
      recommend, _ := m["Recommend"]
      fmt.Println(recommend)
    }
  }
  return _err
}


func main() {
  err := _main(tea.StringSlice(os.Args[1:]))
  if err != nil {
    panic(err)
  }
}
