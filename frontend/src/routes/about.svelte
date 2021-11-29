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

<div class="content">
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
	.content {
		width: 100%;
		max-width: var(--column-width);
		margin: var(--column-margin-top) auto 0 auto;
	}

	p {
		font-size: 18px;
		color: #ffffff;
		font-weight: 800;
	}
</style>
