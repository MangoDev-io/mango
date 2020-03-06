export class Token {
    assetId: string
    creatorAddr: string
    assetName: string
    unitName: string
    totalIssuance: number
    decimals: number
    defaultFrozen: boolean
    url: string
    metadataHash: string
    managerAddr: string
    reserveAddr: string
    freezeAddr: string
    clawbackAddr: string

    constructor(init?: Partial<Token>) {
        Object.assign(this, init)
    }
}
