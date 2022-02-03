export default class ApiClient {
    constructor() {}

    public async shortenUrl(url: string) : Promise<{ short_url: string }> {
        return { short_url: 'http://hpfr.ch/foo_bar_code' }
    }
}