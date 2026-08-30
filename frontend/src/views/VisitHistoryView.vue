<script setup lang="ts">
import { CATEGORY_LABELS } from '@/api/types';
import PhotoUploader from '@/components/PhotoUploader.vue';
import { useVisitsStore } from '@/stores/visits';
import { onMounted, ref } from 'vue';


const visits = useVisitsStore()
const expandedId = ref<number | null>(null)

function toggleExpanded(id: number) {
  expandedId.value = expandedId.value === id ? null : id;
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString("en-GB", {
    day: "numeric",
    month: "short",
    year: "numeric"
  })
}

onMounted(() => visits.load())

</script>

<template>
  <div class="container">
    <header class="page-header">
      <div>
        <h1>My visits</h1>
        <p class="muted">Every place you've marked as visited.</p>
      </div>

      <RouterLink :to="{name: 'add-place'}">
        <button class="primary">Add a place</button>
      </RouterLink>
    </header>

    <p v-if="visits.loading" class="muted">Loading...</p>
    <p v-else-if="visits.error" class="error">{{ visits.error }}</p>
    <p v-else-if="visits.visits.length ===0" class="muted">
      No visits recorded yet. Browse
      <RouterLink :to="{name: 'home'}"> recommendations</RouterLink> or
      <RouterLink :to="{name: 'add-place'}">add your own place</RouterLink>
    </p>

    <ul v-else class="visit-list">
      <li v-for="visit in visits.visits" :key="visit.id" class="card">
        <div class="visit-head">
          <div>
            <h3 v-if="visit.place">
              <RouterLink :to="{name: 'place-detail', params: {id: visit.place_id}}">{{ visit.place.name }}</RouterLink>
            </h3>
            <span v-if="visit.place" class="badge">{{  CATEGORY_LABELS[visit.place.category] }}</span>
            <p class="muted date">Visited {{ formatDate(visit.visited_at) }}</p>
            <p v-if="visit.place" class="muted address">{{ visit.place?.address }}</p>
          </div>

          <div class="visit-actions">
            <button @click="toggleExpanded(visit.id)">
              {{ expandedId === visit.id ? "Close" : "Add photo" }}
            </button>
            <button @click="visits.remove(visit.id)">Remove</button>
          </div>
        </div>

        <p v-if="visit.notes">{{ visit.notes }}</p>

        <div v-if="expandedId === visit.id" class="uploader-block">
          <PhotoUploader
            :upload="(file, caption) => visits.uploadPhoto(visit.id, file, caption)"
          />
        </div>

        <div v-if="visit.photos && visit.photos.length > 0" class="photo-strip">
          <img
            v-for="photo in visit.photos"
            :key="photo.id"
            :src="photo.url"
            :alt="photo.caption ?? 'Visit photo'"
            loading="lazy"/>
        </div>
      </li>
    </ul>
  </div>

</template>

<style scoped>

.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
  flex-wrap: wrap;
}

.page-header h1 {
  margin-bottom: 0.25rem;
}

.visit-list {
  list-style: none;
  padding: 0;
  margin: 1.5rem 0 0;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.visit-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
}

h3 {
  color: var(--ink);
}

.date,
.address {
  margin: 0.25rem 0 0;
  font-size: 0.87rem;
}

.visit-actions {
  display: flex;
  gap: 0.5rem;
}

.visit-actions button {
  font-size: 0.85rem;
  padding: 0.35rem 0.7rem;
}

.notes {
  margin: 0.7rem 0 0;
  font-size: 0.92rem;
}

.uploader-block {
  margin-top: 0.9rem;
  padding-top: 0.9rem;
  border-top: 1px solid var(--line);
}

.photo-strip {
  display: flex;
  gap: 0.6rem;
  margin-top: 0.9rem;
  flex-wrap: wrap;
}

.photo-strip img {
  width: 120px;
  height: 120px;
  object-fit: cover;
  border-radius: 8px;
}

</style>
