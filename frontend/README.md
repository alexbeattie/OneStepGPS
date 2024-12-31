
# OneStepGPS / frontend

A full-stack GPS tracking application with Vue.js frontend and Go backend. This project enables real-time device tracking and management using the OneStepGPS API.

## Project Overview
This monorepo contains:
- `frontend/` - Vue.js web application (frontvueproj)
- `backend/` - Go API server

## Frontend (frontend)

### Features
#### Sidebar
- **Togglable Sidebar**: Show/hide the sidebar for managing devices and settings
- **Device Filters**:
  - Sort by activity, speed, distance, signal strength, and battery voltage
  - Filter by status (Moving, Stopped, Offline) and metrics

#### Map
- **Google Maps Integration**:
  - Real-time device location display
  - Interactive markers with device information
  - Customizable map types (Roadmap, Satellite, Hybrid, Terrain)
- **Real-Time Updates**:
  - Automatic data polling
  - Smooth marker animations

#### Device List
- Status-based categorization (Active/Moving, Stopped, Offline)
- Comprehensive device metrics (Speed, battery, GPS signal)
- Visual indicators for device status

### Project Setup
```bash
# Install dependencies
yarn install

# Development server with hot-reload
yarn serve

# Production build
yarn build

# Linting
yarn lint
```

### Environment Setup
1. Create `.env` file in frontend directory:
```env
VUE_APP_GOOGLE_MAPS_API_KEY=your_google_maps_api_key
VUE_APP_API_URL=http://localhost:8080
```

## File Structure
```
frontend/
└── src/
   ├── assets/
   ├── components/
   │   ├── device/
   │   │   ├── DeviceFilters.vue
   │   │   ├── DeviceListItem.vue
   │   │   └── DriveStopModal.vue
   │   └── map/
   │       ├── AdvancedMarker.vue
   │       └── CustomInfoWindow.vue
   ├── composables/
   │   ├── useDevices.js
   │   ├── useGeocoding.js
   │   └── useMap.js
   ├── services/
   │   ├── api.js
   │   ├── deviceMetrics.js
   │   └── preferencesServices.js
   ├── store/
   │   ├── modules/
   │   │   ├── devices.js
   │   │   └── mapSettings.js
   │   └── index.js
   ├── utils/
   │   ├── deviceStatus.js
   │   ├── formatters.js
   │   ├── index.js
   │   └── unitConversion.js
   ├── views/
   │   └── MapSettings.vue
   ├── App.vue
   └── main.js
```

## Contributing
1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License
MIT

## Acknowledgements
This project was created for the team at One Step GPS demonstrating full-stack development capabilities with Vue.js and Go
