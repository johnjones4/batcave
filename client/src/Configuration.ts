export interface Coordinate {
  latitude: number
  longitude: number
}

export interface Configuration {
  liveFeeds: [string]
  weatherLocations: [Coordinate]
}