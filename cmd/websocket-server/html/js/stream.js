const Stream = {
  data() {
    return {
      count: 0,
      items: []
    }
  },
  methods: {
    call() {
      doStream(this)
    }
  }
}

Vue.createApp(Stream).mount('#stream')

function doStream(self) {
  self.items = []
  const url = 'ws://127.0.0.1:8080/realworld/' + 'Realworld/Stream'
  const ws = new WebSocket(url)
  ws.onopen = () => {
    ws.send(JSON.stringify({
      count: self.count
    }))
  }
  ws.onmessage = ev => {
    const res = JSON.parse(ev.data)
    self.items.push(res)
  }
  ws.onclose = ev => {
    self.items.push(`ws close: ${ev.code} (clean=${ev.wasClean})`)
  }
  ws.onerror = ev => {
    self.items.push(`ws error: ${ev}`)
  }
}
