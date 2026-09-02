<script setup lang="ts">
import type { Place } from '@/api/types';
import { onBeforeUnmount, onMounted, ref, watch } from 'vue';
import L from "leaflet";
import "leaflet/dist/leaflet.css";


const props = defineProps<{
  places: Place[]
  height?: string
}>()

const container = ref<HTMLDivElement | null>(null)
let map: L.Map | null = null;
let markerLayer: L.LayerGroup | null = null;

// Jakarta city centre — the fallback view when there's nothing to fit to.
const JAKARTA_CENTER: L.LatLngExpression = [-6.2088, 106.8456];

// Leaflet's default marker icons resolve to bundler-relative URLs that break
// under Vite; point them at the bundled asset URLs explicitly.
const icon = L.icon({
  iconUrl: new URL("leaflet/dist/images/marker-icon.png", import.meta.url).href,
  iconRetinaUrl: new URL("leaflet/dist/images/marker-icon-2x.png", import.meta.url).href,
  shadowUrl: new URL("leaflet/dist/images/marker-shadow.png", import.meta.url).href,
  iconSize: [25, 41],
  iconAnchor: [12, 41],
  popupAnchor: [1, -34],
  shadowSize: [41, 41],
})

function renderMarkers() {
  if (!map) return
  markerLayer?.remove()
  markerLayer = L.layerGroup().addTo(map)

  const points: L.LatLngExpression[] = []
  for (const place of props.places) {
    L.marker([place.lat, place.lng], {icon})
      .bindPopup(`<strong>${place.name}</strong><br>${place.address}`)
      .addTo(markerLayer);
    points.push([place.lat, place.lng]);
  }

  if (points.length > 0) {
    map.fitBounds(L.latLngBounds(points), {padding: [30, 30], maxZoom: 15})
  }
}

onMounted(() => {
  if (!container.value) return
  map = L.map(container.value).setView(JAKARTA_CENTER, 11)
  L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
    maxZoom: 19,
  }).addTo(map);
  renderMarkers();
})

onBeforeUnmount(() => {
  map?.remove()
  map = null
})

watch(() => props.places, renderMarkers, {deep: true})

</script>

<template>
  <div ref="container" class="map" :style="{height: height ?? '360px'}"></div>
</template>

<style scoped>
.map {
  width: 100%;
  border-radius: 10px;
  border: 1px solid var(--line);
  z-index: 0;
}
</style>
