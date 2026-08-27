<script setup lang="ts">
import { useAuthStore } from '@/stores/auth';
import { ref } from 'vue';
import { useRouter } from 'vue-router';


const auth = useAuthStore()
const router = useRouter()

const name = ref("")
const email = ref("")
const password = ref("")
const error = ref<string | null>(null)
const submitting = ref(false)

async function submit() {
  error.value = null;
  submitting.value = true;
  try {
    await auth.register(name.value, email.value, password.value)
    router.push({name: "home"})
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Registration failed."
  } finally {
    submitting.value = false;
  }

}

</script>
<template>
  <div class="container narrow">
    <h1>Create an account</h1>

    <form class="card" @submit.prevent="submit">
      <div class="field">
        <label for="name">Name</label>
        <input id="name" v-model="name" required autocomplete="name" />
      </div>

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
          minlength="8"
          autocomplete="new-password"
        />
        <small class="muted">At least 8 characters.</small>
      </div>

      <p v-if="error" class="error">{{ error }}</p>

      <button class="primary" type="submit" :disabled="submitting">
       {{ submitting ? "Creating account..." : "Sign up" }}
      </button>
    </form>

    <p class="muted">
      Already have an account? <RouterLink :to="{ name: 'login' }">Log in</RouterLink>
    </p>
  </div>
</template>

<style scoped>
.narrow {
  max-width: 420px;
}
</style>
