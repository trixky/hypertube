<!-- ========================= SCRIPT -->
<script lang="ts">
	import BlackBox from '../../components/containers/black-box.svelte';
	import Logo from '../../components/layouts/logo.svelte';
	import External from './external.svelte';
</script>

<!-- ========================= HTML -->
<BlackBox>
	<Logo alone />
	<h1 class="mt-2 mb-1 text-2xl text-white">Login</h1>
	<External />
	<div>
		<hr />
		<p class="text-white inline-block">or</p>
		<hr />
	</div>
	<form action="" class="pt-1">
		<label for="email">Email</label>
		<input type="email" placeholder="Email" name="email" />
		<label for="password">Password</label>
		<input type="password" placeholder="Password" name="password" minlength="8" />
		<p class="extra-link mt-2 float-right"><a href="/register">Forgot your password ?</a></p>
		<button class="mt-12 py-2 px-6 bg-blue-300 block m-auto rounded-sm">Login</button>
	</form>
	<p class="extra-link mt-4">
		<a href="/register">Not on Hypertube yet ? <span class="underline">Sign up</span></a>
	</p>
</BlackBox>

<!-- ========================= CSS -->
<style lang="postcss">
	.extra-link {
		@apply text-slate-400 text-sm;
	}

	input {
		@apply w-full p-2 rounded-sm;
	}

	hr {
		@apply border-white/60 inline-block w-28 -translate-y-1 mx-1;
	}

	label {
		@apply block p-2 text-white;
	}
</style>
