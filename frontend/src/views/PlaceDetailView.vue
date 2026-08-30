<script setup lang="ts">
import { placesApi } from '@/api';
import { CATEGORY_LABELS, type Place, type PlacePhoto } from '@/api/types';
import { useAuthStore } from '@/stores/auth';
import { useSavedStore } from '@/stores/saved';
import { useVisitsStore } from '@/stores/visits';
import { computed, onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';


const route = useRoute()
const auth = useAuthStore()
const saved = useSavedStore()
const visits = useVisitsStore()

const place = ref<Place | null>(null)
const photos = ref<PlacePhoto[]>([])
const loading = ref(true)
const error = ref<string | null>(null)

const placeId = computed(() => Number(route.params.id))

const openingHours = computed(() => {
  const hours = place.value?.opening_hours
  if (!hours || Object.keys(hours).length === 0) return null
  return Object.entries(hours) as [string, string][]
})

const menu = computed(() => {
  const m = place.value?.menu
  if (!m || Object.keys(m).length === 0) return null
  return Object.entries(m) as [string, string][]
})

async function load() {
  loading.value = true
  error.value = null
  try {
    const data = await placesApi.get(placeId.value)
    place.value = data.place
    photos.value = data.photos ?? []
  } catch(err) {
    error.value = err instanceof Error ? err.message : "Failed to load place"
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await load()
  if (auth.isAuthenticated) {
    await Promise.all([saved.load(), visits.load()])
  }
})

</script>

<template>
  <div class="container">
    <p v-if="loading" class="muted">Loading...</p>
    <p v-else-if="error" class="error">{{ error }}</p>

    <template v-else-if="place">
      <RouterLink :to="{name: 'home'}" class="muted back">← Back to discover</RouterLink>

      <header class="place-header">
        <div>
          <h1>{{ place.name }}</h1>
          <span class="badge">{{ CATEGORY_LABELS[place.category] }}</span>
          <span v-if="place.source === 'user'" class="badge contributed">Member contribution</span>
        </div>

        <div v-if="auth.isAuthenticated" class="actions">
          <button :class="{active: saved.savedPlaceIds.has(place.id)}" @click="saved.toggle(place.id)">
            {{ saved.savedPlaceIds.has(place.id) ? "★ Saved" : "☆ Save"  }}
          </button>

          <button :disabled="visits.visitedPlaceIds.has(place.id)" @click="visits.markVisited(place.id)">
            {{ visits.visitedPlaceIds.has(place.id) ? "✓ Visited" : "Mark visited" }}
          </button>
        </div>
      </header>

      <p v-if="place.description" class="description">{{ place.description }}</p>

      <section class="card details-card">
        <h2>Details</h2>
        <dl>
          <dt>Address</dt>
          <dd>{{ place.address }}</dd>

          <template v-if="place.district">
            <dt>Area</dt>
            <dd>{{ place.district }}</dd>
          </template>

          <template v-if="place.price_range">
            <dt>Price range</dt>
            <dd>{{ place.price_range }}</dd>
          </template>
        </dl>

        <template v-if="menu">
          <h3>Menu highlights</h3>
          <dl>
            <template v-for="[item, price] in menu" :key="item">
              <dt>{{ item }}</dt>
              <dd>{{ price }}</dd>
            </template>
          </dl>
        </template>
      </section>

      <section v-if="photos.length > 0" class="photos">
        <h2>Photos</h2>
        <div class="photo-grid">
          <figure v-for="photo in photos" :key="photo.id">
            <img :src="photo.url" :alt="photo.caption ?? place.name" loading="lazy"/>
            <figcaption v-if="photo.caption">{{ photo.caption }}</figcaption>
          </figure>
        </div>
      </section>
    </template>
  </div>
</template>

<style scoped>
.back {
  display: inline-block;
  margin-bottom: 1rem;
  font-size: 0.9rem;
}

.place-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
  flex-wrap: wrap;
}

.place-header h1 {
  margin: 0 0 0.4rem;
}

.badge.contributed {
  background: var(--line);
  color: var(--muted);
  margin-left: 0.35rem;
}

.actions {
  display: flex;
  gap: 0.5rem;
}

button.active {
  border-color: var(--accent);
  color: var(--accent);
}

.description {
  max-width: 65ch;
}

.details-card {
  margin: 1.5rem 0;
  max-width: 500px;
}

h2 {
  font-size: 1.05rem;
  margin-top: 0;
}

h3 {
  font-size: 0.95rem;
  margin-bottom: 0.4rem;
}

dl {
  display: grid;
  grid-template-columns: auto 1fr;
  gap: 0.3rem 1rem;
  margin: 0 0 1rem;
  font-size: 0.9rem;
}

dt {
  color: var(--muted);
}

dd {
  margin: 0;
}

.photo-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 1rem;
}

.photo-grid img {
  width: 100%;
  border-radius: 8px;
  display: block;
}

figure {
  margin: 0;
}

figcaption {
  font-size: 0.85rem;
  margin-top: 0.3rem;
}
</style>
