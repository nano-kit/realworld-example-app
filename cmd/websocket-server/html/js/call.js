const Call = {
  data() {
    return {
      name: '',
      age: 0,
      msg: '',
      err: ''
    }
  },
  methods: {
    call() {
      doCall(this)
    }
  }
}

Vue.createApp(Call).mount('#call')

function doCall(self) {
  postData('Realworld/Call', { name: self.name, age: self.age })
    .then(res => {
      self.msg = res.msg
      self.err = ''
    }).catch(err => {
      self.err = err
      self.msg = ''
    })
}
