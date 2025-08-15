/*
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-07-31 14:30:20
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-08-13 09:38:31
 * @FilePath: /frontend/src/api/https/errorCodeType.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
export const errorCodeType = function(status:Number):string{
    let errMessage:string = "未知错误"
    switch (status) {
        case 400:
            errMessage = '请求参数无效'
            break
        case 401:
            errMessage = '未认证'
            break
        case 403:
            errMessage = '当前账号无权限访问！'
            break
        case 404:
            errMessage = '你所访问的资源不存在！'
            break
        case 405:
            errMessage = '请求方式错误！请您稍后重试'
            break
        case 408:
            errMessage = '请求超时！请您稍后重试'
            break
        case 500:
            errMessage = '服务器内部错误'
            break
        case 501:
            errMessage = '网络未实现'
            break
        case 502:
            errMessage = '网络错误'
            break
        case 503:
            errMessage = '服务不可用'
            break
        case 504:
            errMessage = '网络超时'
            break
        case 505:
            errMessage = 'http版本不支持该请求'
            break
        default:
            errMessage = `其他连接错误 --${status}`
    }
    return errMessage
}
