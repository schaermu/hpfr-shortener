const shortenUrl = async (url: string) => fetch('/api/shorten', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
    },
    body: JSON.stringify({ url: url })
})

export {
    shortenUrl
}