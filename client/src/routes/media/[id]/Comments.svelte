<!-- ========================= SCRIPT -->
<script lang="ts">
	import { _ } from 'svelte-i18n';
	import { session } from '$app/stores';
	import { fly } from 'svelte/transition';
	import type { MediaComment } from '$types/Media';
	import Warning from '$components/inputs/warning.svelte';
	import { apiMedia } from '$utils/api';

	$: self = $session.user!;

	export let mediaId: number;
	export let list: MediaComment[];

	let cleanComments = list.map((comment) => ({
		...comment,
		id: Number(comment.id),
		user: {
			...comment.user,
			id: Number(comment.user.id)
		},
		date: new Date(comment.date)
	}));

	// Comment
	let loadingComment = false;
	let commentError: string | undefined;
	let commentContent: string | null | undefined;
	$: commentLength = commentContent ? commentContent.length : 0;
	async function postComment() {
		if (!commentContent) {
			commentError = $_('sanitizer.missing');
			return;
		}
		if (commentContent.length < 2) {
			commentError = $_('sanitizer.too_short', { values: { amount: 2 } });
			return;
		}
		if (commentContent.length > 500) {
			commentError = $_('sanitizer.too_long', { values: { amount: 500 } });
			return;
		}
		commentError = undefined;
		loadingComment = true;
		const res = await fetch(apiMedia(`/v1/media/${mediaId}/comment`), {
			method: 'POST',
			credentials: 'include',
			headers: { accept: 'application/json', 'content-type': 'application/json' },
			body: JSON.stringify({ content: commentContent })
		});
		if (res.ok) {
			const body = (await res.json()) as MediaComment;
			const cleanComment = {
				...body,
				date: new Date(body.date)
			};
			commentContent = undefined;
			cleanComments.unshift(cleanComment);
			cleanComments = cleanComments;
		} else {
			commentError = $_('media.comment_fail');
		}
		loadingComment = false;
	}
</script>

<!-- ========================= HTML -->

<div class="my-4">
	<h1 class="text-2xl mb-4">
		{$_('media.comments')}
		{#if cleanComments.length > 0}
			({cleanComments.length})
		{/if}
	</h1>
	{#if cleanComments.length == 0}
		<div>{$_('media.no_comments')}</div>
	{/if}
	<form class="mt-4 mb-6" on:submit|preventDefault={postComment}>
		<textarea
			class="border border-white rounded-md w-full bg-transparent text-white p-4 disabled:opacity-50 transition-all"
			name="comment"
			id="comment"
			rows="4"
			disabled={loadingComment}
			bind:value={commentContent}
			placeholder={$_('media.comment_placeholder')}
		/>
		<div class="flex justify-between">
			<div>
				<Warning
					content={`${commentLength} / 500`}
					color={commentLength > 500 ? 'red' : commentLength > 400 ? 'orange' : 'gray'}
				/>
				{#if commentError}
					<Warning content={commentError} color="red" />
				{/if}
			</div>
			<div class="text-right">
				<button
					class="py-2 px-4 bg-blue-300 text-black mt-1 rounded-sm hover:bg-blue-400 duration-[0.35s] disabled:opacity-50 transition-all"
					disabled={commentLength < 2 || commentLength > 500 || loadingComment}
				>
					{$_('media.post_comment')}
				</button>
			</div>
		</div>
	</form>
	{#each cleanComments as comment (comment.id)}
		<div
			id="comment-{comment.id}"
			class="comment"
			class:self={comment.user.id == self.id}
			in:fly|local
		>
			{#if comment.user.id == self.id}
				<div class="bordered" />
			{/if}
			<div class="comment-header">
				<div>
					<a href="#comment-{comment.id}" class="opacity-60 text-sm mr-2">#{comment.id}</a>
					<a href="/users/{comment.user.id}" class="font-bold">{comment.user.name}</a>
				</div>
				<div class="text-sm">{comment.date.toLocaleString()}</div>
			</div>
			<div class="comment-content">{comment.content}</div>
		</div>
	{/each}
</div>

<!-- ========================= CSS -->
<style lang="postcss">
	.comment {
		@apply mb-4 p-2 border border-stone-400 rounded-md bg-stone-900 relative;
	}

	.comment.self {
		@apply border-transparent overflow-hidden;
		padding: 1px;
	}

	.comment.self .comment-header {
		@apply p-2 pb-0 bg-stone-900 rounded-t-md;
	}
	.comment.self .comment-content {
		@apply p-2 pt-0 bg-stone-900 rounded-b-md;
	}

	.comment.self .bordered {
		@apply absolute top-0 right-0 bottom-0 left-0;
		background: rgb(170, 50, 201);
		background: linear-gradient(to bottom right, rgb(170, 50, 201) 0%, rgba(107, 139, 176, 1) 100%);
	}

	.comment-header {
		@apply flex justify-between w-full relative;
	}

	.comment-content {
		@apply relative;
	}

	.comment-content::before {
		@apply block w-full mb-1;
		content: '';
		height: 1px;
		background: linear-gradient(
			to right,
			rgba(0, 0, 0, 0) 25%,
			rgba(255, 255, 255, 0.8) 50%,
			rgba(0, 0, 0, 0) 75%
		);
	}
</style>
