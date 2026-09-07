<script setup lang="ts">
import type { Place } from '@/api/types';
import { onBeforeUnmount, onMounted, ref } from 'vue';
import placeholderImage from "@/assets/place-placeholder.svg"


const props = defineProps<{
  place: Place[]
}>()

const currentIndex = ref(0)

function next() {
  currentIndex.value  = (currentIndex.value + 1) % props.place.length
}

let timer: ReturnType<typeof setInterval> | undefined

onMounted(() => {
  timer = setInterval(next, 4000)
})

onBeforeUnmount(() => {
  clearInterval(timer)
})
</script>

<template>
  <div class="slideshow">
    <div class="polaroid">
      <div class="photo-frame">
        <img :src="props.place[currentIndex]?.cover_photo_url ?? placeholderImage" :alt="props.place[currentIndex]?.name">
        <span class="badge-newest">Newest!</span>
        <div class="dots">
          <span
            v-for="(p, i) in props.place"
            :key="i"
            class="dot"
            :class="{ active: i === currentIndex }"
          />
        </div>
      </div>
      <p class="slide-name">{{ props.place[currentIndex]?.name }}</p>
    </div>
  </div>
</template>

<style scoped>
.slideshow {
  width: 100%;
  height: 100%;
}

.polaroid {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  background: white;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 10px 30px rgb(0 0 0 / 0.15);
}

.photo-frame {
  position: relative;
  flex: 1;
  overflow: hidden;
}

.photo-frame img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.badge-newest {
  position: absolute;
  top: 12px;
  left: 12px;
  background: var(--accent);
  color: var( --accent-soft);
  padding: 0.3rem 0.8rem;
  border-radius: 999px;
  font-size: 0.8rem;
  font-weight: 600;
}

.dots {
  position: absolute;
  top: 16px;
  right: 16px;
  display: flex;
  gap: 0.4rem;
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  border: 1.5px solid black;
  background: transparent;
}

.dot.active {
  background: var(--accent);
}

.slide-name {
  margin: 0;
  padding: 1rem;
  font-weight: 600;
  font-size: 1.05rem;
}
</style>
