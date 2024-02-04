import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import Login from '../views/Login.vue'
import MyProfile from '../views/MyProfile.vue'
import SearchUser from '../views/SearchUser.vue'
import UserProfile from '../views/UserProfile.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: HomeView},
		{path: '/login', component: Login},
		{path: '/my-profile/', component: MyProfile},
		{path: '/users/', component: SearchUser},
		{path: '/users/:username/', component: UserProfile},
	]
});

router.beforeEach((to, from, next) => {
	if (sessionStorage.getItem('auth-token') == null && to.path != '/login'){
		next({
			path: '/login',
			replace: true
		  });
	} else {
		next();
	}
})

export default router
