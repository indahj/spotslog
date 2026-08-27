import { defineStore } from "pinia";
import { ref } from "vue";
import { placesApi, type PlaceFilters } from "@/api";
import type { Place } from "@/api/types";

export const usePlacesStore = defineStore("places", () => {
  const places = ref<Place[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);

  async function loadHomepage() {
    loading.value = true;
    error.value = null;
    try {
      places.value = await placesApi.homepage();
    } catch (e) {
      error.value = e instanceof Error ? e.message : "Failed to load places";
    } finally {
      loading.value = false;
    }
  }

  async function search(filters: PlaceFilters) {
    loading.value = true;
    error.value = null;
    try {
      places.value = await placesApi.list(filters);
    } catch (e) {
      error.value = e instanceof Error ? e.message : "Failed to search places";
    } finally {
      loading.value = false;
    }
  }

  return { places, loading, error, loadHomepage, search };
});
