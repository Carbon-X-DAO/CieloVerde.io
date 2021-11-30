<script>
	import { page } from '$app/stores';
	import logo from '../../images/logos/w-logo4.png'
</script>

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
					<a href="#" class="ghost-btn">About Us</a>
				</div>
			</div>
		</div>
	</div>
</header>

<style>
	.navbar__nav nav {
		display: none;
	}
	@media (min-width: 600px) {
		.navbar__nav nav {
			display: block;
		}
	}

	.navbar,
	.navbar__nav {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.navbar__nav nav + div {
		margin-left: 1.5rem;
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

	nav {
		display: flex;
		justify-content: center;
		--background: rgba(255, 255, 255, 0.7);
	}

	ul {
		position: relative;
		padding: 0;
		margin: 0;
		display: flex;
		justify-content: center;
		align-items: center;
		list-style: none;
	}

	li {
		position: relative;
		height: 100%;
	}

	li.active::before {
		--size: 6px;
		content: '';
		width: 0;
		height: 0;
		position: absolute;
		top: 0;
		left: calc(50% - var(--size));
		border: var(--size) solid transparent;
		border-top: var(--size) solid var(--accent-color);
	}

	nav a {
		display: flex;
		height: 100%;
		align-items: center;
		padding: 1rem 1.5rem;
		color: var(--heading-color);
		font-weight: 700;
		font-size: 0.9rem;
		text-transform: uppercase;
		letter-spacing: 0.1em;
		text-decoration: none;
		transition: color 0.2s linear;
	}

	a:hover {
		color: var(--accent-color);
	}
</style>
