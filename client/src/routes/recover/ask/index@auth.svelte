<!-- ========================= SCRIPT -->
<script lang="ts">
	import BlackBox from '$components/containers/black-box.svelte';
	import Logo from '$components/layouts/logo.svelte';
	import ConfirmationButton from '$components/buttons/confirmation-button.svelte';
	import Warning from '$components/inputs/warning.svelte';
	import * as sanitizer from '$utils/sanitizer';
	import { browser } from '$app/env';
	import { _ } from 'svelte-i18n';
	import { apiAuth } from '$utils/api';

	const demo_mode = import.meta.env.VITE_DEMO_MODE;

	let loading = false;

	let login_attempts = 0;

	let email_blur = false;

	let email = '';
	let email_warning = '';

	let response_warning = '';
	let response_update_success = '';

	$: disabled = email_warning.length > 0 || !email.length || !email_blur;

	if (browser) {
		document.onkeypress = function (event) {
			if (event.keyCode == 13) {
				event.preventDefault();

				email_blur = true; // email
				check_email();

				if (!disabled) {
					loading = true;
					handle_ask();
				}
			}
		};
	}

	async function handle_ask() {
		return new Promise((resolve) => {
			login_attempts++;
			let inputs_corrupted = false;

			if (check_email()) inputs_corrupted = true;

			if (inputs_corrupted) {
				loading = false;
				return resolve(false);
			}
			
			if (demo_mode === 'false')
				setTimeout(async () => {
					const res = await fetch(apiAuth('/v1/internal/recover-password'), {
						method: 'POST',
						headers: {
							'content-type': 'application/json',
							accept: 'application/json'
						},
						body: JSON.stringify({
							email
						})
					});

					if (res.ok) {
						await res
							.json()
							.then(() => {
								notifies_response_success($_('auth.recover_mail_sent'));
								loading = false;
								resolve(true);
							})
							.catch(() => {
								notifies_response_warning($_('auth.server_error'));
							});
					} else {
						if (res.status == 404) notifies_response_warning($_('auth.no_user_mail'));
						else notifies_response_warning($_('auth.server_error'));
					}
					loading = false;
					resolve(false);
				}, 1000);
			else {
				setTimeout(async () => {
					loading = false;
					notifies_response_warning($_('auth.register_demo_error'))
					resolve(false);
				}, 1000);
			}
		});
	}

	function notifies_response_warning(warning: string) {
		response_warning = warning;
	}

	function notifies_response_success(success: string) {
		response_update_success = success;
		setTimeout(() => {
			response_update_success = '';
		}, 10000);
	}

	// ----------------------------------------------------------------- sanitizing
	function check_email(): boolean {
		response_warning = '';

		let warning = sanitizer.email(email);

		if (email.length && !warning.length) email_blur = true;
		if (login_attempts || email_blur) email_warning = sanitizer.email(email);

		return email_warning.length > 0;
	}
</script>

<!-- ========================= HTML -->
<svelte:head>
	<title>{$_('title.recover_password')}</title>
</svelte:head>
<BlackBox title={$_('auth.forgot_password_header')}>
	<Logo alone />
	<form action="" class="pt-1 w-full">
		<label for="email" class="required truncate">{$_('auth.email')}</label>
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
		<ConfirmationButton
			name={$_('auth.recover_by_mail')}
			handler={handle_ask}
			bind:loading
			bind:disabled
		/>
		<Warning centered content={response_warning} color="red" />
		<Warning centered content={response_update_success} color="green" />
	</form>
	<p class="extra-link mt-4">
		<a href="/login">{$_('auth.back_to')} <span class="underline">{$_('auth.to_login')}</span></a>
	</p>
</BlackBox>

<!-- ========================= CSS -->
<style lang="postcss">
	label {
		@apply ml-2;
	}
</style>
