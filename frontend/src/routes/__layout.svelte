<script>
	import Footer from '$lib/footer/Footer.svelte';
	import Header from '$lib/header/Header.svelte';
	import '../app.css';
</script>

<Header />

<main>
	<slot />
</main>

<Footer />

<style>
	main {
		position: relative;
		flex: 1;
		display: flex;
		flex-direction: column;
		width: 100%;
		box-sizing: border-box;
	}
</style>
