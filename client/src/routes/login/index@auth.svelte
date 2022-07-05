<!-- ========================= SCRIPT -->
<script lang="ts">
	import BlackBox from '../../components/containers/black-box.svelte';
	import Logo from '../../components/layouts/logo.svelte';
	import External from './external.svelte';
	import ConfirmationButton from '../../components/buttons/confirmation-button.svelte';
	import Eye from '../../components/inputs/eye.svelte';
	import Warning from '../../components/inputs/warning.svelte';
	import * as cookies from '../../utils/cookies';
	import * as sanitzer from '../../utils/sanitizer';
	import { uppercase_first_character } from '../../utils/str';
	import { encrypt_password } from '../../utils/password';

	let loading = false;

	let login_attempts = 0;

	let email_blur = false;
	let password_blur = false;

	let email = '';
	let email_warning = '';
	let password = '';
	let password_warning = '';

	let response_warning = '';

	let show_password = false;
	$: password_input_type = show_password ? 'text' : 'password';

	async function handle_login() {
		return new Promise((resolve) => {
			login_attempts++;
			let inputs_corrupted = false;

			if (check_email()) inputs_corrupted = true;
			if (check_password()) inputs_corrupted = true;

			if (inputs_corrupted) return resolve(false);

			setTimeout(async () => {
				const res = await fetch('http://localhost:7070/v1/internal/login', {
					method: 'POST',
					headers: {
						'content-type': 'application/json',
						accept: 'application/json'
					},
					body: JSON.stringify({
						email,
						password: await encrypt_password(password)
					})
				});

				if (res.ok) {
					await res
						.json()
						.then((body) => {
							if (body.hasOwnProperty('token')) {
								cookies.add_a_cookie(cookies.labels.token, body.token);
								resolve(true);
								window.location.href = window.location.origin + '/';
							} else {
								notifies_response_warning(
									'An error occured on server side with your token, please try again'
								);
							}
							body.token;
						})
						.catch(() => {
							notifies_response_warning(
								'An error occured in the response from the server side, please try again'
							);
						});
				} else {
					if (res.status == 403) notifies_response_warning('Incorrect email and/or password');
					else notifies_response_warning('An error occured on server side, please try again');
				}
				resolve(false);
			}, 1000);
		});
	}

	function notifies_response_warning(warning: string) {
		response_warning = uppercase_first_character(warning);
	}

	// ----------------------------------------------------------------- sanitizing
	function check_email(): boolean {
		response_warning = '';
		if (login_attempts || email_blur) email_warning = sanitzer.email(email);

		return email_warning.length > 0;
	}
	function check_password(event: any = null): boolean {
		response_warning = '';
		if (event) password = event.target.value;

		if (login_attempts || password_blur) password_warning = sanitzer.password(password);

		return password_warning.length > 0;
	}
</script>

<!-- ========================= HTML -->
<BlackBox>
	<Logo alone />
	<h1 class="mt-2 mb-1 text-2xl text-white">Login</h1>
	<External disabled={loading} />
	<div>
		<hr />
		<p class="text-white inline-block">or</p>
		<hr />
	</div>
	<form action="" class="pt-1">
		<label for="email" class="required">Email</label>
		<input
			type="email"
			placeholder="Email"
			name="email"
			bind:value={email}
			on:input={check_email}
			on:blur={() => {
				email_blur = true;
				check_email();
			}}
			disabled={loading}
		/>
		<Warning content={email_warning} />
		<label for="password" class="required">Password</label>
		<input
			type={password_input_type}
			placeholder="Password"
			name="password"
			value={password}
			on:input={check_password}
			on:blur={() => {
				password_blur = true;
				check_password();
			}}
			disabled={loading}
		/>
		<Eye bind:open={show_password} />
		<Warning content={password_warning} />
		<p class="extra-link mt-2 pl-28 mb-4 float-right">
			<a href="/register">Forgot your password ?</a>
		</p>
		<ConfirmationButton name="login" handler={handle_login} bind:loading />
		<Warning centered content={response_warning} />
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
</style>
