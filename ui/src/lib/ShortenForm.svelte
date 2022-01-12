<script lang="ts">
    import { createForm } from 'svelte-forms-lib'
    import * as yup from 'yup'

    let shortUrl = null

    const delay = (ms: number) => new Promise(r => setTimeout(r, ms));
    const { form, errors, touched, isValid, isModified, isSubmitting, handleChange, handleSubmit, handleReset } = createForm({
        initialValues: {
            url: ''
        },
        validationSchema: yup.object().shape({
            url: yup.string().required().url()
        }),
        onSubmit: async values => {
            try {
                const res = await fetch('/api/shorten', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ url: values.url })
                })
                await delay(1500);

                shortUrl = (await res.json()).short_url
                handleReset()
            } catch (err) {
                console.error(err)
            }
        }
    })
</script>

<form on:submit={handleSubmit}>
    <input id="url" name="url" placeholder="URL to shorten"
        on:keyup={handleChange}
        bind:value={$form.url}
        aria-invalid={!$touched.url ? null : $errors.url ? true : false}>
    <button type="submit" disabled={!$isValid || !$isModified} aria-busy={$isSubmitting}>Shorten</button>
</form>
<a href={shortUrl} target="_blank" aria-busy={$isSubmitting}>
    {#if $isSubmitting}
    Generating link, please wait...
    {:else if shortUrl}
    {shortUrl}
    {/if}
</a>
