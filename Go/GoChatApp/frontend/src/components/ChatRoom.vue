<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import WebSocketService from '../services/websocket'
import { API_URL, WS_URL } from '../config/api'

const route = useRoute()
const router = useRouter()
const username = route.query.username
const room = route.query.room || 'general'
const messages = ref([])
const newMessage = ref('')
const messagesContainer = ref(null)
let wsService = null

const scrollToBottom = async () => {
  await nextTick()
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

onMounted(() => {
  if (!username) {
    router.push('/')
    return
  }

  // Fetch existing messages
  fetch(`${API_URL}/chatrooms/${encodeURIComponent(room)}/messages`, {
      headers: { 'Authorization': `Bearer ${localStorage.getItem('chat_token')}` }
  })
  .then(res => res.json())
  .then(data => {
      if (data && data.messages) {
          messages.value = data.messages.map(m => ({
              id: m.id || Date.now(),
              sender: m.sender,
              content: m.content,
              room: m.room,
              type: 'message'
          }))
          scrollToBottom()
      }
  })
  .catch(err => console.error("Failed to load history", err))

  // Pass room and username in query
  wsService = new WebSocketService(`${WS_URL}/ws?room=${encodeURIComponent(room)}&username=${encodeURIComponent(username)}`)
  
  wsService.connect(
    () => {
        // Connected
    },
    () => {
        // Closed
        messages.value.push({
            id: Date.now(),
            sender: 'System',
            content: 'Connection closed.',
            type: 'system'
        })
    }
  )

  wsService.addListener((rawMessage) => {
    // Check if it's JSON (expected now)
    try {
        const msgObj = JSON.parse(rawMessage)
        
        messages.value.push({
          id: Date.now(), 
          sender: msgObj.sender, 
          content: msgObj.content,
          room: msgObj.room,
          type: 'message'
        })
    } catch(e) {
        // Fallback for non-JSON (if any legacy)
        messages.value.push({
            id: Date.now(),
            sender: 'Unknown',
            content: rawMessage,
            type: 'message'
        })
    }
    
    scrollToBottom()
  })
})

onUnmounted(() => {
  if (wsService) {
    wsService.close()
  }
})

const sendMessage = () => {
  if (newMessage.value.trim() && wsService) {
    // We send just the content text string, the client.go handles wrapping it into JSON with sender info
    // Wait, implementation plan says user `writePump` logic sends JSON, but frontend sends text? 
    // Let's verify: In client.go, readPump reads messageBytes and wraps into Message struct.
    // So sending raw text IS correct.
    wsService.sendMessage(newMessage.value)
    newMessage.value = ''
  }
}

const exitRoom = async () => {
    const roomId = route.query.id
    const token = localStorage.getItem('chat_token')
    
    if (roomId && token) {
        try {
            await fetch(`${API_URL}/chatrooms/${roomId}/leave`, {
                method: 'POST',
                headers: { 'Authorization': `Bearer ${token}` }
            })
        } catch (e) {
            console.error("Failed to leave", e)
        }
    }
    
    router.push('/rooms')
}

const isMe = (sender) => {
    return sender === username
}
</script>

<template>
  <div class="chat-layout">
    <header class="chat-header">
      <div class="header-content">
        <div class="left-head">
            <button @click="exitRoom" class="back-btn">← Lobby</button>
            <h1 class="room-title">{{ room }}</h1>
        </div>
        <div class="user-badge">{{ username }}</div>
      </div>
    </header>

    <main class="chat-field" ref="messagesContainer">
      <div v-if="messages.length === 0" class="empty-state">
        Start the conversation...
      </div>
      
      <div 
        v-for="msg in messages" 
        :key="msg.id" 
        class="message-row"
        :class="{ 
            'my-msg-row': isMe(msg.sender),
            'other-msg-row': !isMe(msg.sender) && msg.type !== 'system',
            'system-msg-row': msg.type === 'system'
        }"
      >
        <div v-if="msg.type === 'system'" class="system-bubble">
            {{ msg.content }}
        </div>

        <div v-else class="msg-group">
            <span v-if="!isMe(msg.sender)" class="sender-name">{{ msg.sender }}</span>
            <div class="message-bubble" :class="{ 'my-bubble': isMe(msg.sender), 'other-bubble': !isMe(msg.sender) }">
                {{ msg.content }}
            </div>
        </div>
      </div>
    </main>

    <footer class="chat-input-area">
      <form @submit.prevent="sendMessage" class="input-form">
        <input 
          v-model="newMessage" 
          type="text" 
          placeholder="Message..." 
          class="msg-input"
        />
        <button type="submit" class="send-btn">
            Send
        </button>
      </form>
    </footer>
  </div>
