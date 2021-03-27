const Nav = {
  data() {
    return {
      isActive: false,
      account: null
    }
  },
  created() {
    getAccessToken()
      .then(token => {
        const t = checkAccessToken(token)
        if (!t.metadata.Name) {
          t.metadata.Name = t.sub
        }
        this.account = t.metadata
      })
      .catch(_ => this.account = null)
  },
  methods: {
    toggle() {
      this.isActive = !this.isActive
    }
  }
}

Vue.createApp(Nav).mount('#nav')
