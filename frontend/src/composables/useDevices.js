// composables/useDevices.js
import { ref, onMounted, onBeforeUnmount } from 'vue';
import { useStore } from 'vuex';

export function useDevices() {
  const store = useStore();
  const pollingInterval = ref(null);

  const startPolling = (interval = 10000) => {
    pollingInterval.value = setInterval(async () => {
      try {
        await store.dispatch('devices/fetchDevices');
      } catch (error) {
        console.error('Polling error:', error);
      }
    }, interval);
  };

  const stopPolling = () => {
    if (pollingInterval.value) {
      clearInterval(pollingInterval.value);
      pollingInterval.value = null;
    }
  };

  onMounted(() => {
    startPolling();
  });

  onBeforeUnmount(() => {
    stopPolling();
  });

  return {
    startPolling,
    stopPolling
  };
}