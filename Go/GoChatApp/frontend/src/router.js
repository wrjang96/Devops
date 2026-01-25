import { createRouter, createWebHistory } from 'vue-router'
import Login from './components/Login.vue'
import ChatRoom from './components/ChatRoom.vue'
import RoomList from './components/RoomList.vue'

const routes = [
    { path: '/', component: Login },
    { path: '/rooms', component: RoomList },
    { path: '/chat', component: ChatRoom },
]

const router = createRouter({
    history: createWebHistory(), //  Vue Router에서 사용되는 History API를 사용하여 URL 변경을 관리합니다.
    routes,
})

export default router
