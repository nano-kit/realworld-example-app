const Subscribe = {
  data() {
    return {
      connected: false,
      items: []
    }
  },
  methods: {
    call() {
      doSubscribe(this)
    }
  }
}

Vue.createApp(Subscribe).mount('#subscribe')

function doSubscribe(self) {
  if (self.connected) {
    return
  }
  self.items = []
  getAccessToken()
    .then(token => {
      const url = 'ws://127.0.0.1:8080/realworld/' + 'Clubhouse/Subscribe'
      const ws = new WebSocket(url)
      ws.onopen = () => {
        ws.send(JSON.stringify({
          token: token
        }))
        self.connected = true
      }
      ws.onmessage = ev => {
        const res = JSON.parse(ev.data)
        self.items.push(res)
      }
      ws.onclose = ev => {
        self.connected = false
        self.items.push(`ws close: ${ev.code} (clean=${ev.wasClean})`)
      }
      ws.onerror = ev => {
        self.connected = false
        self.items.push(`ws error: ${ev}`)
      }
    })
    .catch(err => {
      self.items.push(`getAccessToken: ${err}`)
    })
}
