<script lang="ts">
    import { createForm } from 'svelte-forms-lib'
    import * as yup from 'yup'
    import ApiClient from './api';

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
                const res = await new ApiClient().shortenUrl(values.url).then(async (res) => {
                    await delay(1500)
                    return res.json()
                })

                shortUrl = res.short_url
                handleReset()
            } catch (err) {
                console.error(err)
            }
        }
    })
</script>

<form on:submit={handleSubmit}>
    <input name="url" placeholder="URL to shorten" role="textbox" aria-label="url"
        on:keyup={handleChange}
        bind:value={$form.url}
        aria-invalid={!$touched.url ? null : $errors.url ? true : false}>
    <button type="submit" aria-label="shorten" disabled={!$isValid || !$isModified} aria-busy={$isSubmitting}>Shorten</button>
</form>
<pre>{$isValid}</pre>
<a href={shortUrl} target="_blank" aria-busy={$isSubmitting}>
    {#if $isSubmitting}
    Generating link, please wait...
    {:else if shortUrl}
    {shortUrl}
    {/if}
</a>