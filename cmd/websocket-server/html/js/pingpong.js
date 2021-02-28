const PingPong = {
  data() {
    return {
      limit: 5,
      items: []
    }
  },
  methods: {
    call() {
      doPingPong(this)
    }
  }
}

Vue.createApp(PingPong).mount('#pingpong')

function doPingPong(self) {
  self.items = []
  const url = 'ws://127.0.0.1:8080/realworld/' + 'Realworld/PingPong'
  const ws = new WebSocket(url)
  let stroke = 1
  ws.onopen = () => {
    ws.send(JSON.stringify({
      stroke: stroke
    }))
  }
  ws.onmessage = ev => {
    const res = JSON.parse(ev.data)
    self.items.push(res)
    stroke++
    if (stroke <= self.limit) {
      ws.send(JSON.stringify({
        stroke: stroke
      }))
    } else {
      ws.close()
    }
  }
  ws.onclose = ev => {
    self.items.push(`ws close: ${ev.code} (clean=${ev.wasClean})`)
  }
  ws.onerror = ev => {
    self.items.push(`ws error: ${ev}`)
  }
}
