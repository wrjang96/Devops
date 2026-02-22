<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { API_URL } from '../config/api'

const mode = ref('login') // 'login' or 'register'
const username = ref('')
const password = ref('')
const errorMsg = ref('')
const router = useRouter()

const toggleMode = () => {
  mode.value = mode.value === 'login' ? 'register' : 'login'
  errorMsg.value = ''
  password.value = ''
}

const handleSubmit = async () => {
  errorMsg.value = ''
  if (!username.value.trim() || !password.value.trim()) {
    errorMsg.value = 'Please enter both ID and password.'
    return
  }

  const endpoint = mode.value === 'login' ? '/login' : '/register'
  const url = `${API_URL}${endpoint}`

  try {
    const res = await fetch(url, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id: username.value, password: password.value })
    })

    const data = await res.json()

    if (!res.ok) {
      errorMsg.value = data.error || 'Action failed'
      return
    }

    if (mode.value === 'register') {
      // Auto login or switch to login
      mode.value = 'login'
      errorMsg.value = 'Registration successful! Please log in.'
      password.value = ''
    } else {
      // Login successful
      localStorage.setItem('chat_token', data.accessToken)
      localStorage.setItem('chat_username', username.value)
      router.push('/rooms')
    }

  } catch (e) {
    console.error(e)
    errorMsg.value = 'Network error'
  }
}
</script>

<template>
  <div class="login-container">
    <div class="login-card">
      <h1 class="logo">GOCHAT<span class="red-dot">.</span></h1>
      <p class="subtitle">{{ mode === 'login' ? 'Enter the Stadium' : 'Join the League' }}</p>
      
      <form @submit.prevent="handleSubmit" class="login-form">
        <input 
          v-model="username" 
          type="text" 
          placeholder="ID" 
          class="input-field"
          autofocus
        />
        <input 
          v-model="password" 
          type="password" 
          placeholder="Password" 
          class="input-field"
        />
        <div v-if="errorMsg" class="error-msg">{{ errorMsg }}</div>
        
        <button type="submit" class="enter-btn">
          {{ mode === 'login' ? 'LOGIN' : 'REGISTER' }}
        </button>
      </form>

      <div class="toggle-link">
        <span v-if="mode === 'login'">New player? <a @click="toggleMode">Create account</a></span>
        <span v-else>Already have a pass? <a @click="toggleMode">Log in</a></span>
      </div>
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

.input-field {
  width: 100%;
  padding: 1rem;
  font-size: 1.1rem;
  border: 2px solid #ddd;
  border-radius: 4px;
  margin-bottom: 1rem;
  box-sizing: border-box; 
  transition: border-color 0.2s;
}

.input-field:focus {
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
  margin-bottom: 1rem;
}

.enter-btn:hover {
  background: #a00b34;
}

.error-msg {
  color: red;
  margin-bottom: 1rem;
  font-size: 0.9rem;
}

.toggle-link {
  font-size: 0.9rem;
  color: #666;
}

.toggle-link a {
  color: var(--primary-blue);
  cursor: pointer;
  text-decoration: underline;
  font-weight: 600;
}
</style>
