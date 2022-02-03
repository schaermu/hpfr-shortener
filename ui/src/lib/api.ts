export default class ApiClient {
    baseUrl: string = ''

    constructor() {
        this.baseUrl = '/api'
    }

    private async call<T>(url: string, method: string, body: any = null) : Promise<T> {
        let request: RequestInit = {
            method: method,
            headers: {
                'Content-Type': 'application/json',
            }
        };

        if (body != null) {
            request.body = JSON.stringify(body)
        }

        return fetch(url, request)
            .then(response => {
                if (!response.ok) {
                    throw new Error(response.statusText)
                }
                return response.json() as Promise<T>
            })
    }

    public async shortenUrl(url: string) : Promise<{ short_url: string }> {
        return this.call(`${this.baseUrl}/shorten`, 'POST', { url })
    }
}