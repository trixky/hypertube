<!-- ========================= SCRIPT -->
<script lang="ts">
	import BlackBox from '$components/containers/black-box.svelte';
	import Logo from '$components/layouts/logo.svelte';
	import ConfirmationButton from '$components/buttons/confirmation-button.svelte';
	import Warning from '$components/inputs/warning.svelte';
	import Eye from '$components/inputs/eye.svelte';
	import * as cookies from '$utils/cookies';
	import * as sanitizer from '$utils/sanitizer';
	import { uppercase_first_character } from '$utils/str';
	import { encrypt_password } from '$utils/password';
	import { goto } from '$app/navigation';
	import { _ } from 'svelte-i18n';

	let loading = false;

	let registration_attempts = 0;

	let username_blur = false;
	let firstname_blur = false;
	let lastname_blur = false;
	let email_blur = false;
	let password_blur = false;
	let confirm_password_blur = false;

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

	$: warnings =
		username_warning.length ||
		firstname_warning.length ||
		lastname_warning.length ||
		email_warning.length ||
		password_warning.length ||
		confirm_password_warning.length;

	$: disabled =
		warnings > 0 ||
		!username.length ||
		!firstname.length ||
		!lastname.length ||
		!email.length ||
		!password.length ||
		!confirm_password.length;

	let show_password = false;
	$: password_input_type = show_password ? 'text' : 'password';

	let emails_already_in_use: Array<string> = [];

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
			show_password = false;

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
						password: await encrypt_password(password)
					})
				});

				if (res.ok) {
					await res
						.json()
						.then((body) => {
							if (cookies.labels.token in body && cookies.labels.user_info in body) {
								cookies.add_a_cookie(cookies.labels.token, body.token);
								cookies.add_a_cookie(cookies.labels.user_info, body[cookies.labels.user_info]);
								resolve(true);
								goto('/');
							} else {
								notifies_response_warning($_('auth.server_error'));
							}
						})
						.catch(() => {
							notifies_response_warning($_('auth.server_error'));
						});
				} else {
					if (res.status == 409) {
						emails_already_in_use.push(email);
						check_email();
					} else notifies_response_warning($_('auth.server_error'));
				}
				resolve(false);
			}, 1000);
		});
	}

	function notifies_response_warning(warning: string) {
		response_warning = uppercase_first_character(warning);
	}

	// ----------------------------------------------------------------- sanitizing
	function check_username(): boolean {
		response_warning = '';
		if (registration_attempts || username_blur) username_warning = sanitizer.name(username);

		return username_warning.length > 0;
	}
	function check_firstname(): boolean {
		response_warning = '';
		if (registration_attempts || firstname_blur) firstname_warning = sanitizer.name(firstname);

		return firstname_warning.length > 0;
	}
	function check_lastname(): boolean {
		response_warning = '';
		if (registration_attempts || lastname_blur) lastname_warning = sanitizer.name(lastname);

		return lastname_warning.length > 0;
	}
	function check_email(): boolean {
		response_warning = '';
		if (registration_attempts || email_blur) {
			if (emails_already_in_use.includes(email)) email_warning = $_('auth.email_already_in_use');
			else email_warning = sanitizer.email(email);
		}

		return email_warning.length > 0;
	}
	function check_password(
		event: (Event & { currentTarget: EventTarget & HTMLInputElement }) | null = null
	): boolean {
		response_warning = '';
		if (event) password = event.currentTarget.value;

		if (registration_attempts || password_blur) password_warning = sanitizer.password(password);

		return password_warning.length > 0;
	}
	function check_confirm_password(
		event: (Event & { currentTarget: EventTarget & HTMLInputElement }) | null = null
	): boolean {
		response_warning = '';
		if (event) confirm_password = event.currentTarget.value;

		if (registration_attempts || confirm_password_blur)
			confirm_password_warning = sanitizer.confirm_password(password, confirm_password);

		return confirm_password_warning.length > 0;
	}
</script>

<!-- ========================= HTML -->
<BlackBox title={$_('auth.register_header')}>
	<Logo alone />
	<form action="" class="pt-1 w-full">
		<label for="username" class="required">{$_('auth.username')}</label>
		<input
			type="text"
			placeholder={$_('auth.username')}
			name="username"
			bind:value={username}
			on:input={check_username}
			on:blur={() => {
				username_blur = true;
				check_username();
			}}
			disabled={loading}
		/>
		<Warning content={username_warning} color="red" />
		<div class="flex justify-between">
			<div class="pr-2">
				<label for="firstname" class="required">{$_('auth.first_name')}</label>
				<input
					type="text"
					placeholder={$_('auth.first_name')}
					name="firstname"
					bind:value={firstname}
					on:input={check_firstname}
					on:blur={() => {
						firstname_blur = true;
						check_firstname();
					}}
					disabled={loading}
				/>
				<Warning content={firstname_warning} color="red" />
			</div>
			<div class="pl-2">
				<label for="lastname" class="required">{$_('auth.last_name')}</label>
				<input
					type="text"
					placeholder={$_('auth.last_name')}
					name="lastname"
					bind:value={lastname}
					on:input={check_lastname}
					on:blur={() => {
						lastname_blur = true;
						check_lastname();
					}}
					disabled={loading}
				/>
				<Warning content={lastname_warning} color="red" />
			</div>
		</div>
		<label for="email" class="required">{$_('auth.email')}</label>
		<input
			type="email"
			placeholder={$_('auth.email')}
			name="email"
			bind:value={email}
			on:input={check_email}
			on:blur={() => {
				email_blur = true;
				check_email();
			}}
			disabled={loading}
		/>
		<Warning content={email_warning} color="red" />
		<label for="password" class="required">{$_('auth.password')}</label>
		<div class="relative">
			<input
				type={password_input_type}
				placeholder={$_('auth.password')}
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
		</div>
		<Warning content={password_warning} color="red" />
		<label for="confirm password" class="required">{$_('auth.confirm_password')}</label>
		<div class="relative">
			<input
				type={password_input_type}
				placeholder={$_('auth.confirm_password')}
				name="confirm password"
				value={confirm_password}
				on:input={check_confirm_password}
				on:blur={() => {
					confirm_password_blur = true;
					check_confirm_password();
				}}
				disabled={loading}
			/>
			<Eye bind:open={show_password} />
		</div>
		<Warning content={confirm_password_warning} color="red" />
		<ConfirmationButton
			name={$_('auth.register_action')}
			handler={handle_register}
			bind:loading
			bind:disabled
		/>
		<Warning centered content={response_warning} color="red" />
	</form>
	<p class="extra-link mt-4">
		<a href="/login"
			>{$_('auth.member_question')} <span class="underline">{$_('auth.login')}</span></a
		>
	</p>
</BlackBox>

<!-- ========================= CSS -->
<style lang="postcss">
	label {
		@apply ml-2;
	}
</style>
