import { createApp } from 'vue';
import App from './App.vue';
import store from './store';
import "core-js/stable";
import "regenerator-runtime/runtime";
import './assets/main.css';

const app = createApp(App);

// Initialize store settings
Promise.all([
  store.dispatch('mapSettings/loadSettings'),
  store.dispatch('devices/loadGeocodeCache')
]).catch(console.error);

app.use(store);
app.mount('#app');

// Persist cache before unloading
window.addEventListener('beforeunload', () => {
  store.dispatch('devices/persistGeocodeCache');
});