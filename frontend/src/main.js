import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import axios from 'axios';
import VueAxios from 'vue-axios';
import store from './store/index.js';
import VueTableLite from 'vue3-table-lite'


axios.defaults.baseURL = process.env.VUE_APP_API_BASE_URL;

const app = createApp(App).use(router);
app.use(router);
app.use(store);
app.use(VueAxios, axios);
app.component('VueTable', VueTableLite);

createApp(App).use(router).mount('#app')
