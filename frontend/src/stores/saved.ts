import { defineStore } from "pinia";
import { ref, computed } from "vue";
import type { SavedPlace } from "@/api/types";
import { savedApi } from "@/api";


export const useSavedStore = defineStore("saved", () => {
  const saved = ref<SavedPlace[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);

  const savedPlaceIds = computed(() => new Set(saved.value.map((s) => s.place_id)));

  async function load() {
    loading.value = true;
    error.value = null;
    try {
      saved.value = await savedApi.list();
    } catch (e) {
      error.value = e instanceof Error ? e.message : "Falied to load wishlist";
    } finally {
      loading.value = false;
    }
  }

  async function toggle(placeId: number) {
    if (savedPlaceIds.value.has(placeId)) {
      await savedApi.remove(placeId);
      saved.value = saved.value.filter((s) => s.place_id !== placeId);
    } else {
      await savedApi.add(placeId);
      await load();
    }
  }

  return { saved, loading, error, savedPlaceIds, load, toggle }

})
