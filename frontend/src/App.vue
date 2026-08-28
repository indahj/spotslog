<script setup lang="ts">
import { RouterLink, RouterView, useRouter } from 'vue-router'
import { useAuthStore } from './stores/auth';

const auth = useAuthStore()
const router = useRouter()

function handleLogout() {
  auth.logout();
  router.push({name: "home"})
}


</script>

<template>
  <header class="site-header">
    <div class="header-inner">
      <RouterLink :to="{ name: 'home' }" class="brand">Spotslog</RouterLink>

      <nav>
        <RouterLink :to="{ name: 'home' }">Discover</RouterLink>
        <template v-if="auth.isAuthenticated">
          <RouterLink :to="{ name: 'visits'}">My Visits </RouterLink>
          <RouterLink :to="{ name: 'wishlist'}">Wishlist </RouterLink>
          <RouterLink :to="{ name: 'add-place'}">Add Place </RouterLink>
          <button @click="handleLogout">Log out</button>
        </template>
        <template v-else>
          <RouterLink :to="{ name: 'login' }">Log in</RouterLink>
          <RouterLink :to="{ name: 'register' }">Sign up</RouterLink>
        </template>
      </nav>
    </div>
  </header>

  <main>
    <RouterView />
  </main>
</template>

<style scoped>
.site-header {
  background: var(--surface);
  border-bottom: 1px solid var(--line);
}

.header-inner {
  max-width: 1100px;
  margin: 0 auto;
  padding: 0.9rem 1rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.brand {
  font-weight: 700;
  font-size: 1.15rem;
  color: var(--ink);
}

nav {
  display: flex;
  align-items: center;
  gap: 1rem;
  font-size: 0.95rem;
}

nav a.router-link-active {
  font-weight: 600;
}
</style>
