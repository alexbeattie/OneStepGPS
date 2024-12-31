// store/modules/devices.js
import { apiService } from '@/services/api';
const state = () => ({
  devices: [],
  selectedDevice: null,
  addressMap: {},
  geocodingCache: {}, // New cache object
  error: null,
  odometerMap: {}, // Add this
  loading: false,
  driveStopRoutes: {},

});

const mutations = {
  SET_DEVICES(state, payload) {
    console.log('SET_DEVICES mutation called with:', payload);
    try {
      // Ensure we always have an array
      const deviceArray = Array.isArray(payload) ? payload : [payload];
      console.log('Processing devices array:', deviceArray.length, 'devices');
      state.devices = deviceArray;
      console.log('Devices set in state:', state.devices.length);
    } catch (error) {
      console.error('Error in SET_DEVICES:', error);
      state.devices = [];
    }
  },
  SET_DRIVE_STOP_ROUTE(state, { deviceId, routeData }) {
    state.driveStopRoutes[deviceId] = routeData;
  },

  SET_SELECTED_DEVICE(state, deviceId) {
    state.selectedDevice = deviceId;
  },
  SET_LOADING(state, loading) {
    state.loading = loading;
  },
  SET_ERROR(state, error) {
    state.error = error;
  },
  UPDATE_DEVICE(state, updatedDevice) {
    const index = state.devices.findIndex(d => d.device_id === updatedDevice.device_id);
    if (index !== -1) {
      state.devices.splice(index, 1, updatedDevice);
    } else {
      state.devices.push(updatedDevice);
    }
  },
  UPDATE_ADDRESS(state, { deviceId, address }) {
    state.addressMap[deviceId] = address;
  },

  // New mutation to update geocoding cache
  UPDATE_GEOCODING_CACHE(state, { coordinates, address, timestamp }) {
    const key = `${coordinates.lat},${coordinates.lng}`;
    state.geocodingCache[key] = {
      address,
      timestamp
    };
    // Optional: Limit cache size
    if (Object.keys(state.geocodingCache).length > 100) {
      const oldestKey = Object.keys(state.geocodingCache)
        .reduce((oldest, current) =>
          state.geocodingCache[current].timestamp <
            state.geocodingCache[oldest].timestamp ? current : oldest
        );
      delete state.geocodingCache[oldestKey];
    }
  },
  SET_GEOCODING_CACHE(state, cache) {
    state.geocodingCache = cache;
  },
  UPDATE_DEVICE_ODOMETER(state, { deviceId, odometer }) {
    state.odometerMap[deviceId] = odometer;
  },
};

