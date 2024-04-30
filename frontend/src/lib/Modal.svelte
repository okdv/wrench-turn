<script lang="ts">
	export let showModal: boolean // boolean

	let dialog: HTMLDialogElement; // HTMLDialogElement

	$: if (dialog && showModal) dialog.showModal();

    /* Example usage
    let showModal = false
    ...
    <Modal bind:showModal>
        <h2 slot="header">This is a header</h2>
        <div>
            <p>This is a body</p>
        </div>
        <a href="https://github.com/okdv/wrench-turn">WrenchTurn Github</a>
    </Modal> 
    */
</script>

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-noninteractive-element-interactions -->
<dialog
	bind:this={dialog}
	on:close={() => (showModal = false)}
	on:click|self={() => dialog.close()}
    class="max-w-lg border-r-2 border-none p-0 scale-95 transform-gpu duration-300 ease-in-out backdrop:backdrop-blur-sm backdrop:bg-white/30 open:scale-100"
>
	<!-- svelte-ignore a11y-no-static-element-interactions -->
	<div class="p-4" on:click|stopPropagation>
		<slot name="header" />
		<hr />
		<slot />
		<hr />
		<!-- svelte-ignore a11y-autofocus -->
		<button class="block" autofocus on:click={() => dialog.close()}>close modal</button>
	</div>
</dialog>