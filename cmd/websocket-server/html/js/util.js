const serverURL = 'http://127.0.0.1:8080/realworld/'

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

var accessToken = ''

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

function checkAccessToken(tokenString) {
    const parts = tokenString.split('.')
    if (parts.length != 3) {
        return null
    }
    try {
        const token = JSON.parse(window.atob(parts[1]))
        if (token.exp * 1000 > Date.now()) {
            return token
        }
        return null
    } catch (e) {
        return null
    }
}
