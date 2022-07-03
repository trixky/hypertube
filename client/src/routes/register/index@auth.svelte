<!-- ========================= SCRIPT -->
<script lang="ts">
	import BlackBox from '../../components/containers/black-box.svelte';
	import Logo from '../../components/layouts/logo.svelte';
	import ConfirmationButton from '../../components/buttons/confirmation-button.svelte';
	import Warning from '../../components/inputs/warning.svelte';
	import { uppercase_first_character } from '../../utils/str';
	import Eye from '../../components/inputs/eye.svelte';

	let registration_attempts = 0;

	let username = '';
	let username_warning = '';
	let firstname = '';
	let firstname_warning = '';
	let lastname = '';
	let lastname_warning = '';
	let email = '';
	let email_warning = '';
	let password = '';
	let password_warning = '';
	let confirm_password = '';
	let confirm_password_warning = '';

	let response_warning = '';

	let show_password = false;
	$: password_input_type = show_password ? 'text' : 'password';

	let emails_already_in_use = [];

	function handle_register() {
		return new Promise((resolve) => {
			registration_attempts++;
			let inputs_corrupted = false;

			if (check_username()) inputs_corrupted = true;
			if (check_firstname()) inputs_corrupted = true;
			if (check_lastname()) inputs_corrupted = true;
			if (check_email()) inputs_corrupted = true;
			if (check_password()) inputs_corrupted = true;
			if (check_confirm_password()) inputs_corrupted = true;

			if (inputs_corrupted) return resolve(false);

			setTimeout(async () => {
				const res = await fetch('http://localhost:7070/v1/internal/register', {
					method: 'POST',
					headers: {
						'content-type': 'application/json',
						accept: 'application/json'
					},
					body: JSON.stringify({
						username,
						firstname,
						lastname,
						email,
						password
					})
				});

				if (res.ok) {
					await res
						.json()
						.then((body) => {
							if (body.hasOwnProperty('token')) {
								console.log('------- d:', body.token);
								document.cookie = 'token=' + body.token;
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
					if (res.status == 409) {
						emails_already_in_use.push(email);
						// notifies_response_warning('Email is already in use, please choose another');
					} else {
						notifies_response_warning('An error occured on server side, please try again');
					}
				}

				resolve(false);
			}, 1000);
		});
	}

	function notifies_response_warning(warning: string) {
		response_warning = warning;

		setTimeout(() => {
			response_warning = '';
		}, 5000);
	}

	// ------------------------------------------- checks
	// ---------------------------- check username
	function check_username(): boolean {
		if (registration_attempts > 0) {
			if (username.length == 0) username_warning = 'is missing';
			else if (username.length < 3) username_warning = 'is too short, needs at least 3 characters';
			else if (username.length > 20)
				username_warning = 'is too long, must contain a maximum of 20 characters';
			else username_warning = '';
		}

		return username_warning.length > 0;
	}
	// ---------------------------- check firstname
	function check_firstname(): boolean {
		if (registration_attempts > 0) {
			if (firstname.length == 0) firstname_warning = 'is missing';
			else if (firstname.length < 3)
				firstname_warning = 'is too short, needs at least 3 characters';
			else if (firstname.length > 20)
				firstname_warning = 'is too long, must contain a maximum of 20 characters';
			else firstname_warning = '';
		}

		return firstname_warning.length > 0;
	}
	// ---------------------------- check lastname
	function check_lastname(): boolean {
		if (registration_attempts > 0) {
			if (lastname.length == 0) lastname_warning = 'is missing';
			else if (lastname.length < 3) lastname_warning = 'is too short, needs at least 3 characters';
			else if (lastname.length > 20)
				lastname_warning = 'is too long, must contain a maximum of 20 characters';
			else lastname_warning = '';
		}

		return lastname_warning.length > 0;
	}
	// ---------------------------- check email
	function check_email(): boolean {
		if (registration_attempts > 0) {
			// https://stackoverflow.com/questions/46155/how-can-i-validate-an-email-address-in-javascript
			const regex =
				/[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?/;

			if (email.length == 0) email_warning = 'is missing';
			else if (!regex.test(email)) email_warning = 'is bad formatted';
			else email_warning = '';
		}

		return email_warning.length > 0;
	}
	// ---------------------------- check password
	function check_password(event: any = null): boolean {
		if (event) password = event.target.value;

		if (registration_attempts > 0) {
			let password_warnings = [];

			if (password.length == 0) password_warnings.push('is missing');
			else {
				if (confirm_password == password) confirm_password_warning = '';
				if (password.length < 8)
					password_warnings.push('is too short, needs at least 3 characters');
				else if (password.length > 30)
					password_warnings.push('is too long, must contain a maximum of 20 characters');
				if (!/[a-z]/.test(password)) {
					password_warnings.push('must contain at least one lowercase character (a-z)');
				}
				if (!/[A-Z]/.test(password)) {
					password_warnings.push('must contain at least one uppercase character (A-Z)');
				}
				if (!/\d/.test(password)) {
					password_warnings.push('must contain at least one numeric character (0-9)');
				}
				if (!/[ !@#$%^&*()-=_+[\]{}\\|'\";:/?.>,<`~]/.test(password)) {
					password_warnings.push('must contain at least one specific character (!@#...)');
				}
			}

			let password_warnings_concatenation = password_warnings
				.map((str) => uppercase_first_character(str))
				.join('\n- ');

			if (password_warnings.length > 1)
				password_warnings_concatenation = '- ' + password_warnings_concatenation;

			password_warning = password_warnings_concatenation;
		}

		return password_warning.length > 0;
	}
	// ---------------------------- check confirm password
	function check_confirm_password(event: any = null): boolean {
		if (event) confirm_password = event.target.value;

		if (registration_attempts > 0) {
			if (confirm_password.length == 0) confirm_password_warning = 'is missing';
			else if (confirm_password != password)
				confirm_password_warning = 'passwords must be the same';
			else confirm_password_warning = '';
		}

		return confirm_password_warning.length > 0;
	}
</script>

<!-- ========================= HTML -->
<BlackBox>
	<Logo alone />
	<h1 class="mt-2 mb-1 text-2xl text-white">Register</h1>
	<form action="" class="pt-1">
		<label for="username">Username</label>
		<input
			type="text"
			placeholder="Username"
			name="username"
			bind:value={username}
			on:input={check_username}
		/>
		<Warning content={username_warning} />
		<div class="flex justify-between">
			<div class="pr-2">
				<label for="firstname">Firstname</label>
				<input
					type="text"
					placeholder="Firstname"
					name="firstname"
					bind:value={firstname}
					on:input={check_firstname}
				/>
				<Warning content={firstname_warning} />
			</div>
			<div class="pl-2">
				<label for="lastname">Lastname</label>
				<input
					type="text"
					placeholder="Lastname"
					name="lastname"
					bind:value={lastname}
					on:input={check_lastname}
				/>
				<Warning content={lastname_warning} />
			</div>
		</div>
		<label for="email">Email</label>
		<input
			type="email"
			placeholder="Email"
			name="email"
			bind:value={email}
			on:input={check_email}
		/>
		<Warning content={email_warning} />
		<label for="password">Password</label>
		<input
			type={password_input_type}
			placeholder="Password"
			name="password"
			value={password}
			on:input={check_password}
		/>
		<Eye bind:open={show_password} />
		<Warning content={password_warning} />
		<label for="confirm password">Confirm password</label>
		<input
			type={password_input_type}
			placeholder="Confirm password"
			name="confirm password"
			value={confirm_password}
			on:input={check_confirm_password}
			aria-hidden="false"
		/>
		<Eye bind:open={show_password} />
		<Warning content={confirm_password_warning} />
		<ConfirmationButton name="register" handler={handle_register} />
		<Warning centered content={response_warning} />
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

	label::after {
		@apply content-['*'] text-blue-300;
	}
</style>
