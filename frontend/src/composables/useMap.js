// composables/useMap.js
import { ref, onMounted } from 'vue';

export function useMap(apiKey) {
  const map = ref(null);
  const markerRefs = ref(new Map());

  const initializeMap = (mapInstance) => {
    map.value = mapInstance;
    map.value.setOptions({
      mapTypeControl: false,
      disableDefaultUI: true,
      gestureHandling: "greedy",
      streetViewControl: false,
      zoomControl: false,
      fullscreenControl: false,
      scaleControl: false,
      rotateControl: false,
      clickableIcons: false,
      keyboardShortcuts: false,
    });
  };

  const registerMarker = (deviceId, marker) => {
    markerRefs.value.set(deviceId, marker);
  };

  const unregisterMarker = (deviceId) => {
    markerRefs.value.delete(deviceId);
  };

  return {
    map,
    markerRefs,
    initializeMap,
    registerMarker,
    unregisterMarker
  };
}