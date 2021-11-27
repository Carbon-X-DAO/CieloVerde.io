<script>
	import Header from '$lib/header/Header.svelte';
	import '../app.css';
	let year = new Date().getFullYear()  // returns the current year
</script>

<Header />

<main>
	<slot />
</main>

<footer>
	<p class="copyright">© CIELO VERDE {year}</p>
</footer>

<style>
	main {
		position: relative;
		flex: 1;
		display: flex;
		flex-direction: column;
		padding: 1rem;
		width: 100%;
		box-sizing: border-box;
	}

	footer {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		padding: 40px;
	}

	footer a {
		font-weight: bold;
	}

	@media (min-width: 480px) {
		footer {
			padding: 40px 0;
		}
	}
</style>
