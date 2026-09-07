import { ratingsApi } from "@/api";
import { defineStore } from "pinia";
import { ref } from "vue";

export const useRatingsStore = defineStore("ratings", () => {
  const myRatings = ref<Map<number, number>>(new Map())
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchMine(placeId: number) {
    loading.value = true
    error.value = null
    try {
      const result = await ratingsApi.getMine(placeId)
      if (result.rating !== null) {
        myRatings.value.set(placeId, result.rating)
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to load rating"
    } finally {
      loading.value = false
    }
  }

  async function rate(placeId: number, rating: number) {
    error.value = null
    try {
      const result = await ratingsApi.rate(placeId, rating)
      myRatings.value.set(placeId, result.rating)
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to save rating"
    }
  }

  return { myRatings, loading, error, fetchMine, rate }
})
