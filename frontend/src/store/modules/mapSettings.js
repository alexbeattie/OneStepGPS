// store/modules/mapSettings.js
export default {
  namespaced: true,
  state: () => ({
    settings: {
      mapType: 'roadmap',
      showTraffic: false,
      showDevices: true, // Make sure this is true
      showWeather: false,
      showDriverNames: false,
      showMoving: true,
      showStopped: true,
      showOffline: true

    },
    isMapInitialized: false
  }),

  mutations: {
    SET_MAP_INITIALIZED(state, value) {
      state.isMapInitialized = value;
    },
    UPDATE_SETTING(state, { key, value }) {
      state.settings[key] = value;
    },
    LOAD_SETTINGS(state, settings) {
      state.settings = { ...state.settings, ...settings };
    },
    updateSetting(state, { key, value }) {
      state.settings[key] = value;
    }

  },

  actions: {
    initializeMap({ commit }) {
      commit('SET_MAP_INITIALIZED', true);
    },
    loadSettings({ commit }) {
      try {
        const savedSettings = localStorage.getItem('mapSettings');
        if (savedSettings) {
          commit('LOAD_SETTINGS', JSON.parse(savedSettings));
        }
      } catch (error) {
        console.error('Error loading map settings:', error);
      }
    },
    // updateSetting({ commit, state }, { key, value }) {
    //   commit('UPDATE_SETTING', { key, value });
    //   // Save to localStorage
    //   localStorage.setItem('mapSettings', JSON.stringify(state.settings));
    // },
    updateSetting({ commit }, payload) {
      commit('updateSetting', payload);
    }

  }
};