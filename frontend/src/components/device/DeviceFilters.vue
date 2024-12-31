<template>
  <div class="space-y-4">
    <!-- Sort Options -->
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-2">Sort by</label>
      <select 
        v-model="localValue.sortBy"
        class="w-full border border-gray-300 rounded-md shadow-sm p-2"
        @change="emitUpdate"
      >
        <option v-for="option in sortOptions" 
                :key="option.value" 
                :value="option.value"
        >
          {{ option.label }}
        </option>
      </select>
    </div>

    <!-- Status Filters -->
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-2">Status</label>
      <div class="space-y-2">
        <label class="flex items-center">
          <input 
            type="checkbox" 
            v-model="localValue.filters.driving"
            @change="emitUpdate"
            class="rounded text-blue-600"
          >
          <span class="ml-2">Driving</span>
        </label>
        <label class="flex items-center">
          <input 
            type="checkbox" 
            v-model="localValue.filters.stopped"
            @change="emitUpdate"
            class="rounded text-blue-600"
          >
          <span class="ml-2">Stopped</span>
        </label>
        <label class="flex items-center">
          <input 
            type="checkbox" 
            v-model="localValue.filters.ignitionOn"
            @change="emitUpdate"
            class="rounded text-blue-600"
          >
          <span class="ml-2">Ignition On</span>
        </label>
      </div>
    </div>

    <!-- Alert Filters -->
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-2">Alerts</label>
      <div class="space-y-2">
        <label class="flex items-center">
          <input 
            type="checkbox" 
            v-model="localValue.filters.overSpeedLimit"
            @change="emitUpdate"
            class="rounded text-blue-600"
          >
          <span class="ml-2">Over Speed Limit</span>
        </label>
        <label class="flex items-center">
          <input 
            type="checkbox" 
            v-model="localValue.filters.lowBattery"
            @change="emitUpdate"
            class="rounded text-blue-600"
          >
          <span class="ml-2">Low Battery</span>
        </label>
        <label class="flex items-center">
          <input 
            type="checkbox" 
            v-model="localValue.filters.poorSignal"
            @change="emitUpdate"
            class="rounded text-blue-600"
          >
          <span class="ml-2">Poor Signal</span>
        </label>
      </div>
    </div>

    <!-- Online Status -->
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-2">Connectivity</label>
      <div class="space-y-2">
        <label class="flex items-center">
          <input 
            type="checkbox" 
            v-model="localValue.filters.offline"
            @change="emitUpdate"
            class="rounded text-blue-600"
          >
          <span class="ml-2">Show Offline Devices</span>
        </label>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'DeviceFilters',
  
  props: {
    modelValue: {
      type: Object,
      required: true
    }
  },

  emits: ['update:modelValue'],
  
  data() {
    return {
      localValue: { ...this.modelValue }, // Create a local copy of the prop
      sortOptions: [
        { value: 'active', label: 'Most Active' },
        { value: 'inactive', label: 'Least Active' },
        { value: 'speed', label: 'Current Speed' },
        { value: 'distance', label: 'Trip Distance' },
        { value: 'voltage', label: 'Battery Voltage' },
        { value: 'signal', label: 'Signal Strength' },
        { value: 'name', label: 'Device Name' }
      ]
    };
  },

  watch: {
    modelValue: {
      handler(newValue) {
        this.localValue = { ...newValue }; // Sync with the parent when the prop changes
      },
      deep: true
    }
  },

  methods: {
    emitUpdate() {
      this.$emit('update:modelValue', { ...this.localValue }); // Emit the updated value
    }
  }
};
</script>
