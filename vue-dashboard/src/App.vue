<template>
	<main class="dashboard" :class="{'no-side-bar': !showSideBar}">

		<section class="sidebar" v-if="showSideBar">
			<TheSidebar></TheSidebar>
		</section>

		<section class="page">
			<router-view v-slot="{ Component }">
				<component :is="Component" />
			</router-view>
		</section>

	</main>

</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { RouterLink, RouterView, useRoute } from 'vue-router'
import TheSidebar from './components/dashboard/TheSidebar.vue';
import { computed } from 'vue';
const route = useRoute();


const showSideBar = computed(() => {
	return !['loginView', 'forgotPasswordView'].includes(route.name)
});





</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>

.dashboard {
    background-color: $col-slate-50;
    display: grid;
    grid-template-columns: 4rem 1fr;
    min-height: 100vh;
    min-width: 100vw;

    // @include respondMobile($bp-medium) {
    //     grid-template-columns: 17rem 1fr;
    // }
}

.no-side-bar {
	grid-template-columns: 1fr
}
</style>
