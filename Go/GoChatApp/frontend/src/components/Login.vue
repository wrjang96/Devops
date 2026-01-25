<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const username = ref('')
const router = useRouter()

const enterLobby = () => {
  if (username.value.trim()) {
    localStorage.setItem('chat_username', username.value.trim())
    router.push('/rooms')
  }
}
</script>

<template>
  <div class="login-container">
    <div class="login-card">
      <h1 class="logo">GOCHAT<span class="red-dot">.</span></h1>
      <p class="subtitle">Enter the Stadium</p>
      
      <form @submit.prevent="enterLobby" class="login-form">
        <input 
          v-model="username" 
          type="text" 
          placeholder="Enter your username" 
          class="username-input"
          autofocus
        />
        <button type="submit" class="enter-btn">ENTER LOBBY</button>
      </form>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background: var(--primary-blue);
  padding: 1rem;
}

.login-card {
  background: white;
  padding: 3rem;
  border-radius: 8px;
  box-shadow: 0 10px 25px rgba(0,0,0,0.2);
  text-align: center;
  width: 100%;
  max-width: 400px;
}

.logo {
  font-size: 2.5rem;
  font-weight: 900;
  color: var(--primary-blue);
  letter-spacing: -1px;
  margin-bottom: 0.5rem;
}

.red-dot {
  color: var(--primary-red);
}

.subtitle {
  color: #666;
  margin-bottom: 2rem;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 1px;
  font-size: 0.9rem;
}

.username-input {
  width: 100%;
  padding: 1rem;
  font-size: 1.1rem;
  border: 2px solid #ddd;
  border-radius: 4px;
  margin-bottom: 1rem;
  box-sizing: border-box; 
  transition: border-color 0.2s;
}

.username-input:focus {
  outline: none;
  border-color: var(--primary-blue);
}

.enter-btn {
  width: 100%;
  padding: 1rem;
  background: var(--primary-red);
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1.1rem;
  font-weight: 700;
  cursor: pointer;
  transition: background 0.2s;
  text-transform: uppercase;
}

.enter-btn:hover {
  background: #a00b34;
}
</style>
