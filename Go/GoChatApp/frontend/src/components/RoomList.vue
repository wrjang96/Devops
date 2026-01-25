<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const rooms = ref([])
const newRoomName = ref('')
const username = localStorage.getItem('chat_username') || 'Anonymous'

const fetchRooms = async () => {
    try {
        const res = await fetch('http://localhost:3434/rooms')
        if (res.ok) {
            rooms.value = await res.json()
        }
    } catch (e) {
        console.error("Failed to fetch rooms", e)
    }
}

onMounted(() => {
    if (!localStorage.getItem('chat_username')) {
        router.push('/')
        return
    }
    fetchRooms()
    setInterval(fetchRooms, 5000)
})

const joinRoom = (roomName) => {
    router.push({ 
        path: '/chat', 
        query: { 
            username: username,
            room: roomName 
        } 
    })
}

const createRoom = () => {
    if (newRoomName.value.trim()) {
        joinRoom(newRoomName.value.trim())
    }
}
</script>

<template>
    <div class="room-list-layout">
        <header class="main-header">
            <div class="header-inner">
                <h1>GOCHAT LOBBY</h1>
                <div class="user-info">Signed in as: <strong>{{ username }}</strong></div>
            </div>
        </header>

        <main class="content-area">
            <div class="create-section">
                <h2>Create New Stadium</h2>
                <form @submit.prevent="createRoom" class="create-form">
                    <input 
                        v-model="newRoomName" 
                        type="text" 
                        placeholder="Enter Room Name" 
                        class="room-input"
                    />
                    <button type="submit" class="create-btn">CREATE</button>
                </form>
            </div>

            <div class="rooms-section">
                <h2>Live Stadiums</h2>
                <div v-if="rooms.length === 0" class="no-rooms">
                    No active games. Create one above!
                </div>
                
                <div class="rooms-grid">
                    <div v-for="room in rooms" :key="room.name" class="room-card">
                        <div class="room-info">
                            <h3>{{ room.name }}</h3>
                            <div class="room-stat">
                                <span class="count">{{ room.count }} / 50</span>
                                <span class="label">Fans</span>
                            </div>
                        </div>
                        <button 
                            @click="joinRoom(room.name)" 
                            class="join-btn"
                            :disabled="room.count >= 50"
                        >
                            {{ room.count >= 50 ? 'FULL' : 'JOIN' }}
                        </button>
                    </div>
                </div>
            </div>
        </main>
    </div>
</template>

<style scoped>
.room-list-layout {
    min-height: 100vh;
    background: #f4f6f8;
}

.main-header {
    background: var(--primary-blue);
    color: white;
    padding: 1rem;
    box-shadow: 0 2px 5px rgba(0,0,0,0.1);
}

.header-inner {
    max-width: 1200px;
    margin: 0 auto;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

h1 {
    font-size: 1.5rem;
    margin: 0;
}

.content-area {
    max-width: 1200px;
    margin: 2rem auto;
    padding: 0 1rem;
}

.create-section {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0,0,0,0.05);
    margin-bottom: 2rem;
}

.create-form {
    display: flex;
    gap: 1rem;
    margin-top: 1rem;
    max-width: 600px;
}

.room-input {
    flex: 1;
    padding: 0.8rem;
    border: 2px solid #ddd;
    border-radius: 4px;
    font-size: 1rem;
}

.create-btn {
    padding: 0 2rem;
    background: var(--primary-red);
    color: white;
    border: none;
    border-radius: 4px;
    font-weight: 700;
    cursor: pointer;
}

.rooms-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 1.5rem;
    margin-top: 1rem;
}

.room-card {
    background: white;
    padding: 1.5rem;
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0,0,0,0.05);
    display: flex;
    justify-content: space-between;
    align-items: center;
    transition: transform 0.2s;
}

.room-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 5px 15px rgba(0,0,0,0.1);
}

.room-info h3 {
    margin: 0 0 0.5rem 0;
    color: var(--primary-blue);
}

.room-stat {
    color: #666;
    font-size: 0.9rem;
}

.join-btn {
    padding: 0.6rem 1.5rem;
    background: var(--primary-blue);
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-weight: 600;
}

.join-btn:disabled {
    background: #ccc;
    cursor: not-allowed;
}

@media (max-width: 768px) {
    .create-form {
        flex-direction: column;
    }
    
    .create-btn {
        width: 100%;
        padding: 1rem;
    }
}
</style>
