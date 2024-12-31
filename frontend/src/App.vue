<template>
  <div class="flex h-screen">
    <!-- Sidebar -->
    <button
      @click="toggleSidebar"
      class="fixed top-4 left-4 z-50 bg-white p-2 rounded-md shadow-lg md:hidden"
    >
      <span v-if="isSidebarOpen">✕</span>
      <span v-else>☰</span>
    </button>

    <div
      :class="{
        'fixed inset-y-0 left-0 transform': true,
        'translate-x-0': isSidebarOpen,
        '-translate-x-full': !isSidebarOpen,
        'md:relative md:translate-x-0': true,
      }"
      class="w-80 bg-white border-r border-gray-200 overflow-y-auto transition-transform duration-300 ease-in-out z-40"
    >
      <div class="p-4">
        <h2 class="text-lg font-semibold mb-4">Devices</h2>

        <!-- Filter Controls -->
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-2"
            >Sort by</label
          >
          <select
            v-model="sortBy"
            class="w-full border border-gray-300 rounded-md shadow-sm p-2"
          >
            <option value="active">Most Active</option>
            <option value="inactive">Least Active</option>
            <option value="speed">Speed</option>
            <option value="distance">Distance</option>
            <option value="signal">Signal Strength</option>
          </select>
        </div>

        <div class="mb-4 space-y-4">
          <!-- Map Settings -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2"
              >Map Settings</label
            >
            <div class="space-y-2">
              <!-- Map Type -->
              <div class="mb-2">
                <select
                  v-model="mapSettings.mapType"
                  class="w-full border border-gray-300 rounded-md shadow-sm p-2"
                  @change="
                    updateSettings({
                      key: 'mapType',
                      value: mapSettings.mapType,
                    })
                  "
                >
                  <option value="roadmap">Road Map</option>
                  <option value="satellite">Satellite</option>
                  <option value="hybrid">Hybrid</option>
                  <option value="terrain">Terrain</option>
                </select>
              </div>
            </div>
          </div>

          <!-- Device Display Settings -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2"
              >Device Display</label
            >
            <div class="space-y-2">
              <label
                class="flex items-center justify-between p-2 hover:bg-gray-50 rounded"
              >
                <span class="text-sm text-gray-600">Show Devices on Map</span>
                <input
                  type="checkbox"
                  v-model="mapSettings.showDevices"
                  class="rounded text-blue-600"
                  @change="
                    updateSettings({
                      key: 'showDevices',
                      value: mapSettings.showDevices,
                    })
                  "
                />
              </label>

              <!-- Device Status Settings -->
              <label
                class="flex items-center justify-between p-2 hover:bg-gray-50 rounded"
              >
                <span class="text-sm text-gray-600">Moving Devices</span>
                <input
                  type="checkbox"
                  v-model="mapSettings.showMoving"
                  class="rounded text-blue-600"
                  @change="
                    updateSettings({
                      key: 'showMoving',
                      value: mapSettings.showMoving,
                    })
                  "
                />
              </label>

              <label
                class="flex items-center justify-between p-2 hover:bg-gray-50 rounded"
              >
                <span class="text-sm text-gray-600">Stopped Devices</span>
                <input
                  type="checkbox"
                  v-model="mapSettings.showStopped"
                  class="rounded text-blue-600"
                  @change="
                    updateSettings({
                      key: 'showStopped',
                      value: mapSettings.showStopped,
                    })
                  "
                />
              </label>

              <label
                class="flex items-center justify-between p-2 hover:bg-gray-50 rounded"
              >
                <span class="text-sm text-gray-600">Offline Devices</span>
                <input
                  type="checkbox"
                  v-model="mapSettings.showOffline"
                  class="rounded text-blue-600"
                  @change="
                    updateSettings({
                      key: 'showOffline',
                      value: mapSettings.showOffline,
                    })
                  "
                />
              </label>
            </div>
          </div>
        </div>
        <!-- Device List -->
        <div class="space-y-2">
          <!-- Active/Moving Devices -->
          <div v-if="deviceGroups.active?.length" class="mb-4">
            <h3 class="text-sm font-medium text-gray-500 px-3 mb-2">
              Moving ({{ deviceGroups.active.length }})
            </h3>
            <div class="bg-green-50 rounded-lg p-2 space-y-2">
              <div
                v-for="device in deviceGroups.active"
                :key="device.device_id"
                @click="handleSidebarDeviceClick(device)"
                class="p-3 rounded-lg cursor-pointer border"
                :class="{
                  'bg-blue-50 border-blue-500': isDeviceSelected(
                    device.device_id
                  ),
                  'hover:bg-green-100 border-transparent': !isDeviceSelected(
                    device.device_id
                  ),
                }"
              >
                <div class="flex justify-between items-center">
                  <span class="font-medium">{{ device.display_name }}</span>
                  <span
                    class="px-2 py-1 rounded-full text-xs"
                    :class="getDeviceStatus(device).className"
                  >
                    {{ getDeviceStatus(device).status }}
                  </span>
                </div>
                <div class="text-sm text-gray-500 mt-1">
                  <div class="truncate max-w-xs">
                    {{ getAddressForDevice(device.device_id) }}
                  </div>
                  <div class="truncate max-w-xs">
                    Last seen: {{ getDeviceStatus(device).time }}
                  </div>
                  <div class="flex items-center">
                    <span class="text-xs">
                      Speed: {{ kmToMph(device.latest_device_point.speed) }}
                    </span>
                  </div>
                  <div class="flex items-center space-x-2">
                    <template
                      v-if="device?.latest_device_point?.params?.gpslev"
                    >
                      <div class="flex items-center">
                        <span class="text-xs"
                          >GPS:
                          {{ device.latest_device_point.params.gpslev }}</span
                        >
                        <div class="ml-1 w-4 h-2.5 bg-gray-200 rounded">
                          <div
                            class="h-full rounded"
                            :class="{
                              'bg-green-500':
                                parseInt(
                                  device.latest_device_point.params.gpslev
                                ) >= 8,
                              'bg-yellow-500':
                                parseInt(
                                  device.latest_device_point.params.gpslev
                                ) >= 4 &&
                                parseInt(
                                  device.latest_device_point.params.gpslev
                                ) < 8,
                              'bg-red-500':
                                parseInt(
                                  device.latest_device_point.params.gpslev
                                ) < 4,
                            }"
                            :style="{
                              width: `${
                                (parseInt(
                                  device.latest_device_point.params.gpslev
                                ) /
                                  12) *
                                100
                              }%`,
                            }"
                          ></div>
                        </div>
                      </div>
                    </template>
                  </div>
                  <!-- Add this inside the device info div where you show GPS and other metrics -->
                  <div class="flex items-center space-x-2">
                    <span class="text-xs">Battery:</span>
                    <div class="flex items-center">
                      <div class="relative inline-flex">
                        <span
                          class="text-xs"
                          :class="{
                            'text-red-600':
                              device.latest_device_point?.params
                                ?.obd_battery_voltage <= 12,
                            'text-green-600':
                              device.latest_device_point?.params
                                ?.obd_battery_voltage > 12,
                          }"
                        >
                          {{
                            device.latest_device_point?.params
                              ?.obd_battery_voltage || "N/A"
                          }}
                        </span>
                        <!-- Battery Icon -->
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          class="h-4 w-4 ml-1"
                          :class="{
                            'text-red-500':
                              device.latest_device_point?.params
                                ?.obd_battery_voltage <= 12,
                            'text-yellow-500':
                              device.latest_device_point?.params
                                ?.obd_battery_voltage <= 12.5,
                            'text-green-500':
                              device.latest_device_point?.params
                                ?.obd_battery_voltage > 12.5,
                          }"
                          fill="none"
                          viewBox="0 0 24 24"
                          stroke="currentColor"
                        >
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M3 5h14a2 2 0 012 2v10a2 2 0 01-2 2H3a2 2 0 01-2-2V7a2 2 0 012-2z"
                          />
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M20 10h2v4h-2"
                          />
                        </svg>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Stopped Devices -->
          <div v-if="deviceGroups.stopped?.length" class="mb-4">
            <h3 class="text-sm font-medium text-gray-500 px-3 mb-2">
              Stopped ({{ deviceGroups.stopped.length }})
            </h3>
            <div class="bg-gray-50 rounded-lg p-2 space-y-2">
              <div
                v-for="device in deviceGroups.stopped"
                :key="device.device_id"
                @click="handleSidebarDeviceClick(device)"
                class="p-3 rounded-lg cursor-pointer border"
                :class="{
                  'bg-blue-50 border-blue-500': isDeviceSelected(
                    device.device_id
                  ),
                  'hover:bg-gray-100 border-transparent': !isDeviceSelected(
                    device.device_id
                  ),
                }"
              >
                <div class="flex justify-between items-center">
                  <span class="font-medium">{{ device.display_name }}</span>
                  <span
                    class="px-2 py-1 rounded-full text-xs"
                    :class="getDeviceStatus(device).className"
                  >
                    {{ getDeviceStatus(device).status }}
                  </span>
                </div>
                <div class="text-sm text-gray-500 mt-1">
                  <div class="truncate max-w-xs">
                    {{ getAddressForDevice(device.device_id) }}
                  </div>

                  <div class="truncate max-w-xs">
                    Last seen: {{ getDeviceStatus(device).time }}
                  </div>

                  <div class="flex items-center space-x-2">
                    <template
                      v-if="device?.latest_device_point?.params?.gpslev"
                    >
                      <div class="flex items-center">
                        <span class="text-xs"
                          >GPS:
                          {{ device.latest_device_point.params.gpslev }}</span
                        >
                        <div class="ml-1 w-4 h-2.5 bg-gray-200 rounded">
                          <div
                            class="h-full rounded"
                            :class="{
                              'bg-green-500':
                                parseInt(
                                  device.latest_device_point.params.gpslev
                                ) >= 8,
                              'bg-yellow-500':
                                parseInt(
                                  device.latest_device_point.params.gpslev
                                ) >= 4 &&
                                parseInt(
                                  device.latest_device_point.params.gpslev
                                ) < 8,
                              'bg-red-500':
                                parseInt(
                                  device.latest_device_point.params.gpslev
                                ) < 4,
                            }"
                            :style="{
                              width: `${
                                (parseInt(
                                  device.latest_device_point.params.gpslev
                                ) /
                                  12) *
                                100
                              }%`,
                            }"
                          ></div>
                        </div>
                      </div>
                      <!-- Add this inside the device info div where you show GPS and other metrics -->
                    </template>
                  </div>
                  <div class="flex items-center space-x-2">
                    <span class="text-xs">Battery:</span>
                    <div class="flex items-center">
                      <div class="relative inline-flex">
                        <span
                          class="text-xs"
                          :class="{
                            'text-red-600':
                              device.latest_device_point?.params
                                ?.obd_battery_voltage <= 12,
                            'text-green-600':
                              device.latest_device_point?.params
                                ?.obd_battery_voltage > 12,
                          }"
                        >
                          {{
                            device.latest_device_point?.params
                              ?.obd_battery_voltage || "N/A"
                          }}
                        </span>
                        <!-- Battery Icon -->
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          class="h-4 w-4 ml-1"
                          :class="{
                            'text-red-500':
                              device.latest_device_point?.params
                                ?.obd_battery_voltage <= 12,
                            'text-yellow-500':
                              device.latest_device_point?.params
                                ?.obd_battery_voltage <= 12.5,
                            'text-green-500':
                              device.latest_device_point?.params
                                ?.obd_battery_voltage > 12.5,
                          }"
                          fill="none"
                          viewBox="0 0 24 24"
                          stroke="currentColor"
                        >
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M3 5h14a2 2 0 012 2v10a2 2 0 01-2 2H3a2 2 0 01-2-2V7a2 2 0 012-2z"
                          />
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M20 10h2v4h-2"
                          />
                        </svg>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Inactive/Offline Devices -->
          <div v-if="deviceGroups.inactive?.length" class="mb-4">
            <h3 class="text-sm font-medium text-gray-500 px-3 mb-2">
              Offline ({{ deviceGroups.inactive.length }})
            </h3>
            <div class="bg-red-50 rounded-lg p-2 space-y-2">
              <div
                v-for="device in deviceGroups.inactive"
                :key="device.device_id"
                @click="handleSidebarDeviceClick(device)"
                class="p-3 rounded-lg cursor-pointer border"
                :class="{
                  'bg-blue-50 border-blue-500': isDeviceSelected(
                    device.device_id
                  ),
                  'hover:bg-red-100 border-transparent': !isDeviceSelected(
                    device.device_id
                  ),
                }"
              >
                <div class="flex justify-between items-center">
                  <span class="font-medium">{{ device.display_name }}</span>
                  <span
                    class="px-2 py-1 rounded-full text-xs"
                    :class="getDeviceStatus(device).className"
                  >
                    {{ getDeviceStatus(device).status }}
                  </span>
                </div>

                <div class="text-sm text-gray-500 mt-1">
                  <div class="truncate max-w-xs">
                    {{ getAddressForDevice(device.device_id) }}
                  </div>

                  <div class="truncate max-w-xs">
                    Last seen: {{ getDeviceStatus(device).time }}
                  </div>

                  <div class="flex items-center space-x-2">
                    <template
                      v-if="device?.latest_device_point?.params?.gpslev"
                    >
                      <div class="flex items-center">
                        <span class="text-xs"
                          >GPS:
                          {{ device.latest_device_point.params.gpslev }}</span
                        >
                        <div class="ml-1 w-4 h-2.5 bg-gray-200 rounded">
                          <div
                            class="h-full rounded"
                            :class="{
                              'bg-green-500':
                                parseInt(
                                  device.latest_device_point.params.gpslev
                                ) >= 8,
                              'bg-yellow-500':
                                parseInt(
                                  device.latest_device_point.params.gpslev
                                ) >= 4 &&
                                parseInt(
                                  device.latest_device_point.params.gpslev
                                ) < 8,
                              'bg-red-500':
                                parseInt(
                                  device.latest_device_point.params.gpslev
                                ) < 4,
                            }"
                            :style="{
                              width: `${
                                (parseInt(
                                  device.latest_device_point.params.gpslev
                                ) /
                                  12) *
                                100
                              }%`,
                            }"
                          ></div>
                        </div>
                      </div>
                    </template>
                  </div>
                  <!-- Add this inside the device info div where you show GPS and other metrics -->
                  <div class="flex items-center space-x-2">
                    <span class="text-xs">Battery:</span>
                    <div class="flex items-center">
                      <div class="relative inline-flex">
                        <span
                          class="text-xs"
                          :class="{
                            'text-red-600':
                              device.latest_device_point?.params
                                ?.obd_battery_voltage <= 12,
                            'text-green-600':
                              device.latest_device_point?.params
                                ?.obd_battery_voltage > 12,
                          }"
                        >
                          {{
                            device.latest_device_point?.params
                              ?.obd_battery_voltage || "N/A"
                          }}V
                        </span>
                        <!-- Battery Icon -->
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          class="h-4 w-4 ml-1"
                          :class="{
                            'text-red-500':
                              device.latest_device_point?.params
                                ?.obd_battery_voltage <= 12,
                            'text-yellow-500':
                              device.latest_device_point?.params
                                ?.obd_battery_voltage <= 12.5,
                            'text-green-500':
                              device.latest_device_point?.params
                                ?.obd_battery_voltage > 12.5,
                          }"
                          fill="none"
                          viewBox="0 0 24 24"
                          stroke="currentColor"
                        >
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M3 5h14a2 2 0 012 2v10a2 2 0 01-2 2H3a2 2 0 01-2-2V7a2 2 0 012-2z"
                          />
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M20 10h2v4h-2"
                          />
                        </svg>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <!-- End of Device List -->
      </div>
    </div>
    <!-- Main Content -->

    <!-- Main Content -->
    <div class="flex-1 relative">
      <GoogleMap
        ref="map"
        :api-key="apiKey"
        :center="defaultCenter"
        :zoom="8"
        :mapTypeId="mapSettings.mapType"
        :options="{
          mapTypeControl: false,
          streetViewControl: false,
          fullscreenControl: false,
          zoomControl: false,
        }"
        class="w-full h-full"
        @ready="initializeLayers"
        @click="handleMapClick"
      >
        <div>
          <!-- Always show MapSettings -->
          <MapSettings
            :modelValue="mapSettings"
            @update:modelValue="updateSettings"
          />

          <!-- Only render markers if showDevices is true -->
          <template v-if="hasDevices && mapSettings.showDevices">
            <template v-for="device in devices" :key="device.device_id">
              <AdvancedMarker
                :position="getDevicePosition(device)"
                :title="device.display_name"
                :device="device"
                @click="handleMarkerClick(device)"
                @register-marker="registerMarker"
                @unregister-marker="unregisterMarker"
              />

              <CustomInfoWindow
                v-if="
                  isDeviceSelected(device.device_id) &&
                  clickedDeviceId === device.device_id
                "
                :device="device"
                :position="getDevicePosition(device)"
                @closeclick="handleInfoWindowClose"
              />
            </template>
          </template>
        </div>
      </GoogleMap>
    </div>
  </div>
