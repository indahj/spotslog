import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { visitsApi } from "@/api";
import type { Visit } from "@/api/types";

export const useVisitsStore = defineStore("visits", () => {
  const visits = ref<Visit[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);

  const visitedPlaceIds = computed(() => new Set(visits.value.map((v) => v.place_id)));

  async function load() {
    loading.value = true;
    error.value = null;
    try {
      visits.value = await visitsApi.list();
    } catch (e) {
      error.value = e instanceof Error ? e.message : "Failed to load visits";
    } finally {
      loading.value = false;
    }
  }

  async function markVisited(placeId: number, notes?: string, visitedAt?: string) {
    await visitsApi.create(placeId, notes, visitedAt);
    await load();
  }

  async function remove(id: number) {
    await visitsApi.remove(id);
    visits.value = visits.value.filter((v) => v.id !== id);
  }

  async function uploadPhoto(visitId: number, file: File, caption?: string) {
    await visitsApi.uploadPhoto(visitId, file, caption);
    await load();
  }

  return { visits, loading, error, visitedPlaceIds, load, markVisited, remove, uploadPhoto };
});
