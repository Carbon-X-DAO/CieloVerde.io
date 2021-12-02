<script>
	import { page } from '$app/stores';
	import logo from '/static/images/logos/w-logo4.png'
</script>

<svelte:head>
	<link rel="stylesheet" href="/styles/font-awesome.css">
	<link rel="apple-touch-icon" sizes="180x180" href="/images/favicons/apple-touch-icon.png">
	<link rel="icon" type="image/png" sizes="32x32" href="/images/favicons/favicon-32x32.png">
	<link rel="icon" type="image/png" sizes="16x16" href="/images/favicons/favicon-16x16.png">
	<link rel="manifest" href="/images/favicons/site.webmanifest">
	<link rel="mask-icon" href="/images/favicons/safari-pinned-tab.svg" color="#318048">
	<link rel="shortcut icon" href="/images/favicons/favicon.ico">
	<meta name="msapplication-TileColor" content="#da532c">
	<meta name="msapplication-config" content="/images/favicons/browserconfig.xml">
	<meta name="theme-color" content="#ffffff">
</svelte:head>

<header>
	<div class="container">
		<div class="navbar">
			<div class="corner navbar__logo">
				<a href="https://cieloverde.io">
					<img src={logo} alt="Cielo Verde" />
				</a>
			</div>
		
			<div class="navbar__nav">
				<!-- <nav>
					<ul>
						<li class:active={$page.path === '/'}><a sveltekit:prefetch href="/">Home</a></li>
						<li class:active={$page.path === '/about'}><a sveltekit:prefetch href="/about">About</a></li>
						<li class:active={$page.path === '/form'}><a sveltekit:prefetch href="/form">Form</a></li>
					</ul>
				</nav> -->
			
				<div>
					<a href="/about" class="ghost-btn">About Us</a>
				</div>
			</div>
		</div>
	</div>
</header>

<style>
	header {
		position: relative;
		z-index: 2;
	}
	.navbar {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding-top: 3px;
	}

	.corner {
		width: 3em;
		height: 3em;
	}

	.corner a {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 100%;
		height: 100%;
	}

	.navbar__logo img {
		width: 35px;
		height: 35px;
		object-fit: contain;
	}

	a:hover {
		color: var(--accent-color);
	}
</style>
