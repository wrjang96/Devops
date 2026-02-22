<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { API_URL } from '../config/api'

const router = useRouter()
const rooms = ref([])
const newRoomName = ref('')
const username = localStorage.getItem('chat_username') || 'Anonymous'
const token = localStorage.getItem('chat_token')

const fetchRooms = async () => {
    try {
        // Use rh.List endpoint to get detailed info (creator, member counts)
        const res = await fetch(`${API_URL}/chatrooms`, {
            headers: { 'Authorization': `Bearer ${token}` }
        })
        if (res.ok) {
            rooms.value = await res.json()
        } else if (res.status === 401) {
            router.push('/')
        }
    } catch (e) {
        console.error("Failed to fetch rooms", e)
    }
}

onMounted(() => {
    if (!token) {
        router.push('/')
        return
    }
    fetchRooms()
    setInterval(fetchRooms, 3000)
})

const joinRoom = async (room) => {
    try {
        const res = await fetch(`${API_URL}/chatrooms/${room.id}/join`, {
            method: 'POST',
            headers: { 'Authorization': `Bearer ${token}` }
        })
        if (res.ok) {
            router.push({ 
                path: '/chat', 
                query: { 
                    username: username,
                    room: room.name,
                    id: room.id
                } 
            })
        }
    } catch (e) {
        console.error("Failed to join", e)
    }
}

const createRoom = async () => {
    if (!newRoomName.value.trim()) return

    try {
        const res = await fetch(`${API_URL}/chatrooms`, {
            method: 'POST',
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}` 
            },
            body: JSON.stringify({ name: newRoomName.value.trim() })
        })
        
        if (res.ok) {
            const room = await res.json()
            newRoomName.value = ''
            await fetchRooms() // Refresh list
            // Optionally auto-join
            await joinRoom(room)
        }
    } catch (e) {
        console.error("Failed to create", e)
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
                    <div v-for="room in rooms" :key="room.id" class="room-card">
                        <div class="room-info">
                            <h3>{{ room.name }}</h3>
                            <div class="room-meta">
                                <span class="creator">Host: {{ room.creatorName }}</span>
                            </div>
                            <div class="room-stat">
                                <span class="count">{{ room.members }}</span>
                                <span class="label">Players</span>
                            </div>
                        </div>
                        <button 
                            @click="joinRoom(room)" 
                            class="join-btn"
                        >
                            JOIN
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

.room-meta {
    font-size: 0.85rem;
    color: #888;
    margin-bottom: 0.5rem;
}

.room-stat {
    color: #666;
    font-size: 0.9rem;
    font-weight: bold;
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

.no-rooms {
    text-align: center;
    color: #888;
    margin-top: 2rem;
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
