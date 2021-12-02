<script context="module">
	import { browser, dev } from '$app/env';

	// we don't need any JS on this page, though we'll load
	// it in dev so that we get hot module replacement...
	export const hydrate = dev;

	// ...but if the client-side router is already loaded
	// (i.e. we came here from elsewhere in the app), use it
	export const router = browser;

	// since there's no dynamic data here, we can prerender
	// it so that it gets served as a static asset in prod
	export const prerender = true;
</script>

<svelte:head>
	<title>About</title>
</svelte:head>

<section class="banner">
	<div class="container">
		<h1>Cielo Verde</h1>
	</div>
</section>

<div class="content">
	<h2>About Us</h2>
	<p>
	Cielo Verde is the premier distributor of medical cannabis products in
	Colombia. With the most loyal customer base and reliable supply chains,
	we ensure the highest quality of product from dirt to door.
	</p>

	<p>
	Started in 2019, Cielo Verde worked closely with regulatory agencies and
	testing facilities to confirm the purity and economic viability of our operation.
	Now, we look to make this product available to you for all of your medical needs.
	</p>

</div>

<style>
	.banner {
		position: relative;
		top: -51px;
		margin-bottom: -51px;
		min-height: 220px;
		padding-top: 48px;
		padding-bottom: 48px;
		background-size: cover;
		background-position: bottom center;
		background-repeat: no-repeat;
		background-image: url('/static/images/background.jpg');
	}
	.banner::before {
		content: "";
		position: absolute;
		top: 0;
		display: block;
		height: 100%;
		width: 100%;
		background-color: rgb(10 8 8 / 45%);
	}
	.banner .container {
		position: absolute;
		top: 50%;
		transform: translateY(-50%);
		left: 0;
		right: 0;
		width: 100%;
	}

	
	.content {
		width: 100%;
		max-width: 660px;
		margin: 0 auto;
		padding: 0 15px;
    	margin: 1rem auto 1rem;
	}

	h2 {
		font-size: 1.7rem;
	}

	p {
		font-size: 18px;
		color: #ffffff;
	}

	@media (min-width: 620px) {
		.banner {
			min-height: 420px;
		}
		h2 {
			font-size: 2rem;
		}
	}
</style>
