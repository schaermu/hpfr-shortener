export default class ApiClient {
    baseUrl: string = ''

    constructor() {
        this.baseUrl = '/api'
    }

    async shortenUrl(url: string) {
        return fetch(`${this.baseUrl}/shorten`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ url: url })
        })
    }
}