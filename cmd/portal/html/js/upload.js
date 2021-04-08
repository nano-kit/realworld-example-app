const Upload = {
  data() {
    return {
      input: "",
      output: ""
    }
  },
  methods: {
    call() {
      doUpload(this)
    }
  }
}

Vue.createApp(Upload).mount('#upload')

function doUpload(self) {
  self.output = ""
  const lines = self.input.split(/\r?\n/)
  const url = 'ws://127.0.0.1:8080/realworld/' + 'Realworld/Upload'
  const ws = new WebSocket(url)
  ws.onopen = () => {
    lines.forEach(line => ws.send(JSON.stringify({ line: line })))
    ws.send(JSON.stringify({ done: true }))
  }
  ws.onmessage = ev => {
    const res = JSON.parse(ev.data)
    self.output = res.file
    self.output += `(uploaded total ${res.total_lines} lines)\n`
  }
  ws.onclose = ev => {
    self.output += `(ws close: ${ev.code} (clean=${ev.wasClean}))\n`
  }
  ws.onerror = ev => {
    self.output += `(ws error: ${ev})\n`
  }
}
