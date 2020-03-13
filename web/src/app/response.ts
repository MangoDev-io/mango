export class Response {
    status: string
    message: string
    assetId: number
    txHash: string

    constructor(init?: Partial<Response>) {
        Object.assign(this, init)
    }
}
