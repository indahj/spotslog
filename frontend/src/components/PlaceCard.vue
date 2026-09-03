<script setup lang="ts">
import { CATEGORY_LABELS, type Place } from '@/api/types';
import { RouterLink } from 'vue-router';
import placeholderImage from "@/assets/place-placeholder.svg";


defineProps<{
  place: Place
  saved?: boolean
  visited?: boolean
  showActions?: boolean
}>()

const emit = defineEmits<{
  toggleSaved: [placeId: number]
  markVisited: [placeId: number]
}>()

</script>
<template>
  <article class="card place-card">
    <img
      :src="place.cover_photo_url ?? placeholderImage"
      :alt="place.name"
      class="cover-photo"
    />

    <header>
      <RouterLink :to="{ name: 'place-detail', params: { id: place.id} }">
        <h3>{{ place.name }}</h3>
      </RouterLink>
      <span class="badge">{{ CATEGORY_LABELS[place.category] }}</span>
    </header>

    <p class="muted address">{{ place.address }}</p>
    <p v-if="place.district" class="muted district">{{ place.district }}</p>
    <p v-if="place.price_range" class="price">{{ place.price_range }}</p>

    <footer v-if="showActions">
      <button :class="{ active: saved }" @click="emit('toggleSaved', place.id)">
        {{ saved? "★ Saved" : "☆ Save"  }}
      </button>
      <button :disabled="visited" @click="emit('markVisited', place.id)">
        {{visited ? "✓ Visited" : "Mark visited"  }}
      </button>
    </footer>

  </article>
</template>

<style scoped>
.place-card {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.cover-photo {
  width: 100%;
  aspect-ratio: 16 / 10;
  object-fit: cover;
  border-radius: 8px;
  display: block;
}


header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 0.5rem;
}

h3 {
  margin: 0;
  font-size: 1.05rem;
  color: var(--ink);
}

.address, .district {
  margin: 0;
  font-size: 0.87rem;
}

.price {
  margin: 0;
  font-size: 0.87rem;
  font-weight: 600;
}

footer {
  margin-top: auto;
  padding-top: 0.7rem;
  display: flex;
  gap: 0.5rem;
}

footer button {
  font-size: 0.85rem;
  padding: 0.35rem 0.7rem;
}

button.active {
  border-color: var(--accent);
  color: var(--accent);
}

</style>
