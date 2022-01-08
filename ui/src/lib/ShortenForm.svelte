<script lang="ts">
    let url = ''
    let result = null

    const handleSubmit = async (e) => {
        const res = await fetch('/api/shorten', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ url })
        })

        const json = await res.json()
        result = JSON.stringify(json)
    }
</script>
<form on:submit|preventDefault={handleSubmit}>
    <label for="url">URL</label>
    <input type="text" bind:value={url} placeholder="URL to shorten" required>
    <button type="submit">Shorten</button>
    <pre>{result}</pre>
</form>