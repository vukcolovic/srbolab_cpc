import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import axios from 'axios';
import VueAxios from 'vue-axios';
import store from './store/index.js';
import VueTableLite from 'vue3-table-lite'
import Toast from "vue-toastification";
import "vue-toastification/dist/index.css";
import vSelect from "vue-select";
import "vue-select/dist/vue-select.css";


axios.defaults.baseURL = process.env.VUE_APP_API_BASE_URL;

const app = createApp(App);
app.use(router);
app.use(store);
app.use(VueAxios, axios);
const options = {
    timeout: 2000
};

app.use(Toast, options);
app.component('VueTable', VueTableLite);
app.component('v-select', vSelect);

axios.interceptors.request.use(request => {
    request.headers.Authorization = store.getters.token;
    return request;
})

axios.interceptors.response.use(function (response) {
    if (response.config.url !== "/api/login" && response.headers.authorization) {
        store.dispatch('setTokenAction', response.headers.authorization);
    }

    return response
}, function(error) {
    if (error.response != null && error.response.status === 403) {
        store.dispatch('setTokenAction', "");
        localStorage.removeItem('roles');
    }
    return Promise.reject(error.response)
})

//create v-can directive
app.directive('can', (el, binding) => {
    if (!JSON.parse(localStorage.getItem('roles')).includes(binding.value)) {
        el.style.display = "none";
    }
})

app.mount('#app')
