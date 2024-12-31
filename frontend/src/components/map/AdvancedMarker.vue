<template>
  <div>
    <!-- Route History Modal -->
    <DriveStopModal
      v-if="showRouteModal"
      :device-id="device.device_id"
      @close="handleRouteModalClose"
    />
  </div>
</template>

<script>
import DriveStopModal from '@/components/device/DriveStopModal.vue';


export default {
  name: 'AdvancedMarker',
  
  components: {
    DriveStopModal
  },

  props: {
    position: {
      type: Object,
      required: true,
      validator: function(value) {
        return value.lat !== undefined && value.lng !== undefined;
      }
    },
    title: {
      type: String,
      default: ''
    },
    device: {
      type: Object,
      required: true,
      validator: function(value) {
        return value.latest_device_point &&
          value.latest_device_point.speed !== undefined &&
          value.latest_device_point.angle !== undefined;
      }
    }
  },

  inject: ['getDeviceStatus', 'mapSettings'],

  data() {
    return {
      marker: null,
      showRouteModal: false
    };
  },

  watch: {
    position: {
      handler(newPos) {
        if (this.marker) {
          this.marker.setPosition(newPos);
        }
      },
      deep: true
    },
    'mapSettings.showDevices': {
      handler(newValue) {
        if (this.marker) {
          this.marker.setVisible(newValue);
        }
      },
      immediate: true
    }
  },

  mounted() {
    this.$nextTick(() => {
      this.initializeMarker();
    });
  },

  methods: {
    initializeMarker() {
      if (!window.google?.maps || !this.$parent.map) {
        console.warn('Google Maps not ready');
        return;
      }

      this.marker = new window.google.maps.Marker({
        position: this.position,
        title: this.title,
        map: this.$parent.map,
        icon: this.getMarkerIcon(),
        visible: this.mapSettings.showDevices
      });

      this.$emit('register-marker', this.device.device_id, this.marker);
      
      this.marker.addListener('click', (e) => {
        if (e && e.stop) {
          e.stop();
        }
        this.$emit('click', this.device);
      });
    },

    getMarkerIcon() {
      const speed = this.device.latest_device_point.speed;
      const heading = this.device.latest_device_point.angle;
      return {
        path: "M 12,0 L 24,24 L 12,16 L 0,24 Z",
        fillColor: speed > 0 ? "#00FF00" : "#FF0000",
        fillOpacity: 1,
        strokeWeight: 1,
        strokeColor: "#000000",
        rotation: heading || 0,
        scale: 1,
        anchor: { x: 12, y: 12 }
      };
    },

    handleRouteModalClose() {
      this.showRouteModal = false;
    },

    handleHistoryClick() {
      if (this.marker) {
        this.marker.setAnimation(window.google.maps.Animation.BOUNCE);
        setTimeout(() => {
          this.marker.setAnimation(null);
        }, 1500);
      }
      this.showRouteModal = true;
    }
  },

  beforeUnmount() {
    if (this.marker) {
      this.marker.setMap(null);
      this.$emit('unregister-marker', this.device.device_id);
    }
  }
};
</script>