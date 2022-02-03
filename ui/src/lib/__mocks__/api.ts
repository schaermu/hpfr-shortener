export const MOCK_SHORT_URL = 'http://hpfr.ch/foo_bar_code'

export default class ApiClient {
    constructor() {}

    public async shortenUrl(url: string) : Promise<{ short_url: string }> {
        return { short_url: MOCK_SHORT_URL }
    }
}