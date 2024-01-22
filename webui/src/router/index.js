import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import Login from '../views/Login.vue'
import MyProfile from '../views/MyProfile.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: HomeView},
		{path: '/login', component: Login},
		{path: '/users/:username/', component: MyProfile},
		{path: '/some/:id/link', component: HomeView},
	]
})

export default router
