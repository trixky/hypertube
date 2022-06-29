<!-- ========================= SCRIPT -->
<script lang="ts">
	import BlackBox from '../../components/containers/black-box.svelte';
	import Logo from '../../components/layouts/logo.svelte';
	import ConfirmationButton from '../../components/buttons/confirmation-button.svelte';
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
		<ConfirmationButton name="register"/>
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
