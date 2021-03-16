const Nav = {
    data() {
      return {
        isActive: false
      }
    },
    methods: {
      toggle() {
        this.isActive = !this.isActive
      }
    }
  }

  Vue.createApp(Nav).mount('#nav')