</template>
<script>
// Import required dependencies and components
import { mapState, mapActions, mapGetters } from "vuex";
import { computed } from "vue";
import { GoogleMap } from "vue3-google-map";
import MapSettings from "./views/MapSettings.vue"; // Make sure path is correct
import AdvancedMarker from "./components/map/AdvancedMarker.vue";
import CustomInfoWindow from "./components/map/CustomInfoWindow.vue";
import { kmToMph } from "@/utils/unitConversion";
import {
  formatTimeAgo,
  getDeviceMetrics,
  getBatteryStatus,
} from "@/services/deviceMetrics";
import { getDeviceStatus } from "./utils/deviceStatus";

import _ from "lodash";

export default {
  name: "App",

  components: {
    GoogleMap,
    MapSettings,
    AdvancedMarker,
    CustomInfoWindow,
  },
  provide() {
    return {
      getDeviceStatus: this.getDeviceStatus,
      mapSettings: computed(() => this.mapSettings),
      getAddressForDevice: this.getAddressForDevice,
      kmToMph,
      getBatteryStatus, // Add this
      getDeviceOdometer: (id) =>
        this.$store.getters["devices/getDeviceOdometer"](id),
      getDeviceEngineHours: (id) =>
        this.$store.getters["devices/getDeviceEngineHours"](id),
    };
  },
  data() {
    return {
      // Keep only component-specific state here
      isSidebarOpen: false,
      sortBy: "active",
      sortOptions: [
        { value: "active", label: "Most Active" },
        { value: "inactive", label: "Least Active" },
        { value: "speed", label: "Current Speed" },
        { value: "distance", label: "Trip Distance" },
        { value: "voltage", label: "Battery Voltage" },
        { value: "signal", label: "Signal Strength" },
      ],

      filters: {
        ignitionOn: false,
        ignitionOff: false,
        overSpeedLimit: false,
        lowBattery: false, // < 12V
        poorSignal: false, // < 8 satellites
      },
      apiKey: process.env.VUE_APP_GOOGLE_MAPS_API_KEY,
      defaultCenter: { lat: 34.0522, lng: -118.2437 },
      showMapSettings: false,

 

      // Instance-specific controls
      timer: null,
      pollingInterval: null, // Add this line

      animationQueue: new Map(),
      markerRefs: new Map(),
      map: null, // Add this
      clickedDeviceId: null,
    };
  },
  // Map Vuex state and getters to component
  computed: {
    // Map Settings state
    ...mapState("mapSettings", {
      mapSettings: (state) => state.settings,
      isMapInitialized: (state) => state.isMapInitialized,
    }),
    ...mapGetters("devices", [
      "hasDevices",
      "getAddressForDevice",
      "isDeviceSelected",
    ]),

    // Devices state
    ...mapState("devices", {
      devices: (state) => {
        console.log("Devices from state:", state.devices);
        return state.devices;
      },
    }),

    // Filtered and sorted devices
    deviceStatuses() {
      return (device) => {
        const status = getDeviceStatus(device);
        return {
          isOffline: status.status === "Offline",
          isMoving: status.status === "Moving",
          isStopped: status.status === "Stopped",
          hasSignal: device?.latest_device_point?.params?.gpslev >= 4,
        };
      };
    },
    deviceGroups() {
      const filteredDevices = this.sortedAndFilteredDevices.filter((device) => {
        const status = this.getDeviceStatus(device);
        if (status.status === "Moving" && !this.mapSettings.showMoving)
          return false;
        if (status.status === "Stopped" && !this.mapSettings.showStopped)
          return false;
        if (status.status === "Offline" && !this.mapSettings.showOffline)
          return false;
        return true;
      });

      return _.groupBy(filteredDevices, (device) => {
        const status = this.getDeviceStatus(device);
        if (status.status === "Moving") return "active";
        if (status.status === "Stopped") return "stopped";
        return "inactive";
      });
    },
    sortedAndFilteredDevices() {
      let filteredDevices = [...this.devices];

      // Apply filters
      filteredDevices = filteredDevices.filter((device) => {
        const metrics = this.getDeviceMetrics(device);
        const status = this.getDeviceStatus(device);

        // First check device status filters using mapSettings
        switch (status.status) {
          case "Moving":
            if (!this.mapSettings.showMoving) return false;
            break;
          case "Stopped":
            if (!this.mapSettings.showStopped) return false;
            break;
          case "Offline":
            if (!this.mapSettings.showOffline) return false;
            break;
        }

        if (this.filters.ignitionOn && !metrics.ignitionOn) return false;
        if (this.filters.ignitionOff && metrics.ignitionOn) return false;
        if (this.filters.overSpeedLimit && !metrics.isOverSpeedLimit)
          return false;
        if (this.filters.lowBattery && metrics.batteryVoltage >= 12)
          return false;
        if (this.filters.poorSignal && metrics.signalStrength >= 8)
          return false;
        if (this.filters.driving && metrics.driveStatus !== "driving")
          return false;
        if (this.filters.stopped && metrics.driveStatus !== "stopped")
          return false;

        return true;
      });

      // Apply sorting
      filteredDevices.sort((a, b) => {
        // Define variables outside switch
        let timeA, timeB, statusA, statusB;

        switch (this.sortBy) {
          case "active": {
            statusA = this.getDeviceStatus(a);
            statusB = this.getDeviceStatus(b);

            if (statusA.status !== "Offline" && statusB.status === "Offline")
              return -1;
            if (statusA.status === "Offline" && statusB.status !== "Offline")
              return 1;

            timeA = new Date(a.latest_device_point?.dt_tracker || 0);
            timeB = new Date(b.latest_device_point?.dt_tracker || 0);
            return timeB - timeA;
          }

          case "inactive": {
            statusA = this.getDeviceStatus(a);
            statusB = this.getDeviceStatus(b);

            if (statusA.status !== "Offline" && statusB.status === "Offline")
              return 1;
            if (statusA.status === "Offline" && statusB.status !== "Offline")
              return -1;

            timeA = new Date(a.latest_device_point?.dt_tracker || 0);
            timeB = new Date(b.latest_device_point?.dt_tracker || 0);
            return timeA - timeB; // Changed from inactiveTimeA/B to timeA/B
          }

          case "speed": {
            const speedA = parseFloat(a.latest_device_point?.speed || 0);
            const speedB = parseFloat(b.latest_device_point?.speed || 0);
            return speedB - speedA;
          }

          case "distance": {
            const distanceA = parseFloat(
              a.latest_device_point?.params?.odo_trip_meter || 0
            );
            const distanceB = parseFloat(
              b.latest_device_point?.params?.odo_trip_meter || 0
            );
            return distanceB - distanceA;
          }

          case "signal": {
            const signalA = parseInt(
              a.latest_device_point?.params?.gpslev || 0
            );
            const signalB = parseInt(
              b.latest_device_point?.params?.gpslev || 0
            );
            return signalB - signalA;
          }

          default:
            return 0;
        }
      });

      return filteredDevices;
    },
  },
  watch: {
    "mapSettings.showMoving"() {
      this.updateMarkersVisibility();
    },
    "mapSettings.showStopped"() {
      this.updateMarkersVisibility();
    },
    "mapSettings.showOffline"() {
      this.updateMarkersVisibility();
    },
    "mapSettings.showDevices"() {
      this.updateMarkersVisibility();
    },
  },
  methods: {
    onMapReady(map) {
      this.map = map;
      // Set map options
      map.setOptions({
        mapTypeControl: false, // Disable the Map/Satellite toggle
        disableDefaultUI: true,
        gestureHandling: "greedy",
        streetViewControl: false,
        zoomControl: false,
        fullscreenControl: false,
        scaleControl: false,
        rotateControl: false,
        clickableIcons: false,
        keyboardShortcuts: false,
        draggable: true,
        scrollwheel: true,
        disableDoubleClickZoom: false,
        
      });
      map.setOptions({
        mapTypeControlOptions: {
          mapTypeIds: [], // Hides all map type options
        },
      });
    },

    toggleSidebar() {
      this.isSidebarOpen = !this.isSidebarOpen;
    },

    registerMarker(deviceId, marker) {
      this.markerRefs.set(deviceId, marker);
    },
    handleMapClick(event) {
      // Only close if clicking on the map itself, not on markers
      // console.log("Map clicked at:", event.latLng.toJSON());

      if (!event.placeId && !event.vertex && !event.edge) {
        this.clickedDeviceId = null;
        this.selectDevice(null);
        // Close sidebar on mobile
        if (window.innerWidth < 768) {
          this.isSidebarOpen = false;
        }
      }
    },
    unregisterMarker(deviceId) {
      this.markerRefs.delete(deviceId);
    },
    updateMarkersVisibility() {
      this.devices.forEach((device) => {
        const marker = this.markerRefs.get(device.device_id);
        if (!marker) return;

        const status = this.getDeviceStatus(device);
        let shouldShow = false;

        // First check if devices are shown at all
        if (this.mapSettings.showDevices) {
          // Then check individual status toggles
          if (status.status === "Moving" && this.mapSettings.showMoving)
            shouldShow = true;
          if (status.status === "Stopped" && this.mapSettings.showStopped)
            shouldShow = true;
          if (status.status === "Offline" && this.mapSettings.showOffline)
            shouldShow = true;
        }

        marker.setVisible(shouldShow);
      });
    },
    updateVisibility(visible) {
      if (this.marker) {
        this.marker.setVisible(visible);
      }
    },
    // Map Vuex actions to methods
    ...mapActions("mapSettings", [
      "initializeMap",
      "loadSettings",
      "updateSetting",
    ]),

    ...mapActions("devices", [
      "fetchDevices",
      "selectDevice",
      "updateAddress",
      "updateDeviceSettings",
    ]),

    getDeviceStatus,

    initializeMarker() {
      if (!window.google?.maps || !this.$parent.map) return;

      const status = this.$parent.getDeviceStatus(this.device);
      const filters = this.$parent.filters;

      let isVisible = false;
      if (status.status === "Moving") isVisible = filters.showMoving;
      if (status.status === "Stopped") isVisible = filters.showStopped;
      if (status.status === "Offline") isVisible = filters.showOffline;

      this.marker = new window.google.maps.Marker({
        position: this.position,
        title: this.title,
        map: this.$parent.map,
        icon: this.getMarkerIcon(),
        visible: isVisible,
      });
    },
    isDeviceActive(device) {
      return (
        parseFloat(
          device.latest_device_point.device_point_detail.speed.display
        ) > 0
      );
    },

    formatTimeAgo,
    kmToMph,
    getDeviceMetrics,
    getBatteryStatus,

    // Map initialization methods
    async initializeLayers() {
      try {
        if (!window.google?.maps) {
          throw new Error("Google Maps not loaded");
        }

        await this.initializeMap();
        this.map = this.$refs.map.map;
      } catch (error) {
        console.error("Layer initialization error:", error);
      }
    },

    getBatteryColor(voltage) {
      const status = this.getBatteryStatus(voltage);
      switch (status) {
        case "critical":
          return "text-red-600";
        case "low":
          return "text-orange-500";
        case "high":
          return "text-yellow-500";
        case "normal":
          return "text-green-500";
        default:
          return "text-gray-500";
      }
    },

    getLastStoppedTime(device) {
      const driveStatus =
        device?.latest_device_point?.device_state?.drive_status;
      const statusTime =
        device?.latest_device_point?.device_state?.drive_status_begin_time;
      const duration =
        device?.latest_device_point?.device_state?.drive_status_duration
          ?.display;

      if (driveStatus === "stopped") {
        return {
          status: "Stopped",
          since: statusTime,
          duration: duration,
        };
      } else {
        // For currently driving vehicles, we'd need to track their last stop in state
        return null;
      }
    },
    handleSidebarDeviceClick(device) {
      // Close sidebar on mobile
      if (window.innerWidth < 768) {
        this.isSidebarOpen = false;
      }

      this.clickedDeviceId = device.device_id;
      this.selectDevice(device.device_id);

      const position = this.getDevicePosition(device);
      this.defaultCenter = position;

      this.reverseGeocode(position.lat, position.lng, device.device_id);
    },

    // Device interaction methods
    handleMarkerClick(device) {
      if (this.clickedDeviceId === device.device_id) {
        this.clickedDeviceId = null;
        this.selectDevice(null);
      } else {
        this.clickedDeviceId = device.device_id;
        this.selectDevice(device.device_id);

        const position = this.getDevicePosition(device);
        this.defaultCenter = position;
        this.reverseGeocode(position.lat, position.lng, device.device_id);
      }
    },

    /**
     * Closes the info window for the currently selected device
     *
     * @see {@link selectDevice}
     */
    handleInfoWindowClose() {
      this.clickedDeviceId = null;
      this.selectDevice(null);
    },
    updateSettings(newSettings) {
      // Update all settings at once through Vuex
      Object.entries(newSettings).forEach(([key, value]) => {
        this.updateSetting({ key, value });
      });
    },

    // Utility methods
    async reverseGeocode(lat, lng, deviceId) {
      if (!lat || !lng || !deviceId) {
        console.warn("Missing parameters for geocoding:", {
          lat,
          lng,
          deviceId,
        });
        return;
      }

      try {
        console.log("Starting reverse geocoding for device:", deviceId);

        // First check if we already have the address
        if (this.$store.state.devices.addressMap[deviceId]) {
          console.log("Address already cached for device:", deviceId);
          return;
        }

        // Set a temporary loading state
        await this.updateAddress({
          deviceId,
          address: "Loading address...",
        });

        const response = await fetch(
          `https://maps.googleapis.com/maps/api/geocode/json?latlng=${lat},${lng}&key=${this.apiKey}`
        );

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        console.log("Geocoding response for device:", deviceId, data.status);

        if (data.status === "OK" && data.results?.[0]) {
          await this.updateAddress({
            deviceId,
            address: data.results[0].formatted_address,
          });
          console.log("Address updated for device:", deviceId);
        } else {
          throw new Error(`Geocoding failed: ${data.status}`);
        }
      } catch (error) {
        console.error("Geocoding error for device:", deviceId, error);
        await this.updateAddress({
          deviceId,
          address: "Address lookup failed",
        });
      }
    },

    // Helper method to check if coordinates are valid
    isValidCoordinate(lat, lng) {
      const validLat =
        typeof lat === "number" && !isNaN(lat) && lat >= -90 && lat <= 90;
      const validLng =
        typeof lng === "number" && !isNaN(lng) && lng >= -180 && lng <= 180;
      return validLat && validLng;
    },

    kelvinToFahrenheit(kelvin) {
      return (((kelvin - 273.15) * 9) / 5 + 32).toFixed(1);
    },

    getMarkerIcon(device) {
      const speed = device.latest_device_point.speed;
      const heading = device.latest_device_point.angle;

      return {
        path: "M 12,0 L 24,24 L 12,16 L 0,24 Z",
        fillColor: speed > 0 ? "#00FF00" : "#FF0000",
        fillOpacity: 1,
        strokeWeight: 1,
        strokeColor: "#000000",
        rotation: heading || 0,
        scale: 1,
        anchor: { x: 12, y: 12 },
      };
    },

    /**
     * Converts a duration in hours into a human-readable format.
     *
     * @param {string|number} durationDisplay - The duration to format,
     *     expressed in hours.
     * @returns {string} A formatted string representing the duration in
     *     days and hours, e.g., "3 days, 5 hours", "7 days", or "15 hours".
     *     If the input is not a valid number, it returns the input as is.
     */

    formatDuration(durationDisplay) {
      const hours = parseInt(durationDisplay);
      if (isNaN(hours)) return durationDisplay;

      const days = Math.floor(hours / 24);
      const remainingHours = hours % 24;

      if (days === 0) return `${hours} hours`;
      if (remainingHours === 0) return `${days} days`;
      return `${days} days, ${remainingHours} hours`;
    },
    getTimeElapsed(timestamp) {
      const now = new Date();
      const startTime = new Date(timestamp);
      const diffInMs = now - startTime;

      // Convert to hours, minutes, seconds
      const hours = Math.floor(diffInMs / (1000 * 60 * 60));
      const minutes = Math.floor((diffInMs % (1000 * 60 * 60)) / (1000 * 60));
      const seconds = Math.floor((diffInMs % (1000 * 60)) / 1000);

      // Format the output
      if (hours > 0) {
        return `${hours}h ${minutes}m`;
      } else if (minutes > 0) {
        return `${minutes}m ${seconds}s`;
      } else {
        return `${seconds}s`;
      }
    },
    //   startTime() {
    //   // Update every second
    //   this.timer = setInterval(() => {
    //     // Force template update
    //     this.$forceUpdate();
    //   }, 1000);
    //  },
    startTimer() {
      // Update every second
      this.timer = setInterval(() => {
        // Force template update
        this.$forceUpdate();
      }, 5000);
    },

    getDevicePosition(device) {
      if (device?.latest_device_point) {
        return {
          lat: Number(device.latest_device_point.lat),
          lng: Number(device.latest_device_point.lng),
        };
      }
      return this.defaultCenter;
    },

    startPolling() {
      // Poll every 5 seconds
      this.pollingInterval = setInterval(async () => {
        try {
          await this.fetchDevices(); // This uses your existing Vuex action
          console.log("Devices updated from polling");
        } catch (error) {
          console.error("Polling error:", error);
        }
      }, 10000); // 5 second interval
    },

    stopPolling() {
      if (this.pollingInterval) {
        clearInterval(this.pollingInterval);
        this.pollingInterval = null;
      }
    },

    animateMarkerMovement(device, start, end, updatedDevice) {
      const duration = 1000;
      const startTime = performance.now();

      const animate = async (currentTime) => {
        const elapsed = currentTime - startTime;
        const progress = Math.min(elapsed / duration, 1);

        const easeProgress =
          progress < 0.5
            ? 2 * progress * progress
            : 1 - Math.pow(-2 * progress + 2, 2) / 2;

        // Calculate current position
        const currentPosition = {
          lat: start.lat + (end.lat - start.lat) * easeProgress,
          lng: start.lng + (end.lng - start.lng) * easeProgress,
        };

        // Update position through Vuex
        await this.$store.dispatch("devices/updateDevicePosition", {
          deviceId: device.device_id,
          position: currentPosition,
        });

        if (progress < 1) {
          const frameId = requestAnimationFrame(animate);
          this.animationQueue.set(device.device_id, frameId);
        } else {
          this.animationQueue.delete(device.device_id);
          // Final update with complete device data
          await this.$store.dispatch("devices/updateDevice", updatedDevice);
        }
      };

      const frameId = requestAnimationFrame(animate);
      this.animationQueue.set(device.device_id, frameId);
    },

    updateMapSettings(key, value) {
      if (this.mapSettings[key] !== value) {
        this.mapSettings[key] = value;
        this.saveMapSettings();
      }
    },

    /**
     * Lifecycle hook: called after the component has been created
     * Loads stored map settings
     */
    created() {
      this.loadMapSettings();
    },
  },

  // Lifecycle hooks
  async mounted() {
    try {
      await this.loadSettings();
      await this.fetchDevices();
      this.startPolling();
      this.startTimer();
    } catch (error) {
      console.error("Mount error:", error);
    }
  },
  beforeUnmount() {
    if (this.timer) {
      clearInterval(this.timer);
    }
    this.stopPolling(); // Stop polling when component unmounts
    // Clean up map listeners
    if (this.map) {
      window.google.maps.event.clearListeners(this.map, "click");
    }
  },
};
</script>

<style>
.gm-style button[title="Show satellite imagery"] {
  display: none !important;
}

.gm-style button[title="Show street map"] {
  display: none !important;
}
</style>
