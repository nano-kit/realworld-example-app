// API 网关的 URL 和 “服务路径”
const serverURL = 'http://127.0.0.1:8080/realworld/'

// postData 调用一次后端服务的 API
// 输入参数：
//     - servicePath 服务的 API 接口全名，如 'Realworld/Call'
//     - data 服务的 API 的输入参数，是一个 JSON 对象，具体规格参考服务的 protobuf 协议
// 成功返回：
//     服务的 API 的响应参数，是一个 JSON 对象，具体规格参考服务的 protobuf 协议
// 出错返回：
//     错误可以是一个 JSON 对象，具体规格为 {"id":"", "code":1, "detail":"", "status":""}
//     也可以是其它任何类型，如 "TypeError: Failed to fetch"
async function postData(servicePath, data = {}) {
    const url = serverURL + servicePath
    const token = await getAccessToken()
    const response = await fetch(url, {
        method: 'POST',
        mode: 'cors',
        cache: 'no-cache',
        credentials: 'same-origin',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + token
        },
        redirect: 'follow',
        referrerPolicy: 'no-referrer',
        body: JSON.stringify(data)
    })
    try {
        const body = await response.json()
        return new Promise((resolve, reject) => {
            if (isError(body)) {
                reject(body)
            } else {
                resolve(body)
            }
        })
    } catch (e) { // the response body to parse is not valid JSON
        return Promise.reject({ id: response.status, detail: response.statusText })
    }
}

function isError(body) {
    if (body.hasOwnProperty('id')
        && body.hasOwnProperty('code')
        && body.code > 299) {
        return true
    } else {
        return false
    }
}

// 在内存中保存的 accessToken
// https://dev.to/cotter/localstorage-vs-cookies-all-you-need-to-know-about-storing-jwt-tokens-securely-in-the-front-end-15id
var accessToken = ''

// 更新 accessToken
// 成功返回：
//     accessToken 类型为 string
// 出错返回：
//     其它任何类型
async function getAccessToken() {
    const token = checkAccessToken(accessToken)
    if (token != null) {
        return Promise.resolve(accessToken)
    }

    const response = await fetch('/token', {
        method: 'POST',
        mode: 'same-origin',
        cache: 'no-cache',
        credentials: 'same-origin',
        redirect: 'follow',
        referrerPolicy: 'no-referrer'
    })
    const body = await response.text()
    return new Promise((resolve, reject) => {
        const token = checkAccessToken(body)
        if (token == null) {
            reject(body)
        } else {
            accessToken = body
            resolve(accessToken)
        }
    })
}

// 校验 accessToken 是否在有效期内
// 成功返回：
//     token 对象，规格为
// {
//   "type": "user",
//   "scopes": [
//     "basic"
//   ],
//   "metadata": {
//     "AvatarUrl": "https://avatars.githubusercontent.com/u/5406765?v=4",
//     "Company": "",
//     "Email": "aclisp@gmail.com",
//     "Location": "Guangzhou, China",
//     "Name": "Homer Huang"
//   },
//   "exp": 1616413818,
//   "iss": "com.example",
//   "sub": "aclisp"
// }
// 出错返回：
//     null
function checkAccessToken(tokenString) {
    const parts = tokenString.split('.')
    if (parts.length != 3) {
        return null
    }
    try {
        const token = JSON.parse(window.atob(base64DecodeUrl(parts[1])))
        if (token.exp * 1000 > Date.now()) {
            return token
        }
        return null
    } catch (e) {
        return null
    }
}

/**
 * use this to make a base64 encoded string URL friendly,
 * i.e. '+' and '/' are replaced with '-' and '_' also any trailing '='
 * characters are removed
 *
 * @param {String} str the encoded string
 * @returns {String} the URL friendly encoded String
 */
 function base64EncodeUrl(str){
    return str.replace(/\+/g, '-').replace(/\//g, '_').replace(/\=+$/, '');
}

/**
 * Use this to recreate a base64 encoded string that was made URL friendly
 * using base64EncodeurlFriendly.
 * '-' and '_' are replaced with '+' and '/' and also it is padded with '+'
 *
 * @param {String} str the encoded string
 * @returns {String} the URL friendly encoded String
 */
function base64DecodeUrl(str){
    str = (str + '===').slice(0, str.length + (str.length % 4));
    return str.replace(/-/g, '+').replace(/_/g, '/');
}
