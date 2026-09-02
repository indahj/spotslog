<script setup lang="ts">
import { CATEGORY_LABELS, type PlaceCategory } from "@/api/types";
import PlaceCard from "@/components/PlaceCard.vue";
import PlaceMap from "@/components/PlaceMap.vue";
import { useAuthStore } from "@/stores/auth";
import { usePlacesStore } from "@/stores/places";
import { useSavedStore } from "@/stores/saved";
import { useVisitsStore } from "@/stores/visits";
import { onMounted, ref } from "vue";

const places = usePlacesStore()
const auth = useAuthStore()
const saved = useSavedStore()
const visits = useVisitsStore()

const category = ref<PlaceCategory | "">("")
const area = ref("")
const showMap = ref(false)

const categories = Object.keys(CATEGORY_LABELS) as PlaceCategory[]

async function applyFilters() {
  if (!category.value && !area.value) {
    await places.loadHomepage()
    return
  }
  await places.search({
    category: category.value || undefined,
    area: area.value || undefined
  })
}

function resetFilters() {
  category.value = ""
  area.value = ""
  places.loadHomepage()
}

onMounted(async () => {
  await places.loadHomepage()
  if (auth.isAuthenticated) {
    await Promise.all([saved.load(), visits.load()])
  }
})

</script>

<template>
  <div class="container">
    <section class="hero">
      <div class="hero-copy">
        <p class="eyebrow">Jakarta, one place at a time</p>
        <h1>Discover new places. Remember the ones you've been.</h1>
        <p class="lead">
          Cafés, restaurants, museums, and libraries around Jakarta — browse
          recommendations, add your own favorites, and keep a visual record of
          everywhere you've actually been.
        </p>
        <div class="hero-actions">
          <a href="#places" class="primary-link">Browse places</a>
          <RouterLink :to="{ name: 'add-place' }" class="secondary-link">
            Add your own
          </RouterLink>
        </div>
      </div>

      <div class="hero-art">
        <img src="../assets/hero.svg" alt="">

      </div>
    </section>

    <form id="places" class="filters card" @submit.prevent="applyFilters">
      <!-- <div class="field">
        <label for="q">Search</label>
        <input id="q" v-model="query" placeholder="Place name…" />
      </div> -->

      <div class="field">
        <label for="category">Category</label>
        <select id="category" v-model="category" >
          <option value="">All categories</option>
          <option v-for="cat in categories" :key="cat" :value="cat">
            {{CATEGORY_LABELS[cat]}}
          </option>
        </select>
      </div>

      <div class="field">
        <label for="area">Area</label>
        <input id="area" v-model="area" placeholder="Kemang, Menteng…" />
      </div>

      <div class="filter-actions">
        <button class="primary" type="submit">Filter</button>
        <button type="button" @click="resetFilters">Reset</button>
      </div>
    </form>

    <div class="toolbar">
      <p class="muted">{{ places.places.length }} places</p>
      <button @click="showMap = !showMap">
        {{ showMap ? "Hide map" : "Show map" }}
      </button>
    </div>

    <PlaceMap v-if="showMap" :places="places.places" class="map-block"/>

    <p v-if="places.loading" class="muted">Loading...</p>
    <p v-else-if="places.error" class="error">{{ places.error }}</p>
    <p v-else-if="places.places.length===0" class="muted">No places match this filter yet.</p>

    <div v-else class="grid">
      <PlaceCard
        v-for="place in places.places"
        :key="place.id"
        :place="place"
        :saved="saved.savedPlaceIds.has(place.id)"
        :visited="visits.visitedPlaceIds.has(place.id)"
        :show-actions="auth.isAuthenticated"
        @toggle-saved="saved.toggle"
        @mark-visited="visits.markVisited"

      />
    </div>

  </div>
</template>

<style scoped>
.filters {
  display: grid;
  grid-template-columns: 1fr 1fr auto;
  gap: 0.9rem;
  align-items: end;
  margin: 1.5rem 0 1rem;
}

.filters .field {
  margin-bottom: 0;
}

.filter-actions {
  display: flex;
  gap: 0.5rem;
}

.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1rem;
}

.map-block {
  margin-bottom: 1.5rem;
}

.hero {
  display: grid;
  grid-template-columns: 1.05fr 0.95fr;
  gap: 2.5rem;
  align-items: center;
  padding: 2.5rem 0 1rem;
}

.eyebrow {
  margin: 0 0 0.6rem;
}

.hero-copy h1 {
  margin: 0 0 0.9rem;
  font-size: 2.3rem;
  line-height: 1.15;
  max-width: 18ch;
}

.hero-copy .lead {
  margin: 0 0 1.5rem;
  max-width: 46ch;
  font-size: 1.02rem;
}

.hero-actions {
  display: flex;
  gap: 0.9rem;
  flex-wrap: wrap;
}

.primary-link,
.secondary-link {
  display: inline-flex;
  align-items: center;
  border-radius: 8px;
  padding: 0.6rem 1.15rem;
  font-weight: 600;
  font-size: 0.95rem;
  text-decoration: none;
}

.primary-link {
  background: var(--accent);
  color: white;
}

.primary-link:hover {
  opacity: 0.92;
  text-decoration: none;
}

.secondary-link {
  border: 1px solid var(--line);
  color: var(--ink);
}

.secondary-link:hover {
  border-color: var(--accent);
  color: var(--accent);
  text-decoration: none;
}

.hero-art svg {
  width: 100%;
  height: auto;
  display: block;
}

.art-card {
  fill: var(--surface);
  stroke: var(--line);
  stroke-width: 1.5;
}

.art-route {
  fill: none;
  stroke: var(--line);
  stroke-width: 2.5;
  stroke-linecap: round;
  stroke-dasharray: 1 9;
}

.pin {
  stroke: var(--surface);
  stroke-width: 0.6;
}

.pin-cafe {
  fill: #b45309;
}
.pin-restaurant {
  fill: #be5b4a;
}
.pin-museum {
  fill: #4a6b63;
}
.pin-library {
  fill: #4a5b8c;
}

.icon-stroke {
  stroke: white;
  stroke-width: 1.1;
  stroke-linecap: round;
  fill: none;
}

.visited-badge {
  fill: #2f7a4f;
  stroke: var(--surface);
  stroke-width: 2;
}

.visited-check {
  fill: none;
  stroke: white;
  stroke-width: 2;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.art-polaroid {
  fill: var(--surface);
  stroke: var(--line);
  stroke-width: 1.2;
  filter: drop-shadow(0 6px 10px rgb(0 0 0 / 0.12));
}

.art-photo {
  fill: #efe3d3;
}

.art-photo-hill {
  fill: #c4906a;
  opacity: 0.85;
}

.art-photo-sun {
  fill: #e8b95c;
}

.grid {
  display: grid;
  gap: 1rem;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
}

@media (max-width: 860px) {
  .hero {
    grid-template-columns: 1fr;
    gap: 1.75rem;
    padding-top: 1.5rem;
    text-align: left;
  }

  .hero-copy h1 {
    font-size: 1.9rem;
    max-width: none;
  }

  .hero-art {
    max-width: 380px;
    margin: 0 auto;
  }
}

@media (max-width: 780px) {
  .filters {
    grid-template-columns: 1fr;
  }
}
</style>
