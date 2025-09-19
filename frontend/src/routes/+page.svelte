<script lang="ts">
	import { onMount, tick } from 'svelte';
	import { marked } from 'marked';
	import {
		collection,
		getDocs,
		doc,
		updateDoc,
		increment,
		writeBatch,
		getDoc,
		setDoc,
		deleteDoc,
		type DocumentData
	} from 'firebase/firestore';
	import { db } from '$lib/firebase';
	import { authUser } from '$lib/stores/auth';
	import * as Card from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { Input } from '$lib/components/ui/input';
	import { Button } from '$lib/components/ui/button';
	import * as ToggleGroup from '$lib/components/ui/toggle-group';
	import { ThumbsUp, ThumbsDown, Copy, X } from 'lucide-svelte';

	let snippets: DocumentData[] = [];
	let searchTerm = '';
	let sortBy = 'newest';
	let activeFilters = new Set<string>();
	let copiedState: { [key: string]: boolean } = {};
	let userVotes: { [key: string]: 'thumbs_up' | 'thumbs_down' | null } = {};

	onMount(async () => {
		const querySnapshot = await getDocs(collection(db, 'snippets'));
		snippets = querySnapshot.docs.map((doc) => {
			const data = doc.data();
			return {
				id: doc.id,
				...data,
				createdAt: data.createdAt ? data.createdAt.toDate() : new Date() // Convert Firestore Timestamp to JS Date
			};
		});

		if ($authUser) {
			for (const snippet of snippets) {
				const voteDocRef = doc(db, 'snippets', snippet.id, 'votes', $authUser.uid);
				const voteDoc = await getDoc(voteDocRef);
				if (voteDoc.exists()) {
					userVotes[snippet.id] = voteDoc.data().vote;
				}
			}
		}
	});

	$: sortedAndFilteredSnippets = (() => {
		let result = snippets.filter((snippet) => {
			const lowerCaseSearch = searchTerm.toLowerCase();
			const inText = snippet.content.toLowerCase().includes(lowerCaseSearch);
			const inLabels = snippet.labels.some((label: string) =>
				label.toLowerCase().includes(lowerCaseSearch)
			);

			const hasActiveLabels =
				activeFilters.size === 0 ||
				[...activeFilters].every((filter) => snippet.labels.includes(filter));

			return (inText || inLabels) && hasActiveLabels;
		});

		result.sort((a, b) => {
			switch (sortBy) {
				case 'newest':
					return b.createdAt - a.createdAt;
				case 'oldest':
					return a.createdAt - b.createdAt;
				case 'most_liked':
					return b.thumbs_up - b.thumbs_down - (a.thumbs_up - a.thumbs_down);
				case 'least_liked':
					return a.thumbs_up - a.thumbs_down - (b.thumbs_up - b.thumbs_down);
				default:
					return 0;
			}
		});

		return result;
	})();

	async function rateSnippet(id: string, rating: 'thumbs_up' | 'thumbs_down') {
		if (!$authUser) {
			// maybe show a message to login
			return;
		}

		const snippetRef = doc(db, 'snippets', id);
		const voteDocRef = doc(db, 'snippets', id, 'votes', $authUser.uid);
		const existingVote = userVotes[id];

		const snippet = snippets.find((s) => s.id === id);
		if (!snippet) return;

		if (existingVote === rating) {
			// Revoke vote
			await deleteDoc(voteDocRef);
			await updateDoc(snippetRef, { [rating]: increment(-1) });
			userVotes[id] = null;
			snippet[rating]--;
		} else if (existingVote) {
			// Change vote
			const batch = writeBatch(db);
			batch.update(snippetRef, {
				[existingVote]: increment(-1),
				[rating]: increment(1)
			});
			batch.set(voteDocRef, { vote: rating });
			await batch.commit();

			userVotes[id] = rating;
			snippet[existingVote]--;
			snippet[rating]++;
		} else {
			// New vote
			await setDoc(voteDocRef, { vote: rating });
			await updateDoc(snippetRef, { [rating]: increment(1) });
			userVotes[id] = rating;
			snippet[rating]++;
		}
		snippets = [...snippets];
	}

	function copyToClipboard(id: string, text: string) {
		navigator.clipboard.writeText(text);
		copiedState[id] = true;
		setTimeout(() => {
			copiedState[id] = false;
			// Use tick to ensure the UI updates before we potentially re-render
			tick();
		}, 2000);
	}

	function toggleFilter(label: string) {
		if (activeFilters.has(label)) {
			activeFilters.delete(label);
		} else {
			activeFilters.add(label);
		}
		activeFilters = new Set(activeFilters);
	}