</template>

<style scoped>
.chat-layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: white;
}

.chat-header {
  background: white;
  border-bottom: 1px solid #ddd;
  padding: 0.8rem 1rem;
  flex-shrink: 0;
  z-index: 10;
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.left-head {
    display: flex;
    align-items: center;
    gap: 1rem;
}

.back-btn {
    background: none;
    border: none;
    font-size: 1rem;
    color: var(--primary-blue);
    cursor: pointer;
    font-weight: 600;
}

.room-title {
  font-size: 1.2rem;
  margin: 0;
  font-weight: 700;
}

.user-badge {
    font-size: 0.9rem;
    color: #666;
}

.chat-field {
  flex: 1;
  overflow-y: auto;
  padding: 1rem;
  background: #fff;
  display: flex;
  flex-direction: column;
  /* Max width for readability on large screens */
  max-width: 1200px; 
  width: 100%;
  margin: 0 auto;
  box-sizing: border-box;
}

.empty-state {
    text-align: center;
    color: #999;
    margin-top: 2rem;
}

.message-row {
    margin-bottom: 0.5rem;
    width: 100%;
    display: flex;
}

.my-msg-row {
    justify-content: flex-end;
}

.other-msg-row {
    justify-content: flex-start;
}

.system-msg-row {
    justify-content: center;
}

.system-bubble {
    background: #f0f0f0;
    padding: 0.5rem 1rem;
    border-radius: 20px;
    font-size: 0.8rem;
    color: #666;
}

.msg-group {
    display: flex;
    flex-direction: column;
    max-width: 70%;
}

.my-msg-row .msg-group {
    align-items: flex-end;
}

.other-msg-row .msg-group {
    align-items: flex-start;
}

.sender-name {
    font-size: 0.75rem;
    color: #888;
    margin-bottom: 0.2rem;
    margin-left: 0.5rem; /* Indent slightly for other's name */
}

.message-bubble {
    padding: 0.8rem 1rem;
    border-radius: 18px;
    font-size: 1rem;
    line-height: 1.4;
    word-break: break-word;
}

.other-bubble {
    background: #efefef;
    color: black;
    border-bottom-left-radius: 4px; /* Instagram style tweak */
}

.my-bubble {
    background: #3797f0; /* Instagram Blueish */
    color: white;
    border-bottom-right-radius: 4px;
}

.chat-input-area {
  background: white;
  padding: 1rem;
  border-top: 1px solid #ddd;
  flex-shrink: 0;
}

.input-form {
    max-width: 1200px;
    margin: 0 auto;
    display: flex;
    gap: 0.5rem;
    align-items: center;
}

.msg-input {
    flex: 1;
    padding: 0.8rem 1.2rem;
    border: 1px solid #ddd;
    border-radius: 24px;
    font-size: 1rem;
    background: #f0f0f0;
}

.msg-input:focus {
    outline: none;
    background: white;
    border-color: #bbb;
}

.send-btn {
    padding: 0 1rem;
    background: none;
    border: none;
    color: #3797f0;
    font-weight: 600;
    font-size: 1rem;
    cursor: pointer;
}

.send-btn:hover {
    color: #0056b3;
}
</style>
