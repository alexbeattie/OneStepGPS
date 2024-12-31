# One Step GPS Project

A GPS tracking application with Vue.js frontend and Go backend.

## Project Structure
- `/frontend` - Vue.js application
- `/backend` - Go server

## Getting Started
1. Clone the repository
2. Set up environment variables (see below)
3. Follow setup instructions in frontend and backend directories

## Environment Setup
- Copy `.env.example` to `.env` in both frontend and backend directories
- Fill in your own environment variables

# One Step GPS - Project

This is my submission for the project detail. The project contains 2 folders in a [monorepo on Github](https://github.com/alexbeattie/OneStepProj). One is labeled **[frontend](https://github.com/alexbeattie/OneStepProj/tree/main/frontend)** and the other is labeled [backend](https://github.com/alexbeattie/OneStepProj/tree/main/backend).

The **backend** is written in Go and the **frontend** is written in Vue.js. The backend is intended to connected to connect to a [PostgreSQL](https://www.postgresql.org/) database in order to demonstrate persistence. 

Persistence / state is shared through the **frontend** through the use of Vue's [Vuex state management package](https://vuex.vuejs.org/installation.html) interfacing with PostgreSQL database.

The local repo setup and configuration of the Go web server differs from the production server in the initial configuration. 

For the localhost - I put both the Go web server files under a folder entitled backend and beside it I created another folder called frontend. The frontend contains the Vue.js project in its entirety.

The variance mostly is in the main.go file. 

When you put the local version up on an [EC2](https://aws.amazon.com/pm/ec2/) instance and connect to an [AWS RDS PostgreSQL](https://aws.amazon.com/free/database/) database you need to account for the ports, url strings, CORS and SSL. 

I created an AMI (Amazon Machine Image) of an Ubuntu instance Amazon EC2. The container uses Security Groups, which need to be configured to allow incoming connections for ports 443 & 80, https & http respectively. 

This will also need to be connected to your PostgreSQL database instance, which, in this case is also hosted on AWS their RDS.

Once you do that, you run into the problem of serving an app from a domain such as 34.207.185.237, which is, obviously, not good. This also presents the TLS/SSL certificate validation warning, which is overcome through creating the proper certs & keys.

In this case, I mapped the domain onestepgpsdemo.com to the domain 34.207.185.237 in order to appropriately initiate the TLS protocol handshake / authentication of the domain name.  I did this with [Lets Encrypt](https://letsencrypt.org) and [AWS R53](https://aws.amazon.com/route53/)

# Go Web Server with Gin, GORM, and Middleware
 

This project demonstrates a Go web server using the following packages and features:

-  **Gin** for HTTP routing and middleware.
-  **CORS** middleware from `github.com/gin-contrib/cors`.
-  **GORM** for database interactions with PostgreSQL.
-  **dotenv** for environment variable management.

- A modular architecture with configuration, handlers, models, and services.


## Backend Features - (Local & Production)

-  **Database Integration**: PostgreSQL with automatic migrations.
-  **API Endpoints**:
-  `/api/v1/preferences/:userId` (GET/PUT) for managing user preferences.
-  `/api/v1/devices` (GET) for fetching device information.
-  `/v3/api/device-info` (GET) for retrieving detailed device information.
-  `/v3/api/route/drive-stop` (GET) for driving and stopping routes.
-  **CORS**: Configured for both development and production environments.
-  **Environment Variables**: API keys and sensitive data managed via `.env`.

## Complete File Structure
 ```
OneStepProj-main/
├── frontend/
│   └── src/
│       ├── assets/
│       ├── components/
│       │   ├── device/
│       │   │   ├── DeviceFilters.vue
│       │   │   ├── DeviceListItem.vue 
│       │   │   └── DriveStopModal.vue
│       │   └── map/
│       │       ├── AdvancedMarker.vue
│       │       └── CustomInfoWindow.vue
│       ├── composables/
│       │   ├── useDevices.js
│       │   ├── useGeocoding.js
│       │   └── useMap.js
│       ├── services/
│       │   ├── api.js
│       │   ├── deviceMetrics.js
│       │   └── preferencesServices.js
│       ├── store/
│       │   ├── modules/
│       │   │   ├── devices.js
│       │   │   └── mapSettings.js
│       │   └── index.js
│       ├── utils/
│       │   ├── deviceStatus.js
│       │   ├── formatters.js
│       │   ├── index.js
│       │   └── unitConversion.js
│       ├── views/
│       │   └── MapSettings.vue
│       ├── App.vue
│       └── main.js
└── backend/
   ├── handlers/
   ├── models/
   └── services/
   ```