<script setup lang="ts">
import PlaceCard from '@/components/PlaceCard.vue';
import { useSavedStore } from '@/stores/saved';
import { useVisitsStore } from '@/stores/visits';
import { onMounted } from 'vue';


const saved = useSavedStore()
const visits = useVisitsStore()

onMounted(async () => {
  await Promise.all([saved.load(), visits.load()])
})
</script>

<template>

  <div class="container">
    <h1>Wishlist</h1>
    <p class="muted">Places you've saved to visit later.</p>

    <p v-if="saved.loading" class="muted">Loading...</p>
    <p v-else-if="saved.error" class="error">{{ saved.error }}</p>
    <p v-else-if="saved.saved.length === 0" class="muted">
      Nothing saved yet. Browse
      <RouterLink :to="{name: 'home'}">recommendations</RouterLink> and tap Save.
    </p>

    <div v-else class="grid">
      <PlaceCard
        v-for="entry in saved.saved"
        :key="entry.id"
        :place="entry.place!"
        :saved="true"
        :visited="visits.visitedPlaceIds.has(entry.place_id)"
        :show-actions="true"
        @toggle-saved="saved.toggle"
        @mark-visited="visits.markVisited"
      />
    </div>
  </div>
</template>

<style scoped>
h1 {
  margin-bottom: 0.25rem;
}

.grid {
  margin-top: 1.5rem;
  display: grid;
  gap: 1rem;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
}
</style>
