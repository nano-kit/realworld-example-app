const serverURL = 'http://127.0.0.1:8080/realworld/'

async function postData(servicePath, data = {}) {
    const url = serverURL + servicePath
    const response = await fetch(url, {
        method: 'POST',
        mode: 'cors',
        cache: 'no-cache',
        credentials: 'same-origin',
        headers: {
            'Content-Type': 'application/json'
        },
        redirect: 'follow',
        referrerPolicy: 'no-referrer',
        body: JSON.stringify(data)
    })
    const body = await response.json()
    return new Promise((resolve, reject) => {
        if (isError(body)) {
            reject(body)
        } else {
            resolve(body)
        }
    })
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
