import { createRouter, createWebHistory } from 'vue-router'
import MapView from '../views/MapView.vue'
import DeviceView from '@/views/DeviceView.vue'
import UserView from '@/views/UserView.vue'
import AuthView from '@/views/AuthView.vue'
import { useMessageStore } from '@/stores/messageStore'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/authStore'
import { useJwtComposable } from '@/composables/useJwtComposable'
import LogoutView from '@/views/LogoutView.vue'
import SvgSpriteView from '@/views/helpers/SvgSpriteView.vue'
import { useAppStore } from '@/stores/appStore'
import SettingsView from '@/views/SettingsView.vue'
import DebugView from '@/views/DebugView.vue'

const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	routes: [
		{ path: '/', name: 'mapView', component: MapView, },
		{ path: '/user', name: 'userView', component: UserView, },
		{ path: '/user/:userID', name: 'userEditView', component: UserView, props: true},
		{ path: '/device', name: 'deviceView', component: DeviceView, },
		{ path: '/device/:deviceID', name: 'deviceEditView', component: DeviceView, props: true},
		{ path: '/login', name: 'loginView', component: AuthView, },
		{ path: '/forgot-password', name: 'forgotPasswordView', component: AuthView, },
		{ path: '/logout', name: 'logoutView', component: LogoutView },
		{ path: '/settings', name: 'settingsView', component: SettingsView },
		{ path: '/debug/:id?', name: 'debugView', component: DebugView, props: true },
		{ path: '/helpers/svg', name: 'viewHelpersSvg', component: SvgSpriteView },		

		// {
		//   path: '/about', name: 'about',
		//   // route level code-splitting
		//   // this generates a separate chunk (About.[hash].js) for this route
		//   // which is lazy-loaded when the route is visited.
		//   component: () => import('../views/AboutView.vue'),
		// },
	],
});

router.beforeEach(async (to, from, next) => {
	const { isAuthenticated } = storeToRefs(useAuthStore());
	const { checkJwtExpiration } = useJwtComposable();

	// Note: We are access the Pinia store within the navigation guard by:
	// useAuthStore() & useMessageStore() to ensure Pinia is properly initialized 

	const authStore = useAuthStore();
	const appStore = useAppStore();
	const messageStore = useMessageStore();

	// Clear flash message
	messageStore.getPersistFlashMessage ?
		messageStore.decreasePersistFlashMessage() :
		messageStore.clearFlashMessage();

	// Redirect to login if not authenticated
	if (!isAuthenticated.value && !['loginView', 'forgotPasswordView'].includes(to.name)) {
		authStore.setRedirectTo(to);
		return next({ name: 'loginView' });
	}

	// Token expiration check
	if (!checkJwtExpiration()) {
		authStore.resetAuthStore();
		appStore.resetAppStore();		
	
		if (!isAuthenticated.value && !['loginView', 'forgotPasswordView'].includes(to.name)) {
			// TODO: Log token has expired - logout
			authStore.setRedirectTo(to);	
			return next({ name: 'loginView' });
		}
	}
	// Continue with the navigation
	next();

});

export default router
