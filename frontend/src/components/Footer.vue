<script setup lang="ts">
import githubIcon from "@/assets/github-icon.svg";
import linkedinIcon from "@/assets/linkedin-icon.svg";
</script>

<template>
  <footer class="site-footer">
    <div class="footer-inner">
      <p class="muted">Est. 2026 - Belle's portofolio</p>

      <div class="footer-links">
        <span class="muted">find me on:</span>
        <a href="#" target="'_blank'" rel="noopener noreferrer" aria-label="GitHub">
          <img :src="githubIcon" alt="GitHub" width="20" height="20">
        </a>
        <a href="#" target="'_blank'" rel="noopener noreferrer" aria-label="Linkedln">
          <img :src="linkedinIcon" alt="linkedln" width="20" height="20">
        </a>
      </div>
    </div>
  </footer>

</template>

<style scoped>
.site-footer {
  border-top: 1px solid var(--line);
  margin-top: 3rem;
}

.footer-inner {
  max-width: 1100px;
  margin: 0 auto;
  padding: 1.5rem 1rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  flex-wrap: wrap;
}

.footer-links {
  display: flex;
  align-items: center;
  gap: 0.9rem;
  font-size: 0.9rem;
}

.footer-links a {
  display: inline-flex;
  opacity: 0.75;
}

.footer-links a:hover {
  opacity: 1;
}
</style>
