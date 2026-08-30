<script setup lang="ts">
import { ref } from 'vue';


const props = defineProps<{
  upload: (file: File, caption?: string) => Promise<unknown>
}>()

const file = ref<File | null>(null)
const caption = ref("")
const uploading = ref(false)
const error = ref<string | null>(null)

function onFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  file.value = input.files?.[0] ?? null
}

async function submit() {
  if (!file.value) return
  uploading.value = true
  error.value = null
  try {
    await props.upload(file.value, caption.value || undefined)
    file.value = null
    caption.value = ""
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Upload failed"
  } finally {
    uploading.value = false
  }
}

</script>

<template>
  <form class="uploader" @submit.prevent="submit">
    <input type="file" accept="image/*" @change="onFileChange"/>
    <input v-model="caption" placeholder="Caption (optional)"/>
    <button type="submit" :disabled="!file || uploading">
      {{ uploading ? "Uploading..." : "Upload photo" }}
    </button>
    <p v-if="error" class="error">{{ error }}</p>
  </form>

</template>

<style scoped>
.uploader {
  display: flex;
  gap: 0.5rem;
  align-items: center;
  flex-wrap: wrap;
  font-size: 0.9rem;
}

.uploader input[type="file"] {
  border: none;
  padding: 0;
  width: auto;
}

.uploader input:not([type="file"]) {
  width: auto;
  flex: 1;
  min-width: 160px;
}
</style>
