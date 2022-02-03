<script lang="ts">
    import { createForm } from 'svelte-forms-lib'
    import * as yup from 'yup'
    import ApiClient from './api';
    import { delayFor } from './utils'

    let shortUrl = null
    let error = null
    const client = new ApiClient();

    const { form, errors, touched, isValid, isModified, isSubmitting, handleChange, handleSubmit, handleReset } = createForm({
        initialValues: {
            url: ''
        },
        validationSchema: yup.object().shape({
            url: yup.string().required().url()
        }),
        onSubmit: async values => {
            try {
                const res = await client.shortenUrl(values.url).then(async (res) => {
                    await delayFor(1500)
                    return res
                })

                shortUrl = res.short_url
                handleReset()
            } catch (err) {
                error = err
            }
        }
    })
</script>

<style>
    div[role=alert] {
        color: #d63b3b;
    }
</style>

<form on:submit={handleSubmit}>
    <input name="url" placeholder="URL to shorten" role="textbox" aria-label="url"
        on:keyup={handleChange}
        bind:value={$form.url}
        aria-invalid={!$touched.url ? null : $errors.url ? true : false}>
    <button type="submit" aria-label="shorten" disabled={!$isValid || !$isModified} aria-busy={$isSubmitting}>Shorten</button>
</form>
{#if error}
<div role="alert" aria-label="error message">{error}</div>
{/if}
<a href={shortUrl} target="_blank" aria-busy={$isSubmitting}>
    {#if $isSubmitting}
    Generating link, please wait...
    {:else if shortUrl}
    {shortUrl}
    {/if}
</a>
