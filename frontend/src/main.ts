import { createApp } from 'vue';
import { createPinia } from 'pinia';
import naive from 'naive-ui';

import App from './App.vue';
import router from './router';

import './assets/main.css';

const app = createApp(App);

app.use(naive);
app.use(createPinia());
app.use(router);

app.mount('#app');