const actions = {
  async fetchDevices({ commit, dispatch }) {
    commit('SET_LOADING', true);
    try {
      console.log('Fetching devices from Go server...');
     
      const data = await apiService.getDevices();


      let devicesToProcess = [];

      // Handle different data structures
      if (data && !Array.isArray(data) && data.device_id) {
        devicesToProcess = [data];
        commit('SET_DEVICES', [data]);
      } else if (Array.isArray(data)) {
        devicesToProcess = data;
        commit('SET_DEVICES', data);
      } else if (data && data.devices) {
        devicesToProcess = data.devices;
        commit('SET_DEVICES', data.devices);
      } else {
        commit('SET_DEVICES', []);
        devicesToProcess = [];
      }

      // Process geocoding for each device
      for (const device of devicesToProcess) {
        if (device?.latest_device_point?.lat && device?.latest_device_point?.lng) {
          // Add a small delay between geocoding requests
          await new Promise(resolve => setTimeout(resolve, 100));
          dispatch('updateDeviceAddress', device);
        }
      }

    } catch (error) {
      console.error('Error fetching devices:', error);
      commit('SET_ERROR', error.message);
    } finally {
      commit('SET_LOADING', false);
    }
    
  },
  

  selectDevice({ commit }, deviceId) {
    commit('SET_SELECTED_DEVICE', deviceId);
  },

  updateAddress({ commit }, { deviceId, address }) {
    commit('UPDATE_ADDRESS', { deviceId, address });
  },

  async fetchDriveStopRoute({ commit }, { deviceId, fromDate, toDate, stopDuration }) {
    try {
      const data = await apiService.getDriveStopRoute(deviceId, fromDate, toDate, stopDuration);
      commit('SET_DRIVE_STOP_ROUTE', { deviceId, routeData: data });
      return data;
    } catch (error) {
      console.error('Error fetching drive-stop route:', error);
      throw error;
    }
  },
  // Add new action for geocoding
  async updateDeviceAddress({ commit, state }, device) {
    const lat = device.latest_device_point.lat;
    const lng = device.latest_device_point.lng;
    const cacheKey = `${lat},${lng}`;
    const now = Date.now();

    // Check cache first
    const cachedResult = state.geocodingCache[cacheKey];
    if (cachedResult) {
      // Check if cache is fresh (e.g., less than 24 hours old)
      const twentyFourHoursAgo = now - (24 * 60 * 60 * 1000);
      if (cachedResult.timestamp > twentyFourHoursAgo) {
        commit('UPDATE_ADDRESS', {
          deviceId: device.device_id,
          address: cachedResult.address
        });
        return;
      }
    }

    try {
      // Set a loading state
      commit('UPDATE_ADDRESS', {
        deviceId: device.device_id,
        address: 'Loading address...'
      });

      const response = await fetch(
        `https://maps.googleapis.com/maps/api/geocode/json?latlng=${lat},${lng}&key=${process.env.VUE_APP_GOOGLE_MAPS_API_KEY}`
      );

      if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);

      const data = await response.json();
      if (data.status === "OK" && data.results?.[0]) {
        const address = data.results[0].formatted_address;

        // Update both address map and geocoding cache
        commit('UPDATE_ADDRESS', {
          deviceId: device.device_id,
          address: address
        });

        commit('UPDATE_GEOCODING_CACHE', {
          coordinates: { lat, lng },
          address: address,
          timestamp: now
        });
      }
    } catch (error) {
      console.error('Geocoding error:', error);
      commit('UPDATE_ADDRESS', {
        deviceId: device.device_id,
        address: 'Address lookup failed'
      });
    }
  },

  // Optional: Method to persist cache to localStorage
  persistGeocodeCache({ state }) {
    try {
      localStorage.setItem('geocodingCache', JSON.stringify(state.geocodingCache));
    } catch (error) {
      console.error('Failed to persist geocoding cache', error);
    }
  },

  // Optional: Method to load cache from localStorage on app init
  loadGeocodeCache({ commit }) {
    try {
      const cachedData = localStorage.getItem('geocodingCache');
      if (cachedData) {
        const parsedCache = JSON.parse(cachedData);
        commit('SET_GEOCODING_CACHE', parsedCache);
      }
    } catch (error) {
      console.error('Failed to load geocoding cache', error);
    }
  },
  async fetchSingleDevice({ commit }, deviceId) {
    try {
      const device = await apiService.getDevice(deviceId);
      commit('UPDATE_DEVICE', device);
      return device;
    } catch (error) {
      console.error(`Error fetching device ${deviceId}:`, error);
      commit('SET_ERROR', error.message);
    }
  },

};



const getters = {
  hasDevices: (state) => {
    const result = state.devices && state.devices.length > 0;
    console.log('hasDevices getter:', {
      hasDevices: result,
      deviceCount: state.devices?.length
    });
    return result;
  },
  getDriveStopRoute: (state) => (deviceId) => {
    return state.driveStopRoutes[deviceId];
  },
  getDeviceById: (state) => (id) => {
    return state.devices.find(device => device.device_id === id);
  },
  getAddressForDevice: (state) => (id) => {
    return state.addressMap[id] || 'Loading address...';
  },
  getDeviceOdometer: (state) => (id) => {
    const odometer = state.odometerMap[id];
    if (!odometer) return { value: 'N/A', unit: 'mi' };
    return {
      value: odometer.value,
      unit: odometer.unit,
      display: odometer.display
    };
  },
  getDeviceEngineHours: (state) => (id) => {
    const device = state.devices.find(d => d.device_id === id);
    return device?.latest_device_point?.params?.engine_hours || null;
  },

  isDeviceSelected: (state) => (id) => {
    return state.selectedDevice === id;
  }
};

export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters
};