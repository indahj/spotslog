<script setup lang="ts">
import { useAuthStore } from '@/stores/auth';
import { ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';


const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const email = ref("")
const password = ref("")
const error = ref<string | null>(null)
const submitting = ref(false)

async function submit() {
  error.value = null;
  submitting.value = true;
  try {
    await auth.login(email.value, password.value)
    const redirect = route.query.redirect as string | undefined
    router.push(redirect ?? {name: "home"})
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Login failed"
  } finally {
    submitting.value = false;
  }
}

</script>

<template>
  <div class="container narrow">
    <h1>Log in</h1>

    <form class="card" @submit.prevent="submit">
      <div class="field">
        <label for="email">Email</label>
        <input id="email" v-model="email" type="email" required autocomplete="email" />
      </div>

      <div class="field">
        <label for="password">Password</label>
        <input
          id="password"
          v-model="password"
          type="password"
          required
          autocomplete="current-password"
        />
      </div>

      <p v-if="error" class="error"> {{ error }}</p>

      <button class="primary" type="submit" :disabled="submitting">
        {{ submitting ? "Logging in..." : "Log in" }}
      </button>
    </form>

    <p class="muted">
      No account yet? <RouterLink :to="{ name: 'register' }">Sign up</RouterLink>
    </p>
  </div>
</template>

<style scoped>
.narrow {
  max-width: 420px;
}
</style>
