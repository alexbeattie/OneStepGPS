<template>
  <div class="fixed top-4 right-4 bg-white rounded-lg shadow-lg w-80 z-50 map-settings">
  
  </div>
</template>

<script>
import { mapState } from 'vuex';

export default {
  name: 'MapSettings',
  
  props: {
    modelValue: {
      type: Object,
      required: true
    }
  },

  emits: ['update:modelValue', 'close'],

  data() {
    return {
      // localSettings: { ...this.modelValue }
    }
  },

  computed: {
    ...mapState('mapSettings', {
      // mapSettings: state => state.settings
    })
  },

  watch: {
    modelValue: {
      // handler(newValue) {
      //   this.localSettings = { ...newValue };
      // },
      deep: true
    }
  },

  methods: {
    handleChange() {
      // Emit the entire settings object when any setting changes
      // this.$emit('update:modelValue', { ...this.localSettings });
      
      // Also update Vuex store
      // this.$store.dispatch('mapSettings/updateSetting', {
      //   key: 'showDevices',
      //   value: this.localSettings.showDevices
      // });
    }
  }
}
</script>