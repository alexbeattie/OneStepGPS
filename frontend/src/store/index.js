import { createStore } from 'vuex';
import mapSettings from './modules/mapSettings';
import devices from './modules/devices';

// The createStore function sets up our central Vuex store.
export default createStore({
  modules: {
    mapSettings,  // Handles map-related settings
    devices       // Handles device-related data
  }
});