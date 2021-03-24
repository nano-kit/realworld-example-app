const Publish = {
  data() {
    return {
      topic: '',
      test: '',
      err: ''
    }
  },
  methods: {
    call() {
      doPublish(this)
    }
  }
}

Vue.createApp(Publish).mount('#publish')

function doPublish(self) {
  postData('Clubhouse/Publish',
    {
      publish_note: {
        topic: self.topic,
        text: self.text
      }
    })
    .then(res => {
      self.err = ''
    })
    .catch(err => {
      self.err = err
    })
}
