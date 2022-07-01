<!-- ========================= SCRIPT -->
<script lang="ts">
	import BlackBox from '../../components/containers/black-box.svelte';
	import Logo from '../../components/layouts/logo.svelte';
	import ConfirmationButton from '../../components/buttons/confirmation-button.svelte';

	let result = null;

	async function handle_registerd() {
		const response = await fetch('https://jsonplaceholder.typicode.com/posts', {
			mode: 'no-cors'
		})
		let posts = await response.json()
		console.log(posts)
	}

	async function handle_register4() {
		const response = await fetch('http://localhost:7070/v1/internal/register', {
			method: 'POST',
			mode: 'no-cors'
		})
		let posts = await response.json()
		console.log(posts)
	}

	async function handle_register() {
		console.log('----- 0');
		const res = await fetch('http://localhost:7070/v1/internal/register', {
			method: 'POST',
			headers: {
				'content-type': 'application/json',
				// 'authauth': 'yolo',
				accept: 'application/json'
				// like application/json or text/xml
			},
			body: JSON.stringify({
				toto: 'foo'
			})
		});

		console.log(res);
		console.log(await res.text());
	}
</script>

<!-- ========================= HTML -->
<BlackBox>
	<Logo alone />
	<h1 class="mt-2 mb-1 text-2xl text-white">Register</h1>
	<form action="" class="pt-1">
		<label for="username">Username</label>
		<input type="text" placeholder="Username" name="username" />
		<div class="flex justify-between">
			<div class="pr-2">
				<label for="firstname">Firstname</label>
				<input type="text" placeholder="Firstname" name="firstname" />
			</div>
			<div class="pl-2">
				<label for="lastname">Lastname</label>
				<input type="text" placeholder="Lastname" name="lastname" />
			</div>
		</div>
		<label for="email">Email</label>
		<input type="email" placeholder="Email" name="email" />
		<label for="password">Password</label>
		<input type="password" placeholder="Password" name="password" minlength="8" />
		<label for="confirm password">Confirm password</label>
		<input type="password" placeholder="Confirm password" name="confirm password" minlength="8" />
		<ConfirmationButton name="register" handler={handle_register} />
	</form>
	<p class="extra-link mt-4">
		<a href="/login">Already a member ? <span class="underline">Log in</span></a>
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

	label {
		@apply block p-2 text-white;
	}
</style>