</script>

<div class="container mx-auto p-4">
	<div class="mb-4 flex items-center justify-between">
		<h1 class="text-2xl font-bold">Instruction Snippets</h1>
	</div>

	<div class="mb-4 flex space-x-4">
		<div class="flex-grow">
			<Input placeholder="Search snippets..." bind:value={searchTerm} />
		</div>
		<ToggleGroup.Root type="single" bind:value={sortBy} class="flex-shrink-0">
			<ToggleGroup.Item value="newest">Newest</ToggleGroup.Item>
			<ToggleGroup.Item value="oldest">Oldest</ToggleGroup.Item>
			<ToggleGroup.Item value="most_liked">Most Liked</ToggleGroup.Item>
			<ToggleGroup.Item value="least_liked">Least Liked</ToggleGroup.Item>
		</ToggleGroup.Root>
	</div>

	<div class="mb-4 flex items-center gap-2">
		{#if activeFilters.size > 0}
			<span class="font-semibold">Includes:</span>
		{/if}
		{#each [...activeFilters] as filter}
			<div
				on:click={() => toggleFilter(filter)}
				on:keydown={(e) => e.key === 'Enter' && toggleFilter(filter)}
				role="button"
				tabindex="0"
				class="cursor-pointer"
			>
				<Badge class="bg-blue-500 text-white hover:bg-blue-600">
					{filter}
					<X class="ml-2 h-4 w-4" />
				</Badge>
			</div>
		{/each}
	</div>

	<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
		{#each sortedAndFilteredSnippets as snippet}
			<Card.Root class="relative">
				<Card.Header>
					<Card.Title>{snippet.title}</Card.Title>
					<div
						class="absolute top-2 right-2"
						on:click={() => copyToClipboard(snippet.id, snippet.content)}
						on:keydown={(e) => e.key === 'Enter' && copyToClipboard(snippet.id, snippet.content)}
						role="button"
						tabindex="0"
					>
						<Button variant="ghost" size="icon">
							<Copy class="h-4 w-4" />
						</Button>
						{#if copiedState[snippet.id]}
							<span
								class="absolute -bottom-8 left-1/2 -translate-x-1/2 rounded-md bg-gray-900 px-2 py-1 text-xs text-white"
								>Copied!</span
							>
						{/if}
					</div>
				</Card.Header>
				<Card.Content>
					<div class="prose dark:prose-invert">{@html marked(snippet.content)}</div>
				</Card.Content>
				<Card.Footer class="flex-col items-start">
					<div class="mt-2 flex flex-wrap">
						{#each snippet.labels as label}
							<div
								on:click={() => toggleFilter(label)}
								on:keydown={(e) => e.key === 'Enter' && toggleFilter(label)}
								role="button"
								tabindex="0"
								class="inline-block cursor-pointer"
							>
								<Badge
									class="mr-2 mb-2"
									variant={activeFilters.has(label) ? 'secondary' : 'default'}>{label}</Badge
								>
							</div>
						{/each}
					</div>
					<div class="flex items-center space-x-2 self-end">
						<Button
							variant="ghost"
							size="icon"
							onclick={() => rateSnippet(snippet.id, 'thumbs_up')}
						>
							<ThumbsUp
								class="h-4 w-4 {userVotes[snippet.id] === 'thumbs_up' ? 'fill-current' : ''}"
							/>
						</Button>
						<span>{snippet.thumbs_up}</span>
						<Button
							variant="ghost"
							size="icon"
							onclick={() => rateSnippet(snippet.id, 'thumbs_down')}
						>
							<ThumbsDown
								class="h-4 w-4 {userVotes[snippet.id] === 'thumbs_down' ? 'fill-current' : ''}"
							/>
						</Button>
						<span>{snippet.thumbs_down}</span>
					</div>
				</Card.Footer>
			</Card.Root>
		{/each}
	</div>
</div>
