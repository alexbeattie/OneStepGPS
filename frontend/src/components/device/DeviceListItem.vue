<template>
  <div
    class="p-3 rounded-lg cursor-pointer border transition-colors duration-150"
    :class="{
      'bg-blue-50 border-blue-500': isSelected,
      'hover:bg-gray-50 border-transparent': !isSelected,
    }"
    @click="$emit('click', device)"
  >
    <!-- Header: Name and Status -->
    <div class="flex justify-between items-center">
      <span class="font-medium">{{ metrics.displayName }}</span>
      <span class="px-2 py-1 rounded-full text-xs" :class="statusClass">
        {{ deviceStatus }}
      </span>
    </div>

    <!-- Device Info -->
    <div class="text-sm text-gray-500 mt-1 space-y-1">
      <!-- Location -->
      <div v-if="address">{{ address }}</div>

      <!-- Speed -->
      <div class="flex justify-between items-center">
        <span>Speed:</span>
        <span :class="{ 'text-red-600': metrics.isOverSpeedLimit }">
          {{ formatSpeed(metrics.currentSpeed) }}
          <span v-if="metrics.speedLimit" class="text-xs">
            (Limit: {{ formatSpeed(metrics.speedLimit) }})
          </span>
        </span>
      </div>

      <!-- Battery -->
      <div class="flex justify-between items-center">
        <span>Battery:</span>
        <span
          :class="{
            'text-red-600': metrics.batteryVoltage <= 12,
            'text-green-600': metrics.batteryVoltage > 13,
          }"
        >
          {{ metrics.batteryVoltage.toFixed(1) }}V
        </span>
      </div>

      <!-- Signal -->
      <div class="flex justify-between items-center">
        <span>GPS Signal:</span>
        <div class="flex items-center">
          <span
            :class="{
              'text-red-600': metrics.signalStrength < 8,
              'text-yellow-600':
                metrics.signalStrength >= 8 && metrics.signalStrength < 12,
              'text-green-600': metrics.signalStrength >= 12,
            }"
          >
            {{ metrics.signalStrength }} satellites
          </span>
        </div>
      </div>

      <!-- Status Duration -->
      <div class="flex justify-between items-center text-xs">
        <span>{{ deviceStatus }} for:</span>
        <span>{{ metrics.driveTimeFormatted }}</span>
      </div>

      <!-- Last Update -->
      <div class="text-xs text-right text-gray-400">
        Last seen: {{ formatTimeAgo(metrics.lastUpdated) }}
      </div>
    </div>
  </div>
</template>

<script>
import {
  getDeviceMetrics,
  getDeviceStatusClass,
} from "@/services/deviceMetrics";

export default {
  name: "DeviceListItem",

  props: {
    device: {
      type: Object,
      required: true,
    },
    isSelected: {
      type: Boolean,
      default: false,
    },
    address: {
      type: String,
      default: "Loading address...",
    },
  },

  emits: ["click"],

  computed: {
    metrics() {
      return getDeviceMetrics(this.device);
    },

    statusClass() {
      return getDeviceStatusClass(this.metrics);
    },

    deviceStatus() {
      if (!this.metrics.isOnline) return "Offline";
      if (this.metrics.signalStrength < 8) return "No Signal";
      return (
        this.metrics.driveStatus.charAt(0).toUpperCase() +
        this.metrics.driveStatus.slice(1)
      );
    },
  },

  methods: {
    formatSpeed(speed) {
      return `${speed} km/h`;
    },

    formatTimeAgo(timestamp) {
      if (!timestamp) return "Unknown";

      const date = new Date(timestamp);
      const now = new Date();
      const seconds = Math.floor((now - date) / 1000);

      if (seconds < 60) return "just now";
      if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`;
      if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`;
      return `${Math.floor(seconds / 86400)}d ago`;
    },
  },
};
</script>
