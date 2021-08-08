Vue.component('transactions-table', {
    data: {
    },
    template: `
<table role="grid">
  <thead>
    <tr>
      <th scope="col">#</th>
      <th scope="col">From Address</th>
      <th scope="col">To Address</th>
      <th scope="col">Asset</th>
      <th scope="col">Amount</th>
    </tr>
  </thead>
  <tbody>
    <tr v-for="item in items">
      <th scope="row">1</th>
      <td>{{ item.from_address }}</td>
      <td>{{ item.to_address }}</td>
      <td>{{ item.asset }}</td>
      <td>{{ item.amount }}</td>
    </tr>
  </tbody>
</table>
  `,
    computed: {
        items: []
    }
})

Vue.component('add-hd-wallet', {
    props: [
        "x_pub"
    ],
    template: `
    <form>
        <label for="xpub">XPub</label>
        <input v-model="x_pub" placeholder="XPub" required>
        <button type="submit" onclick="event.preventDefault()" v-on:click="addHDWallet({'x_pub': x_pub })">Submit</button>
    </form>
    `,
    methods: {
        addHDWallet: function (data) {
            api.addHDWallet(data)
        }
    }
})

Vue.component('hd-wallet', {
    props: [
        "x_pub"
    ],
    data () {
        return {
            items: [],
        }
    },
    template: `
    <table role="grid">
      <thead>
        <tr>
          <th scope="col">#</th>
          <th scope="col">XPub</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in items">
          <th scope="row">{{ item.id }}</th>
          <td>{{ item.x_pub }}</td>
        </tr>
      </tbody>
    </table>
    `,
    methods: {
        getHDWalletList: function () {
            api.getHDWalletList()
                .then(response => {
                    this.items =  response
                })
        }
    },
    created () {
        this.getHDWalletList()
    },
    watch: {
        '$route': 'fetchData'
    }
})

const routes = [
    { path: '/hd-wallet/add', component: Vue.component('add-hd-wallet') },
    { path: '/hd-wallet', component: Vue.component('hd-wallet') },
]

const router = new VueRouter({
    router
})

new Vue({
    router,
    delimiters: ['${', '}'],
    data: {
        "email": "",
        "password": "",
    },
    methods: {
        auth: function (data) {
            api.auth(data)
                .then(response => {
                    console.log(response)
                    localStorage.setItem('auth_token', response.token)
                    window.location.href = '/admin/dashboard';
                })
                .catch(err => console.error(err));
        }
    }
}).$mount('#app')