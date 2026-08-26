import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { authApi } from "@/api";
import { clearToken, getToken, setToken } from "@/api/client";
import type { User } from "@/api/types";

export const useAuthStore = defineStore("auth", () => {
  const user = ref<User | null>(null);
  const loading = ref(false);

  const isAuthenticated = computed(() => user.value !== null);
  const isAdmin = computed(() => user.value?.role === "admin");

  async function register(name: string, email: string, password: string) {
    const { user: u, token } = await authApi.register(name, email, password);
    setToken(token);
    user.value = u;
  }

  async function login(email: string, password: string) {
    const { user: u, token } = await authApi.login(email, password);
    setToken(token);
    user.value = u;
  }

  function logout() {
    clearToken();
    user.value = null;
  }

  async function restore() {
    if (!getToken()) return;
    loading.value = true;
    try {
      user.value = await authApi.me();
    } catch {
      clearToken();
      user.value = null;
    } finally {
      loading.value = false;
    }
  }

  return { user, loading, isAuthenticated, isAdmin, register, login, logout, restore };
});
