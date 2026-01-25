<script setup>
</script>

<template>
  <div class="app-container">
    <router-view></router-view>
  </div>
</template>

<style>
:root {
  --primary-blue: #041E42; /* MLB Navy */
  --primary-red: #BF0D3E;  /* MLB Red */
  --bg-color: #f4f4f4;
  --text-main: #333;
}

body {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  background-color: var(--bg-color);
  color: var(--text-main);
}

.app-container {
  width: 100%;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}
</style>
