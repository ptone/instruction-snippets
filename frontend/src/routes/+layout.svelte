<script lang="ts">
	import { onMount } from 'svelte';
	import { auth, provider } from '$lib/firebase';
	import { signInWithPopup, onAuthStateChanged, signOut } from 'firebase/auth';
	import { authUser } from '$lib/stores/auth';
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import Button from '$lib/components/ui/button/button.svelte';

	let { children } = $props();

	onMount(() => {
		const unsubscribe = onAuthStateChanged(auth, (user) => {
			$authUser = user;
		});
		return unsubscribe;
	});

	async function signIn() {
		try {
			await signInWithPopup(auth, provider);
		} catch (error) {
			console.error('Error signing in', error);
		}
	}

	async function signOutUser() {
		try {
			await signOut(auth);
		} catch (error) {
			console.error('Error signing out', error);
		}
	}
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<header class="flex justify-between items-center p-4 bg-gray-100 border-b">
	<h1 class="text-2xl font-bold">Instruction Snippets</h1>
	<div>
		{#if $authUser}
			<span class="mr-4">Welcome, {$authUser.displayName}</span>
			<a href="/contribute" class="mr-4">
				<Button variant="outline">Add Contribution</Button>
			</a>
			<Button on:click={signOutUser}>Logout</Button>
		{:else}
			<Button on:click={signIn}>Login with Google</Button>
		{/if}
	</div>
</header>

<main class="p-4">
	{@render children?.()}
</main>

