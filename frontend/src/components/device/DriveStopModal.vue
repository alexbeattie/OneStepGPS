<!-- # DriveStopModal.vue -->
<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-start justify-center">
    <div class="bg-white w-full h-full md:rounded-lg md:shadow-xl md:w-11/12 md:max-w-6xl md:max-h-[90vh] md:my-4 overflow-hidden">
      <!-- Modal Header -->
      <div class="px-4 md:px-6 py-4 border-b border-gray-200 flex justify-between items-center sticky top-0 bg-white z-10">
        <h3 class="text-xl font-semibold text-gray-900">Route History</h3>
        <button @click="$emit('close')" class="text-gray-500 hover:text-gray-700 p-2">
          <span class="sr-only">Close</span>
          <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Date Range Selector - Mobile Optimized -->
      <div class="px-4 md:px-6 py-4 border-b border-gray-200 space-y-4">
        <!-- Date Controls -->
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div class="flex flex-col">
            <label class="text-sm font-medium text-gray-700">From Date</label>
            <input 
              type="datetime-local" 
              v-model="fromDate"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 text-sm"
            />
          </div>
          <div class="flex flex-col">
            <label class="text-sm font-medium text-gray-700">To Date</label>
            <input 
              type="datetime-local" 
              v-model="toDate"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 text-sm"
            />
          </div>
          <div class="flex flex-col">
            <label class="text-sm font-medium text-gray-700">Stop Duration</label>
            <select 
              v-model="stopDuration"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 text-sm"
            >
              <option value="5m0s">5 minutes</option>
              <option value="10m0s">10 minutes</option>
              <option value="15m0s">15 minutes</option>
              <option value="30m0s">30 minutes</option>
            </select>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex flex-wrap gap-2">
          <button 
            @click="fetchRouteData"
            class="flex-1 md:flex-none px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 text-sm"
          >
            Fetch Data
          </button>
          <button 
            @click="downloadData"
            :disabled="!routeData?.events?.length"
            class="flex-1 md:flex-none px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed text-sm flex items-center justify-center gap-2"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"/>
            </svg>
            Download CSV
          </button>
        </div>
      </div>

      <!-- Main Content Area -->
      <div class="flex flex-col h-[calc(100vh-240px)] md:h-[calc(90vh-220px)]">
        <!-- Loading/Error/No Data States -->
        <div v-if="loading" class="flex-1 flex items-center justify-center">
          <div class="text-center">
            <div class="inline-block animate-spin rounded-full h-8 w-8 border-4 border-blue-500 border-t-transparent"></div>
            <p class="mt-2 text-gray-600">Loading route data...</p>
          </div>
        </div>

        <div v-else-if="error" class="flex-1 flex items-center justify-center p-4">
          <div class="text-center text-red-500">
            <p>{{ error }}</p>
          </div>
        </div>

        <div v-else-if="!routeData?.events?.length" class="flex-1 flex items-center justify-center p-4">
          <div class="text-center">
            <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <h3 class="mt-2 text-sm font-medium text-gray-900">No route data found</h3>
            <p class="mt-1 text-sm text-gray-500">Try adjusting the date range or stop duration</p>
          </div>
        </div>

        <template v-else>
          <!-- Summary Stats - Scrollable on mobile -->
          <div class="px-4 md:px-6 py-4 bg-gray-50 border-b border-gray-200 overflow-x-auto">
            <div class="min-w-[640px]">
              <!-- First Row -->
              <div class="grid grid-cols-3 gap-4">
                <div class="col-span-1">
                  <div class="text-sm font-medium text-gray-500">Total Distance</div>
                  <div class="mt-1 text-lg font-semibold text-gray-900">{{ routeData.summary.distance }}</div>
                </div>
                <div class="col-span-1">
                  <div class="text-sm font-medium text-gray-500">Average Speed</div>
                  <div class="mt-1 text-lg font-semibold text-gray-900">{{ routeData.summary.avgSpeed }}</div>
                </div>
                <div class="col-span-1">
                  <div class="text-sm font-medium text-gray-500">Top Speed</div>
                  <div class="mt-1 text-lg font-semibold text-gray-900">{{ routeData.summary.topSpeed }}</div>
                </div>
              </div>
              <!-- Second Row -->
              <div class="grid grid-cols-3 gap-4 mt-4">
                <div class="col-span-1">
                  <div class="text-sm font-medium text-gray-500">Total Duration</div>
                  <div class="mt-1 text-lg font-semibold text-gray-900">{{ routeData.summary.duration }}</div>
                </div>
                <div class="col-span-1">
                  <div class="text-sm font-medium text-gray-500">Idle Time</div>
                  <div class="mt-1 text-lg font-semibold text-gray-900">{{ routeData.summary.idleTime }}</div>
                </div>
                <div class="col-span-1">
                  <div class="text-sm font-medium text-gray-500">Stop Duration</div>
                  <div class="mt-1 text-lg font-semibold text-gray-900">{{ routeData.summary.stopTime }}</div>
                </div>
              </div>
            </div>
          </div>

          <!-- Map View -->
          <div class="h-1/2 border-b border-gray-200" ref="mapContainer"></div>

          <!-- Table View - Mobile Optimized -->
          <div class="h-1/2 overflow-auto">
            <div class="min-w-[640px]">
              <table class="w-full divide-y divide-gray-200">
                <thead class="bg-gray-50 sticky top-0">
                  <tr>
                    <th class="px-4 md:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Type</th>
                    <th class="px-4 md:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Start Time</th>
                    <th class="px-4 md:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">End Time</th>
                    <th class="px-4 md:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Duration</th>
                    <th class="px-4 md:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Location</th>
                    <th class="px-4 md:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Odometer</th>
                  </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                  <tr v-for="(event, index) in routeData.events" :key="index" class="hover:bg-gray-50">
                    <td class="px-4 md:px-6 py-4 whitespace-nowrap">
                      <span :class="getEventTypeClass(event.type)">
                        {{ event.type }}
                      </span>
                    </td>
                    <td class="px-4 md:px-6 py-4 whitespace-nowrap">{{ formatDate(event.start_time) }}</td>
                    <td class="px-4 md:px-6 py-4 whitespace-nowrap">{{ formatDate(event.end_time) }}</td>
                    <td class="px-4 md:px-6 py-4 whitespace-nowrap">{{ event.duration_display }}</td>
                    <td class="px-4 md:px-6 py-4">{{ event.location }}</td>
                    <td class="px-4 md:px-6 py-4 whitespace-nowrap">{{ event.odometer_from }} → {{ event.odometer_to }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script>
import { apiService } from '@/services/api';


export default {
  name: 'DriveStopModal',
  props: {
    deviceId: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      fromDate: this.getDefaultFromDate(),
      toDate: this.getDefaultToDate(),
      stopDuration: '5m0s',
      routeData: {
        events: [],
        points: [],
        summary: {
          distance: '0 mi',
          avgSpeed: '0 mph',
          topSpeed: '0 mph',
          duration: '0h 0m',
          idleTime: '0m',
          stopTime: '0m'
        }
      },
      loading: false,
      error: null,
      map: null,
      markers: [],
      polyline: null
    };
  },
  methods: {
     downloadData() {
      if (!this.routeData?.events?.length) return;

      // Prepare CSV data
      const headers = [
        'Type',
        'Start Time',
        'End Time',
        'Duration',
        'Location',
        'Odometer Start',
        'Odometer End'
      ];

      const rows = this.routeData.events.map(event => [
        event.type,
        this.formatDate(event.start_time),
        this.formatDate(event.end_time),
        event.duration_display,
        event.location,
        event.odometer_from,
        event.odometer_to
      ]);

      // Add summary data at the top
      const summaryRows = [
        ['Summary'],
        ['Total Distance', this.routeData.summary.distance],
        ['Average Speed', this.routeData.summary.avgSpeed],
        ['Top Speed', this.routeData.summary.topSpeed],
        ['Total Duration', this.routeData.summary.duration],
        ['Idle Time', this.routeData.summary.idleTime],
        ['Stop Duration', this.routeData.summary.stopTime],
        [],  // Empty row between summary and details
      ];

      // Combine all rows
      const csvContent = [
        ...summaryRows,
        headers,
        ...rows
      ]
        .map(row => row.join(','))
        .join('\n');

      // Create and trigger download
      const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
      const link = document.createElement('a');
      if (navigator.msSaveBlob) {
        // IE 10+
        navigator.msSaveBlob(blob, 'route_history.csv');
      } else {
        link.href = URL.createObjectURL(blob);
        link.setAttribute('download', 'route_history.csv');
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
      }
    },
    getDefaultFromDate() {
      const date = new Date();
      date.setHours(0, 0, 0, 0);
      return date.toISOString().slice(0, 16);
    },
    getDefaultToDate() {
      const date = new Date();
      date.setHours(23, 59, 59, 999);
      return date.toISOString().slice(0, 16);
    },
    async fetchRouteData() {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await apiService.getDriveStopRoute(
          this.deviceId,
          this.fromDate,
          this.toDate,
          this.stopDuration
        );
        console.log('Raw API response:', response);

        // Process the response data
        const processedData = {
          events: [],
          points: [],
          summary: {
            distance: response.distance?.display || '0 mi',
            avgSpeed: response.average_speed?.display || '0 mph',
            topSpeed: response.top_speed?.display || '0 mph',
            duration: response.duration?.display || '0h 0m',
            idleTime: response.idle_duration?.display || '0m',
            stopTime: response.stop_duration?.display || '0m'
          }
        };

        // Process drive_stop_list into events
        if (response.drive_stop_list && Array.isArray(response.drive_stop_list)) {
          processedData.events = response.drive_stop_list.map(event => ({
            type: event.type,
            start_time: event.time_from,
            end_time: event.time_to,
            duration_display: event.duration?.display,
            location: `${event.first_valid_lat_lng?.lat.toFixed(6)}, ${event.first_valid_lat_lng?.lng.toFixed(6)}`,
            lat: event.first_valid_lat_lng?.lat,
            lng: event.first_valid_lat_lng?.lng,
            odometer_from: event.odometer_from?.display,
            odometer_to: event.odometer_to?.display
          }));

          // Add points for the map
          processedData.points = response.drive_stop_list.flatMap(event => ([
            event.first_valid_lat_lng,
            event.last_valid_lat_lng
          ]).filter(Boolean));
        }

        console.log('Processed route data:', processedData);
        this.routeData = processedData;

        this.$nextTick(() => {
          this.initializeMap();
          this.plotRouteOnMap();
        });
      } catch (error) {
        this.error = error.message;
        console.error('Error fetching route data:', error);
      } finally {
        this.loading = false;
      }
    },
     formatDate(dateString) {
      return new Date(dateString).toLocaleString();
    },

    getEventTypeClass(type) {
      const classes = {
        'drive': 'bg-green-100 text-green-800',
        'stop': 'bg-red-100 text-red-800'
      };
      return `px-2 py-1 text-xs rounded-full ${classes[type.toLowerCase()] || 'bg-gray-100 text-gray-800'}`;
    },

    initializeMap() {
      if (!window.google?.maps) {
        console.warn('Google Maps not loaded');
        return;
      }

      if (!this.$refs.mapContainer) {
        console.warn('Map container not found');
        return;
      }

      if (!this.map && this.routeData.events.length > 0) {
        this.map = new window.google.maps.Map(this.$refs.mapContainer, {
          zoom: 12,
          mapTypeId: 'roadmap',
          center: { lat: 0, lng: 0 } // Default center
        });
      }
    },

    clearMapOverlays() {
      if (this.polyline) {
        this.polyline.setMap(null);
      }
      this.markers.forEach(marker => marker.setMap(null));
      this.markers = [];
    },

    plotRouteOnMap() {
      if (!this.map || !this.routeData?.points?.length) return;

      this.clearMapOverlays();

      // Create path for polyline
      const path = this.routeData.points.map(point => ({
        lat: point.lat,
        lng: point.lng
      }));

      // Create and set polyline
      this.polyline = new window.google.maps.Polyline({
        path,
        geodesic: true,
        strokeColor: '#FF0000',
        strokeOpacity: 1.0,
        strokeWeight: 2,
        map: this.map
      });

      // Add markers for stops
      this.routeData.events
        .filter(event => event.type.toLowerCase() === 'stop')
        .forEach(stop => {
          const marker = new window.google.maps.Marker({
            position: { lat: stop.lat, lng: stop.lng },
            map: this.map,
            title: `Stop: ${stop.duration_display}`
          });
          this.markers.push(marker);
        });

      // Fit bounds to show all points
      const bounds = new window.google.maps.LatLngBounds();
      path.forEach(point => bounds.extend(point));
      this.map.fitBounds(bounds);
    }

    // ... rest of your methods remain the same
  }
};
</script>