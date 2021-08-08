"use strict";

(function() {
    var xfetch = function(resource, init) {
        init = init || {}
        if (['post', 'put', 'delete'].indexOf(init.method) !== -1) {
            init['headers'] = init['headers'] || {}
            init['headers']['x-requested-by'] = 'ethpayserver'
        }
        return fetch(resource, init)
    }
    var api = function(method, endpoint, data) {
        var headers = {
            'Content-Type': 'application/json',
            'Token': localStorage.getItem('auth_token'),
        }

        return xfetch(endpoint, {
            method: method,
            headers: headers,
            body: JSON.stringify(data),
        })
    }

    var json = function(res) {
        return res.json()
    }

    var param = function(query) {
        if (!query) return ''
        return '?' + Object.keys(query).map(function(key) {
            return encodeURIComponent(key) + '=' + encodeURIComponent(query[key])
        }).join('&')
    }
    window.api = {
        auth: function (data) {
            return api('post', '/api/v1/auth', data).then(json)
        },
        addHDWallet: function (data) {
            return api('post', '/api/v1/hd-wallet', data).then(json)
        },
        getHDWalletList: function () {
            return api('get', '/api/v1/hd-wallet').then(json)
        }
    }
})()