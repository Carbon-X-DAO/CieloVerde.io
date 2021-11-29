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
	<div class="container">
		<p class="copyright">© CIELO VERDE {year}</p>
	</div>
</footer>

<style>
	main {
		position: relative;
		flex: 1;
		display: flex;
		flex-direction: column;
		width: 100%;
		box-sizing: border-box;
	}

	footer {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		background: #232323;
    	border-top: 2px solid #55736c;
	}

	.copyright {
		font-size: 0.7rem;
	}

	@media (min-width: 480px) {
		footer {
			padding: 10px 0;
		}
	}
</style>
