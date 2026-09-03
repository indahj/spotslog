<script setup lang="ts">
import { placesApi } from '@/api';
import { CATEGORY_LABELS, type PlaceCategory } from '@/api/types';
import { useVisitsStore } from '@/stores/visits';
import { OpenStreetMapProvider } from 'leaflet-geosearch';
import { ref } from 'vue';
import { useRouter } from 'vue-router';


const router = useRouter()
const visits = useVisitsStore()

const categories = Object.keys(CATEGORY_LABELS) as PlaceCategory[]

const name = ref("")
const category = ref<PlaceCategory>("cafe")
const address = ref("")
const district = ref("")
const lat = ref<number | null>(null)
const lng = ref<number | null>(null)
const description = ref("")
const priceRange = ref("")
const visibility = ref<"public" | "private">("private")
const alsoMarkVisited = ref(true)
const photo = ref<File | null>(null)

const searchResults = ref<{ label: string; lat: number; lng: number }[]>([])
const searching = ref(false)
let searchTimeout:ReturnType<typeof setTimeout> | undefined;

const provider = new OpenStreetMapProvider({
  params: { countrycodes: "id"}
})

const error = ref<string | null>(null)
const submitting = ref(false)

function onAddressInput() {
  clearTimeout(searchTimeout)
  if (!address.value.trim()) {
    searchResults.value = []
    return
  }
  searchTimeout = setTimeout(async() => {
    searching.value = true
    try {
      const results = await provider.search({query: address.value})
      searchResults.value = results.map((r) => ({
        label: r.label,
        lat: r.y,
        lng: r.x,
      }))
    } finally {
      searching.value = false
    }
  }, 500)
}

function selectResult(result: {label: string; lat: number; lng: number}) {
  lat.value = Number(result.lat.toFixed(6))
  lng.value = Number(result.lng.toFixed(6))
  address.value = result.label
  address.value = result.label
  searchResults.value = []
}

function onPhotoChange(event: Event) {
  const input = event.target as HTMLInputElement
  photo.value = input.files?.[0] ?? null
}

async function submit() {
  if (lat.value === null || lng.value === null) {
    error.value = "Latitude and longitude are required."
    return
  }

  error.value = null
  submitting.value = true
  try {
    const place = await placesApi.create({
      name: name.value,
      category: category.value,
      address: address.value,
      district: district.value || undefined,
      lat: lat.value,
      lng: lng.value,
      description: description.value || undefined,
      price_range: priceRange.value || undefined,
      source: "user",
      visibility: visibility.value
    })

    let photoError: string | null = null;
    if (photo.value) {
      try {
        await placesApi.uploadPhoto(place.id, photo.value)
      } catch (err) {
        photoError = err instanceof Error ? err.message : "unknown error"
      }
    }

    if (alsoMarkVisited.value) {
      await visits.markVisited(place.id)
    }

    if (photoError) {
      error.value = `Place was added, but the photo failed to upload: ${photoError}`
      submitting.value = false
      return
    }

    if (alsoMarkVisited.value) {
      router.push({name: "visits"})
    } else {
      router.push({name: "place-detail", params: {id: place.id}})
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Failed to add place"
  } finally {
    submitting.value = false
  }
}

</script>

<template>

  <div class="container narrow">
    <h1>Add a place</h1>
    <p class="muted">
      Somewhere that isn't in the recommendations yet. Keep it private for your own records, or make it public so it shows up on the homepage for everyone.
    </p>

    <form class="card" @submit.prevent="submit">
      <div class="field">
        <label for="name">Name</label>
        <input id="name" v-model="name"/>
      </div>

      <div class="field">
        <label for="category">Category</label>
        <select id="category" v-model="category">
          <option v-for="cat in categories" :key="cat" :value="cat">
            {{ CATEGORY_LABELS[cat] }}
          </option>
        </select>
      </div>

      <div class="field">
        <label for="address">Address</label>
        <input id="address" v-model="address" required @input="onAddressInput">
        <p v-if="searching" class="muted">Searching...</p>
        <ul v-if="searchResults.length > 0" class="search-results">
          <li v-for="(result, i) in searchResults" :key="i" @click="selectResult(result)">
            {{ result.label }}
          </li>
        </ul>
      </div>

      <div class="field">
        <label for="district">Area / district</label>
        <input id="district" v-model="district" placeholder="Kemang, Menteng..."/>
      </div>

       <div class="coords">
        <div class="field">
          <label for="lat">Latitude</label>
          <input id="lat" v-model.number="lat" type="number" step="any" disabled/>
        </div>
        <div class="field">
          <label for="lng">Longitude</label>
          <input id="lng" v-model.number="lng" type="number" step="any" disabled/>
        </div>
      </div>

      <div class="field">
        <label for="price">Price range</label>
        <input id="price" v-model="priceRange" placeholder="Rp 50–100k" />
      </div>

      <div class="field">
        <label for="description">Notes / description</label>
        <textarea id="description" v-model="description" rows="3" />
      </div>

      <div class="field">
        <label for="photo">Photo (optional)</label>
        <input id="photo" type="file" accept="image/*" @change="onPhotoChange" />
      </div>

      <div class="field">
        <label for="visibility">Visibility</label>
        <select id="visibility" v-model="visibility">
          <option value="private">Private — only I can see it</option>
          <option value="public">Public — show it on the homepage</option>
        </select>
      </div>

      <label class="checkbox">
        <input v-model="alsoMarkVisited" type="checkbox" />
        Also add this to my visit history
      </label>

      <p v-if="error" class="error">{{ error }}</p>

      <button class="primary" type="submit" :disabled="submitting">
        {{ submitting ? "Saving…" : "Add place" }}
      </button>
    </form>

  </div>

</template>

<style scoped>
.narrow {
  max-width: 560px;
}

.search-results {
  list-style: none;
  margin: 0.4rem 0 0;
  padding: 0;
  border: 1px solid var(--line);
  border-radius: 6px;
  max-height: 200px;
  overflow-y: auto;
}

.search-results li {
  padding: 0.5rem 0.7rem;
  cursor: pointer;
  font-size: 0.9rem;
  border-bottom: 1px solid var(--line);
}

.search-results li:last-child {
  border-bottom: none;
}
.search-results li:hover {
  background: var(--paper);
}

.coords {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.7rem;
  align-items: end;
}

.coords .field {
  margin-bottom: 0.9rem;
}

.checkbox {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.9rem;
  margin-bottom: 1rem;
}

.checkbox input {
  width: auto;
}

textarea {
  resize: vertical;
}
</style>
